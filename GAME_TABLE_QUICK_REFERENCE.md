# 🎮 GAME TABLE - QUICK REFERENCE CARD

**Updated:** January 24, 2026  
**Status:** ✅ Live and Working

---

## 📍 TABLE LAYOUT

```
        Opponent(s) at Top
        ┌────────────────┐
        │   [Card] ...   │
        │                │
        │  [Revealed]    │
        │   [Deck 🎴]    │
        │                │
        │ [Your Card]    │
        └────────────────┘
         You at Bottom
```

---

## 🎯 GAME FLOW

| Step | Action | What Happens |
|------|--------|--------------|
| 1 | Click deck | Card drawn & revealed |
| 2 | Click opponent pile | Card placed there |
| 3 | Wait | Next player's turn |
| Repeat | ... | Game continues |

---

## 🎨 COLOR MEANINGS

| Color | Meaning |
|-------|---------|
| 🟢 Green | Current player's turn |
| 🔵 Blue | Your pile |
| 🟡 Yellow | Selected target |
| ⚫ Black | Inactive pile |

---

## 🖱️ HOW TO PLAY

### Your Turn
```
1. See green "YOUR TURN!" indicator
2. Click deck in center
3. Card appears above deck
4. Click opponent pile OR your pile
5. Card is placed
6. Turn passes to next player
```

### Not Your Turn
```
1. See opponent's name with green glow
2. Watch them draw card
3. See card revealed
4. Watch them place it
5. Wait for your turn
```

---

## 🎮 INTERACTIVE ELEMENTS

### Deck Button (Center)
- **When ENABLED:** Blue, says "TAP", bounces
- **When DISABLED:** Grayed out
- **Click:** Draw a card

### Player Piles
- **When ACTIVE (can place):** Highlight yellow on hover
- **When SELECTED:** Gold border + ring
- **Click:** Place drawn card there

### Drawn Card
- **Appears:** When you draw from deck
- **Shows:** Big, animated, everyone sees it
- **Disappears:** When placed on a pile

---

## 💻 TECHNICAL INFO

| Item | Details |
|------|---------|
| Component | `GameBoard.jsx` |
| Layout | Circular (`rounded-full`) |
| Table Bg | Green felt poker table |
| Border | Yellow-gold (8px) |
| Responsive | ✅ Mobile/Tablet/Desktop |
| Backend | GO server on :8080 |
| Frontend | React on :3000 |

---

## ✅ FEATURES

- [x] Circular table
- [x] Clickable deck
- [x] Card revealed
- [x] Place on opponent
- [x] Place on self
- [x] Visual feedback
- [x] Animations
- [x] Responsive design
- [x] Turn indicators
- [x] Error handling

---

## 🚀 QUICK START

### Backend
```powershell
cd cardGame
go run cmd/server/main.go
```

### Frontend (Already Running)
```
http://localhost:3000
```

### Test
1. Login
2. Create game
3. Join with 2nd player
4. Start game
5. **Click deck & place cards!**

---

## 🎯 GAME RULES (Phase 1)

### Card Placement (+1 Rule)
- Can place on opponent if card is **+1** higher
  - Example: 3 on 2, Jack on 10, Ace on King
- Can place on yourself if **+1** to your previous card
- If can't place anywhere: **Must draw**

### Turn Flow
- **A Phase:** Try to place on others
- **B Phase:** If can't, must draw
- **C Phase:** Auto-place if possible
- **D Phase:** Turn ends

### Winning
- First player to empty their hand wins!

---

## 🎨 VISUAL GUIDE

### Colors Used
```
Green:   #1F7A3A (table), #10B981 (player turn), #34D399 (highlight)
Yellow:  #FBBF24 (deck), #FCD34D (drawn card), #FDE047 (instruction)
Blue:    #3B82F6 (your pile), #1E3A8A (deck bg), #1F2937 (player area)
Red:     #DC2626 (errors), #EF4444 (highlight)
Gray:    #6B7280 (inactive), #9CA3AF (text), #D1D5DB (borders)
```

### Size Ratios
```
Deck Card:      20px × 28px
Player Card:    20px × 28px
Drawn Card:     24px × 32px (large)
Table Min:      600px × 600px
Player Area:    200px width
Gap (desktop):  32px (gap-8)
Gap (mobile):   8px (gap-2)
```

---

## 🐛 TROUBLESHOOTING

| Problem | Solution |
|---------|----------|
| Deck not clickable | Not your turn or card already drawn |
| Can't place card | Doesn't follow +1 rule |
| Card not visible | Check if it's your turn |
| Wrong positions | Refresh page |
| No table visible | Ensure backend is running |
| Mobile layout broken | Refresh page |

---

## 📱 RESPONSIVE BEHAVIOR

| Device | Display |
|--------|---------|
| Desktop (1200px+) | Full-size 600px table |
| Tablet (768px) | Medium 400px table |
| Mobile (<768px) | Compact 300px table |
| All | Maintains circular shape |

---

## 🔗 RELATED DOCS

- `GAME_TABLE_COMPLETE.md` - Full documentation
- `GAME_TABLE_VISUAL_GUIDE.md` - Visual breakdown
- `GAME_RULES.md` - Complete game rules
- `GAME_START_EXPLANATION.md` - Game flow explanation

---

## 🎊 SUMMARY

**You have a complete, beautiful, functional circular game table!**

Click deck → Place card → Next player → Repeat! 🎮

---

**Last Updated:** January 24, 2026  
**Status:** ✅ Live  
**Frontend:** http://localhost:3000  
**Backend:** http://localhost:8080

