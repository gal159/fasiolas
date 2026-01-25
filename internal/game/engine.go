package game

import (
	"cardgame/internal/models"
	"cardgame/pkg/cards"
	"errors"
	"fmt"
)

var (
	ErrGameNotFound       = errors.New("game not found")
	ErrGameFull           = errors.New("game is full")
	ErrGameAlreadyStarted = errors.New("game already started")
	ErrNotYourTurn        = errors.New("not your turn")
	ErrInvalidMove        = errors.New("invalid move")
	ErrMustPlaceFirst     = errors.New("must place card before drawing")
	ErrNoValidPlacement   = errors.New("no valid placement available")
	ErrCannotCallCheat    = errors.New("cannot call cheat")
	ErrPlayerNotFound     = errors.New("player not found")
)

// Engine handles game logic
type Engine struct{}

// NewEngine creates a new game engine
func NewEngine() *Engine {
	return &Engine{}
}

// InitializeGame sets up a new game with shuffled deck
func (e *Engine) InitializeGame(game *models.Game, players []models.GamePlayer) error {
	// Create and shuffle deck
	deck := cards.NewDeck()
	deck.Shuffle()

	// Deal one card to each player
	for i := range players {
		if card := deck.Draw(); card != nil {
			players[i].Cards = []cards.Card{*card}
			players[i].TopCard = card
			players[i].CardCount = 1
		}
	}

	// Store remaining deck
	game.DeckCards = deck.Cards
	game.TableCards = []cards.Card{}
	game.Phase = 1
	game.State = models.StatePhase1

	// Find player with lowest card to start
	lowestValue := 15
	lowestPosition := 0
	for _, player := range players {
		if player.TopCard != nil && player.TopCard.Value < lowestValue {
			lowestValue = player.TopCard.Value
			lowestPosition = player.Position
		}
	}
	game.CurrentPlayerPosition = &lowestPosition

	return nil
}

// CanPlaceCard checks if a player can place their top card on any other player's pile (+1 rule)
// Returns whether can place, the target player position, and whether it's on self
func (e *Engine) CanPlaceCard(currentPlayer models.GamePlayer, allPlayers []models.GamePlayer) (bool, int, bool) {
	if currentPlayer.TopCard == nil {
		return false, -1, false
	}

	// FIRST: Check if can place on any OTHER player's top card
	for _, targetPlayer := range allPlayers {
		if targetPlayer.Position == currentPlayer.Position {
			continue // Skip self for now
		}
		if targetPlayer.TopCard != nil && currentPlayer.TopCard.IsOnePlus(*targetPlayer.TopCard) {
			return true, targetPlayer.Position, false
		}
	}

	// SECOND: If no other player can receive, check if can place on OWN pile (self)
	// Can place on self only if card is +1 to current top card of own pile
	if currentPlayer.TopCard != nil && len(currentPlayer.Cards) > 1 {
		cardBelow := currentPlayer.Cards[len(currentPlayer.Cards)-2]
		if currentPlayer.TopCard.IsOnePlus(cardBelow) {
			return true, currentPlayer.Position, true
		}
	}

	return false, -1, false
}

// CanPlaceCardOnTarget checks if a player can place their top card on a SPECIFIC target player
// Returns true if the placement is valid according to +1 rule
func (e *Engine) CanPlaceCardOnTarget(currentPlayer models.GamePlayer, targetPosition int, allPlayers []models.GamePlayer) bool {
	if currentPlayer.TopCard == nil {
		return false
	}

	// NEW RULE: Player can place card on ANY player (self or opponent)
	// No +1 validation required for placement
	// The +1 rule is ONLY used to determine turn continuation

	// Check if target player exists
	for _, targetPlayer := range allPlayers {
		if targetPlayer.Position == targetPosition {
			// Target exists, placement is always valid
			return true
		}
	}

	return false // Target player doesn't exist
}

// DoesPlacementApplyPlusOneRule checks if placing currentPlayer's top card on target follows +1 rule
// This determines if the turn should continue (must draw again) or end
func (e *Engine) DoesPlacementApplyPlusOneRule(currentPlayer models.GamePlayer, targetPosition int, allPlayers []models.GamePlayer) bool {
	if currentPlayer.TopCard == nil {
		return false
	}

	// Find target player
	for _, targetPlayer := range allPlayers {
		if targetPlayer.Position == targetPosition {
			// Check if current card is +1 from target's top card
			if targetPlayer.TopCard != nil && currentPlayer.TopCard.IsOnePlus(*targetPlayer.TopCard) {
				return true
			}
			break
		}
	}

	return false
}

