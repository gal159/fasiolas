# Visual Guide - New Game Table Layout

## 🎰 Table Layout Diagram

### 4-Player Game View
```
                    ┌─────────────────────────────────┐
                    │     PLAYER 2 (Top)              │
                    │                                 │
                    │     ┌─────────────────┐        │
                    │     │  Username 2     │        │
                    │     │  ┌───────────┐  │        │
                    │     │  │   ♠ K     │  │        │
                    │     │  │           │  │        │
                    │     │  │   ♠       │  │        │
                    │     │  └───────────┘  │        │
                    │     │  🂠 4 cards   │        │
                    │     └─────────────────┘        │
                    └─────────────────────────────────┘
                           (green border)
                           (it's their turn)

┌──────────────────────┐                ┌──────────────────────┐
│  PLAYER 1 (Left)     │                │  PLAYER 3 (Right)    │
│                      │                │                      │
│  ┌────────────────┐  │                │  ┌────────────────┐  │
│  │  Username 1    │  │                │  │  Username 3    │  │
│  │  ┌──────────┐  │  │                │  │  ┌──────────┐  │  │
│  │  │  ♥ 5    │  │  │                │  │  │  ♣ 7    │  │  │
│  │  │          │  │  │                │  │  │          │  │  │
│  │  │   ♥     │  │  │                │  │  │   ♣     │  │  │
│  │  └──────────┘  │  │                │  │  └──────────┘  │  │
│  │  🂠 6 cards  │  │                │  │  🂠 5 cards  │  │
│  └────────────────┘  │                │  └────────────────┘  │
│                      │                │                      │
└──────────────────────┘                └──────────────────────┘
(gray border)                            (gray border)


                    ┌─────────────────────────────────┐
                    │                                 │
                    │       ┌─────────────────┐      │
                    │       │  🂠              │      │
                    │       │ Deck             │      │
                    │       │ 32 cards         │      │
                    │       │  (clickable)     │      │
                    │       └─────────────────┘      │
                    │                                 │
                    │   (Drawn card appears here)    │
                    │                                 │
                    └─────────────────────────────────┘
                    (center of table)


                    ┌─────────────────────────────────┐
                    │    PLAYER 4 (Bottom - YOU)      │
                    │                                 │
                    │     ┌─────────────────┐        │
                    │     │  You            │        │
                    │     │  ┌───────────┐  │        │
                    │     │  │  ♦ 6     │  │        │
                    │     │  │           │  │        │
                    │     │  │  ♦        │  │        │
                    │     │  └───────────┘  │        │
                    │     │  🂠 7 cards   │        │
                    │     └─────────────────┘        │
                    └─────────────────────────────────┘
                    (blue border - your player)
```

## 🎬 Game Flow Visualization

### Turn Sequence Diagram
```
┌─────────────────────────────────────────────────────────────┐
│                    GAME STATE: Phase 1                       │
│              Current Player: Player 2 (Top)                  │
└─────────────────────────────────────────────────────────────┘
                              ↓
                    ┌──────────────────────┐
                    │  IT'S PLAYER 2'S TURN │
                    │  (green border)       │
                    │  ACTION: Draw Card    │
                    │  "TAP DECK"           │
                    └──────────────────────┘
                              ↓
                    ┌──────────────────────┐
                    │  PLAYER 2 DRAWS      │
                    │  Card = ♥ 5          │
                    │                      │
                    │  [Above Deck]        │
                    │     ┌──────────┐     │
                    │     │  ♥ 5    │     │
                    │     │          │     │
                    │     │   ♥     │     │
                    │     └──────────┘     │
                    │   Drawn Card         │
                    └──────────────────────┘
                              ↓
                    ┌──────────────────────┐
                    │  SELECT TARGET       │
                    │                      │
                    │  Click on Player who │
                    │  can receive this    │
                    │  card (♥ 5)          │
                    │                      │
                    │  Player 1 has ♥ 4   │
                    │  (5 can go on 4)     │
                    │                      │
                    │  →  Select Player 1  │
                    │      (yellow border) │
                    └──────────────────────┘
                              ↓
                    ┌──────────────────────┐
                    │  PLACE CARD          │
                    │                      │
                    │  Click "Place Card"  │
                    │  Button              │
                    │                      │
                    │  ♥ 5 is placed on    │
                    │  Player 1's pile     │
                    └──────────────────────┘
                              ↓
                    ┌──────────────────────┐
                    │  TURN COMPLETE       │
                    │                      │
                    │  Player 2's turn     │
                    │  ends                │
                    │                      │
                    │  State updates       │
                    │  (auto-refresh)      │
                    └──────────────────────┘
                              ↓
                    ┌──────────────────────┐
                    │  NEXT PLAYER'S TURN  │
                    │                      │
                    │  Passes to Player 3  │
                    │  (next in order)     │
                    │                      │
                    │  Player 3 now has    │
                    │  green border        │
                    │  Action panel active │
                    └──────────────────────┘
```

