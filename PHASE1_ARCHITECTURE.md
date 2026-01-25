# Fasiolas Card Game - Phase 1 Architecture Document

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (React/Vite)                   │
│                   http://localhost:3000                      │
│                                                              │
│  Components:                                                 │
│  - Login (Google OAuth)                                     │
│  - Dashboard (Game List)                                    │
│  - Game Room (Lobby)                                        │
│  - Game Table (Phase 1 & 2)                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/WebSocket
                              │
┌─────────────────────────────────────────────────────────────┐
│              Backend API (Go/Gin)                            │
│              http://localhost:8080                           │
│                                                              │
│  Routes:                                                    │
│  POST   /api/v1/games              - CreateGame            │
│  POST   /api/v1/games/join         - JoinGame              │
│  POST   /api/v1/games/{id}/start   - StartGame             │
│  GET    /api/v1/games/{id}         - GetGameState          │
│  POST   /api/v1/games/{id}/place   - PlaceCard             │
│  POST   /api/v1/games/{id}/draw    - DrawCard              │
│  GET    /api/v1/auth/google        - GoogleAuth            │
│  GET    /api/v1/auth/google/callback - OAuthCallback       │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ SQL Queries
                              │
┌─────────────────────────────────────────────────────────────┐
│          PostgreSQL Database                                 │
│          postgres://postgres:5432/fasiolas_game             │
│                                                              │
│  Tables:                                                    │
│  - users                                                    │
│  - games                                                    │
│  - game_players                                             │
│  - game_actions                                             │
└─────────────────────────────────────────────────────────────┘
```

## Game Logic Layer Structure

```
pkg/cards/
├── cards.go
│   ├── Card          - Single card representation
│   ├── Deck          - 52-card deck
│   ├── IsOnePlus()   - +1 rule check (with cyclic A→2)
│   ├── Shuffle()     - Fisher-Yates shuffle
│   └── Draw()        - Remove top card

internal/game/
├── engine.go
│   ├── Engine        - Game rule engine
│   ├── CanPlaceCard() - Check if card can be placed (A1/A2/A3 logic)
│   ├── PlaceCard()   - Execute card placement
│   ├── ExecutePhase1Turn() - Full turn A phase (loop)
│   ├── DrawPhase1Card() - Full turn B phase (B1/B2/B3)
│   └── StartPhase2() - Initialize phase 2

internal/service/
├── game_service.go
│   ├── GameService   - Business logic
│   ├── CreateGame()  - Create new game
│   ├── JoinGame()    - Add player to game
│   ├── StartGame()   - Initialize and deal cards
│   ├── PlaceCard()   - Service-level placement
│   └── DrawCard()    - Service-level drawing

internal/models/
├── models.go
│   ├── Game          - Game state
│   ├── GamePlayer    - Player in game
│   ├── User          - User profile
│   └── GameAction    - Action history
```

## Phase 1 Turn Execution Flow

### Detailed State Machine

```
START_TURN (currentPlayer)
    │
    ├─→ PHASE_A: Try to Place
    │   │
    │   └─→ LOOP:
    │       ├─→ CanPlaceCard(currentPlayer, allPlayers)
    │       │   │
    │       │   ├─→ [false] → EXIT_LOOP, goto PHASE_B_MUST_DRAW
    │       │   │
    │       │   ├─→ [true, targetPos != self] → PlaceCard(current, target)
    │       │   │   └─→ Update topCard, Update cardCount
    │       │   │   └─→ CONTINUE LOOP
    │       │   │
    │       │   └─→ [true, targetPos == self] → PlaceCard(current, self)
    │       │       └─→ Update topCard, Update cardCount
    │       │       └─→ Mark: MustDraw = true
    │       │       └─→ EXIT_LOOP
    │       │
    │       └─→ END_LOOP
    │
    ├─→ PHASE_B: Draw (if MustDraw)
    │   │
    │   ├─→ Draw 1 card from deck
    │   │
    │   ├─→ B1: Can place on other players?
    │   │   ├─→ YES → PlaceCard(currentPlayer, targetPlayer)
    │   │   │        → TURN_ENDS
    │   │   │
    │   │   └─→ NO → goto B2
    │   │
    │   ├─→ B2: Can place on self?
    │   │   ├─→ YES → PlaceCard(currentPlayer, currentPlayer)
    │   │   │        → TURN_ENDS
    │   │   │
    │   │   └─→ NO → goto B3
    │   │
    │   └─→ B3: Place on self (keep card)
    │       → TURN_ENDS
    │
    └─→ NEXT_PLAYER = (currentPlayer + 1) % totalPlayers
        GOTO START_TURN
```

## Data Flow for a Single Action

### Example: Place Card Action

```
Frontend Request:
POST /api/v1/games/{id}/place
{
  "target_player_position": 1
}

↓

HTTP Handler (handler/game_handler.go)
├─ Extract userID from JWT
├─ Parse request body
└─ Call gameService.PlaceCard()

↓

