# 🎮 COMPLETE FIX - Card Drawing and Placement

## Problem in Screenshot
- You pressed "Draw Card" button
- Deck shows "Deck is empty!" 
- No card appeared in the middle
- Cannot place the drawn card

## Root Cause Analysis

### Issue 1: Deck Empty Check
The frontend might be showing "Deck is empty!" even when deck has cards because the game state polling hasn't updated yet.

### Issue 2: Card Not Appearing in Middle
After drawing, the card should appear in the middle but it's not showing up. This could be because:
1. The backend returns the card but frontend doesn't set it correctly
2. The polling overwrites the drawnCard state
3. The refs aren't being set properly

## Complete Fix Applied

### Backend Changes (game_service.go)

#### 1. DrawCard Function
```go
// Store drawn card in TableCards[0] instead of player's pile
g.TableCards = []cards.Card{*card}

// Update game in database
if err := s.gameRepo.Update(g); err != nil {
    return nil, fmt.Errorf("failed to update game: %w", err)
}

// Return card to frontend
return card, nil
```

**What this does:**
- ✅ Draws card from deck
- ✅ Stores in `game.TableCards[0]` (NOT in player's pile)
- ✅ Saves to database
- ✅ Returns card to frontend

#### 2. PlaceCard Function
```go
// Check if there's a drawn card waiting to be placed
isPlacingDrawnCard := len(g.TableCards) > 0

if isPlacingDrawnCard {
    // Use the drawn card from TableCards[0]
    cardToPlace = &g.TableCards[0]
    
    // Add it to player's pile temporarily for validation
    currentPlayer.Cards = append(currentPlayer.Cards, *cardToPlace)
    currentPlayer.TopCard = cardToPlace
} else {
    // Use player's starting hand TopCard
    cardToPlace = currentPlayer.TopCard
}

// Validate and place card...

// Clear TableCards after placing
if isPlacingDrawnCard {
    g.TableCards = []cards.Card{}
}
```

**What this does:**
- ✅ Checks if there's a drawn card in TableCards
- ✅ If yes, uses that card for placement
- ✅ If no, uses player's TopCard (starting hand)
- ✅ Validates placement
- ✅ Clears TableCards after successful placement

### Frontend Changes (GameBoard.jsx)

#### 1. Added useRef for Persistent State
```javascript
const isPlacingCardRef = useRef(false);
const drawnCardRef = useRef(null);
```

**Why:** To persist state across re-renders caused by polling.

#### 2. handleDrawCard Updates
```javascript
const handleDrawCard = async () => {
    setWaitingForPlacement(true);
    isPlacingCardRef.current = true;
    
    const res = await axios.post(`/api/v1/games/${game.id}/draw`);
    const drawn = res.data?.card;
    
    setDrawnCard(drawn || null);
    drawnCardRef.current = drawn || null;
    
    // Don't call onUpdate - prevents showing card in pile
};
```

**What this does:**
- ✅ Sets state flags
- ✅ Draws card from backend
- ✅ Stores in both state and ref
- ✅ Does NOT call onUpdate (prevents pile update)

#### 3. TopCard Display Logic
```javascript
const shouldHideCard = isYou && (drawnCard || isPlacingCardRef.current || drawnCardRef.current);

if (shouldHideCard) {
    return <div>(placing above)</div>;
} else if (player.top_card) {
    return <PlayingCard card={player.top_card} />;
}
```

**What this does:**
- ✅ Checks multiple conditions (state + ref)
- ✅ Hides TopCard while placing
- ✅ Shows "(placing above)" message
- ✅ Persists across polling re-renders

## How to Test

### Step 1: Restart Docker
Run this PowerShell script:
```powershell
.\restart-and-test.ps1
```

Or manually:
```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose down
docker compose up -d --build
```

### Step 2: Create NEW Game
**IMPORTANT:** Use a NEW game, not the old one from your screenshot!

1. Go to http://localhost:3000/login
2. Log in with Google
3. Click "Create Game"
4. Set max players to 2
5. Copy room code
6. In incognito window, log in with another account
7. Join with room code
8. Start game

### Step 3: Test Drawing
1. Wait for your turn (green border, "YOUR TURN!" message)
2. Click "Draw Card" button (NOT the deck in the middle)
3. **Expected:** Card appears in the middle with yellow border
4. **Expected:** Your pile shows "(placing above)" instead of the card
5. **Expected:** Deck count decreases by 1

### Step 4: Test Placement
1. Drag the card to a player OR click on a player
2. **Expected:** Card moves to their pile
3. **Expected:** Your turn ends
4. **Expected:** Next player's turn begins

## Troubleshooting

### If "Deck is empty!" shows
- This means `game.deck_count === 0`
- Check if you're using an OLD game from before the fix
- Create a NEW game and test again

### If card doesn't appear in middle
1. Open browser console (F12)
2. Look for errors
3. Check if `drawnCard` state is set:
   - Should see: `drawnCard: {rank: "X", suit: "Y"}`
4. Check if backend returned card:
   - Look at Network tab
   - Find POST `/api/v1/games/{id}/draw`
   - Check response: should have `{card: {...}}`

### If card appears in both places
1. Check console logs:
   - Should see: `refPlacing: true`
   - Should see: `shouldHide: true`
2. If `shouldHide: false`, the ref isn't working
3. Try hard refresh: Ctrl+Shift+R

### If placement doesn't work
1. Check console for errors
2. Verify POST `/api/v1/games/{id}/place` request
3. Check if `target_player_position` is sent correctly
4. Look at backend logs: `docker compose logs app --tail 50`

## Game Flow Summary

### Correct Flow:
1. **Draw:** Click "Draw Card" button
   - Card removed from deck ✅
   - Card stored in `game.TableCards[0]` ✅
   - Card returned to frontend ✅
   - Card shows in middle ✅
   - Your pile shows "(placing above)" ✅

2. **Place:** Drag or click to place
   - Card taken from `TableCards[0]` ✅
   - Validation checked ✅
   - Card added to target pile ✅
   - `TableCards` cleared ✅
   - Turn ends if needed ✅

### Incorrect Flow (Before Fix):
1. Draw → Card added to YOUR pile immediately ❌
2. Card shows in BOTH middle AND your pile ❌
3. Polling updates make it worse ❌

## Files Modified

1. **internal/service/game_service.go**
   - `DrawCard()` - stores in TableCards
   - `PlaceCard()` - reads from TableCards

2. **frontend/src/components/GameBoard.jsx**
   - Added `useRef` for persistent state
   - Updated `handleDrawCard()` 
   - Updated `handlePlaceCard()`
   - Updated TopCard display logic

## Next Steps

1. Run `.\restart-and-test.ps1`
2. Create a NEW game
3. Test the complete flow
4. If issues persist, check browser console + backend logs
5. Report any errors you see

---

**All fixes are complete and ready to test!** 🎉

