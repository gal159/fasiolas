package repository

import (
	"cardgame/internal/models"
	"encoding/json"
)

// GameActionRepository handles game action data operations
type GameActionRepository struct {
	db *DB
}

// NewGameActionRepository creates a new game action repository
func NewGameActionRepository(db *DB) *GameActionRepository {
	return &GameActionRepository{db: db}
}

// Create creates a new game action
func (r *GameActionRepository) Create(action *models.GameAction) error {
	query := `
		INSERT INTO game_actions (game_id, user_id, action_type, action_data, phase)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, timestamp
	`

	var actionDataJSON []byte
	if action.ActionData != nil {
		actionDataJSON, _ = json.Marshal(action.ActionData)
	}

	return r.db.QueryRow(
		query,
		action.GameID,
		action.UserID,
		action.ActionType,
		actionDataJSON,
		action.Phase,
	).Scan(&action.ID, &action.Timestamp)
}

// GetByGame retrieves all actions for a game
func (r *GameActionRepository) GetByGame(gameID int, limit int) ([]models.GameAction, error) {
	query := `
		SELECT id, game_id, user_id, action_type, action_data, phase, timestamp
		FROM game_actions
		WHERE game_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, gameID, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var actions []models.GameAction
	for rows.Next() {
		var action models.GameAction
		var actionDataJSON []byte

		err := rows.Scan(
			&action.ID,
			&action.GameID,
			&action.UserID,
			&action.ActionType,
			&actionDataJSON,
			&action.Phase,
			&action.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if len(actionDataJSON) > 0 {
			json.Unmarshal(actionDataJSON, &action.ActionData)
		}

		actions = append(actions, action)
	}

	return actions, rows.Err()
}

// GetByUser retrieves all actions by a user
func (r *GameActionRepository) GetByUser(userID int, limit int) ([]models.GameAction, error) {
	query := `
		SELECT id, game_id, user_id, action_type, action_data, phase, timestamp
		FROM game_actions
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var actions []models.GameAction
	for rows.Next() {
		var action models.GameAction
		var actionDataJSON []byte

		err := rows.Scan(
			&action.ID,
			&action.GameID,
			&action.UserID,
			&action.ActionType,
			&actionDataJSON,
			&action.Phase,
			&action.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		if len(actionDataJSON) > 0 {
			json.Unmarshal(actionDataJSON, &action.ActionData)
		}

		actions = append(actions, action)
	}

	return actions, rows.Err()
}