Service Layer (service/game_service.go)
├─ Load game from DB
├─ Load all players from DB
├─ Find current player
├─ Verify it's their turn
├─ Verify phase == 1
├─ Call engine.CanPlaceCard() to validate
├─ Call engine.PlaceCard() to execute
├─ Update currentPlayer in DB
├─ Update targetPlayer in DB
├─ Log action in game_actions table
├─ Check if can place again
└─ Return success

↓

Repository Layer (repository/)
├─ gameRepo.Update(game)
├─ playerRepo.Update(currentPlayer)
├─ playerRepo.Update(targetPlayer)
└─ actionRepo.Create(action)

↓

Database (PostgreSQL)
├─ UPDATE game_players SET cards = ? WHERE id = currentPlayer.id
├─ UPDATE game_players SET cards = ? WHERE id = targetPlayer.id
├─ INSERT INTO game_actions (game_id, user_id, action_type, ...) VALUES (...)
└─ ✓ Persisted

↓

Response to Frontend:
200 OK
{
  "message": "card placed",
  "game_state": { ... }  // Optional
}

↓

Frontend Update
├─ Refresh game state
├─ Update UI to show new pile positions
├─ Check if turn changed
├─ Show notifications
└─ Display updated top cards
```

## Card Representation & Storage

### In Memory (During Game)
```go
type Card struct {
    Suit  Suit  // "hearts", "diamonds", "clubs", "spades"
    Rank  Rank  // "2", "3", ..., "10", "J", "Q", "K", "A"
    Value int   // 2-14 (for ranking)
}

type GamePlayer struct {
    Position  int
    Cards     []Card      // All cards, ordered bottom→top
    TopCard   *Card       // Pointer to Cards[len-1]
    CardCount int         // len(Cards)
}
```

### In Database
```sql
-- Cards are stored as JSON in game_players.cards
{
  "cards": [
    {"suit": "hearts", "rank": "5", "value": 5},
    {"suit": "clubs", "rank": "8", "value": 8},
    {"suit": "diamonds", "rank": "9", "value": 9}
  ],
  "top_card": {"suit": "diamonds", "rank": "9", "value": 9},
  "card_count": 3
}
```

## Turn Order & Position System

```
Game with 4 players:

Position 0 (Player A) --- dealer
Position 1 (Player B)
Position 2 (Player C)
Position 3 (Player D)

Turn order: 0 → 1 → 2 → 3 → 0 → 1 → ...

Next position formula:
nextPos = (currentPosition + 1) % totalPlayers

Example:
Current: 2 (Player C)
Total: 4
Next: (2 + 1) % 4 = 3 (Player D)

Current: 3 (Player D)
Total: 4
Next: (3 + 1) % 4 = 0 (Player A)
```

## Phase 1 End Condition

```
PHASE_1_ENDS when:
├─ Deck is empty (all cards drawn)
├─ AND current player can't place on others
├─ AND current player can't place on self
└─ (i.e., must draw but deck empty)