// Phase1TurnState represents the state during a Phase 1 turn
type Phase1TurnState struct {
	MustDraw    bool // If true, player must draw a card (placed on self, or can't place)
	PlacedCards []cards.Card
	Action      string // "place_on_other", "place_on_self", "draw_card", "take_own_card"
}

// ExecutePhase1Turn executes a complete Phase 1 turn for the current player
// This implements the +1 cyclic rule mechanic:
// A) Player tries to place their top card on others (+1 rule)
// B) If can't place on others but can on self, must place on self then draw
// C) If can't place anywhere, must draw
// Returns the updated game state and action taken
func (e *Engine) ExecutePhase1Turn(game *models.Game, players []models.GamePlayer) (*Phase1TurnState, error) {
	turnState := &Phase1TurnState{
		PlacedCards: []cards.Card{},
	}

	currentPos := *game.CurrentPlayerPosition
	var currentPlayer *models.GamePlayer
	for i := range players {
		if players[i].Position == currentPos {
			currentPlayer = &players[i]
			break
		}
	}

	if currentPlayer == nil {
		return turnState, ErrPlayerNotFound
	}

	// STEP A: Try to place on other players' piles (repeat until can't place on others)
	for {
		canPlace, targetPos, isOnSelf := e.CanPlaceCard(*currentPlayer, players)

		if !canPlace {
			// Can't place anywhere - must draw
			turnState.Action = "draw_card"
			turnState.MustDraw = true
			break
		}

		if !isOnSelf {
			// Can place on another player - MUST do it
			// Find target player
			var targetPlayer *models.GamePlayer
			for i := range players {
				if players[i].Position == targetPos {
					targetPlayer = &players[i]
					break
				}
			}

			if targetPlayer != nil {
				// Execute placement
				e.PlaceCard(currentPlayer, targetPlayer)
				turnState.PlacedCards = append(turnState.PlacedCards, *currentPlayer.TopCard)
				turnState.Action = "place_on_other"
			}
			// Loop back to check if can place again
		} else {
			// Can only place on self - MUST do it, then MUST draw
			e.PlaceCard(currentPlayer, currentPlayer)
			turnState.PlacedCards = append(turnState.PlacedCards, *currentPlayer.TopCard)
			turnState.Action = "place_on_self_then_draw"
			turnState.MustDraw = true
			break
		}
	}

	return turnState, nil
}

// DrawPhase1Card draws a card and handles placement according to Phase 1 rules
// If card can be placed on others, MUST place it
// If card can't be placed on others but can be placed on self, place it
// If card can't be placed anywhere, keep it (add to own pile)
func (e *Engine) DrawPhase1Card(game *models.Game, player *models.GamePlayer, allPlayers []models.GamePlayer) (*Phase1TurnState, error) {
	turnState := &Phase1TurnState{
		PlacedCards: []cards.Card{},
	}

	if len(game.DeckCards) == 0 {
		// Deck is empty, turn ends
		return turnState, nil
	}

	// Draw one card
	card := game.DeckCards[0]
	game.DeckCards = game.DeckCards[1:]
	player.Cards = append(player.Cards, card)
	player.TopCard = &card
	player.CardCount = len(player.Cards)

	// B1: Check if drawn card can be placed on other players
	for _, targetPlayer := range allPlayers {
		if targetPlayer.Position == player.Position {
			continue
		}
		if targetPlayer.TopCard != nil && card.IsOnePlus(*targetPlayer.TopCard) {
			// MUST place on other player
			e.PlaceCard(player, &targetPlayer)
			turnState.Action = "place_drawn_on_other"
			return turnState, nil
		}
	}

	// B2: Check if drawn card can be placed on self
	if len(player.Cards) > 1 {
		cardBelow := player.Cards[len(player.Cards)-2]
		if card.IsOnePlus(cardBelow) {
			// MUST place on self
			e.PlaceCard(player, player)
			turnState.Action = "place_drawn_on_self"
			return turnState, nil
		}
	}

	// B3: Card can't be placed anywhere - keep it (turn ends)
	turnState.Action = "keep_drawn_card"
	return turnState, nil
}