## 🎨 Color Coding

### Border Colors
```
GREEN BORDER   = ┌───────────────┐    Current player's turn
               │  🎯 Playing   │    (can take actions)
               └───────────────┘

BLUE BORDER    = ┌───────────────┐    Your player
               │  You          │    (always at bottom)
               └───────────────┘

YELLOW BORDER  = ┌───────────────┐    Selected target
               │  Selected     │    (ready to place card on)
               └───────────────┘

GRAY BORDER    = ┌───────────────┐    Other players
               │  Player 2     │    (waiting their turn)
               └───────────────┘
```

### Card Display
```
WHITE CARD     = ┌───────────────┐    Face-up card
               │   ♠ K        │    (visible to all)
               │               │
               │    ♠          │
               └───────────────┘

BLUE CARD      = ┌───────────────┐    Face-down deck
               │    🂠         │    (showing back)
               │   Deck        │
               │  32 cards     │
               └───────────────┘
```

## 📍 Position Reference

### 2-Player Game
```
      ┌─────────┐
      │ Player 1│  (Top)
      └─────────┘
         (Deck)
      ┌─────────┐
      │ You     │  (Bottom)
      └─────────┘
```

### 3-Player Game
```
        ┌─────────┐
        │ Player 1│  (Top)
        └─────────┘
    (Deck)
┌─────────┐     ┌─────────┐
│ Player 2│     │ Player 3│  (Left/Right)
└─────────┘     └─────────┘
    ┌─────────┐
    │ You     │  (Bottom)
    └─────────┘
```

### 4-Player Game (as shown in main diagram above)
```
        ┌─────────┐
        │ Player 2│  (Top)
        └─────────┘
    (Deck)
┌─────────┐     ┌─────────┐
│ Player 1│     │ Player 3│  (Left/Right)
└─────────┘     └─────────┘
    ┌─────────┐
    │ You     │  (Bottom)
    └─────────┘
```

## 🖱️ Interactive Elements

### Deck Button
```
Normal State:
┌──────────────────┐
│     🂠          │  ← Clickable
│    Deck         │     (when your turn)
│   32 cards      │
│    [TAP]        │
└──────────────────┘

Disabled State:
┌──────────────────┐
│     🂠          │  ← Greyed out
│    Deck         │     (when not your turn
│   32 cards      │      or card drawn)
│   (Disabled)    │
└──────────────────┘
```

### Drawn Card Display
```
Before Drawing:
(Nothing shown)

After Drawing:
    ┌──────────────┐
    │  ♥ 5       │  ← Drawn card
    │             │     visible above deck
    │   ♥        │
    └──────────────┘
    ┌──────────────┐
    │     🂠      │
    │    Deck     │
    │   32→31     │
    └──────────────┘

Message:
"Select a player to place the drawn card"
```

### Player Selection
```
Before Selection:
┌──────────────────┐
│  Player Name     │
│  [Card shown]    │  ← Click to select
│  🂠 5 cards    │
└──────────────────┘

After Selection:
┌──────────────────┐
│  Player Name     │
│  [Card shown]    │
│  🂠 5 cards    │
└──────────────────┘
(Yellow glow)

Selected Indicator:
"✓ Target selected: Player Name"
```

## 📊 State Transitions

### During Your Turn
```
IDLE STATE
    ↓
You see: "TAP DECK"
Click deck button
    ↓
DRAWING STATE
    ↓
Drawn card appears above deck
You see: "SELECT PLAYER"
    ↓
SELECTING STATE
    ↓
Click on a player (yellow border)
You see: "✓ Target selected: [Name]"
    ↓
CONFIRMING STATE
    ↓
Click "PLACE CARD" button
    ↓
ACTION SENT TO SERVER
    ↓
TURN ENDS
    ↓
Next player's turn begins
```

### Info Panel
```
┌────────────────────────────────┐
│ Phase: 1                       │  ← Current phase
│ Current Turn: Player 2         │  ← Whose turn
│ Deck: 31 cards                 │  ← Cards left
│                                │
│ 🎯 YOUR TURN!                  │  ← Shows when your turn
│    (pulsating animation)       │
└────────────────────────────────┘
```

## 🎭 Example Game Sequence

### Turn 1: Player 1
```
Player 1 draws ♥ 5
Deck: 52 → 51 cards
Can place on Player 2 (who has ♥ 4)
Places ♥ 5 on Player 2
Turn ends
```

### Turn 2: Player 2
```
Player 2 now has top card: ♥ 5
Draws ♣ 6
Deck: 51 → 50 cards
Cannot place ♣ 6 on anyone
Places on self
Turn ends
```

### Turn 3: Player 3
```
Player 3 draws ♠ 4
Deck: 50 → 49 cards
Can place on Player 1 (who has ♠ 3)
Places ♠ 4 on Player 1
Can check again? (New top card?)
Turn ends
```

And the game continues...

---

**This visual guide helps understand the new table layout and game flow!**

