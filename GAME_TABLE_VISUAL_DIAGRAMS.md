# 🎮 GAME TABLE - VISUAL DIAGRAMS

**Date:** January 24, 2026  
**Status:** ✅ Implementation Complete

---

## 1️⃣ TABLE LAYOUT DIAGRAM

### Circular Table (Your Diagram)
```
        ┌─────────────────────────────┐
        │  Opponents at Top            │
        │  [Card] [Card] [Card] ...   │
        │                             │
        │      Revealed Card           │
        │         [Big Card]           │
        │      (visible to all)        │
        │                             │
        │       [Deck 🎴]             │
        │   (Click to draw)           │
        │                             │
        │       You at Bottom          │
        │        [Your Card]           │
        └─────────────────────────────┘
```

### Real Table Comparison
```
REAL POKER TABLE:
        North (Opponent 2)
              [Card]
West                    East
[Card]      [Deck]    [Card]
        South (You)
              [Card]

FASIOLAS TABLE:
       Top (All Opponents)
        [Card] [Card]
              [Deck]
      Bottom (You)
              [Card]
```

---

## 2️⃣ GAME FLOW DIAGRAM

### Turn Sequence
```
START YOUR TURN
    ↓
[YOU SEE "YOUR TURN!" - Green highlight]
    ↓
    CLICK DECK (in center)
    ↓
    CARD IS DRAWN
    ↓
    [CARD REVEALED ABOVE DECK - BIG, ANIMATED]
    ↓
    [ALL PLAYER PILES LIGHT UP]
    ↓
    CLICK A PLAYER PILE
    ↓
    [CARD ANIMATES TO PILE]
    ↓
    CARD IS PLACED
    ↓
    [YOUR TURN ENDS]
    ↓
    TURN PASSES TO NEXT PLAYER
    ↓
    [OPPONENT GETS "YOUR TURN!" - Green highlight]
    ↓
    WAIT FOR YOUR NEXT TURN
```

### Card Placement Decision
```
DRAWN CARD

├─ Can place on OPPONENT? (Rule: +1)
│  ├─ YES → Click opponent's pile
│  │       → Card placed on opponent
│  │       → Turn ends
│  │
│  └─ NO → Can place on SELF? (Rule: +1 to own)
│         ├─ YES → Click your pile
│         │        → Card placed on self
│         │        → Turn ends
│         │
│         └─ NO → MUST DRAW
│                 → No choice
│                 → Turn ends
```

---

## 3️⃣ CIRCULAR LAYOUT STRUCTURE

### HTML/CSS Structure
```
┌─── CIRCULAR TABLE (rounded-full) ───┐
│                                      │
│  ┌──────────────────────────────┐  │
│  │ TOP: Opponents Area          │  │
│  │ (absolute top-8)            │  │
│  │ [Opp1] [Opp2] [Opp3] ...   │  │
│  └──────────────────────────────┘  │
│                                      │
│         ┌─────────────────┐         │
│         │  CENTER FLEX    │         │
│         │  [Drawn Card]   │         │
│         │   [Deck 🎴]     │         │
│         └─────────────────┘         │
│         (z-20: on top)              │
│                                      │
│  ┌──────────────────────────────┐  │
│  │ BOTTOM: You Area             │  │
│  │ (absolute bottom-8)          │  │
│  │      [Your Pile]             │  │
│  └──────────────────────────────┘  │
│                                      │
└──────────────────────────────────────┘
     (minHeight: 600px, aspectRatio: 1)
```

### Positioning Math
```
TOP POSITIONING:
<div class="absolute top-8 left-1/2 transform -translate-x-1/2">
                 ↑      ↑     ↑                           ↑
            8px from    center horizontally      centered (half-width)
            top

BOTTOM POSITIONING:
<div class="absolute bottom-8 left-1/2 transform -translate-x-1/2">
                 ↑        ↑     ↑                           ↑
            8px from    center horizontally      centered (half-width)
            bottom

CENTER POSITIONING:
<div class="flex items-center justify-center">
     ↑      ↑                           ↑
  use flex   vertical center        horizontal center
```

---

## 4️⃣ COLOR SCHEME DIAGRAM

