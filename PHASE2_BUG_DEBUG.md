# 🐛 Phase 2 Card Visibility - Debugging Guide

## Issue Report
**Problem:** In Phase 2, players only see their own cards, not other players' cards  
**Expected:** In Phase 2, ALL players should see ALL cards from ALL players  
**Date:** January 30, 2026

---

## How to Reproduce

### Step 1: Start the Game
1. Start backend server: `.\bin\server.exe`
2. Start frontend: `cd frontend && npm start`
3. Open 2 browser windows (or use incognito)

### Step 2: Create a Game
1. **Browser 1:** Login as User 1
2. **Browser 1:** Create a game (2 players)
3. **Browser 2:** Login as User 2
4. **Browser 2:** Join the game with room code
5. **Browser 1:** Start the game

### Step 3: Get to Phase 2
1. Play Phase 1 normally
2. Draw cards until deck is empty
3. Observe phase transition to Phase 2

### Step 4: Check Card Visibility
**Open Browser Console (F12) in Both Windows**

```javascript
// Check game phase
game.phase  // Should be 2

// Check players data
players.forEach((p, idx) => {
  console.log(`Player ${idx+1} (${p.user.username}):`, {
    userId: p.user_id,
    cardCount: p.card_count,
    cardsArrayLength: p.cards.length,
    cards: p.cards
  });
});
```

**Expected Output (Phase 2):**
```javascript
Player 1 (Alice): { 
  userId: 1, 
  cardCount: 5, 
  cardsArrayLength: 5,  // ✅ Should have cards
  cards: [{suit: "hearts", rank: "5"}, ...] 
}
Player 2 (Bob): { 
  userId: 2, 
  cardCount: 7, 
  cardsArrayLength: 7,  // ✅ Should have cards
  cards: [{suit: "clubs", rank: "K"}, ...] 
}
```

**Actual Output (Bug):**
```javascript
Player 1 (Alice): { 
  userId: 1, 
  cardCount: 5, 
  cardsArrayLength: 5,  // ✅ Your cards visible
  cards: [{suit: "hearts", rank: "5"}, ...] 
}
Player 2 (Bob): { 
  userId: 2, 
  cardCount: 7, 
  cardsArrayLength: 0,  // ❌ Empty array! (BUG)
  cards: [] 
}
```

---

## Debugging Steps

### Step 1: Check Backend Logs

With the debug logging I just added, check the backend console output:

```bash
# Look for these logs when calling GetGameState
[GetGameState] Player 1 (Alice) - Phase: 2, CardCount: 5, CardsLen: 5, IsCurrentUser: true
[GetGameState] Phase 2: Showing all cards for player 1 (cards: 5)
[GetGameState] Player 2 (Bob) - Phase: 2, CardCount: 7, CardsLen: 0, IsCurrentUser: false
[GetGameState] Phase 2: Showing all cards for player 2 (cards: 0)
```

**If CardsLen is 0 for other players in Phase 2, this means:**
- The database is not returning the cards
- OR the cards were cleared somewhere earlier
- OR there's a database query issue

### Step 2: Check Database Directly

Connect to PostgreSQL and run:

```sql
-- Check game state
SELECT id, phase, state, deck_cards 
FROM games 
WHERE id = 1;

-- Check players' cards
SELECT id, user_id, card_count, 
       jsonb_array_length(cards) as cards_in_db
FROM game_players 
WHERE game_id = 1;
```

**Expected Result:**
```
 id | user_id | card_count | cards_in_db 
----+---------+------------+-------------
  1 |       1 |          5 |           5
  2 |       2 |          7 |           7
```

**If cards_in_db is 0:**
- Cards are not being saved to database
- This is a different bug (card persistence issue)

**If cards_in_db > 0 but API returns empty:**
- Issue is in the service layer
- Check if GetByGame() is loading cards correctly

---

## Potential Root Causes

### Cause 1: Database Query Not Loading Cards
**Location:** `internal/repository/game_player_repository.go`

