# 🎮 Card Placement Rules - Updated Implementation

## Overview
This document describes the corrected card placement rules for Phase 1 of the Fasiolas card game.

## Core Rules

### 1. Placing Cards on Opponents
- **Requirement**: The +1 rule MUST apply (your card must be exactly 1 rank higher than opponent's top card)
- **Example**: You can place Jack (11) on opponent's 10 ✅
- **Turn Continuation**: When you place a card on an opponent with +1 rule:
  - ✅ Turn CONTINUES
  - 🎲 You MUST draw another card
  - 🔄 Repeat until you place on yourself or can't place anywhere

### 2. Placing Cards on Yourself
- **Requirement**: You can ALWAYS place ANY card on yourself (no +1 validation required)
- **Turn Continuation Rules**:
  - **If +1 rule applies** (your card is +1 from your previous top card):
    - ✅ Turn CONTINUES
    - 🎲 You CAN draw another card
  - **If +1 rule does NOT apply**:
    - ❌ Turn ENDS
    - ➡️ Next player's turn

### 3. Drawing Cards
- You draw when:
  - You can't place on any opponent
  - You placed on an opponent (must draw again)
  - You placed on yourself with +1 rule (can draw again)

## Examples

### Example 1: Placing J♦ on Opponent
```
Your pile: [5♣, 8♦, J♦]  (top card: J♦)
Opponent pile: [7♠, 10♥] (top card: 10♥)

Action: Click opponent's pile to place J♦

Validation:
- Is J♦ (11) = 10♥ (10) + 1?
- YES ✅ (11 = 10 + 1)

Result:
- J♦ is placed on opponent's pile
- Opponent now has: [7♠, 10♥, J♦]
- Your pile: [5♣, 8♦]
- **Turn CONTINUES** - You MUST draw another card
```

### Example 2: Placing Any Card on Yourself (No +1)
```
Your pile: [5♣, 8♦]     (top card: 8♦)
Drawn card: J♦

Action: Click your own pile to place J♦

Validation:
- Placing on self? YES ✅ (always allowed)
- Does +1 rule apply? Is J♦ (11) = 8♦ (8) + 1? 
- NO (11 ≠ 9)

Result:
- J♦ is placed on your pile
- Your pile: [5♣, 8♦, J♦]
- **Turn ENDS** - Next player's turn
```

### Example 3: Placing Card on Yourself (With +1)
```
Your pile: [5♣, 8♦]     (top card: 8♦)
Drawn card: 9♥

Action: Click your own pile to place 9♥

Validation:
- Placing on self? YES ✅ (always allowed)
- Does +1 rule apply? Is 9♥ (9) = 8♦ (8) + 1?
- YES ✅ (9 = 8 + 1)

Result:
- 9♥ is placed on your pile
- Your pile: [5♣, 8♦, 9♥]
- **Turn CONTINUES** - You CAN draw another card
```

### Example 4: Chain Placement on Opponents
```
Setup:
- Your pile: [5♣, 6♦, 7♥, 8♠] (top: 8♠)
- Opponent A: [9♣]
- Opponent B: [10♥]

Turn sequence:
1. Place 8♠ on opponent's 7? Can't find any opponent with 7
2. You have 8♠, Opponent A has 9♣
3. Can't place 8 on 9 (need +1 rule: 8→7 not 8→9)
4. Must draw from deck
5. Draw 10♦
6. Can place 10♦ on Opponent A's 9♣? YES ✅ (10 = 9 + 1)
7. Place 10♦ on Opponent A → Turn continues
8. Must draw again...
```

## Implementation Details

### Backend Changes

#### 1. New Function: `DoesPlacementApplyPlusOneRule`
**File**: `internal/game/engine.go`

```go
// Checks if placing currentPlayer's top card on target follows +1 rule
func (e *Engine) DoesPlacementApplyPlusOneRule(
    currentPlayer models.GamePlayer, 
    targetPosition int, 
    allPlayers []models.GamePlayer
) bool
```

**Purpose**: Determines if the +1 rule applies to a placement, which controls turn continuation.

#### 2. Updated Function: `CanPlaceCardOnTarget`
**File**: `internal/game/engine.go`

**Key Changes**:
- Placing on **self**: ALWAYS returns `true` (no validation)
- Placing on **opponent**: Returns `true` only if +1 rule applies

#### 3. Updated Function: `PlaceCard`
**File**: `internal/service/game_service.go`

**Key Changes**:
- Checks `plusOneApplies` BEFORE placing card
- Determines `isPlacingOnSelf` to decide turn continuation
- Turn continuation logic:
  ```go
  if isPlacingOnSelf {
      if !plusOneApplies {
          shouldEndTurn = true  // Turn ends
      }
      // else: turn continues, can draw again
  } else {
      // Placed on opponent with +1 rule
      // Turn continues, must draw again
  }
  ```

### API Response Changes

The `/api/v1/games/:id/place` endpoint now logs additional information:
```json
{
  "target_position": 1,
  "plus_one_rule": true,
  "on_self": false
}
```

## Testing Scenarios

### Test 1: Place on Opponent (Turn Continues)
1. Create game with 2 players
2. Player 1 draws card that is +1 from opponent's top card
3. Player 1 places on opponent
4. **Verify**: Turn does NOT end, player can draw again

### Test 2: Place on Self Without +1 (Turn Ends)
1. Player draws card that is NOT +1 from their top card
2. Player places on self
3. **Verify**: Turn ends, next player's turn begins

### Test 3: Place on Self With +1 (Turn Continues)
1. Player draws card that IS +1 from their top card
2. Player places on self
3. **Verify**: Turn continues, player can draw again

### Test 4: Cannot Place Invalid Card on Opponent
1. Player tries to place card on opponent where +1 rule doesn't apply
2. **Verify**: API returns error "invalid placement: card cannot be placed on target"

## Frontend Considerations

### UI Changes Needed
1. **Turn Indicator**: Show "Your turn - Draw another card" when turn continues after placing on opponent
2. **Action Buttons**: 
   - After placing on opponent with +1: Automatically enable "Draw" button
   - After placing on self without +1: Disable all action buttons
3. **Visual Feedback**: Highlight available actions based on game state

### User Experience Flow
```
Player's Turn:
├─ Has card in hand
│  ├─ Click opponent's pile
│  │  ├─ Valid (+1 applies) → Card placed → "Draw again"
│  │  └─ Invalid → Error message
│  └─ Click own pile
│     ├─ Placed with +1 → "You can draw again"
│     └─ Placed without +1 → "Turn ended" → Next player
└─ No card in hand → Click "Draw" button
```

## Migration Notes

### Breaking Changes
⚠️ **IMPORTANT**: This update changes core game mechanics

**Before**:
- Placing any card on self ended turn
- Turn continuation was based on `CanPlaceCard` check

**After**:
- Placing card on self with +1 allows turn continuation
- Placing card on opponent with +1 requires drawing again
- Turn continuation is based on +1 rule application

### Database Schema
No database schema changes required. All changes are in game logic only.

## Summary

| Scenario | Requirement | Turn Continuation |
|----------|-------------|-------------------|
| Place on opponent | +1 rule MUST apply | ✅ Continues (must draw) |
| Place on self (with +1) | Always allowed | ✅ Continues (can draw) |
| Place on self (no +1) | Always allowed | ❌ Ends |
| Draw card | No placement possible | Depends on drawn card |

---

**Last Updated**: January 25, 2026  
**Version**: 2.0  
**Status**: ✅ Implemented
