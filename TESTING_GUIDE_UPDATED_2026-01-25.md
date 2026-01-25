# 🧪 Testing Guide: Card Placement Rules (Updated)

## Quick Start Testing

### Prerequisites
✅ Docker containers running (`docker-compose up -d`)  
✅ Backend on `http://localhost:8080`  
✅ Frontend on `http://localhost:3000`

---

## Test Case 1: Place Card on Opponent (Turn Continues)

### Setup
1. Login to game at `http://localhost:3000`
2. Create a new game
3. Join with 2 players
4. Start the game

### Scenario
```
Player 1 has: 7♠
Player 2 has: 6♥
```

### Steps
1. **Player with 6♥ draws**: 7♦
2. **Click on opponent's pile** (Player with 7♠)
3. **Expected**: 7♦ (7) = 7♠ (7) + 1? NO
4. **Result**: ❌ Should show error "invalid placement"

### Correct Scenario
```
Player 1 has: 7♠
Player 2 has: 6♥
```

1. **Player with 7♠ draws**: 7♦
2. **Click on opponent's pile** (Player with 6♥)
3. **Expected**: 7♦ (7) = 6♥ (6) + 1? YES ✅
4. **Result**: 
   - ✅ Card placed on opponent
   - 🔄 Turn continues
   - 🎲 Can draw another card

---

## Test Case 2: Place Card on Self (No +1 Rule - Turn Ends)

### Scenario
```
Your pile: [5♣, 8♦] (top: 8♦)
Draw: J♦ (value: 11)
```

### Steps
1. **Draw card**: J♦
2. **Click on your own pile**
3. **Expected**: 
   - Is placement on self? YES ✅ (always allowed)
   - Does J♦ (11) = 8♦ (8) + 1? NO (11 ≠ 9)
4. **Result**:
   - ✅ Card placed on your pile
   - ❌ Turn ENDS
   - ➡️ Next player's turn

### Verification
- Check that `current_player_position` changed in game state
- Your pile now has 3 cards: [5♣, 8♦, J♦]
- It's the next player's turn

---

## Test Case 3: Place Card on Self (+1 Rule - Turn Continues)

### Scenario
```
Your pile: [5♣, 8♦] (top: 8♦)
Draw: 9♥ (value: 9)
```

### Steps
1. **Draw card**: 9♥
2. **Click on your own pile**
3. **Expected**:
   - Is placement on self? YES ✅
   - Does 9♥ (9) = 8♦ (8) + 1? YES ✅ (9 = 8 + 1)
4. **Result**:
   - ✅ Card placed on your pile
   - ✅ Turn CONTINUES
   - 🎲 Can draw another card

### Verification
- Check that `current_player_position` DID NOT change
- Your pile now has 3 cards: [5♣, 8♦, 9♥]
- It's still your turn
- Deck button should be clickable

---

## Test Case 4: Multiple Placements on Opponents (Chain)

### Scenario
```
Setup:
- Player 1: [7♠]
- Player 2: [6♥]
- Player 3: [9♦]
```

### Steps (Player 1's Turn)
1. **Draw**: 7♦
2. **Try to place on Player 2** (6♥)
   - 7♦ (7) = 6♥ (6) + 1? YES ✅
3. **Result**: Card placed, turn continues
4. **Draw again**: 10♣
5. **Try to place on Player 3** (9♦)
   - 10♣ (10) = 9♦ (9) + 1? YES ✅
6. **Result**: Card placed, turn continues
7. **Draw again**: 2♠
8. **Try to place on anyone**
   - Player 2 now has 7♦ (need 8)
   - Player 3 now has 10♣ (need J)
   - Can't place on anyone
9. **Click own pile** (7♠)
   - 2♠ (2) = 7♠ (7) + 1? NO
10. **Result**: Turn ends

---

## API Testing with Postman/cURL

### 1. Get Game State
```bash
GET http://localhost:8080/api/v1/games/{gameID}
Headers: Authorization: Bearer {token}
```

**Check**:
- `current_player_position`: Who's turn it is
- `table_cards`: Any drawn card waiting
- `players[].top_card`: Each player's top card

### 2. Draw Card
```bash
POST http://localhost:8080/api/v1/games/{gameID}/draw
Headers: Authorization: Bearer {token}
```

**Response**:
```json
{
  "card": {
    "suit": "hearts",
    "rank": "9",
    "value": 9
  }
}
```

### 3. Place Card
```bash
POST http://localhost:8080/api/v1/games/{gameID}/place
Headers: 
  Authorization: Bearer {token}
  Content-Type: application/json
Body:
{
  "target_player_position": 0
}
```

**Success Response**:
```json
{
  "message": "card placed"
}
```

**Error Response** (Invalid placement):
```json
{
  "error": "invalid placement: card cannot be placed on target"
}
```

---

## Debugging Tips

### Check Action Logs
```bash
docker logs fasiolas_app | grep "ActionPlaceCard"
```

Look for:
```json
{
  "target_position": 1,
  "plus_one_rule": true,
  "on_self": false
}
```

### Check Game State Updates
```bash
docker logs fasiolas_app | grep "current_player_position"
```

### Check Turn Continuation
After placing card:
1. Get game state
2. Check if `current_player_position` changed
3. If same: Turn continues ✅
4. If different: Turn ended ✅

---

## Expected Behaviors Summary

| Action | +1 Rule? | Result | Turn Continues? |
|--------|----------|--------|-----------------|
| Place on opponent | ❌ NO | Error | N/A |
| Place on opponent | ✅ YES | Success | ✅ YES |
| Place on self | ❌ NO | Success | ❌ NO |
| Place on self | ✅ YES | Success | ✅ YES |

---

## Common Issues

### Issue: "invalid placement" when trying to place on self
**Cause**: Should never happen (always allowed)  
**Fix**: Check `CanPlaceCardOnTarget` returns `true` for self

### Issue: Turn ends when it shouldn't
**Cause**: `plusOneApplies` logic incorrect  
**Fix**: Verify `DoesPlacementApplyPlusOneRule` is called correctly

### Issue: Turn continues when it shouldn't
**Cause**: Turn end logic not triggered  
**Fix**: Check `shouldEndTurn` logic in `PlaceCard` service

---

## Browser Console Testing

Open browser console (`F12`) and watch for logs:
```javascript
// Look for these messages:
"✅ Card placed successfully"
"❌ Failed to place card"
"🎯 Turn continues - draw again"
"➡️ Turn ended - next player"
```

---

## Database Verification

### Check Player Cards
```sql
SELECT id, position, cards, top_card 
FROM game_players 
WHERE game_id = 1;
```

### Check Action Logs
```sql
SELECT * FROM game_actions 
WHERE game_id = 1 
ORDER BY created_at DESC 
LIMIT 10;
```

---

**Happy Testing!** 🎮

If you encounter any issues, check:
1. Docker logs: `docker logs fasiolas_app`
2. Browser console: `F12`
3. Network tab: Check API responses
4. Game state: Verify with GET `/api/v1/games/{id}`
