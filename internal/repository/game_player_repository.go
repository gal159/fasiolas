package repository

import (
	"cardgame/internal/models"
	"cardgame/pkg/cards"
	"database/sql"
	"encoding/json"
)

// GamePlayerRepository handles game player data operations
type GamePlayerRepository struct {
	db *DB
}

// NewGamePlayerRepository creates a new game player repository
func NewGamePlayerRepository(db *DB) *GamePlayerRepository {
	return &GamePlayerRepository{db: db}
}

// Create creates a new game player
func (r *GamePlayerRepository) Create(gp *models.GamePlayer) error {
	cardsJSON, _ := json.Marshal(gp.Cards)
	topCardJSON, _ := json.Marshal(gp.TopCard)

	query := `
		INSERT INTO game_players (game_id, user_id, position, cards, top_card,
		                          card_count, status, can_call_cheat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, joined_at
	`
	return r.db.QueryRow(
		query,
		gp.GameID,
		gp.UserID,
		gp.Position,
		cardsJSON,
		topCardJSON,
		gp.CardCount,
		gp.Status,
		gp.CanCallCheat,
	).Scan(&gp.ID, &gp.JoinedAt)
}

// GetByID retrieves a game player by ID
func (r *GamePlayerRepository) GetByID(id int) (*models.GamePlayer, error) {
	gp := &models.GamePlayer{}
	var cardsJSON, topCardJSON []byte

	query := `
		SELECT id, game_id, user_id, position, cards, top_card, card_count,
		       status, can_call_cheat, joined_at
		FROM game_players
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&gp.ID,
		&gp.GameID,
		&gp.UserID,
		&gp.Position,
		&cardsJSON,
		&topCardJSON,
		&gp.CardCount,
		&gp.Status,
		&gp.CanCallCheat,
		&gp.JoinedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cardsJSON, &gp.Cards)
	json.Unmarshal(topCardJSON, &gp.TopCard)

	return gp, nil
}

// GetByGameAndUser retrieves a game player by game ID and user ID
func (r *GamePlayerRepository) GetByGameAndUser(gameID, userID int) (*models.GamePlayer, error) {
	gp := &models.GamePlayer{}
	var cardsJSON, topCardJSON []byte

	query := `
		SELECT id, game_id, user_id, position, cards, top_card, card_count,
		       status, can_call_cheat, joined_at
		FROM game_players
		WHERE game_id = $1 AND user_id = $2
	`
	err := r.db.QueryRow(query, gameID, userID).Scan(
		&gp.ID,
		&gp.GameID,
		&gp.UserID,
		&gp.Position,
		&cardsJSON,
		&topCardJSON,
		&gp.CardCount,
		&gp.Status,
		&gp.CanCallCheat,
		&gp.JoinedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cardsJSON, &gp.Cards)
	json.Unmarshal(topCardJSON, &gp.TopCard)

	return gp, nil
}

// GetByGame retrieves all players in a game
func (r *GamePlayerRepository) GetByGame(gameID int) ([]models.GamePlayer, error) {
	query := `
		SELECT id, game_id, user_id, position, cards, top_card, card_count,
		       status, can_call_cheat, joined_at
		FROM game_players
		WHERE game_id = $1
		ORDER BY position
	`
	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.GamePlayer

	for rows.Next() {
		var gp models.GamePlayer
		var cardsJSON, topCardJSON []byte

		err := rows.Scan(
			&gp.ID,
			&gp.GameID,
			&gp.UserID,
			&gp.Position,
			&cardsJSON,
			&topCardJSON,
			&gp.CardCount,
			&gp.Status,
			&gp.CanCallCheat,
			&gp.JoinedAt,
		)
		if err != nil {
			return nil, err
		}

		gp.Cards = []cards.Card{}
		json.Unmarshal(cardsJSON, &gp.Cards)
		json.Unmarshal(topCardJSON, &gp.TopCard)

		players = append(players, gp)

	}

	return players, rows.Err()
}

// Update updates a game player
func (r *GamePlayerRepository) Update(gp *models.GamePlayer) error {
	cardsJSON, _ := json.Marshal(gp.Cards)
	topCardJSON, _ := json.Marshal(gp.TopCard)

	query := `
		UPDATE game_players
		SET cards = $1, top_card = $2, card_count = $3, status = $4, can_call_cheat = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(
		query,
		cardsJSON,
		topCardJSON,
		gp.CardCount,
		gp.Status,
		gp.CanCallCheat,
		gp.ID,
	)

	return err
}

// Delete deletes a game player
func (r *GamePlayerRepository) Delete(id int) error {
	query := `DELETE FROM game_players WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CountByGame counts players in a game
func (r *GamePlayerRepository) CountByGame(gameID int) (int, error) {
	query := `SELECT COUNT(*) FROM game_players WHERE game_id = $1`
	var count int
	err := r.db.QueryRow(query, gameID).Scan(&count)
	return count, err
}

// GetUserActiveGame retrieves the active game for a user (if any)
// Returns the game player record if user is in an active game (waiting or playing)
func (r *GamePlayerRepository) GetUserActiveGame(userID int) (*models.GamePlayer, error) {
	gp := &models.GamePlayer{}
	var cardsJSON, topCardJSON []byte

	query := `
		SELECT gp.id, gp.game_id, gp.user_id, gp.position, gp.cards, gp.top_card, 
		       gp.card_count, gp.status, gp.can_call_cheat, gp.joined_at
		FROM game_players gp
		JOIN games g ON gp.game_id = g.id
		WHERE gp.user_id = $1 
		  AND g.state IN ('waiting', 'playing')
		  AND gp.status = 'active'
		LIMIT 1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&gp.ID,
		&gp.GameID,
		&gp.UserID,
		&gp.Position,
		&cardsJSON,
		&topCardJSON,
		&gp.CardCount,
		&gp.Status,
		&gp.CanCallCheat,
		&gp.JoinedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cardsJSON, &gp.Cards)
	json.Unmarshal(topCardJSON, &gp.TopCard)

	return gp, nil
}
