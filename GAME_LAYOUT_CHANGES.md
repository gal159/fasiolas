# Game Layout and Mechanics Update

## Changes Made

### 1. Frontend UI Redesign (GameBoard.jsx)
The game board has been completely redesigned to match real-life card game table layout:

#### Table Layout
- **Players positioned around the table**: Current player at bottom, other players distributed around (top, left, right)
- **Deck in the center**: A clickable deck button showing card count in the middle of the table
- **Drawn card display**: When a card is drawn, it's displayed above the deck

#### Player Areas
Each player area now shows:
- Player name (highlighted in blue if it's you)
- Their top card (the card they can play)
- Card count
- Turn indicator (green border and "🎯 Playing" if it's their turn)

#### Game Flow - Phase 1 (Kaupimas - Stacking)
1. **Player's turn begins**: 
   - They see deck in center
   - Can click deck to draw a card

2. **Card drawn**:
   - Card appears above the deck (visible to all players)
   - Player is prompted to select where to place it
   - Player clicks on any other player's area to select target
   - Clicks "Place Card" button

3. **After placement**:
   - Turn ends and passes to next player
   - Board updates automatically

### 2. New State Management
Added new state variables:
- `drawnCard`: Tracks the card that was drawn
- `waitingForPlacement`: Indicates if a card has been drawn and awaiting placement
- `selectedTarget`: Tracks which player the card will be placed on

### 3. Simplified Actions
- **Draw Card**: Click the deck button (only available when it's your turn and no card is drawn)
- **Place Card**: After drawing, click on a player, then click "Place Card" button
- **Turn ends**: Automatically when card is placed

### 4. Visual Indicators
- Green border + "🎯 YOUR TURN!" = It's your turn
- Green table border = Current player's turn
- Yellow border = Selected target player
- Drawn card pulsates above deck = Card is ready to be placed

## Game Rules Implemented

### Phase 1: Kaupimas (Stacking/Accumulation)

**Turn Structure:**
1. If you have cards, you **must** draw from the deck (click deck button)
2. You can see the card drawn (displayed above deck)
3. You must place it on another player (if +1 rule allows)
4. If you can't place on others, place on yourself (adds to your pile)
5. If you can't place anywhere, you keep it (pasiimti) and turn ends

**+1 Rule (Cyclic):**
- 2 can be placed on A
- 3 can be placed on 2
- ... and so on ...
- K can be placed on Q
- A can be placed on K

**Placement Logic:**
- Find a player whose top card is +1 from your card
- Place your card on them
- That becomes their new top card

## Testing Instructions

1. **Start the project**: Docker containers should be running
2. **Open browser**: http://localhost:3000
3. **Login**: Use Google OAuth or test account
4. **Create game**: Create a new game with 2+ players
5. **Join game**: Have another player join using room code
6. **Start game**: Click "START GAME" when 2+ players ready
7. **Play**:
   - When it's your turn, click the deck in the center
   - Drawn card appears above deck
   - Click on another player to select them
   - Click "Place Card" to place the drawn card
   - Turn passes to next player

## Files Modified
- `frontend/src/components/GameBoard.jsx`: Complete UI redesign

## Backward Compatibility
- All existing API endpoints remain unchanged
- Game state response format unchanged
- Backend logic unchanged (Phase 1 mechanics already implemented)

