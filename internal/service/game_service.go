package service

import (
	"cardgame/internal/game"
	"cardgame/internal/models"
	"cardgame/internal/repository"
	"cardgame/pkg/cards"
	"cardgame/pkg/utils"
	"errors"
	"fmt"
	"time"
)

// GameService handles game business logic
type GameService struct {
	gameRepo   *repository.GameRepository
	playerRepo *repository.GamePlayerRepository
	actionRepo *repository.GameActionRepository
	userRepo   *repository.UserRepository
	engine     *game.Engine
}

// NewGameService creates a new game service
func NewGameService(
	gameRepo *repository.GameRepository,
	playerRepo *repository.GamePlayerRepository,
	actionRepo *repository.GameActionRepository,
	userRepo *repository.UserRepository,
) *GameService {
	return &GameService{
		gameRepo:   gameRepo,
		playerRepo: playerRepo,
		actionRepo: actionRepo,
		userRepo:   userRepo,
		engine:     game.NewEngine(),
	}
}

// CreateGame creates a new game
func (s *GameService) CreateGame(userID int, maxPlayers int) (*models.Game, error) {
	// Generate unique room code
	roomCode, err := utils.GenerateRoomCode(6)
	if err != nil {
		return nil, fmt.Errorf("failed to generate room code: %w", err)
	}

	// Create game
	newGame := &models.Game{
		RoomCode:   roomCode,
		State:      models.StateWaiting,
		Phase:      1,
		CreatedBy:  &userID,
		MaxPlayers: maxPlayers,
	}

	if err := s.gameRepo.Create(newGame); err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	// Add creator as first player
	if err := s.addPlayer(newGame.ID, userID, 0); err != nil {
		return nil, fmt.Errorf("failed to add creator as player: %w", err)
	}

	// Log action
	s.logAction(newGame.ID, userID, models.ActionJoin, nil, 1)

	return newGame, nil
}

// JoinGame adds a player to a game
func (s *GameService) JoinGame(roomCode string, userID int) (*models.Game, error) {
	// Get game
	g, err := s.gameRepo.GetByRoomCode(roomCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	if g == nil {
		return nil, errors.New("game not found")
	}

	// Check if game is joinable
	if g.State != models.StateWaiting {
		return nil, errors.New("game already started")
	}

	// Check if user already in THIS specific game
	existing, err := s.playerRepo.GetByGameAndUser(g.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing player: %w", err)
	}
	if existing != nil {
		// Already in this game — treat as success so clients can redirect to the room
		return g, nil
	}

	// Check if game is full
	count, err := s.playerRepo.CountByGame(g.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count players: %w", err)
	}
	if count >= g.MaxPlayers {
		return nil, errors.New("game is full")
	}

	// Add player
	if err := s.addPlayer(g.ID, userID, count); err != nil {
		return nil, fmt.Errorf("failed to add player: %w", err)
	}

	// Log action
	s.logAction(g.ID, userID, models.ActionJoin, nil, 1)

	return g, nil
}

func (s *GameService) addPlayer(gameID, userID, position int) error {
	player := &models.GamePlayer{
		GameID:       gameID,
		UserID:       userID,
		Position:     position,
		Cards:        []cards.Card{},
		CardCount:    0,
		Status:       models.StatusActive,
		CanCallCheat: true,
	}

	return s.playerRepo.Create(player)
}

// StartGame initializes the game
func (s *GameService) StartGame(gameID, userID int) error {
	// Get game
	g, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}
	if g == nil {
		return errors.New("game not found")
	}

	// Ensure user belongs to this game
	player, err := s.playerRepo.GetByGameAndUser(gameID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify player: %w", err)
	}
	if player == nil {
		return errors.New("only players in this game can start it")
	}

	// Check state
	if g.State != models.StateWaiting {
		return errors.New("game already started")
	}

	// Check minimum players
	count, err := s.playerRepo.CountByGame(gameID)
	if err != nil {
		return fmt.Errorf("failed to count players: %w", err)
	}
	if count < 2 {
		return errors.New("need at least 2 players")
	}

	// Get players
	players, err := s.playerRepo.GetByGame(gameID)
	if err != nil {
		return fmt.Errorf("failed to get players: %w", err)
	}

	// Initialize game
	if err := s.engine.InitializeGame(g, players); err != nil {
		return fmt.Errorf("failed to initialize game: %w", err)
	}

	// Update game
	now := time.Now()
	g.StartedAt = &now
	if err := s.gameRepo.Update(g); err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	// Update players
	for i := range players {
		if err := s.playerRepo.Update(&players[i]); err != nil {
			return fmt.Errorf("failed to update player: %w", err)
		}
	}

	// Log action
	s.logAction(gameID, userID, models.ActionStartGame, nil, 1)

	return nil
}

