package repository

import (
	"cardgame/internal/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection with retry logic
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	connStr := cfg.ConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Retry logic to wait for database to be ready
	maxRetries := 30
	retryDelay := time.Second

	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			fmt.Printf("✓ Database connection successful on attempt %d\n", i+1)
			break
		}

		if i < maxRetries-1 {
			fmt.Printf("Database not ready (attempt %d/%d): %v - retrying in %v\n", i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
