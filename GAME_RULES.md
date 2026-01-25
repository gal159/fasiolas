# 🎮 Fasiolas Game Rules & Implementation Guide

This document explains the Fasiolas card game rules and how they are implemented in the API.

---

## 🎴 Game Overview

**Fasiolas** is a two-phase Lithuanian card game for 2-8 players using a standard 52-card deck.

### Game Phases
1. **Phase 1 (Accumulation)**: Players place cards on each other's piles following the +1 rank rule
2. **Phase 2 (Trick-Taking)**: Players get rid of cards by beating previous cards with higher rank or trump

The player who ends up with cards after everyone else is empty **LOSES**.

---

## 🟦 PHASE 1: ACCUMULATION (The "+1" Phase)

### Goal
Accumulate cards on other players' piles while following the +1 rank rule and avoiding cheating penalties.

### Setup

```
1. Shuffle the standard 52-card deck
2. Each player receives 1 face-up card (top card of their pile)
3. Player with the LOWEST card starts (e.g., 2 beats A)
4. Cards ranked: A=1, 2=2, ..., K=13 (for determining starter only)
```

**API Implementation:**
```go
// In engine.go - InitializeGame()
// Finds lowest card value to determine starter
for _, player := range players {
    if player.TopCard.Value < lowestValue {
        lowestValue = player.TopCard.Value
        lowestPosition = player.Position
    }
}
game.CurrentPlayerPosition = &lowestPosition
```

### Player's Turn Sequence

Each player's turn follows this strict order:

#### Step 1️⃣: Mandatory Placement Check
**Before doing anything else, check if you can place your top card!**

A card can be placed if:
- Your top card rank = another player's top card rank + 1
- Example: If opponent has 7, you place 8; if opponent has K, you place A

