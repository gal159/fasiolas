package repository

import (
	"cardgame/internal/models"
	"cardgame/pkg/cards"
	"database/sql"
	"encoding/json"
	"time"
)

// GameRepository handles game data operations
type GameRepository struct {
	db *DB
}

// NewGameRepository creates a new game repository
func NewGameRepository(db *DB) *GameRepository {
	return &GameRepository{db: db}
}

// Create creates a new game
func (r *GameRepository) Create(game *models.Game) error {
	query := `
		INSERT INTO games (room_code, state, phase, current_player_position, trump_suit, 
		                   deck_cards, table_cards, created_by, max_players)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	deckJSON, _ := json.Marshal(game.DeckCards)
	tableJSON, _ := json.Marshal(game.TableCards)

	return r.db.QueryRow(
		query,
		game.RoomCode,
		game.State,
		game.Phase,
		game.CurrentPlayerPosition,
		game.TrumpSuit,
		deckJSON,
		tableJSON,
		game.CreatedBy,
		game.MaxPlayers,
	).Scan(&game.ID, &game.CreatedAt, &game.UpdatedAt)
}

// GetByID retrieves a game by ID
func (r *GameRepository) GetByID(id int) (*models.Game, error) {
	query := `
		SELECT id, room_code, state, phase, current_player_position, trump_suit, 
		       deck_cards, table_cards, created_by, created_at, updated_at, 
		       started_at, finished_at, winner_id, max_players
		FROM games
		WHERE id = $1
	`

	game := &models.Game{}
	var deckJSON, tableJSON []byte
	var currentPos sql.NullInt64
	var createdBy sql.NullInt64
	var winnerID sql.NullInt64
	var trumpSuit sql.NullString
	var startedAt sql.NullTime
	var finishedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&game.ID,
		&game.RoomCode,
		&game.State,
		&game.Phase,
		&currentPos,
		&trumpSuit,
		&deckJSON,
		&tableJSON,
		&createdBy,
		&game.CreatedAt,
		&game.UpdatedAt,
		&startedAt,
		&finishedAt,
		&winnerID,
		&game.MaxPlayers,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Assign nullable/pointer fields
	if currentPos.Valid {
		v := int(currentPos.Int64)
		game.CurrentPlayerPosition = &v
	}
	if trumpSuit.Valid {
		v := trumpSuit.String
		game.TrumpSuit = &v
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		game.CreatedBy = &v
	}
	if winnerID.Valid {
		v := int(winnerID.Int64)
		game.WinnerID = &v
	}
	if startedAt.Valid {
		v := startedAt.Time
		game.StartedAt = &v
	}
	if finishedAt.Valid {
		v := finishedAt.Time
		game.FinishedAt = &v
	}

	if err := json.Unmarshal(deckJSON, &game.DeckCards); err != nil {
		game.DeckCards = []cards.Card{}
	}
	if err := json.Unmarshal(tableJSON, &game.TableCards); err != nil {
		game.TableCards = []cards.Card{}
	}

	return game, nil
}

// GetByRoomCode retrieves a game by room code
func (r *GameRepository) GetByRoomCode(roomCode string) (*models.Game, error) {
	query := `
		SELECT id, room_code, state, phase, current_player_position, trump_suit, 
		       deck_cards, table_cards, created_by, created_at, updated_at, 
		       started_at, finished_at, winner_id, max_players
		FROM games
		WHERE room_code = $1
	`

	game := &models.Game{}
	var deckJSON, tableJSON []byte
	var currentPos sql.NullInt64
	var createdBy sql.NullInt64
	var winnerID sql.NullInt64
	var trumpSuit sql.NullString
	var startedAt sql.NullTime
	var finishedAt sql.NullTime

	err := r.db.QueryRow(query, roomCode).Scan(
		&game.ID,
		&game.RoomCode,
		&game.State,
		&game.Phase,
		&currentPos,
		&trumpSuit,
		&deckJSON,
		&tableJSON,
		&createdBy,
		&game.CreatedAt,
		&game.UpdatedAt,
		&startedAt,
		&finishedAt,
		&winnerID,
		&game.MaxPlayers,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if currentPos.Valid {
		v := int(currentPos.Int64)
		game.CurrentPlayerPosition = &v
	}
	if trumpSuit.Valid {
		v := trumpSuit.String
		game.TrumpSuit = &v
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		game.CreatedBy = &v
	}
	if winnerID.Valid {
		v := int(winnerID.Int64)
		game.WinnerID = &v
	}
	if startedAt.Valid {
		v := startedAt.Time
		game.StartedAt = &v
	}
	if finishedAt.Valid {
		v := finishedAt.Time
		game.FinishedAt = &v
	}

	if err := json.Unmarshal(deckJSON, &game.DeckCards); err != nil {
		game.DeckCards = []cards.Card{}
	}
	if err := json.Unmarshal(tableJSON, &game.TableCards); err != nil {
		game.TableCards = []cards.Card{}
	}

	return game, nil
}

// Update updates a game
func (r *GameRepository) Update(game *models.Game) error {
	query := `
		UPDATE games
		SET state = $1, phase = $2, current_player_position = $3, trump_suit = $4,
		    deck_cards = $5, table_cards = $6, updated_at = $7, started_at = $8,
		    finished_at = $9, winner_id = $10
		WHERE id = $11
	`

	deckJSON, _ := json.Marshal(game.DeckCards)
	tableJSON, _ := json.Marshal(game.TableCards)

	var currentPos *int = game.CurrentPlayerPosition
	var trumpSuit *string = game.TrumpSuit

	_, err := r.db.Exec(
		query,
		game.State,
		game.Phase,
		currentPos,
		trumpSuit,
		deckJSON,
		tableJSON,
		time.Now(),
		game.StartedAt,
		game.FinishedAt,
		game.WinnerID,
		game.ID,
	)

	return err
}

// Delete deletes a game
func (r *GameRepository) Delete(id int) error {
	query := `DELETE FROM games WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ListByState retrieves games by state
func (r *GameRepository) ListByState(state models.GameState, limit, offset int) ([]models.Game, error) {
	query := `
		SELECT id, room_code, state, phase, current_player_position, trump_suit, 
		       deck_cards, table_cards, created_by, created_at, updated_at, 
		       started_at, finished_at, winner_id, max_players
		FROM games
		WHERE state = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, state, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return r.scanGames(rows)
}

// ListAll retrieves all games with pagination
func (r *GameRepository) ListAll(limit, offset int) ([]models.Game, error) {
	query := `
		SELECT id, room_code, state, phase, current_player_position, trump_suit, 
		       deck_cards, table_cards, created_by, created_at, updated_at, 
		       started_at, finished_at, winner_id, max_players
		FROM games
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return r.scanGames(rows)
}

func (r *GameRepository) scanGames(rows *sql.Rows) ([]models.Game, error) {
	var games []models.Game

	for rows.Next() {
		var game models.Game
		var deckJSON, tableJSON []byte
		var currentPos sql.NullInt64
		var createdBy sql.NullInt64
		var winnerID sql.NullInt64
		var trumpSuit sql.NullString
		var startedAt sql.NullTime
		var finishedAt sql.NullTime

		err := rows.Scan(
			&game.ID,
			&game.RoomCode,
			&game.State,
			&game.Phase,
			&currentPos,
			&trumpSuit,
			&deckJSON,
			&tableJSON,
			&createdBy,
			&game.CreatedAt,
			&game.UpdatedAt,
			&startedAt,
			&finishedAt,
			&winnerID,
			&game.MaxPlayers,
		)
		if err != nil {
			return nil, err
		}

		if currentPos.Valid {
			v := int(currentPos.Int64)
			game.CurrentPlayerPosition = &v
		}
		if trumpSuit.Valid {
			v := trumpSuit.String
			game.TrumpSuit = &v
		}
		if createdBy.Valid {
			v := int(createdBy.Int64)
			game.CreatedBy = &v
		}
		if winnerID.Valid {
			v := int(winnerID.Int64)
			game.WinnerID = &v
		}
		if startedAt.Valid {
			v := startedAt.Time
			game.StartedAt = &v
		}
		if finishedAt.Valid {
			v := finishedAt.Time
			game.FinishedAt = &v
		}

		game.DeckCards = []cards.Card{}
		if err := json.Unmarshal(deckJSON, &game.DeckCards); err != nil {
			game.DeckCards = []cards.Card{}
		}
		if err := json.Unmarshal(tableJSON, &game.TableCards); err != nil {
			game.TableCards = []cards.Card{}
		}

		games = append(games, game)
	}

	return games, rows.Err()
}