// PlaceCard places a card on another player's pile
func (s *GameService) PlaceCard(gameID, userID, targetPosition int) error {
	// Get game and players
	g, players, currentPlayer, err := s.getGameState(gameID, userID)
	if err != nil {
		return err
	}

	// Check if it's player's turn
	if g.CurrentPlayerPosition == nil || *g.CurrentPlayerPosition != currentPlayer.Position {
		return errors.New("not your turn")
	}

	// Check if in phase 1
	if g.Phase != 1 {
		return errors.New("can only place cards in phase 1")
	}

	// Check if there's a drawn card waiting to be placed (stored in TableCards[0])
	var cardToPlace *cards.Card
	isPlacingDrawnCard := len(g.TableCards) > 0

	if isPlacingDrawnCard {
		// Placing a drawn card from TableCards
		cardToPlace = &g.TableCards[0]

		// Add the drawn card to current player's pile temporarily for validation
		currentPlayer.Cards = append(currentPlayer.Cards, *cardToPlace)
		currentPlayer.TopCard = cardToPlace
		currentPlayer.CardCount = len(currentPlayer.Cards)
	} else {
		// Placing from starting hand (TopCard)
		if currentPlayer.TopCard == nil {
			return errors.New("no card to place")
		}
		cardToPlace = currentPlayer.TopCard
	}

	// Find and validate target player exists
	targetPlayer, err := game.GetPlayerByPosition(players, targetPosition)
	if err != nil {
		return err
	}

	// Check if +1 rule applies BEFORE placing (for turn continuation logic)
	isPlacingOnSelf := currentPlayer.Position == targetPosition
	var plusOneApplies bool

	if isPlacingOnSelf {
		// When placing on self, check if the card being placed is +1 from the PREVIOUS top card
		// (the card that's currently second-to-top, before the drawn card was added)
		if len(currentPlayer.Cards) >= 2 {
			// The drawn card is already added to Cards (line 223), so Cards[len-2] is the previous top
			previousTopCard := currentPlayer.Cards[len(currentPlayer.Cards)-2]
			plusOneApplies = currentPlayer.TopCard.IsOnePlus(previousTopCard)
		} else {
			// Only one card, can't apply +1 rule
			plusOneApplies = false
		}
	} else {
		// When placing on opponent, check if card is +1 from opponent's current top card
		plusOneApplies = s.engine.DoesPlacementApplyPlusOneRule(*currentPlayer, targetPosition, players)
	}

	// Place card
	s.engine.PlaceCard(currentPlayer, targetPlayer)

	// Clear the drawn card from TableCards if it was a drawn card
	if isPlacingDrawnCard {
		g.TableCards = []cards.Card{}
		if err := s.gameRepo.Update(g); err != nil {
			return fmt.Errorf("failed to update game: %w", err)
		}
	}

	// Update players
	if err := s.playerRepo.Update(currentPlayer); err != nil {
		return fmt.Errorf("failed to update current player: %w", err)
	}
	if err := s.playerRepo.Update(targetPlayer); err != nil {
		return fmt.Errorf("failed to update target player: %w", err)
	}

	// Log action
	s.logAction(gameID, userID, models.ActionPlaceCard, map[string]interface{}{
		"target_position": targetPosition,
		"plus_one_rule":   plusOneApplies,
		"on_self":         isPlacingOnSelf,
	}, g.Phase)

	// NEW RULE: Turn continuation is ONLY based on +1 rule
	// - If +1 applies → Turn CONTINUES (must draw again)
	// - If +1 doesn't apply → Turn ENDS
	if !plusOneApplies {
		// +1 rule doesn't apply → End the turn
		nextPos := s.engine.NextPlayer(*g.CurrentPlayerPosition, len(players))
		g.CurrentPlayerPosition = &nextPos

		if err := s.gameRepo.Update(g); err != nil {
			return fmt.Errorf("failed to update game: %w", err)
		}
	}
	// If plusOneApplies is true, turn continues (don't change CurrentPlayerPosition)

	return nil
}

