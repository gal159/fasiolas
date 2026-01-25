package service

import (
	"cardgame/internal/models"
	"cardgame/internal/repository"
	"fmt"
)

// AdminService handles admin operations
type AdminService struct {
	userRepo   *repository.UserRepository
	gameRepo   *repository.GameRepository
	playerRepo *repository.GamePlayerRepository
}

// NewAdminService creates a new admin service
func NewAdminService(
	userRepo *repository.UserRepository,
	gameRepo *repository.GameRepository,
	playerRepo *repository.GamePlayerRepository,
) *AdminService {
	return &AdminService{
		userRepo:   userRepo,
		gameRepo:   gameRepo,
		playerRepo: playerRepo,
	}
}

// ListUsers retrieves all users with pagination
func (s *AdminService) ListUsers(limit, offset int) ([]models.User, error) {
	return s.userRepo.List(limit, offset)
}

// GetUser retrieves a user by ID
func (s *AdminService) GetUser(userID int) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

// UpdateUserRole updates a user's role
func (s *AdminService) UpdateUserRole(userID int, newRole models.UserRole) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	user.Role = newRole
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user and all their related data
func (s *AdminService) DeleteUser(userID int) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Delete user's game participations first (foreign key constraint)
	if err := s.playerRepo.DeleteByUserID(uint64(userID)); err != nil {
		return fmt.Errorf("failed to delete user's game participations: %w", err)
	}

	// Delete the user
	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetSystemStats returns system statistics
func (s *AdminService) GetSystemStats() (map[string]interface{}, error) {
	// Get total users
	users, err := s.userRepo.List(9999, 0) // Get all users (practical limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Count by role
	var playerCount, spectatorCount, adminCount int
	for _, user := range users {
		switch user.Role {
		case models.RolePlayer:
			playerCount++
		case models.RoleSpectator:
			spectatorCount++
		case models.RoleAdmin:
			adminCount++
		}
	}

	stats := map[string]interface{}{
		"total_users": len(users),
		"players":     playerCount,
		"spectators":  spectatorCount,
		"admins":      adminCount,
	}

	return stats, nil
}
