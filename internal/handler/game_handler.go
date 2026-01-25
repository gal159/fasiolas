package handler

import (
	"cardgame/internal/middleware"
	"cardgame/internal/models"
	"cardgame/internal/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GameHandler handles game endpoints
type GameHandler struct {
	gameService *service.GameService
	extService  *service.ExternalAPIService
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService *service.GameService, extService *service.ExternalAPIService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		extService:  extService,
	}
}

// CreateGame godoc
// @Summary Create a new game
// @Tags games
// @Security BearerAuth
// @Param request body models.CreateGameRequest true "Create Game Request"
// @Success 201 {object} models.Game
// @Router /api/v1/games [post]
func (h *GameHandler) CreateGame(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[Game] Creating game for user %d with maxPlayers %d\n", userID, req.MaxPlayers)
	game, err := h.gameService.CreateGame(userID, req.MaxPlayers)
	if err != nil {
		fmt.Printf("[Game] ❌ CreateGame failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[Game] ✅ CreateGame success: ID=%d, creator=%d\n", game.ID, userID)
	c.JSON(http.StatusCreated, game)
}

// JoinGame godoc
// @Summary Join a game
// @Tags games
// @Security BearerAuth
// @Param request body models.JoinGameRequest true "Join Game Request"
// @Success 200 {object} models.Game
// @Router /api/v1/games/join [post]
func (h *GameHandler) JoinGame(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := h.gameService.JoinGame(req.RoomCode, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

// StartGame godoc
// @Summary Start a game
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/games/{id}/start [post]
func (h *GameHandler) StartGame(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	if err := h.gameService.StartGame(gameID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "game started"})
}

// GetGameState godoc
// @Summary Get game state
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} models.GameStateResponse
// @Router /api/v1/games/{id} [get]
func (h *GameHandler) GetGameState(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	// Get user role
	userRole, _ := middleware.GetUserRole(c)

	state, err := h.gameService.GetGameState(gameID, userID, userRole)
	if err != nil {
		fmt.Printf("[Game] ❌ GetGameState failed for game %d, user %d: %v\n", gameID, userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[Game] ✅ GetGameState success for game %d: %d players (role: %s)\n", gameID, len(state.Players), userRole)
	c.JSON(http.StatusOK, state)
}

// GetGameStateSpectator godoc
// @Summary Get game state as spectator (view-only)
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} models.GameStateResponse
// @Router /api/v1/games/{id}/spectate [get]
func (h *GameHandler) GetGameStateSpectator(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	state, err := h.gameService.GetGameStateForSpectator(gameID)
	if err != nil {
		fmt.Printf("[Game] ❌ GetGameStateSpectator failed for game %d, user %d: %v\n", gameID, userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[Game] ✅ GetGameStateSpectator success for game %d (viewer: %d)\n", gameID, userID)
	c.JSON(http.StatusOK, state)
}

// PlaceCard godoc
// @Summary Place a card
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Param request body models.PlaceCardRequest true "Place Card Request"
// @Success 200 {object} map[string]string
// @Router /api/v1/games/{id}/place [post]
func (h *GameHandler) PlaceCard(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	var req models.PlaceCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body", "details": err.Error()})
		return
	}

	// Validate the target position (position 0 is valid, positions are 0-indexed)
	if req.TargetPlayerPosition < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_player_position must be >= 0"})
		return
	}

	if err := h.gameService.PlaceCard(gameID, userID, req.TargetPlayerPosition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "card placed"})
}

// DrawCard godoc
// @Summary Draw a card
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/games/{id}/draw [post]
func (h *GameHandler) DrawCard(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	card, err := h.gameService.DrawCard(gameID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "card drawn", "card": card})
}

// CallCheat godoc
// @Summary Call out a player for cheating
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Param cheater_id query int true "Cheater User ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/games/{id}/cheat [post]
func (h *GameHandler) CallCheat(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	cheaterID, err := strconv.Atoi(c.Query("cheater_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cheater ID"})
		return
	}

	if err := h.gameService.CallCheat(gameID, userID, cheaterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cheat called"})
}

// SkipTurn godoc
// @Summary Skip turn
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/games/{id}/skip [post]
func (h *GameHandler) SkipTurn(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gameID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	if err := h.gameService.SkipTurn(gameID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "turn skipped"})
}

// ListGames godoc
// @Summary List games
// @Tags games
// @Security BearerAuth
// @Param state query string false "Game State (waiting, phase1, phase2, finished)"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.Game
// @Router /api/v1/games [get]
func (h *GameHandler) ListGames(c *gin.Context) {
	state := models.GameState(c.DefaultQuery("state", ""))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	games, err := h.gameService.ListGames(state, limit, offset)
	if err != nil {
		// Log the error for debugging
		println("[ERROR] ListGames failed:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null if no games
	if games == nil {
		games = []models.Game{}
	}

	c.JSON(http.StatusOK, games)
}

// GetCardImage godoc
// @Summary Get card image from external API
// @Tags external
// @Param suit query string true "Card Suit"
// @Param rank query string true "Card Rank"
// @Success 200 {object} service.CardImage
// @Router /api/v1/external/card-image [get]
func (h *GameHandler) GetCardImage(c *gin.Context) {
	suit := c.Query("suit")
	rank := c.Query("rank")

	if suit == "" || rank == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "suit and rank required"})
		return
	}

	image, err := h.extService.GetCardImage(suit, rank)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// GetGameStats godoc
// @Summary Get game statistics
// @Tags external
// @Success 200 {object} service.CardStats
// @Router /api/v1/external/stats [get]
func (h *GameHandler) GetGameStats(c *gin.Context) {
	stats, err := h.extService.GetGameStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
