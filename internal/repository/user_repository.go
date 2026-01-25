package repository

import (
	"cardgame/internal/models"
	"database/sql"
	"time"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, username, oauth_provider, oauth_id, role, avatar_url, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		user.Email,
		user.Username,
		user.OAuthProvider,
		user.OAuthID,
		user.Role,
		user.AvatarURL,
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, email, username, oauth_provider, oauth_id, role, avatar_url, 
		       created_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.Role,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByOAuth retrieves a user by OAuth provider and ID
func (r *UserRepository) GetByOAuth(provider, oauthID string) (*models.User, error) {
	query := `
		SELECT id, email, username, oauth_provider, oauth_id, role, avatar_url, 
		       created_at, updated_at, last_login
		FROM users
		WHERE oauth_provider = $1 AND oauth_id = $2
	`

	user := &models.User{}
	err := r.db.QueryRow(query, provider, oauthID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.Role,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, username, oauth_provider, oauth_id, role, avatar_url, 
		       created_at, updated_at, last_login
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.Role,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2, role = $3, avatar_url = $4, 
		    last_login = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.Exec(
		query,
		user.Email,
		user.Username,
		user.Role,
		user.AvatarURL,
		user.LastLogin,
		time.Now(),
		user.ID,
	)

	return err
}

// UpdateLastLogin updates the last login time
func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `
		UPDATE users
		SET last_login = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.Exec(query, now, now, userID)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves all users with pagination
func (r *UserRepository) List(limit, offset int) ([]models.User, error) {
	query := `
		SELECT id, email, username, oauth_provider, oauth_id, role, avatar_url, 
		       created_at, updated_at, last_login
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.OAuthProvider,
			&user.OAuthID,
			&user.Role,
			&user.AvatarURL,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}
