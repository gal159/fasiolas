package service

import (
	"cardgame/internal/config"
	"cardgame/internal/models"
	"cardgame/internal/repository"
	"cardgame/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo      *repository.UserRepository
	cfg           *config.Config
	googleConfig  *oauth2.Config
	githubConfig  *oauth2.Config
	discordConfig *oauth2.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		googleConfig: &oauth2.Config{
			ClientID:     cfg.OAuth.Google.ClientID,
			ClientSecret: cfg.OAuth.Google.ClientSecret,
			RedirectURL:  cfg.OAuth.Google.RedirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		},
		githubConfig: &oauth2.Config{
			ClientID:     cfg.OAuth.GitHub.ClientID,
			ClientSecret: cfg.OAuth.GitHub.ClientSecret,
			RedirectURL:  cfg.OAuth.GitHub.RedirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
		discordConfig: &oauth2.Config{
			ClientID:     cfg.OAuth.Discord.ClientID,
			ClientSecret: cfg.OAuth.Discord.ClientSecret,
			RedirectURL:  cfg.OAuth.Discord.RedirectURL,
			Scopes:       []string{"identify", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
	}
}

// GetAuthURL returns the OAuth2 authorization URL for a provider
func (s *AuthService) GetAuthURL(provider, state string) (string, error) {
	var oauthConfig *oauth2.Config

	var opts []oauth2.AuthCodeOption
	switch provider {
	case "google":
		oauthConfig = s.googleConfig
		fmt.Printf("Google OAuth Config - ClientID: %s, RedirectURL: %s\n", oauthConfig.ClientID, oauthConfig.RedirectURL)
		// Force Google to show the account selector instead of auto-login with the last session
		opts = append(opts, oauth2.SetAuthURLParam("prompt", "select_account"))
	case "github":
		oauthConfig = s.githubConfig
	case "discord":
		oauthConfig = s.discordConfig
	default:
		return "", errors.New("unsupported OAuth provider")
	}

	if oauthConfig.ClientID == "" {
		return "", fmt.Errorf("OAuth provider %s is not configured (missing ClientID)", provider)
	}

	authURL := oauthConfig.AuthCodeURL(state, opts...)
	fmt.Printf("Generated OAuth URL: %s\n", authURL)
	return authURL, nil
}

// HandleCallback processes OAuth2 callback
func (s *AuthService) HandleCallback(ctx context.Context, provider, code string) (*models.User, string, error) {
	var oauthConfig *oauth2.Config

	switch provider {
	case "google":
		oauthConfig = s.googleConfig
	case "github":
		oauthConfig = s.githubConfig
	case "discord":
		oauthConfig = s.discordConfig
	default:
		return nil, "", errors.New("unsupported OAuth provider")
	}

	// Exchange code for token
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info
	userInfo, err := s.getUserInfo(ctx, provider, token)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists
	user, err := s.userRepo.GetByOAuth(provider, userInfo.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Create user if doesn't exist
	if user == nil {
		fmt.Printf("[Auth] Creating new user from OAuth: email=%s, provider=%s\n", userInfo.Email, provider)
		user = &models.User{
			Email:         userInfo.Email,
			Username:      userInfo.Username,
			OAuthProvider: provider,
			OAuthID:       userInfo.ID,
			Role:          models.RolePlayer,
			AvatarURL:     userInfo.AvatarURL,
		}

		if err := s.userRepo.Create(user); err != nil {
			fmt.Printf("[Auth] ❌ Failed to create user: %v\n", err)
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}
		fmt.Printf("[Auth] ✅ User created successfully with ID: %d\n", user.ID)
	} else {
		fmt.Printf("[Auth] User already exists: id=%d, email=%s\n", user.ID, user.Email)
		// Update last login
		if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
			fmt.Printf("[Auth] ❌ Failed to update last login: %v\n", err)
			return nil, "", fmt.Errorf("failed to update last login: %w", err)
		}
	}

	// Generate JWT
	jwtToken, err := utils.GenerateJWT(
		user.ID,
		user.Email,
		user.Username,
		string(user.Role),
		s.cfg.JWT.Secret,
		s.cfg.JWT.Expiration,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return user, jwtToken, nil
}

type oauthUserInfo struct {
	ID        string
	Email     string
	Username  string
	AvatarURL *string
}

func (s *AuthService) getUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*oauthUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	var url string
	switch provider {
	case "google":
		url = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "github":
		url = "https://api.github.com/user"
	case "discord":
		url = "https://discord.com/api/users/@me"
	default:
		return nil, errors.New("unsupported provider")
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return s.parseUserInfo(provider, body)
}

func (s *AuthService) parseUserInfo(provider string, data []byte) (*oauthUserInfo, error) {
	var result oauthUserInfo

	switch provider {
	case "google":
		var googleUser struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		if err := json.Unmarshal(data, &googleUser); err != nil {
			return nil, err
		}
		result.ID = googleUser.ID
		result.Email = googleUser.Email
		result.Username = googleUser.Name
		result.AvatarURL = &googleUser.Picture

	case "github":
		var githubUser struct {
			ID        int64  `json:"id"`
			Email     string `json:"email"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
		}
		if err := json.Unmarshal(data, &githubUser); err != nil {
			return nil, err
		}
		result.ID = fmt.Sprintf("%d", githubUser.ID)
		result.Email = githubUser.Email
		if githubUser.Name != "" {
			result.Username = githubUser.Name
		} else {
			result.Username = githubUser.Login
		}
		result.AvatarURL = &githubUser.AvatarURL

		// GitHub may not provide email in the main response
		if result.Email == "" {
			result.Email = fmt.Sprintf("%s@github.user", githubUser.Login)
		}

	case "discord":
		var discordUser struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			Username      string `json:"username"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		}
		if err := json.Unmarshal(data, &discordUser); err != nil {
			return nil, err
		}
		result.ID = discordUser.ID
		result.Email = discordUser.Email
		result.Username = discordUser.Username + "#" + discordUser.Discriminator
		if discordUser.Avatar != "" {
			avatarURL := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordUser.ID, discordUser.Avatar)
			result.AvatarURL = &avatarURL
		}
	}

	return &result, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return utils.ValidateJWT(tokenString, s.cfg.JWT.Secret)
}

// RefreshToken generates a new JWT token
func (s *AuthService) RefreshToken(userID int) (string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.userRepo.Update(user); err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	// Generate new JWT
	return utils.GenerateJWT(
		user.ID,
		user.Email,
		user.Username,
		string(user.Role),
		s.cfg.JWT.Secret,
		s.cfg.JWT.Expiration,
	)
}

// GetUserByID fetches a user from the database by ID
func (s *AuthService) GetUserByID(userID int) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}
