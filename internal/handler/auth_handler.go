package handler

import (
	"cardgame/internal/middleware"
	"cardgame/internal/service"
	"cardgame/pkg/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
	stateStore  map[string]bool // Temporary state storage (in production, use Redis/database)
	stateMutex  sync.Mutex      // Protect state store from concurrent access
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		stateStore:  make(map[string]bool),
	}
}

// GetAuthURL godoc
// @Summary Get OAuth authorization URL
// @Tags auth
// @Param provider path string true "OAuth Provider (google, github, discord)"
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/{provider} [get]
func (h *AuthHandler) GetAuthURL(c *gin.Context) {
	provider := c.Param("provider")

	// Generate state token for CSRF protection
	state, err := utils.GenerateSecureToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	// Store state in memory (primary storage)
	h.stateMutex.Lock()
	h.stateStore[state] = true
	h.stateMutex.Unlock()
	fmt.Printf("[OAuth] State generated and stored in memory: %s\n", state)

	// Also set state as cookie (backup for cross-domain scenarios)
	c.SetCookie("oauth_state", state, 300, "/", "localhost", false, true)

	url, err := h.authService.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return OAuth URL
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// HandleCallback godoc
// @Summary Handle OAuth callback
// @Tags auth
// @Param provider path string true "OAuth Provider (google, github, discord)"
// @Param code query string true "Authorization Code"
// @Param state query string true "State Token"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/{provider}/callback [get]
func (h *AuthHandler) HandleCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	fmt.Printf("[OAuth] Callback received - Provider: %s, State: %s, Code: %s\n", provider, state, code)

	// Validate required parameters
	if code == "" {
		fmt.Println("[OAuth] ❌ Missing authorization code")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing authorization code",
		})
		return
	}

	if state == "" {
		fmt.Println("[OAuth] ❌ Missing state parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing state parameter",
		})
		return
	}

	// Validate state parameter using multiple strategies
	isStateValid := false

	// Strategy 1: Check in-memory state store (most reliable)
	h.stateMutex.Lock()
	if h.stateStore[state] {
		fmt.Printf("[OAuth] ✅ State validated from memory store\n")
		delete(h.stateStore, state) // Use state only once
		isStateValid = true
	}
	h.stateMutex.Unlock()

	// Strategy 2: Fallback to cookie validation (if memory store fails)
	if !isStateValid {
		cookieState, err := c.Cookie("oauth_state")
		if err == nil && cookieState == state {
			fmt.Printf("[OAuth] ✅ State validated from cookie\n")
			isStateValid = true
		}
	}

	// Strategy 3: Check URL parameter consistency
	if !isStateValid {
		// Additional safety check - state must be non-empty and reasonable length
		if len(state) > 10 && len(state) < 1000 {
			fmt.Printf("[OAuth] ⚠️  State parameter exists but not found in store. Attempting validation anyway.\n")
			// For development, we can be more lenient, but log it
		}
	}

	if !isStateValid {
		fmt.Println("[OAuth] ❌ VALIDATION FAILED - Invalid state parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state parameter",
		})
		return
	}

	// Clear state cookie after successful validation
	c.SetCookie("oauth_state", "", -1, "/", "localhost", false, true)

	// Handle OAuth callback
	fmt.Printf("[OAuth] State validated successfully, exchanging code for token\n")
	_, token, err := h.authService.HandleCallback(c.Request.Context(), provider, code)
	if err != nil {
		fmt.Printf("[OAuth] ❌ HandleCallback error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[OAuth] ✅ Token generated successfully, redirecting to frontend\n")
	// Redirect to frontend with token as URL parameter
	// Frontend will store it and redirect to dashboard
	redirectURL := "http://localhost:3000/auth/callback?token=" + token
	c.Redirect(302, redirectURL)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	token, err := h.authService.RefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	// Fetch user from database to ensure they actually exist
	user, err := h.authService.GetUserByID(userID)
	if err != nil || user == nil {
		fmt.Printf("[Auth] ❌ User not found in database: ID=%d, error=%v\n", userID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	fmt.Printf("[Auth] ✅ Returning user profile: id=%d, email=%s\n", user.ID, user.Email)
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
	})
}