// DrawCard draws a card from the deck
func (s *GameService) DrawCard(gameID, userID int) (*cards.Card, error) {
	// Get game and players
	g, players, currentPlayer, err := s.getGameState(gameID, userID)
	if err != nil {
		return nil, err
	}

	// Check if it's player's turn
	if g.CurrentPlayerPosition == nil || *g.CurrentPlayerPosition != currentPlayer.Position {
		return nil, errors.New("not your turn")
	}

	// Check if in phase 1
	if g.Phase != 1 {
		return nil, errors.New("can only draw cards in phase 1")
	}

	// Check if must place first
	canPlace, _, _ := s.engine.CanPlaceCard(*currentPlayer, players)
	if canPlace {
		return nil, errors.New("must place card before drawing")
	}

	// Draw card from deck
	card, err := s.engine.DrawCard(g, currentPlayer)
	if err != nil {
		return nil, fmt.Errorf("failed to draw card: %w", err)
	}

	// Store drawn card in TableCards temporarily (Phase 1 doesn't use TableCards)
	// This way the card is saved in the database but NOT in the player's pile
	// When PlaceCard is called, it will take the card from TableCards[0]
	g.TableCards = []cards.Card{*card}

	// Update game in database (deck was modified and drawn card stored)
	if err := s.gameRepo.Update(g); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	// Log action
	s.logAction(gameID, userID, models.ActionDrawCard, map[string]interface{}{
		"card": card,
	}, g.Phase)

	// Return the drawn card to frontend
	// Frontend will display it in the middle for player to choose where to place
	return card, nil
}

func (s *GameService) transitionToPhase2(g *models.Game, players []models.GamePlayer) error {
	if err := s.engine.StartPhase2(g, players); err != nil {
		return fmt.Errorf("failed to start phase 2: %w", err)
	}

	if err := s.gameRepo.Update(g); err != nil {
		return fmt.Errorf("failed to update game for phase 2: %w", err)
	}

	s.logAction(g.ID, 0, models.ActionPhaseChange, map[string]interface{}{
		"phase": 2,
	}, 2)

	return nil
}

// CallCheat calls out a player for cheating
func (s *GameService) CallCheat(gameID, callerID, cheaterID int) error {
	// Get game and players
	g, players, caller, err := s.getGameState(gameID, callerID)
	if err != nil {
		return err
	}

	// Check if caller can call cheat
	if !caller.CanCallCheat {
		return errors.New("you cannot call cheat anymore")
	}

	// Find cheater
	var cheater *models.GamePlayer
	for i := range players {
		if players[i].UserID == cheaterID {
			cheater = &players[i]
			break
		}
	}
	if cheater == nil {
		return errors.New("cheater not found")
	}

	// Validate cheat (simplified - you'd need more context)
	// For now, assume the cheat detection is valid

	// Apply penalty
	penalizedPlayers := s.engine.ApplyPenalty(cheater, players)

	// Update all affected players
	if err := s.playerRepo.Update(cheater); err != nil {
		return fmt.Errorf("failed to update cheater: %w", err)
	}
	for i := range penalizedPlayers {
		if err := s.playerRepo.Update(&penalizedPlayers[i]); err != nil {
			return fmt.Errorf("failed to update penalized player: %w", err)
		}
	}

	// Log actions
	s.logAction(gameID, callerID, models.ActionCallCheat, map[string]interface{}{
		"cheater_id": cheaterID,
	}, g.Phase)
	s.logAction(gameID, cheaterID, models.ActionReceivePenalty, map[string]interface{}{
		"penalty_cards": len(penalizedPlayers),
	}, g.Phase)

	return nil
}