Transition to PHASE_2:
├─ Trump suit determined (last player's top card suit)
├─ Game state updated to "phase2"
├─ Current player set to holder of 9♠
├─ Table cleared
└─ Ready for Phase 2 play
```

## Key Functions & Their Responsibilities

### `IsOnePlus(other Card) bool`
**Location:** `pkg/cards/cards.go`
**Purpose:** Check if card is exactly +1 rank from other
**Special:** Handles cyclic rule (A+1=2)
**Time Complexity:** O(1)

### `CanPlaceCard(player, allPlayers) (bool, int, bool)`
**Location:** `internal/game/engine.go`
**Purpose:** Determine if and where a card can be placed
**Returns:**
1. Can place? (bool)
2. Target position (int, -1 if can't)
3. Is target self? (bool)
**Logic:** A1 (others) → A2 (self) → A3 (can't)
**Time Complexity:** O(n) where n = number of players

### `PlaceCard(current, target) void`
**Location:** `internal/game/engine.go`
**Purpose:** Execute the physical placement of card
**Side Effects:**
- Removes from current player's hand
- Adds to target player's hand
- Updates both topCards
- Updates both cardCounts
**Time Complexity:** O(n) where n = card count

### `ExecutePhase1Turn(game, players) (*Phase1TurnState, error)`
**Location:** `internal/game/engine.go`
**Purpose:** Execute complete Phase A of turn
**Loop Logic:**
- While can place on others: place and continue
- When can only place on self: place and break
- When can't place: break
**Returns:** State of turn (must draw? cards placed? action?)
**Time Complexity:** O(p*n) where p = players, n = avg cards

### `DrawPhase1Card(game, player, allPlayers) (*Phase1TurnState, error)`
**Location:** `internal/game/engine.go`
**Purpose:** Execute Phase B (draw and auto-place)
**Logic:**
- B1: If matches other player → place, end
- B2: If matches self → place, end
- B3: If matches nothing → keep, end
**Time Complexity:** O(p*n)

## Error Handling Strategy

```go
// Graceful degradation
if player == nil {
    return ErrPlayerNotFound  // Recoverable
}

if len(deck) == 0 {
    return errors.New("deck empty")  // Expected in Phase 1→2 transition
}

// Validation before action
canPlace, _, _ := engine.CanPlaceCard(player, players)
if !canPlace {
    return errors.New("invalid placement")  // User error
}

// Database consistency
if err := playerRepo.Update(player); err != nil {
    return fmt.Errorf("failed to update: %w", err)  // Internal error
}
```

## Testing Strategy

### Unit Tests (cards.go)
```go
func TestIsOnePlus(t *testing.T) {
    // Test A→2 cyclic
    // Test K→A regular
    // Test 5→6 regular
    // Test invalid progressions
}

func TestDeckShuffle(t *testing.T) {
    // Verify all 52 cards present
    // Verify randomization
}
```

### Integration Tests (engine.go)
```go
func TestExecutePhase1Turn(t *testing.T) {
    // Setup 2+ players with cards
    // Execute turn
    // Verify correct placements
    // Verify topCard updates
    // Verify cardCount matches
}

func TestDrawPhase1Card(t *testing.T) {
    // B1: Can place on other
    // B2: Can place on self
    // B3: Can't place anywhere
}
```

### End-to-End Tests (game_service.go)
```go
func TestFullPhase1Game(t *testing.T) {
    // Create game
    // Add 2+ players
    // Start game
    // Execute multiple turns
    // Verify Phase 1 ends correctly
}
```

## Performance Considerations

| Operation | Time | Space | Notes |
|-----------|------|-------|-------|
| CanPlaceCard | O(p*n) | O(1) | p=players, n=cards |
| PlaceCard | O(n) | O(1) | Slice operations |
| ExecutePhase1Turn | O(p²*n) | O(n) | Worst: place all turns |
| DrawPhase1Card | O(p*n) | O(1) | Single draw |
| Database Update | O(1) | O(n) | JSON serialization |
| Game Serialize | O(p*n) | O(p*n) | Full game state |

## Scalability Notes

**Current Limitations:**
- Game state in-memory only (lost on restart)
- Card arrays copied for each update
- No caching layer
- Sequential turn processing

**Improvements for 1000+ concurrent games:**
1. Redis caching for active games
2. Event sourcing instead of state snapshots
3. Async turn processing
4. Read replicas for game state queries

## Security Considerations

✅ **Implemented:**
- JWT authentication for all endpoints
- User ownership verification on game actions
- Turn validation (can't act on other's turn)

⚠️ **TODO for Production:**
- Input validation (card placement bounds)
- Rate limiting on API endpoints
- SQL injection prevention (using parameterized queries)
- Replay attack prevention
- Cheat detection algorithm

## Database Schema (Relevant for Phase 1)

```sql
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    room_code VARCHAR(10) UNIQUE,
    state VARCHAR(20),  -- 'waiting', 'phase1', 'phase2', 'finished'
    phase INTEGER,       -- 1 or 2
    current_player_position INTEGER,
    deck_cards JSONB,   -- Array of Card objects
    table_cards JSONB,  -- For Phase 2
    created_by INTEGER REFERENCES users(id),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    max_players INTEGER
);

CREATE TABLE game_players (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games(id),
    user_id INTEGER REFERENCES users(id),
    position INTEGER,
    cards JSONB,        -- Array of Card objects
    top_card JSONB,     -- Single Card object
    card_count INTEGER,
    status VARCHAR(20), -- 'active', 'eliminated', 'winner'
    joined_at TIMESTAMP
);

CREATE TABLE game_actions (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games(id),
    user_id INTEGER REFERENCES users(id),
    action_type VARCHAR(50),  -- 'place_card', 'draw_card', etc.
    action_data JSONB,
    phase INTEGER,
    timestamp TIMESTAMP
);
```

## Deployment & Environment

```
Development:
├─ Docker Compose for local dev
├─ In-memory game states OK
└─ PostgreSQL in Docker

Production:
├─ Kubernetes pods
├─ RDS PostgreSQL
├─ Redis for session/cache
├─ CloudFront CDN for frontend
└─ Lambda for serverless functions

Environment Variables:
- DB_HOST=postgres (dev), RDS_ENDPOINT (prod)
- DB_PORT=5432
- DB_NAME=fasiolas_game
- DB_USER=postgres
- GOOGLE_CLIENT_ID=...
- GOOGLE_CLIENT_SECRET=...
- JWT_SECRET=...
- GIN_MODE=debug (dev), release (prod)
```

## Documentation Files Reference

| File | Purpose |
|------|---------|
| `GAME_RULES_IMPLEMENTATION.md` | Full game rules in Lithuanian |
| `PHASE1_IMPLEMENTATION_SUMMARY.md` | What code changed |
| `PHASE1_TESTING_GUIDE.md` | How to test the implementation |
| `PHASE1_ARCHITECTURE.md` | This file - System design |

---

**Created:** January 23, 2026
**Status:** ✅ Phase 1 Implementation Complete
**Next:** Phase 2 (Trump suit strategy game)

