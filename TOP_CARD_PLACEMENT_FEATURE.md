# ✅ NEW FEATURE - Place Top Card on Opponents

## 🎯 Feature Added

Players can now place their **top card** directly on opponents' piles if the **+1 rule applies**!

### What This Means

You don't always need to draw a card. If your current top card is +1 from an opponent's top card, you can place it directly without drawing first.

## 📋 Placement Rules (Complete)

### 1. Placing a Drawn Card
**After drawing a card:**
- ✅ Can place on **anyone** (self or opponent)
- ✅ No validation required
- ✅ +1 rule determines turn continuation

### 2. Placing Your Top Card (NEW!)
**Without drawing:**
- ✅ Can place on **yourself** (always allowed)
- ✅ Can place on **opponent** ONLY if +1 rule applies
- ✅ +1 rule determines turn continuation

## 🎮 Examples

### Example 1: Place Top Card on Opponent ✅
```
Your top card: 7♦
Opponent's top: 6♥

Action: Click opponent's pile
Check: 7 = 6 + 1? YES ✅
Result: Card placed successfully!
→ Turn continues if you want
```

### Example 2: Cannot Place Top Card (Invalid) ❌
```
Your top card: 10♠
Opponent's top: 6♥

Action: Click opponent's pile
Check: 10 = 6 + 1? NO ❌
Result: Error - "can only place on opponent if card is +1"
→ Must draw a card first
```

### Example 3: Place Top Card on Yourself ✅
```
Your pile: [5♣, 8♦] (top: 8♦)
Action: Click your own pile

Result: Always allowed ✅
→ Turn ends (no +1 rule)
```

### Example 4: Strategic Chain Placement
```
Your top: 7♦
Opponent A top: 6♥
Opponent B top: 7♠

Turn flow:
1. Place 7♦ on Opponent A (6♥) → 7 = 6 + 1? YES ✅
2. Your new top: (card underneath)
3. Can continue turn or draw
```

## 🎲 Turn Flow Options

### Option A: Draw First
```
1. Draw a card
2. Place on anyone (no validation)
3. If +1 applies → Continue
4. If no +1 → Turn ends
```

### Option B: Place Top Card First (NEW!)
```
1. Check if your top is +1 from opponent's
2. If YES: Place directly (no draw needed)
3. If +1 still applies → Continue
4. If no +1 → Turn ends
```

## 📊 Validation Matrix

| Scenario | Placing | Validation | Allowed? |
|----------|---------|------------|----------|
| **Drawn card → Self** | Any card | None | ✅ YES |
| **Drawn card → Opponent** | Any card | None | ✅ YES |
| **Top card → Self** | Any card | None | ✅ YES |
| **Top card → Opponent** | +1 card | Must be +1 | ✅ YES |
| **Top card → Opponent** | Non-+1 | Must be +1 | ❌ NO |

## 🔧 Implementation Details

### Code Changes

**File**: `internal/service/game_service.go`  
**Function**: `PlaceCard()`

**Added Validation:**
```go
// When placing your top card (not a drawn card) on an opponent, +1 rule MUST apply
if !isPlacingDrawnCard && !isPlacingOnSelf {
    // Check if +1 rule applies
    if targetPlayer.TopCard == nil || !currentPlayer.TopCard.IsOnePlus(*targetPlayer.TopCard) {
        return errors.New("can only place on opponent if card is +1 from their top card")
    }
}
```

### What Changed

**Before:**
- Only drawn cards could be placed
- Top card placement was limited

**After:**
- ✅ Drawn cards: Place on anyone (no validation)
- ✅ Top card on self: Always allowed
- ✅ Top card on opponent: +1 rule required

## 🚀 Strategic Benefits

### More Flexibility
- Don't need to draw if you already have a placeable card
- Can chain multiple placements
- More strategic options

### Faster Gameplay
- Skip unnecessary draws
- Direct placements when possible
- More dynamic turns

### Better Strategy
- Plan ahead with your pile
- Look for +1 opportunities
- Control your turn length

## 🧪 How to Test

### Test Case 1: Place Top Card on Opponent
1. **Setup**: Your top is 7♦, opponent has 6♥
2. **Action**: Click opponent's pile (without drawing)
3. **Expected**: ✅ Card placed successfully

### Test Case 2: Try Invalid Top Card Placement
1. **Setup**: Your top is 10♠, opponent has 6♥
2. **Action**: Click opponent's pile (without drawing)
3. **Expected**: ❌ Error message

### Test Case 3: Draw Then Place
1. **Setup**: Cannot place top card
2. **Action**: Draw a card first
3. **Result**: ✅ Can place drawn card on anyone

## 📝 Error Messages

### "can only place on opponent if card is +1 from their top card"
**Meaning**: You're trying to place your top card on an opponent, but it's not +1 from their top card.

**Solution**: 
- Draw a card instead (drawn cards can be placed anywhere)
- Or place on yourself
- Or wait for your turn to end

## 🎯 Game Flow Summary

```
YOUR TURN:

Option 1: Place Top Card
├─ Can place on self? → YES (always)
├─ Can place on opponent? → Only if +1
└─ Turn continues? → Check +1 rule

Option 2: Draw Card
├─ Draw from deck
├─ Can place on anyone? → YES (always)
└─ Turn continues? → Check +1 rule

Repeat until turn ends (no +1 rule applies)
```

## 🚀 Deployment Status

- ✅ Feature implemented
- ✅ Backend rebuilt
- ✅ Docker restarted
- ✅ Server running
- ✅ Ready to use!

## 🎮 Try It Now!

**URL**: http://localhost:3000

1. Start a game
2. After first card is dealt
3. Check if your top is +1 from opponent's
4. If YES: Click opponent's pile directly ✅
5. If NO: Draw a card first

---

**Status**: 🟢 **LIVE**  
**Date**: January 25, 2026  
**Feature**: Top Card Placement on Opponents
