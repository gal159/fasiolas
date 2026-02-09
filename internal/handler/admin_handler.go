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

// AdminHandler handles admin endpoints
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// ListUsers godoc
// @Summary List all users (admin only)
// @Tags admin
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.User
// @Router /api/v1/admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.adminService.ListUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if users == nil {
		users = []models.User{}
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a specific user (admin only)
// @Tags admin
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /api/v1/admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.adminService.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserRole godoc
// @Summary Update user role (admin only)
// @Tags admin
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body map[string]string true "Role Update Request"
// @Success 200 {object} models.User
// @Router /api/v1/admin/users/{id}/role [put]
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	adminID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, ok := req["role"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role field required"})
		return
	}

	// Validate role
	if role != "player" && role != "spectator" && role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role. must be one of: player, spectator, admin"})
		return
	}

	user, err := h.adminService.UpdateUserRole(userID, models.UserRole(role))
	if err != nil {
		fmt.Printf("[Admin] ❌ UpdateUserRole failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	fmt.Printf("[Admin] ✅ User %d role updated to %s by admin %d\n", userID, role, adminID)
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user (admin only)
// @Tags admin
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	adminID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Prevent admin from deleting themselves
	if uint64(userID) == uint64(adminID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete your own account"})
		return
	}

	err = h.adminService.DeleteUser(userID)
	if err != nil {
		fmt.Printf("[Admin] ❌ DeleteUser failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[Admin] ✅ User %d deleted by admin %d\n", userID, adminID)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// GetStats godoc
// @Summary Get system statistics (admin only)
// @Tags admin
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/stats [get]
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