### Visual States
```
PLAYER PILE STATES:

┌────────────────────────────────────┐
│ Your Pile (Not Current Player)     │
│ ┌──────────────────────────────┐  │
│ │ Blue Border  🔵              │  │
│ │ Blue Glow                    │  │
│ │ [Your Card]                  │  │
│ │ 📚 5 cards                   │  │
│ └──────────────────────────────┘  │
└────────────────────────────────────┘

┌────────────────────────────────────┐
│ Current Player's Pile              │
│ ┌──────────────────────────────┐  │
│ │ Green Border  🟢              │  │
│ │ Green Glow                   │  │
│ │ 🎯 Playing                   │  │
│ │ [Their Card]                 │  │
│ │ 📚 3 cards                   │  │
│ └──────────────────────────────┘  │
└────────────────────────────────────┘

┌────────────────────────────────────┐
│ Opponent Pile (When Card Drawn)    │
│ ┌──────────────────────────────┐  │
│ │ Gray Border (hover → Yellow)  │  │
│ │ [Opponent Card]              │  │
│ │ 📚 7 cards                   │  │
│ │ (Clickable)                  │  │
│ └──────────────────────────────┘  │
└────────────────────────────────────┘

┌────────────────────────────────────┐
│ Selected Target (You clicked here) │
│ ┌──────────────────────────────┐  │
│ │ Yellow Border  🟡             │  │
│ │ Yellow Glow + Ring           │  │
│ │ [Target Card]                │  │
│ │ 📚 4 cards                   │  │
│ │ (About to place here)        │  │
│ └──────────────────────────────┘  │
└────────────────────────────────────┘
```

### Color Meanings
```
🟢 GREEN GLOW     = Current player (your turn)
🔵 BLUE GLOW      = Your pile
🟡 YELLOW GLOW    = Selected target pile
⚫ GRAY BORDER     = Inactive (not involved)
🟡 YELLOW BORDER  = Drawn card (center)
🔵 BLUE BG        = Deck button
🟢 GREEN BG       = Table felt
```

---

## 5️⃣ RESPONSIVE DESIGN DIAGRAM

### Desktop (1200px+)
```
       ┌─────────────────────────────┐
       │    Full-size Table          │
       │    minHeight: 600px         │
       │                             │
       │ [Opp1] [Opp2] [Opp3]       │
       │                             │
       │        [Deck]               │
       │      [Card 🎴]             │
       │                             │
       │      [Your Pile]            │
       │                             │
       └─────────────────────────────┘
       Gap: 32px (gap-8)
       Padding: 32px (top-8, bottom-8)
       Text: Normal size (text-sm)
```

### Tablet (768px)
```
     ┌──────────────────────────┐
     │  Medium Table            │
     │  scales down 70%         │
     │                          │
     │ [Opp1] [Opp2] [Opp3]    │
     │                          │
     │       [Deck]             │
     │     [Card]              │
     │                          │
     │    [Your Pile]           │
     │                          │
     └──────────────────────────┘
     Gap: 16px (sm:gap-4)
     Padding: 16px (sm:top-8)
     Text: Medium size
```

### Mobile (< 768px)
```
   ┌──────────────────┐
   │  Compact Table   │
   │  fits screen     │
   │                  │
   │ [Opp1] [Opp2]   │
   │                  │
   │    [Deck]        │
   │   [Card]        │
   │                  │
   │  [Your Pile]     │
   │                  │
   └──────────────────┘
   Gap: 8px (gap-2)
   Padding: 16px (top-4)
   Text: Small (text-xs)
```

---

## 6️⃣ INTERACTION SEQUENCE DIAGRAM

### Click Deck Sequence
```
┌─────────────────────┐
│  Player clicks      │
│  deck button        │
└──────┬──────────────┘
       ↓
┌──────────────────────────────┐
│ Frontend:                    │
│ - Call handleDrawCard()      │
│ - Set waitingForPlacement=T  │
│ - API POST /draw            │
└──────┬───────────────────────┘
       ↓
┌──────────────────────────────┐
│ Backend:                     │
│ - Verify it's player's turn  │
│ - Draw from deck             │
│ - Validate placement (Phase1)│
│ - Return drawn card          │
└──────┬───────────────────────┘
       ↓
┌──────────────────────────────┐
│ Frontend:                    │
│ - Set drawnCard = response   │
│ - Display card above deck    │
│ - Highlight player piles     │
│ - Show: "Click to place"     │
└──────────────────────────────┘
```

### Click Pile Sequence
```
┌──────────────────────────┐
│  Player clicks           │
│  opponent pile           │
└──────┬───────────────────┘
       ↓
┌────────────────────────────────┐
│ Frontend:                      │
│ - Call handlePlaceCard()       │
│ - targetPosition = pile.pos    │
│ - API POST /place              │
└──────┬─────────────────────────┘
       ↓
┌────────────────────────────────┐
│ Backend:                       │
│ - Verify player & card exists  │
│ - Verify +1 rule              │
│ - Place card on pile          │
│ - Move to next player         │
│ - Update game state           │
└──────┬─────────────────────────┘
       ↓
┌────────────────────────────────┐
│ Frontend:                      │
│ - Clear drawnCard             │
│ - Clear waitingForPlacement   │
│ - Call onUpdate()             │
│ - Re-fetch game state         │
└──────┬─────────────────────────┘
       ↓
┌────────────────────────────────┐
│ Display:                       │
│ - Card is removed from center  │
│ - Piles are no longer lit      │
│ - Next player highlighted      │
│ - Ready for next turn         │
└────────────────────────────────┘
```

---

## 7️⃣ CARD SIZE PROGRESSION

