# Phase 1 Game Mechanics - Complete Implementation

## Overview
Phase 1 is called **"Kaupimas"** (Stacking/Accumulation). The goal is to get rid of all your cards by placing them on other players' piles.

## Game Setup (Before Phase 1)

1. **Deck**: Standard 52 cards (no jokers)
2. **Card Order**: 2 < 3 < 4 < 5 < 6 < 7 < 8 < 9 < 10 < J < Q < K < A
3. **Distribution**: Each player gets 1 card face up (starting pile)
4. **First Player**: The player with the lowest-ranking card starts
5. **Remaining deck**: Placed in center (can be drawn from)

## The +1 Rule (Cyclic)

Cards can only be placed if they are exactly +1 from the target card's rank:

```
2 → 3 → 4 → 5 → 6 → 7 → 8 → 9 → 10 → J → Q → K → A → 2 (cycles)
```

**Examples of valid placements:**
- 5 can go on 4 ✅
- 6 can go on 5 ✅
- 10 can go on 9 ✅
- J can go on 10 ✅
- Q can go on J ✅
- K can go on Q ✅
- A can go on K ✅
- 2 can go on A ✅ (cyclic)

**Invalid placements:**
- 5 cannot go on 3 ❌
- 5 cannot go on 6 ❌
- A cannot go on 2 ❌

## Turn Structure - Phase A (Place Before Draw)

When it's your turn, you **must** follow this sequence:

### Step A1: Check if Can Place on Others
```
Do you have valid placement for ANY other player?
├─ YES → Place on one of them (mandatory)
│        └─ Go back to A1 (check again with new top card)
└─ NO  → Go to A2
```

**What is "valid placement"?**
- Your top card is +1 from another player's top card
- You can place on ANY player who meets this rule
- You CHOOSE which one if multiple options

### Step A2: Check if Can Place on Yourself
```
Can you place on yourself?
├─ YES → Place on yourself (mandatory)
│        └─ MUST draw (go to Phase B)
└─ NO  → Go to A3 (must draw)
```

### Step A3: Must Draw
```
Can't place anywhere?
└─ MUST draw from deck (go to Phase B)
```

## Turn Structure - Phase B (Draw & Auto-Place)

After Phase A, if you must draw:

### Step B1: Draw Card
- Draw 1 card from deck
- Card is visible to ALL players
- **The card is held temporarily** (not in your pile yet)

### Step B2: Check Where Drawn Card Can Go
```
Does drawn card match +1 rule?
├─ B1: Can go on OTHER player?
│      └─ YES → Place on them (mandatory, turn ends)
├─ B2: Can go on YOURSELF?
│      └─ YES → Place on yourself (mandatory, turn ends)
└─ B3: Can't go anywhere?
       └─ Add to your pile (pasiimti - "take for yourself")
          └─ Turn ends
```

## Example Turn Sequence

**Setup:**
- Player A has top card: 5
- Player B has top card: 6
- Player C has top card: 4
- Deck has 32 cards

**Player A's Turn (Phase A):**
1. Check: Can A's 5 go on others?
   - On B? 5→6? NO (needs +1, which is 7)
   - On C? 5→4? NO (needs 3)
   - Answer: NO, can't place on others

2. Check: Can A place on self?
   - A has [5] (only visible card)
   - Answer: NO
   - Action: MUST DRAW

**Player A's Turn (Phase B):**
1. Draw from deck → Gets **6**
2. Check where 6 can go:
   - On B? B has 6, need 7. NO
   - On C? C has 4, need 5. NO
   - On A? A's pile has 5 underneath. 6→5? YES! (+1 rule)
   - Action: Place 6 on self

3. Turn ends. Deck now has 31 cards.

**Updated state:**
- Player A: [5, 6] top card: 6
- Player B: [6] top card: 6
- Player C: [4] top card: 4

## Key Rules Summary

1. **Mandatory Placement**
   - You MUST place on others if possible
   - You MUST place on self if can't place on others
   - You MUST draw if can't place anywhere

2. **Drawn Cards**
   - Drawn card is visible to all players
   - You must place it according to B1/B2/B3
   - You cannot choose to keep it if rules allow placement

3. **Turn Flow**
   - Turn = Phase A (place if possible) + Phase B (draw if needed)
   - Turn ends when:
     - You place on another player, OR
     - You place on yourself OR you keep a drawn card (pasiimti)

4. **Winning Phase 1**
   - You get rid of all cards
   - When you place your last card, you "get out"
   - Last player remaining holds all cards

## Penalties & Cheating

### What is Cheating?
1. Drawing when you can place ❌
2. Placing a card that doesn't follow +1 rule ❌
3. Placing on wrong player ❌

### Consequences
- Called out by another player
- Must draw additional penalty cards
- Could be disqualified from game

## UI Elements Explained

### Deck Button (Center)
```
     🂠
    Deck
    32 left
```
- Shows number of remaining cards in deck
- Click to draw a card (only when it's your turn)
- Only clickable if no card currently drawn

### Drawn Card Display
```
  [Card shown]
  Drawn Card
```
- Appears above deck after drawing
- Shows to all players what was drawn
- Indicates selection required for placement

### Player Areas
```
┌──────────────────┐
│  Player Name     │
│  🎯 Playing      │  (if current turn)
│  [Top Card]      │
│  🂠 5 cards      │  (card count)
└──────────────────┘
```
- Shows each player's top card
- Colored border indicates status:
  - **Green**: Currently playing
  - **Blue**: You (your player)
  - **Yellow**: Selected target (ready to place)

### Action Panel (When Your Turn)
```
🎮 Your Turn
TAP DECK → Draw
SELECT PLAYER → Choose target
PLACE CARD → Confirm placement
```
- Only visible when it's your turn
- Guide for next action
- Updates as game state changes

## Automatic Features

✅ **Automatic turn detection**
- System knows whose turn it is
- Only your controls are enabled

✅ **Automatic turn passing**
- After placement, passes to next player
- Updates every 2 seconds

✅ **Automatic card counting**
- Card count always accurate
- Updates after every placement

✅ **Automatic state validation**
- Can't place invalid cards
- Can't draw if can place
- Can't select invalid targets

## Phase Transition

Phase 1 ends when:
- All players have been processed
- Approximately 5 rounds completed
- OR when game decides transition (based on specific rules)

Phase 2 begins:
- Different rules apply
- Table cards come into play
- Trump suit matters

---

**Remember**: The golden rule of Phase 1 is **"Must place if possible"**
- If you can place on others → you MUST
- If you can place on self → you MUST (after trying others)
- If you can place drawn card → you MUST

