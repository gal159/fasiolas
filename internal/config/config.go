package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	OAuth    OAuthConfig
	JWT      JWTConfig
	CORS     CORSConfig
	External ExternalAPIConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port int
	Env  string // "development" or "production"
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// OAuthConfig holds OAuth2 provider configurations
type OAuthConfig struct {
	Google  OAuthProvider
	GitHub  OAuthProvider
	Discord OAuthProvider
}

// OAuthProvider holds individual OAuth provider config
type OAuthProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scope        []string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// ExternalAPIConfig holds external API configuration
type ExternalAPIConfig struct {
	DeckOfCardsAPIURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvInt("SERVER_PORT", 8080),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "fasiolas_game"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		OAuth: OAuthConfig{
			Google: OAuthProvider{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
				Scope:        []string{"openid", "email", "profile"},
			},
			GitHub: OAuthProvider{
				ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
				ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/api/v1/auth/github/callback"),
				Scope:        []string{"user:email"},
			},
			Discord: OAuthProvider{
				ClientID:     getEnv("DISCORD_CLIENT_ID", ""),
				ClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("DISCORD_REDIRECT_URL", "http://localhost:8080/api/v1/auth/discord/callback"),
				Scope:        []string{"identify", "email"},
			},
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080")),
		},
		External: ExternalAPIConfig{
			DeckOfCardsAPIURL: getEnv("DECK_OF_CARDS_API_URL", "https://deckofcardsapi.com/api/deck"),
		},
	}

	// Validate critical configuration
	// In development mode, allow any non-empty OAuth values (for testing)
	// In production, require at least one real OAuth provider to be configured
	if config.Server.Env == "production" {
		hasValidGoogle := config.OAuth.Google.ClientID != "" && config.OAuth.Google.ClientID != "default-client-id"
		hasValidGitHub := config.OAuth.GitHub.ClientID != "" && config.OAuth.GitHub.ClientID != "default-client-id"
		hasValidDiscord := config.OAuth.Discord.ClientID != "" && config.OAuth.Discord.ClientID != "default-client-id"

		if !hasValidGoogle && !hasValidGitHub && !hasValidDiscord {
			return nil, fmt.Errorf("in production mode, at least one real OAuth provider must be configured")
		}
	}

	return config, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer with a default value
func getEnvInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnvDuration gets an environment variable as a duration with a default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(strValue)
	if err != nil {
		return defaultValue
	}
	return duration
}

// parseOrigins parses comma-separated origins
func parseOrigins(originsStr string) []string {
	origins := strings.Split(originsStr, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	return origins
}

// DSN returns the PostgreSQL connection string
func (db DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.DBName,
		db.SSLMode,
	)
}

// ConnectionString returns the PostgreSQL connection string (alias for DSN)
func (db DatabaseConfig) ConnectionString() string {
	return db.DSN()
}