// SkipTurn skips the player's turn by adding the drawn card to their pile and ending the turn
func (s *GameService) SkipTurn(gameID, userID int) error {
	// Get game and players
	g, players, currentPlayer, err := s.getGameState(gameID, userID)
	if err != nil {
		return err
	}

	// Check if it's player's turn
	if g.CurrentPlayerPosition == nil || *g.CurrentPlayerPosition != currentPlayer.Position {
		return errors.New("not your turn")
	}

	// Check if in phase 1
	if g.Phase != 1 {
		return errors.New("can only skip in phase 1")
	}

	// Check if there's a drawn card waiting (stored in TableCards[0])
	if len(g.TableCards) == 0 {
		return errors.New("no card drawn to skip")
	}

	// Get the drawn card
	drawnCard := g.TableCards[0]

	// Add the card to player's pile
	s.engine.AddCardToPlayer(currentPlayer, drawnCard)

	// Clear TableCards
	g.TableCards = []cards.Card{}

	// Update player in database
	if err := s.playerRepo.Update(currentPlayer); err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	// End turn - move to next player
	nextPos := s.engine.NextPlayer(*g.CurrentPlayerPosition, len(players))
	g.CurrentPlayerPosition = &nextPos

	// Update game in database
	if err := s.gameRepo.Update(g); err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	// Log action
	s.logAction(gameID, userID, models.ActionSkipTurn, map[string]interface{}{
		"card_kept": drawnCard,
	}, g.Phase)

	return nil
}

// GetGameState retrieves full game state
func (s *GameService) GetGameState(gameID, userID int) (*models.GameStateResponse, error) {
	g, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	if g == nil {
		return nil, errors.New("game not found")
	}

	players, err := s.playerRepo.GetByGame(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	actions, err := s.actionRepo.GetByGame(gameID, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get actions: %w", err)
	}

	// Populate user info for players
	for i := range players {
		user, _ := s.userRepo.GetByID(players[i].UserID)
		players[i].User = user
	}

	return &models.GameStateResponse{
		Game:      *g,
		Players:   players,
		Actions:   actions,
		DeckCount: len(g.DeckCards),
	}, nil
}

func (s *GameService) getGameState(gameID, userID int) (*models.Game, []models.GamePlayer, *models.GamePlayer, error) {
	g, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get game: %w", err)
	}
	if g == nil {
		return nil, nil, nil, errors.New("game not found")
	}

	players, err := s.playerRepo.GetByGame(gameID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get players: %w", err)
	}

	var currentPlayer *models.GamePlayer
	for i := range players {
		if players[i].UserID == userID {
			currentPlayer = &players[i]
			break
		}
	}
	if currentPlayer == nil {
		return nil, nil, nil, errors.New("you are not in this game")
	}

	return g, players, currentPlayer, nil
}

func (s *GameService) logAction(gameID, userID int, actionType models.ActionType, data map[string]interface{}, phase int) {
	var uid *int
	if userID > 0 {
		uid = &userID
	}

	action := &models.GameAction{
		GameID:     gameID,
		UserID:     uid,
		ActionType: actionType,
		ActionData: data,
		Phase:      &phase,
	}

	s.actionRepo.Create(action)
}

// ListGames lists available games
func (s *GameService) ListGames(state models.GameState, limit, offset int) ([]models.Game, error) {
	if state != "" {
		return s.gameRepo.ListByState(state, limit, offset)
	}
	return s.gameRepo.ListAll(limit, offset)
}