Check if `GetByGame()` is properly loading the `cards` JSONB field:

```go
func (r *GamePlayerRepository) GetByGame(gameID int) ([]models.GamePlayer, error) {
    var players []models.GamePlayer
    err := r.db.Select(&players, 
        "SELECT * FROM game_players WHERE game_id = $1 ORDER BY position", 
        gameID)
    return players, err
}
```

**Fix:** Ensure the query includes all fields including `cards`

### Cause 2: Cards Array Being Cleared Before Phase Check
**Location:** `internal/service/game_service.go`

Issue: Cards might be cleared in Phase 1 and never repopulated.

**Current Logic:**
```go
// This MODIFIES the players slice from database
if g.Phase == 1 && players[i].UserID != userID {
    players[i].Cards = []cards.Card{}
}
```

**Problem:** If the same players slice is reused/cached, cards stay empty.

**Fix:** Always fetch fresh data, or only clear on response (not in memory).

### Cause 3: Frontend Caching
**Location:** `frontend/src/components/GameBoard.jsx`

React might be caching the old player data.

**Check:** Verify polling is actually fetching new data every 2 seconds.

---

## Quick Fix Options

### Option 1: Remove Card Hiding Logic (Simple)
Since cards are only shown in UI conditionally, we don't need to hide them in the backend:

```go
// Remove this entire block
// if g.Phase == 1 && players[i].UserID != userID {
//     players[i].Cards = []cards.Card{}
// }
```

**Pros:** Simple, no hiding logic needed  
**Cons:** Sends unnecessary data in Phase 1

### Option 2: Clone Player Data Before Modifying (Better)
Don't modify the original players slice:

```go
// Create response players with conditional visibility
responsePlayers := make([]models.GamePlayer, len(players))
for i := range players {
    responsePlayers[i] = players[i]  // Copy
    user, _ := s.userRepo.GetByID(players[i].UserID)
    responsePlayers[i].User = user
    
    if g.Phase == 1 && players[i].UserID != userID {
        responsePlayers[i].Cards = []cards.Card{}  // Hide in copy only
    }
}

return &models.GameStateResponse{
    Game:      *g,
    Players:   responsePlayers,  // Use modified copy
    Actions:   actions,
    DeckCount: len(g.DeckCards),
}, nil
```

### Option 3: Add Visibility Flag (Most Flexible)
Add a field to indicate if cards should be visible:

```go
type GamePlayerResponse struct {
    models.GamePlayer
    CardsVisible bool `json:"cards_visible"`
}
```

Then in frontend:
```javascript
{player.cards_visible && player.cards.map(...)}
```

---

## Testing the Fix

### Test 1: Check API Response Directly
```bash
# Get game state (replace token and game ID)
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/v1/games/1 \
     | jq '.players[] | {user_id, card_count, cards_length: (.cards | length)}'
```

**Expected in Phase 2:**
```json
{"user_id": 1, "card_count": 5, "cards_length": 5}
{"user_id": 2, "card_count": 7, "cards_length": 7}
```

### Test 2: Frontend Console Check
```javascript
// After API response
console.log('Phase:', game.phase);
players.forEach(p => {
  console.log(p.user.username, '- Cards:', p.cards.length, 'Expected:', p.card_count);
});

// All should match in Phase 2
```

### Test 3: Visual Check
In Phase 2, you should see:
- ✅ Your cards in a grid
- ✅ Other players' cards in grids
- ✅ "🃏 Hand Visible" label on all players
- ✅ Small cards displayed (not just top card)

---

## Recommended Fix (Option 2)

I recommend **Option 2** - clone the data before modifying. This ensures:
1. Original database data is never modified
2. Phase 1 hiding still works
3. Phase 2 shows all cards
4. No data leakage issues

Let me implement this fix now...

---

**Status:** 🔍 Investigating  
**Next Step:** Implement Option 2 fix