// PlaceCard executes the placement of a card
func (e *Engine) PlaceCard(currentPlayer *models.GamePlayer, targetPlayer *models.GamePlayer) {
	// Remove top card from current player
	placedCard := *currentPlayer.TopCard
	currentPlayer.Cards = currentPlayer.Cards[:len(currentPlayer.Cards)-1]
	currentPlayer.CardCount = len(currentPlayer.Cards)

	// Update current player's top card
	if len(currentPlayer.Cards) > 0 {
		currentPlayer.TopCard = &currentPlayer.Cards[len(currentPlayer.Cards)-1]
	} else {
		currentPlayer.TopCard = nil
	}

	// Add card to target player
	targetPlayer.Cards = append(targetPlayer.Cards, placedCard)
	// Point to the card in the Cards array, not a local variable!
	targetPlayer.TopCard = &targetPlayer.Cards[len(targetPlayer.Cards)-1]
	targetPlayer.CardCount = len(targetPlayer.Cards)
}

// DrawCard draws a card from the deck
func (e *Engine) DrawCard(game *models.Game, player *models.GamePlayer) (*cards.Card, error) {
	if len(game.DeckCards) == 0 {
		return nil, errors.New("deck is empty")
	}

	// Draw card
	card := game.DeckCards[0]
	game.DeckCards = game.DeckCards[1:]

	return &card, nil
}

// AddCardToPlayer adds a drawn card to player's pile
func (e *Engine) AddCardToPlayer(player *models.GamePlayer, card cards.Card) {
	player.Cards = append(player.Cards, card)
	player.TopCard = &card
	player.CardCount = len(player.Cards)
}

// CheckCheat validates if a cheat occurred
func (e *Engine) CheckCheat(action string, currentPlayer models.GamePlayer, allPlayers []models.GamePlayer) bool {
	switch action {
	case "draw_when_can_place":
		canPlace, _, _ := e.CanPlaceCard(currentPlayer, allPlayers)
		return canPlace
	case "invalid_placement":
		// This would need more context about the actual placement attempted
		return false
	default:
		return false
	}
}

// ApplyPenalty applies penalty to a cheating player
func (e *Engine) ApplyPenalty(cheater *models.GamePlayer, allPlayers []models.GamePlayer) []models.GamePlayer {
	penalizedPlayers := []models.GamePlayer{}

	for i := range allPlayers {
		// Skip the cheater and players with only 1 card
		if allPlayers[i].Position == cheater.Position || allPlayers[i].CardCount <= 1 {
			continue
		}

		// Take one card (not top card) from each player
		if len(allPlayers[i].Cards) > 1 {
			// Take the second-to-last card (not the top)
			cardIndex := len(allPlayers[i].Cards) - 2
			penaltyCard := allPlayers[i].Cards[cardIndex]

			// Remove from player
			allPlayers[i].Cards = append(allPlayers[i].Cards[:cardIndex], allPlayers[i].Cards[cardIndex+1:]...)
			allPlayers[i].CardCount = len(allPlayers[i].Cards)

			// Add to bottom of cheater's pile
			cheater.Cards = append([]cards.Card{penaltyCard}, cheater.Cards...)
			cheater.CardCount = len(cheater.Cards)

			penalizedPlayers = append(penalizedPlayers, allPlayers[i])
		}
	}

	return penalizedPlayers
}

// CheckPhase1End checks if phase 1 should end
func (e *Engine) CheckPhase1End(game *models.Game) bool {
	return len(game.DeckCards) == 0
}

// StartPhase2 initializes phase 2 of the game
func (e *Engine) StartPhase2(game *models.Game, players []models.GamePlayer) error {
	game.Phase = 2
	game.State = models.StatePhase2

	// Determine trump suit (last card drawn in phase 1 that's not spades)
	for i := len(players) - 1; i >= 0; i-- {
		if players[i].TopCard != nil && players[i].TopCard.Suit != cards.Spades {
			trumpSuit := string(players[i].TopCard.Suit)
			game.TrumpSuit = &trumpSuit
			break
		}
	}

	// Find player with 9 of spades
	startPosition := 0
	for _, player := range players {
		if cards.FindNineOfSpades(player.Cards) != -1 {
			startPosition = player.Position
			break
		}
	}
	game.CurrentPlayerPosition = &startPosition

	// Initialize empty table
	game.TableCards = []cards.Card{}

	return nil
}

