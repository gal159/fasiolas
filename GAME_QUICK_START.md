# Quick Start Guide - New Game Layout

## 🎮 Game is Ready to Play!

### System Status
- ✅ Backend Server: http://localhost:8080
- ✅ Frontend: http://localhost:3000
- ✅ Database: Connected and initialized
- ✅ Docker Containers: All running

### How to Test

#### Step 1: Open the Game
1. Open browser and go to: **http://localhost:3000**
2. You should see login page

#### Step 2: Login
- Click "Google Login" button
- Complete Google OAuth flow
- You'll be redirected to Dashboard

#### Step 3: Create a Game
1. Click "Create Game" button
2. Select max players (e.g., 2-4)
3. Click "Create"
4. You'll be redirected to game lobby (waiting room)

#### Step 4: Join as Second Player (IMPORTANT!)
You need 2 players to test. Options:
1. **Use incognito window**:
   - Open incognito/private window
   - Go to http://localhost:3000
   - Login with a different Google account
   - Join game using room code

2. **Or use the Dashboard**:
   - In the game lobby, you see room code (e.g., ABC123)
   - On another browser/account, go to Dashboard
   - Paste code in "Join by Room Code" section
   - Click "Join Game"

#### Step 5: Start Game
When 2+ players are in lobby:
- Click big green "START GAME" button
- Game state changes to "phase1"

#### Step 6: Play the Game
The new table layout shows:
- **Your position**: Bottom of the table
- **Other players**: Around the table
- **Deck**: In the center (with card count)
- **Each player area**: Shows their top card + card count

**During your turn:**
1. Click the **deck button** in the center
   - Card appears above the deck
   - You see the drawn card
2. Click on a **player area** to select where to place
   - Their area gets yellow border
3. Click **"Place Card"** button
   - Card is placed
   - Turn passes to next player

**When it's NOT your turn:**
- You see "Your Hand" info panel
- Watch other players play
- See current player highlighted with green

### What You'll See

```
           ┌─────────────────────┐
           │   PLAYER (TOP)      │
           │  [Card shown]       │
           │  🂠 5 cards         │
           └─────────────────────┘

┌──────────────┐          ┌──────────────┐
│ PLAYER LEFT  │          │ PLAYER RIGHT │
│  [Card]      │          │  [Card]      │
│  🂠 4 cards  │          │  🂠 3 cards  │
└──────────────┘          └──────────────┘

           ┌─────────────────────┐
           │   📚 Deck (32)      │  <- Click to draw
           │                     │
           │   [Drawn card ▲]    │  <- Shows drawn card
           └─────────────────────┘

┌──────────────┐          ┌──────────────┐
│ PLAYER (YOU) │          │              │
│  [Your Card] │          │ ACTION PANEL │
│  🂠 6 cards  │          │ TAP DECK     │
│ (Blue border)│          │ PLACE CARD   │
└──────────────┘          └──────────────┘
```

### Common Actions

**Action**: Draw a card
- Click the **deck button** in center
- Wait for card to appear
- Card shows above deck

**Action**: Place a drawn card
1. Click the **player area** you want to place on (gets yellow border)
2. Click **"Place Card"** button
3. Card moves to that player
4. Turn ends

**Action**: See turn progress
- **Green border** = Currently playing
- **Blue border** = Your player (you)
- **Yellow border** = Selected target
- **Green header** = Your turn indicator

### Troubleshooting

**Can't see the deck?**
- Make sure you're in game (not lobby)
- Check browser console (F12) for errors
- Refresh page if needed

**Drawn card not appearing?**
- Wait a moment after clicking deck
- Card should appear above deck
- If nothing appears, check network tab in dev tools

**Can't place card?**
- First, click a player area (they get yellow border)
- Then click "Place Card" button
- If button is disabled, you may need to select target first

**Turn not passing?**
- Check if "Place Card" button worked
- Refresh the page to see updated state
- Game refreshes every 2 seconds automatically

### Game Rules (Phase 1 Simplified)

1. **When it's your turn:**
   - Draw card from deck
   - Place it on another player (if possible)
   - If not possible, place on yourself
   - If nowhere to place, keep it (pasiimti)

2. **+1 Rule:**
   - You can only place a card if it's +1 from target's top card
   - Cyclic: A→2, K→A, Q→K, etc.
   - Examples:
     - 5 can go on 4
     - 6 can go on 5
     - A can go on K
     - 2 can go on A

3. **Card Placement:**
   - Click deck → draw
   - See card above deck
   - Click player → select target
   - Click "Place Card" → place

4. **Turn ends:**
   - Automatically after placing card
   - Passes to next player

### Need Help?

Check these logs:
```bash
# Backend logs
docker logs fasiolas_app

# Frontend logs  
docker logs fasiolas_frontend

# Database logs
docker logs fasiolas_postgres
```

Restart if needed:
```bash
docker compose restart
```

Rebuild if code changed:
```bash
docker compose down -v
docker compose up -d
```

---

**Enjoy the game! 🎮🃏**

The new table layout matches real-life card game setup where everyone sits around the table!