**Cannot skip this!** If you can place, you MUST place (or it's cheating).

```go
// In engine.go - CanPlaceCard()
// Checks if current player's top card is +1 from any opponent
for _, targetPlayer := range allPlayers {
    if targetPlayer.Position == currentPlayer.Position {
        continue // Skip self
    }
    if currentPlayer.TopCard.IsOnePlus(*targetPlayer.TopCard) {
        return true, targetPlayer.Position // Can place!
    }
}
```

#### Step 2️⃣: Place Cards (if applicable)
If you found a valid placement:
- Place your top card on that player's pile
- That becomes their new top card
- Repeat: Check your new top card - if it can be placed, place it again
- Continue until you have no more valid placements

**API Endpoint:**
```
POST /api/v1/games/{gameID}/place
Body: { "target_player_position": 2 }
```

```go
// In game_service.go - PlaceCard()
// Moves card from current player to target player
s.engine.PlaceCard(currentPlayer, targetPlayer)

// Check if can place again
canPlace, _ := s.engine.CanPlaceCard(*currentPlayer, players)
if canPlace {
    // Client should call place endpoint again
    return nil
}
```

#### Step 3️⃣: Draw from Deck
If you can't place anymore:
- Draw 1 card from the deck
- This card goes into your "hand" temporarily
- Check: Does this drawn card fit (+1 rule)?
  - If YES: Place it immediately (and check if can place again)
  - If NO: Add it to the bottom of your pile, turn ends

**API Endpoint:**
```
POST /api/v1/games/{gameID}/draw
```

```go
// In game_service.go - DrawCard()
card, err := s.engine.DrawCard(g, currentPlayer)

// Check if drawn card can be placed
canPlaceDrawn, targetPos := s.engine.CanPlaceCard(*currentPlayer, players)

if canPlaceDrawn {
    // Client must place it in next action
    return nil
}

// Otherwise, add to pile
s.engine.AddCardToPlayer(currentPlayer, *card)
```

---

## ⚠️ CHEATING & PENALTIES

### What is Cheating?

Cheating occurs when a player:

1. **Draws when they can place** ❌
   - You have a valid +1 placement but drew instead
   - Extremely bad move and disqualifying offense

2. **Places incorrectly** ❌
   - Placed a card that isn't +1 from target's top card
   - Invalid placement attempt

### Detection

```go
// In engine.go - CheckCheat()
// Validates that player followed rules
func (e *Engine) CheckCheat(action string, currentPlayer models.GamePlayer, allPlayers []models.GamePlayer) bool {
    if action == "draw_when_can_place" {
        canPlace, _ := e.CanPlaceCard(currentPlayer, allPlayers)
        return canPlace  // Is cheating if true
    }
    return false
}
```

### Penalty System

When cheating is proven:

1. **All other players** (who have more than 1 card) give:
   - 1 card (not top card) to the cheater
   - Card goes to **bottom** of cheater's pile
   - Cards given in turn order

2. **Cheater loses cheat calling rights** for rest of game

3. **Wrong accusation**:
   - If accused player is innocent
   - Accuser loses ability to call cheats for rest of game

**API Endpoint:**
```
POST /api/v1/games/{gameID}/cheat
Body: { "cheater_id": 5 }
```

```go
// In game_service.go - CallCheat()
// Validates and applies penalty
penalizedPlayers := s.engine.ApplyPenalty(cheater, players)

// Update all players and log action
for i := range penalizedPlayers {
    s.playerRepo.Update(&penalizedPlayers[i])
}
```

---

## 🟧 PHASE 1 END CONDITIONS

### When Does Phase 1 End?

**When the deck becomes empty** (no more cards to draw)

```go
// In engine.go - CheckPhase1End()
func (e *Engine) CheckPhase1End(game *models.Game) bool {
    return len(game.DeckCards) == 0
}
```

### What Happens Next?

1. **Record each player's card count** (stored in database)
2. **Save each player's top card** (important for trump determination)
3. **Transition to Phase 2**

```go
// In game_service.go - transitionToPhase2()
if s.engine.CheckPhase1End(g) {
    return s.transitionToPhase2(g, players)
}
```

---

## 🟨 PHASE 2: TRICK-TAKING (The "Mušimas" Phase)

### Goal
Get rid of all your cards. **Last player with cards LOSES.**

### Trump Suit Determination

**Trump = The suit of the last non-spade card in Phase 1**

- Look at the last player's top card (player who received last card)
- If that card is NOT spades (♠), that suit becomes trump
- If all remaining cards are spades, there is NO trump

**Why:** Trump cards beat everything except spades

```go
// In engine.go - StartPhase2()
// Find trump suit from last card dealt
for i := len(players) - 1; i >= 0; i-- {
    if players[i].TopCard != nil && players[i].TopCard.Suit != cards.Spades {
        trumpSuit := string(players[i].TopCard.Suit)
        game.TrumpSuit = &trumpSuit
        break
    }
}
```

### Starting Player

**Player with 9 of Spades starts Phase 2**

```go
// In engine.go - StartPhase2()
for _, player := range players {
    if cards.FindNineOfSpades(player.Cards) != -1 {
        startPosition = player.Position
        break
    }
}
```

### Playing Rules

Players take turns. Each turn: **Play ONE card to the table**

#### If table is EMPTY:
→ **Play any card** you want

#### If table is NOT empty:
→ **Must follow one of these rules:**

**Rule 1: Same Suit + Higher Rank**
- Last card on table: 7♥
- You have: 9♥
- You can play: 8♥, 9♥, 10♥, ..., A♥

**Rule 2: Trump Card**
- Last card on table: 7♥ (not trump)
- Trump is: ♦
- You can play: Any diamond (2♦ beats 7♥)

**Rule 3: Higher Trump**
- Last card on table: 5♦ (trump)
- You must play: 6♦, 7♦, ..., A♦
- Cannot play lower trump

**Special Rule: Spades (♠)**
- Spades can ONLY be beaten by higher spades
- Trump cannot beat spades (unless trump IS spades)
- Example: If trump is ♥, cannot play 9♥ on A♠

```go
// In engine.go - ValidatePhase2Play()
if lastCard.Suit == Spades && game.TrumpSuit != nil && *game.TrumpSuit != "spades" {
    return errors.New("cannot beat spades with non-spade trump")
}
```

### Cannot Play

If you cannot play any card:
1. Take the **oldest card from the table** (the first one laid down)
2. Add it to the **bottom** of your pile
3. **Skip your turn** (do NOT play a card this round)
4. Pla continues to next player

**This is NOT a turn!** You don't go again.

```go
// In engine.go - TakeOldestTableCard()
card := game.TableCards[0]  // First card laid
game.TableCards = game.TableCards[1:]
player.Cards = append([]Card{card}, player.Cards...)  // Add to bottom
```

### Round End

When **all players** have played one card each:
- Table is cleared
- **Last player to play starts next round**
- Continue until someone has no cards left

```go
// In engine.go - ClearTable()
game.TableCards = []cards.Card{}  // Empty table

// Next round: player who just played goes first
game.CurrentPlayerPosition = &playerWhoPlayed
```

---

## 🏁 GAME END

### Win Condition
Game ends when only **ONE player has cards remaining**

That player **LOSES**
Everyone else **WINS**

```go
// In engine.go - CheckGameEnd()
activePlayers := 0
for _, player := range players {
    if player.CardCount > 0 {
        activePlayers++
        loserPosition = &player.Position
    }
}

if activePlayers == 1 {
    return true, loserPosition  // Game over!
}
```

### Winner Assignment

```go
// In game_service.go
if gameEnded {
    loserID := // Get from CheckGameEnd()
    for _, player := range players {
        if player.ID == loserID {
            player.Status = models.StatusEliminated
        } else {
            player.Status = models.StatusWinner
        }
    }
}
```

---

## 📊 Implementation Reference

### Card Values (Phase 1 starter determination)
```
Ace = 1 (Lowest!)
2-10 = face value
Jack = 11
Queen = 12
King = 13
```

### Card Comparison Methods
```go
// +1 check (Phase 1)
card.IsOnePlus(other)      // 8.IsOnePlus(7) = true

// Same suit check (Phase 2)
card.IsSameSuit(other)     // 8♥.IsSameSuit(7♥) = true

// Higher rank check (Phase 2)
card.IsHigherRank(other)   // 9.IsHigherRank(7) = true

// Trump check
card.IsTrump(trumpSuit)    // Check if card is trump suit
```

### Player Status Values
```go
enum PlayerStatus {
    ACTIVE = "active"       // Currently playing
    ELIMINATED = "eliminated"  // Out of game (cards left)
    WINNER = "winner"       // Game finished, has no cards
    LEFT = "left"          // Disconnected during game
}
```

---

## 🔄 Game Flow Diagram

```
┌─────────────────────────────────┐
│   Players Join Game             │
│   (via /api/v1/games/join)      │
└──────────────┬──────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│   Game Creator Starts Game      │
│   (via /api/v1/games/{id}/start)│
└──────────────┬──────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│          PHASE 1: ACCUMULATION              │
│  ┌──────────────────────────────────────┐   │
│  │ Player's Turn:                       │   │
│  │ 1. Check if can place (+1)           │   │
│  │ 2. If yes: Place cards (loop until   │   │
│  │    no more +1 possible)              │   │
│  │ 3. Draw from deck                    │   │
│  │ 4. Check drawn card:                 │   │
│  │    - If +1: Place it                 │   │
│  │    - If not: Add to bottom of pile   │   │
│  │ 5. Next player                       │   │
│  │ (Can call cheat any time)            │   │
│  └──────────────────────────────────────┘   │
│  Continue until deck is empty              │
└──────────────┬───────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│      PHASE TRANSITION                        │
│  - Determine trump suit                      │
│  - Find player with 9♠                       │
│  - Clear table                               │
└──────────────┬───────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────┐
│       PHASE 2: TRICK-TAKING                 │
│  ┌──────────────────────────────────────┐   │
│  │ Player's Turn:                       │   │
│  │ 1. Play valid card (or take oldest)  │   │
│  │ 2. Validate against rules:           │   │
│  │    - Same suit + higher rank         │   │
│  │    - OR play trump                   │   │
│  │    - OR can't play: take oldest card │   │
│  │ 3. Next player                       │   │
│  │                                      │   │
│  │ When all played:                     │   │
│  │ - Clear table                        │   │
│  │ - Last player goes again             │   │
│  └──────────────────────────────────────┘   │
│  Continue until only 1 player has cards    │
└──────────────┬───────────────────────────────┘
               │
               ▼
┌─────────────────────────────────┐
│   GAME END                       │
│   Last player with cards LOSES   │
│   All others WIN                 │
└─────────────────────────────────┘
```

---

## API Endpoint Reference by Phase

### Phase 1 Endpoints
```
POST /api/v1/games/{id}/place    - Place card (+1 rule)
POST /api/v1/games/{id}/draw     - Draw from deck
POST /api/v1/games/{id}/cheat    - Call cheat
GET  /api/v1/games/{id}          - Check game state & phase
```

### Phase 2 Endpoints
```
POST /api/v1/games/{id}/place    - Play card (trick-taking rules)
GET  /api/v1/games/{id}          - Check game state & trump suit
```

### Shared Endpoints
```
GET  /api/v1/games/{id}          - Get full game state
POST /api/v1/games/{id}/start    - Start the game
```

---

## 🎯 Summary

| Aspect | Phase 1 | Phase 2 |
|--------|---------|---------|
| **Goal** | Accumulate cards | Get rid of cards |
| **Card Placement** | +1 rank only | Suit/trump/rank based |
| **Deck** | Draw from deck | No drawing |
| **Table** | No table | Cards played to table |
| **Trump** | N/A | Determined at transition |
| **Loser** | N/A | Last with cards |

---

## 📚 Further Reading

- [Fasiolas Wikipedia](https://lt.wikipedia.org/wiki/Faziolas) (Lithuanian)
- See `/api/v1/games/{id}` response for current game state details
- Check `GameAction` log for complete action history