### Card Sizes Used
```
Player Card (normal):
┌──────┐
│ A ♥  │     Width: 80px (w-20)
│      │     Height: 112px (h-28)
│ A ♥  │     Used in player piles
└──────┘

Drawn Card (large):
┌────────┐
│ A ♥    │     Width: 96px (w-24)
│        │     Height: 128px (h-32)
│ A ♥    │     Used in center, above deck
│        │     More prominent
└────────┘

Deck Card:
┌──────┐
│ 🂠    │     Width: 80px (w-20)
│ DECK │     Height: 112px (h-28)
│ 30   │     Shows card count
└──────┘
```

---

## 8️⃣ ANIMATION EFFECTS

### Deck Button Animation
```
DISABLED:                    ENABLED:
┌──────────┐               ┌──────────┐
│          │               │          │
│ 🂠 DECK  │ (gray)        │ 🂠 DECK  │ (blue)
│ 30       │ opacity-60    │ 30       │ bright
│          │               │ TAP ↑    │ bouncing
└──────────┘               └──────────┘

HOVER:
┌──────────┐
│          │ scale-105
│ 🂠 DECK  │ (5% larger)
│ 30       │ -translate-y-2
│          │ (8px higher)
└──────────┘
```

### Drawn Card Animation
```
APPEARING:                  CONTINUOUS:
Frame 1:                    Pulse Effect:
[Card] opacity: 0.3        ┌─────────────┐
                           │ opacity: 1  │
Frame 2:                   │ [Card]      │
[Card] opacity: 0.7        │ fade        │
                           │ [Card]      │
Frame 3:                   │ opacity: 0.8│
[Card] opacity: 1          └─────────────┘

Text Bounce:               (animate-pulse)
⬇️ Tap a player ⬇️
↓ (bounces up/down)
```

### Hover Effect Animation
```
NORMAL:                     HOVER:
┌─────────────┐           ┌─────────────┐
│             │           │             │
│ Player Card │           │ Player Card │ scale-105
│ Blue Border │ → hover → │ Blue Border │
│             │           │             │
└─────────────┘           └─────────────┘
                          (shadow added)
```

---

## 9️⃣ STATE TRANSITIONS

### Game States
```
WAITING (Setup)
    ↓
    [Players join]
    ↓
READY TO START
    ↓
    [Click START]
    ↓
IN_PROGRESS (Playing)
    ↓
    ├─ Player 1's turn
    │   ├─ Click Deck (draw)
    │   ├─ Click Pile (place)
    │   └─ Turn ends
    │       ↓
    ├─ Player 2's turn
    │   ├─ Click Deck (draw)
    │   ├─ Click Pile (place)
    │   └─ Turn ends
    │       ↓
    └─ ... cycles through players
    ↓
COMPLETED (Winner)
```

### UI State Machine
```
NO_CARD_DRAWN:
  ├─ Deck: ENABLED (blue, clickable)
  ├─ Piles: INACTIVE (gray)
  └─ Message: "Click deck to draw"

CARD_DRAWN:
  ├─ Deck: DISABLED (grayed out)
  ├─ Drawn Card: VISIBLE (animated)
  ├─ Piles: ACTIVE (highlight on hover)
  └─ Message: "Click pile to place card"

CARD_PLACED:
  ├─ Deck: ENABLED again
  ├─ Drawn Card: HIDDEN
  ├─ Piles: INACTIVE
  └─ Message: "Next player's turn"
```

---

## 🔟 SIZING REFERENCE

### Desktop Dimensions
```
Table:           600px × 600px
Player Pile:     200px wide
Deck Button:     80px × 112px
Card:            80px × 112px
Drawn Card:      96px × 128px (larger)
Gap (opponents): 32px
Top/Bottom:      32px padding
```

### Mobile Dimensions
```
Table:           300px × 300px (50% of desktop)
Player Pile:     140px wide
Deck Button:     60px × 84px (scaled)
Card:            60px × 84px
Drawn Card:      72px × 96px
Gap (opponents): 8px
Top/Bottom:      16px padding
```

### Responsive Breakpoints
```
sm:     640px (Tailwind small)
        ├─ gap-2 → gap-4
        ├─ top-4 → top-8
        ├─ bottom-4 → bottom-8
        └─ text-xs → text-sm

md:     768px (Tailwind medium)
        └─ Full-size layout

lg:     1024px (Tailwind large)
        └─ Full-size layout

xl:     1280px (Tailwind extra large)
        └─ Full-size layout
```

---

## SUMMARY

Your game table has:

✅ **Circular shape** - Perfect poker table  
✅ **Centered deck** - Clear focal point  
✅ **Revealed card** - Big and visible  
✅ **Natural positioning** - You at bottom  
✅ **Responsive design** - All screen sizes  
✅ **Smooth animations** - Professional look  
✅ **Color feedback** - Know what's happening  
✅ **Interactive piles** - Clear click targets  

**It's a beautiful, functional game table!** 🎮

