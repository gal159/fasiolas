package models

import (
	"cardgame/pkg/cards"
	"time"
)

// UserRole represents user access level
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RolePlayer    UserRole = "player"
	RoleSpectator UserRole = "spectator"
)

// User represents a user in the system
type User struct {
	ID            int        `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	Username      string     `json:"username" db:"username"`
	OAuthProvider string     `json:"oauth_provider" db:"oauth_provider"`
	OAuthID       string     `json:"oauth_id" db:"oauth_id"`
	Role          UserRole   `json:"role" db:"role"`
	AvatarURL     *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin     *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// GameState represents the state of a game
type GameState string

const (
	StateWaiting   GameState = "waiting"
	StatePhase1    GameState = "phase1"
	StatePhase2    GameState = "phase2"
	StateFinished  GameState = "finished"
	StateCancelled GameState = "cancelled"
)

// Game represents a game session
type Game struct {
	ID                    int          `json:"id" db:"id"`
	RoomCode              string       `json:"room_code" db:"room_code"`
	State                 GameState    `json:"state" db:"state"`
	Phase                 int          `json:"phase" db:"phase"`
	CurrentPlayerPosition *int         `json:"current_player_position,omitempty" db:"current_player_position"`
	TrumpSuit             *string      `json:"trump_suit,omitempty" db:"trump_suit"`
	DeckCards             []cards.Card `json:"deck_cards" db:"deck_cards"`
	TableCards            []cards.Card `json:"table_cards" db:"table_cards"`
	CreatedBy             *int         `json:"created_by,omitempty" db:"created_by"`
	CreatedAt             time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at" db:"updated_at"`
	StartedAt             *time.Time   `json:"started_at,omitempty" db:"started_at"`
	FinishedAt            *time.Time   `json:"finished_at,omitempty" db:"finished_at"`
	WinnerID              *int         `json:"winner_id,omitempty" db:"winner_id"`
	MaxPlayers            int          `json:"max_players" db:"max_players"`
}

// PlayerStatus represents a player's status in a game
type PlayerStatus string

const (
	StatusActive     PlayerStatus = "active"
	StatusEliminated PlayerStatus = "eliminated"
	StatusWinner     PlayerStatus = "winner"
	StatusLeft       PlayerStatus = "left"
)

// GamePlayer represents a player in a game
type GamePlayer struct {
	ID           int          `json:"id" db:"id"`
	GameID       int          `json:"game_id" db:"game_id"`
	UserID       int          `json:"user_id" db:"user_id"`
	Position     int          `json:"position" db:"position"`
	Cards        []cards.Card `json:"cards" db:"cards"`
	TopCard      *cards.Card  `json:"top_card,omitempty" db:"top_card"`
	CardCount    int          `json:"card_count" db:"card_count"`
	Status       PlayerStatus `json:"status" db:"status"`
	CanCallCheat bool         `json:"can_call_cheat" db:"can_call_cheat"`
	JoinedAt     time.Time    `json:"joined_at" db:"joined_at"`

	// Joined data (not in DB)
	User *User `json:"user,omitempty" db:"-"`
}

// ActionType represents types of game actions
type ActionType string

const (
	ActionJoin           ActionType = "join"
	ActionLeave          ActionType = "leave"
	ActionStartGame      ActionType = "start_game"
	ActionDrawCard       ActionType = "draw_card"
	ActionPlaceCard      ActionType = "place_card"
	ActionCallCheat      ActionType = "call_cheat"
	ActionReceivePenalty ActionType = "receive_penalty"
	ActionPhaseChange    ActionType = "phase_change"
	ActionTakeTableCard  ActionType = "take_table_card"
	ActionSkipTurn       ActionType = "skip_turn"
	ActionWin            ActionType = "win"
	ActionLose           ActionType = "lose"
)

// GameAction represents an action taken in a game
type GameAction struct {
	ID         int                    `json:"id" db:"id"`
	GameID     int                    `json:"game_id" db:"game_id"`
	UserID     *int                   `json:"user_id,omitempty" db:"user_id"`
	ActionType ActionType             `json:"action_type" db:"action_type"`
	ActionData map[string]interface{} `json:"action_data,omitempty" db:"action_data"`
	Phase      *int                   `json:"phase,omitempty" db:"phase"`
	Timestamp  time.Time              `json:"timestamp" db:"timestamp"`
}

// CreateGameRequest represents a request to create a game
type CreateGameRequest struct {
	MaxPlayers int `json:"max_players" binding:"required,min=2,max=8"`
}

// JoinGameRequest represents a request to join a game
type JoinGameRequest struct {
	RoomCode string `json:"room_code" binding:"required"`
}

// PlaceCardRequest represents a request to place a card
type PlaceCardRequest struct {
	TargetPlayerPosition int `json:"target_player_position"`
}

// GameStateResponse represents the full game state for a player
type GameStateResponse struct {
	Game      Game         `json:"game"`
	Players   []GamePlayer `json:"players"`
	Actions   []GameAction `json:"recent_actions"`
	DeckCount int          `json:"deck_count"`
}