// ValidatePhase2Play validates a card play in phase 2
func (e *Engine) ValidatePhase2Play(game *models.Game, card cards.Card, player models.GamePlayer) error {
	// If table is empty, any card is valid
	if len(game.TableCards) == 0 {
		return nil
	}

	lastCard := game.TableCards[len(game.TableCards)-1]

	// Check if last card is trump
	isTrump := game.TrumpSuit != nil && string(card.Suit) == *game.TrumpSuit
	lastIsTrump := game.TrumpSuit != nil && string(lastCard.Suit) == *game.TrumpSuit

	if lastIsTrump {
		// Must play higher trump
		if !isTrump || !card.IsHigherRank(lastCard) {
			return errors.New("must play higher trump card")
		}
	} else {
		// Must play same suit with higher rank, or trump
		if card.IsSameSuit(lastCard) {
			if !card.IsHigherRank(lastCard) {
				return errors.New("must play higher rank of same suit")
			}
		} else if isTrump {
			// Trump can beat non-trump, unless it's spades and trump isn't spades
			if lastCard.Suit == cards.Spades && game.TrumpSuit != nil && *game.TrumpSuit != string(cards.Spades) {
				return errors.New("cannot beat spades with non-spade trump")
			}
		} else {
			return errors.New("must play same suit or trump")
		}
	}

	return nil
}

// CanPlayAnyCard checks if player can play any card
func (e *Engine) CanPlayAnyCard(game *models.Game, player models.GamePlayer) bool {
	for _, card := range player.Cards {
		if err := e.ValidatePhase2Play(game, card, player); err == nil {
			return true
		}
	}
	return false
}

// PlayCardPhase2 plays a card in phase 2
func (e *Engine) PlayCardPhase2(game *models.Game, player *models.GamePlayer, cardIndex int) error {
	if cardIndex < 0 || cardIndex >= len(player.Cards) {
		return errors.New("invalid card index")
	}

	card := player.Cards[cardIndex]

	// Validate the play
	if err := e.ValidatePhase2Play(game, card, *player); err != nil {
		return err
	}

	// Remove card from player
	player.Cards = append(player.Cards[:cardIndex], player.Cards[cardIndex+1:]...)
	player.CardCount = len(player.Cards)

	// Update top card
	if len(player.Cards) > 0 {
		player.TopCard = &player.Cards[len(player.Cards)-1]
	} else {
		player.TopCard = nil
	}

	// Add to table
	game.TableCards = append(game.TableCards, card)

	return nil
}

// TakeOldestTableCard takes the oldest card from table and adds to player
func (e *Engine) TakeOldestTableCard(game *models.Game, player *models.GamePlayer) error {
	if len(game.TableCards) == 0 {
		return errors.New("no cards on table")
	}

	// Take oldest (first) card
	card := game.TableCards[0]
	game.TableCards = game.TableCards[1:]

	// Add to bottom of player's pile
	player.Cards = append([]cards.Card{card}, player.Cards...)
	player.CardCount = len(player.Cards)

	return nil
}

// CheckGameEnd checks if the game should end
func (e *Engine) CheckGameEnd(players []models.GamePlayer) (bool, *int) {
	activePlayers := 0
	var loserPosition *int

	for i, player := range players {
		if player.CardCount > 0 {
			activePlayers++
			pos := players[i].Position
			loserPosition = &pos
		}
	}

	// Game ends when only one player has cards
	if activePlayers == 1 {
		return true, loserPosition
	}

	return false, nil
}

// NextPlayer determines the next player position
func (e *Engine) NextPlayer(currentPosition int, totalPlayers int) int {
	return (currentPosition + 1) % totalPlayers
}

// ClearTable clears all cards from the table
func (e *Engine) ClearTable(game *models.Game) {
	game.TableCards = []cards.Card{}
}

// GetPlayerByPosition finds a player by their position
func GetPlayerByPosition(players []models.GamePlayer, position int) (*models.GamePlayer, error) {
	for i := range players {
		if players[i].Position == position {
			return &players[i], nil
		}
	}
	return nil, fmt.Errorf("player at position %d not found", position)
}
