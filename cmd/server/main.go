package main

import (
	"cardgame/internal/config"
	"cardgame/internal/handler"
	"cardgame/internal/middleware"
	"cardgame/internal/repository"
	"cardgame/internal/service"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// @title Fasiolas Card Game API
// @version 1.0
// @description REST API for Fasiolas card game with OAuth2 authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@fasiolas.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// ...existing code...

	fmt.Println("=== Fasiolas Card Game Server Starting ===")
	fmt.Printf("Current time: %v\n", time.Now())
	fmt.Printf("Environment: %s\n", os.Getenv("SERVER_ENV"))
	fmt.Printf("Database host: %s\n", os.Getenv("DB_HOST"))

	// Load configuration
	fmt.Println("Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("✓ Config loaded - ENV: %s, DB Host: %s, Port: %d\n", cfg.Server.Env, cfg.Database.Host, cfg.Server.Port)

	// Initialize database with timeout
	fmt.Println("Connecting to database...")
	fmt.Printf("Connection string: host=%s port=%d user=%s dbname=%s\n",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.DBName)

	db, err := repository.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✓ Database connected successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)
	playerRepo := repository.NewGamePlayerRepository(db)
	actionRepo := repository.NewGameActionRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	gameService := service.NewGameService(gameRepo, playerRepo, actionRepo, userRepo)
	adminService := service.NewAdminService(userRepo, gameRepo, playerRepo)
	extService := service.NewExternalAPIService(cfg.External.DeckOfCardsAPIURL)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	gameHandler := handler.NewGameHandler(gameService, extService)
	adminHandler := handler.NewAdminHandler(adminService)

	// Setup Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "fasiolas-card-game",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.GET("/:provider", authHandler.GetAuthURL)
			auth.GET("/:provider/callback", authHandler.HandleCallback)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Auth protected routes
			authProtected := protected.Group("/auth")
			{
				authProtected.GET("/profile", authHandler.GetProfile)
				authProtected.POST("/refresh", authHandler.RefreshToken)
			}

			// Game routes - View for all roles (player, admin, spectator)
			games := protected.Group("/games")
			{
				// List and view games - available to all authenticated users
				games.GET("", gameHandler.ListGames)
				games.GET("/:id", gameHandler.GetGameState)

				// Spectator-only viewing endpoint
				games.GET("/:id/spectate", gameHandler.GetGameStateSpectator)
			}

			// Game actions - Only players and admins can play
			gameActions := protected.Group("/games")
			gameActions.Use(middleware.RequireRole("player", "admin"))
			{
				gameActions.POST("", gameHandler.CreateGame)
				gameActions.POST("/join", gameHandler.JoinGame)
				gameActions.POST("/:id/start", gameHandler.StartGame)
				gameActions.POST("/:id/place", gameHandler.PlaceCard)
				gameActions.POST("/:id/draw", gameHandler.DrawCard)
				gameActions.POST("/:id/skip", gameHandler.SkipTurn)
				gameActions.POST("/:id/cheat", gameHandler.CallCheat)
			}

			// External API routes
			external := protected.Group("/external")
			{
				external.GET("/card-image", gameHandler.GetCardImage)
				external.GET("/stats", gameHandler.GetGameStats)
			}

			// Admin routes (admin role required)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", adminHandler.ListUsers)
				admin.GET("/users/:id", adminHandler.GetUser)
				admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
				admin.DELETE("/users/:id", adminHandler.DeleteUser)
				admin.GET("/stats", adminHandler.GetStats)
			}
		}
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("✓ Server starting on %s", addr)
	log.Printf("✓ Environment: %s", cfg.Server.Env)
	log.Printf("✓ API Documentation: http://localhost:%d/api/v1", cfg.Server.Port)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
