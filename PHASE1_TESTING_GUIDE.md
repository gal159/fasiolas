# Phase 1 Implementation Testing Guide

## Server Status ✅
- **Backend API**: Running on `http://localhost:8080`
- **Frontend**: Running on `http://localhost:3000`
- **Database**: PostgreSQL connected and healthy

## What Was Implemented

### 1. ✅ Cyclic +1 Rule (A → 2)
The game now correctly recognizes that an Ace can be followed by a 2 in a cyclic pattern:
- Sequence: 2 < 3 < 4 < ... < K < A < 2 < 3 < ...

**Implementation Location:** `pkg/cards/cards.go` - `IsOnePlus()` function

### 2. ✅ Phase 1 Turn Logic (A + B)
Complete turn execution with two phases:

#### Phase A - Placement Before Drawing
- Player tries to place their top card on any other player's pile
- Must follow the +1 rule
- Can place on self only if can't place on others
- Loops until can't place on others

#### Phase B - Drawing & Automatic Placement  
- If must draw: draws 1 card from deck
- Automatically places if possible (on others or self)
- Keeps card if can't place (pasiimti)

**Implementation Locations:**
- `internal/game/engine.go` - `ExecutePhase1Turn()`, `DrawPhase1Card()`
- `internal/service/game_service.go` - Updated service layer

### 3. ✅ Game State Tracking
New `Phase1TurnState` struct tracks:
- Whether player must draw
- Cards placed during turn
- Action type taken

## Testing Manually

### Step 1: Access Frontend
```
http://localhost:3000
```

### Step 2: Create/Join Game
1. Login with Google OAuth (set up in console)
2. Create a new game or join with room code
3. Wait for minimum 2 players

### Step 3: Start Game
- Game automatically initializes Phase 1
- Each player gets 1 card
- Player with lowest card starts

### Step 4: Observe Phase 1 Mechanics

**Turn Flow:**
1. Current player's card is displayed
2. System checks if can place on others (A1)
3. If yes: Place automatically
4. If no: Check if can place on self (A2)
5. If yes: Place then draw
6. If no: Draw card (enters phase B)

**Drawing (Phase B):**
1. Card drawn from deck
2. System checks if can place on others (B1)
3. If yes: Place, turn ends
4. If no: Check if can place on self (B2)
5. If yes: Place, turn ends
6. If no: Keep card, turn ends (B3)

## API Endpoints (For Testing with Postman/curl)

### Create Game
```
POST /api/v1/games
Authorization: Bearer {token}
Body: { "max_players": 4 }
```

### Join Game
```
POST /api/v1/games/join
Authorization: Bearer {token}
Body: { "room_code": "ABC123" }
```

### Start Game
```
POST /api/v1/games/{id}/start
Authorization: Bearer {token}
```

### Get Game State
```
GET /api/v1/games/{id}
Authorization: Bearer {token}
```

### Place Card
```
POST /api/v1/games/{id}/place
Authorization: Bearer {token}
Body: { "target_player_position": 1 }
```

### Draw Card
```
POST /api/v1/games/{id}/draw
Authorization: Bearer {token}
```

## Key Rules Verification Checklist

### Cyclic +1 Rule
- [ ] A of any suit can be followed by 2 of any suit
- [ ] 2 can be followed by 3
- [ ] K can be followed by A
- [ ] Other +1 progressions work (5→6, etc.)

### Mandatory Placement on Others
- [ ] If player's top card matches any other player, must place on one of them
- [ ] Player cannot draw when can place on others
- [ ] Loop continues until can't place on others

### Self-Placement & Drawing
- [ ] If can't place on others but can on self, must place on self
- [ ] After placing on self, must draw a card
- [ ] Cannot skip drawing after self-placement

### Automatic Placement of Drawn Card
- [ ] Drawn card automatically placed if matches any player (B1)
- [ ] Drawn card automatically placed on self if matches (B2)
- [ ] Drawn card kept if doesn't match (B3)

### Card Counting
- [ ] Card count increases when receiving card
- [ ] Card count decreases when placing on others
- [ ] TopCard always points to last card in pile

## Expected Behavior Examples

### Example 1: Direct Placement
```
Player A has 5♥ on top
Player B has 4♣ on top

A's turn:
- Check: 5♥ is +1 to 4♣? YES (4+1=5)
- Action: Place 5♥ on Player B
- B's pile: [4♣, 5♥]
- A's pile: now one card less
- Check again: Can place remaining top card?
```

### Example 2: Self-Placement Then Draw
```
Player A has [5♠, 6♥] (top: 6♥)
Player B has [3♣, 4♦] (top: 4♦)
Player C has [Q♠, K♣] (top: K♣)

A's turn:
- Check: 6♥ is +1 to 4♦? NO (4+1=5, not 6)
- Check: 6♥ is +1 to K♣? NO (K+1=A, not 6)
- Check: Can place on self? Check: 6♥ is +1 to 5♠? YES!
- Action: Place 6♥ on self
- A's pile: [5♠, 6♥, 6♥]
- Must DRAW from deck
```

### Example 3: Draw & Auto-Place
```
Player A has [7♦, 8♠] (top: 8♠)
Other players can't receive 8♠

A's turn:
- Can't place 8♠ on others → Must draw
- Draw: 9♠
- Check B1: Can place 9♠ on others? NO
- Check B2: Can place 9♠ on self? 9♠ is +1 to 8♠? YES!
- Action: Place on self automatically
- A's pile: [7♦, 8♠, 9♠]
- Turn ends
```

## Debugging

### Check Phase 1 State
The `GET /api/v1/games/{id}` endpoint returns:
```json
{
  "game": {
    "phase": 1,
    "state": "phase1",
    "deck_cards": [...],
    "current_player_position": 0
  },
  "players": [
    {
      "position": 0,
      "card_count": 3,
      "top_card": { "suit": "hearts", "rank": "5", "value": 5 },
      "cards": [...]
    }
  ]
}
```

### Verify TopCard
- Check that `top_card` matches last element in `cards` array
- Verify `card_count` equals length of `cards` array

### Check Current Player
- `current_player_position` shows whose turn it is
- Should cycle through positions 0, 1, 2, ... in order

## Known Limitations (Phase 2 Pending)

- Phase 1 currently doesn't automatically transition to Phase 2
- No trump suit selection yet (for Phase 2)
- No winning logic yet
- No cheat detection yet

## Next Steps After Phase 1 Works

1. **Implement Phase 2 Logic** - Trump suit strategy game
2. **Add WebSocket Events** - Real-time turn updates
3. **Implement Game Ending** - When only 1 player has cards
4. **Add Cheat Detection** - Validate player moves
5. **Implement Elo Rating** - Player ranking system
6. **Add Replay System** - Watch game recordings

## Support

For issues or questions:
1. Check server logs: `docker compose logs fasiolas_app`
2. Check database: `docker compose exec postgres psql -U postgres -d fasiolas_game`
3. Review game rules: See `GAME_RULES_IMPLEMENTATION.md`
4. Review implementation: See `PHASE1_IMPLEMENTATION_SUMMARY.md`

