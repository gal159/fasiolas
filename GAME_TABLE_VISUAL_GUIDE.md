# 🎮 Game Table - Visual Implementation Guide

---

## 📐 Circular Table Design

### HTML Structure
```html
<div className="relative mx-auto bg-gradient-to-br from-green-800 via-green-700 to-green-900 rounded-full shadow-2xl border-8 border-yellow-900 flex items-center justify-center"
     style={{ minHeight: '600px', minWidth: '600px', maxWidth: '100%', aspectRatio: '1' }}>
  
  <!-- Texture Overlay -->
  <div className="absolute inset-0 rounded-full opacity-20"></div>
  
  <!-- CENTER: Deck & Drawn Card -->
  <div className="relative z-20 flex flex-col items-center justify-center gap-4">
    <DeckCard />           <!-- Deck button -->
    {drawnCard && <DrawCard />}  <!-- Revealed card -->
  </div>
  
  <!-- POSITIONS: Players around circle -->
  <div className="absolute inset-0 pointer-events-none">
    
    <!-- TOP: Opponents -->
    <div className="absolute top-8 left-1/2 transform -translate-x-1/2 flex gap-8 pointer-events-auto">
      {opponents.map(player => <PlayerArea player={player} />)}
    </div>
    
    <!-- BOTTOM: You -->
    <div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 pointer-events-auto">
      <PlayerArea player={currentUser} />
    </div>
    
  </div>
  
</div>
```

### Circular Table CSS
```css
/* Make it circular */
rounded-full              /* Perfect circle border */
aspect-ratio: 1          /* Square dimensions (height = width) */
minHeight: 600px         /* Minimum size for playability */
minWidth: 600px          /* Keep it square */
maxWidth: 100%           /* Don't overflow on mobile */

/* Center positioning */
flex items-center justify-center   /* Flex children to center */
relative                           /* Position children absolutely inside */

/* Styling */
bg-gradient-to-br from-green-800 via-green-700 to-green-900
border-8 border-yellow-900
shadow-2xl
```

---

## 🎯 Positioning System

### 1. Center (Deck & Card)
```
absolute positioning = NO (uses flex centering)
z-index = 20 (above background)
layout = flex column
gap = 4 (16px spacing)

Components:
├── DeckCard button (clickable)
└── PlayingCard (if drawnCard exists)
```

### 2. Top (Opponents)
```
absolute top-8              /* 32px from top */
left-1/2 transform -translate-x-1/2  /* Centered horizontally */
flex gap-8                  /* Spread out horizontally */
pointer-events-auto         /* Make clickable */

Mapped from:
relativeIndex !== 0         /* All except current player */
```

### 3. Bottom (You)
```
absolute bottom-8            /* 32px from bottom */
left-1/2 transform -translate-x-1/2  /* Centered horizontally */
pointer-events-auto          /* Make clickable */

Only shows:
currentUser's PlayerArea     /* Just you */
```

---

## 🎨 Component Breakdown

### DeckCard Component
```jsx
<button
  className={`relative w-20 h-28 rounded-lg border-2 border-yellow-500 
              shadow-2xl flex flex-col items-center justify-center transition-all ${
    disabled ? 'opacity-60 cursor-not-allowed bg-blue-700' 
           : 'bg-blue-800 hover:scale-105 hover:-translate-y-2 cursor-pointer'
  }`}
  onClick={handleDrawCard}
  disabled={!isYourTurn || drawnCard !== null}
>
  <div className="text-2xl">🂠</div>           {/* Card back symbol */}
  <div className="text-xs text-white mt-1">Deck</div>
  <div className="text-xs text-yellow-300 font-bold">{count}</div>
  {!disabled && <div className="absolute -top-3 -right-2 bg-yellow-400 text-black text-xs px-2 py-0.5 rounded-full font-bold animate-bounce">TAP</div>}
</button>
```

**Features:**
- Blue background when enabled
- Yellow border (gold border)
- Card count displayed
- "TAP" indicator bounces
- Scale up on hover
- Disabled state when not your turn or card drawn

### PlayingCard Component
```jsx
<div className={`w-20 h-28 bg-white rounded-lg shadow-lg 
                flex flex-col items-center justify-center border-2 
                ${isTopCard ? 'border-yellow-400 ring-2 ring-yellow-400' : 'border-gray-300'}`}
>
  <div className={`font-bold ${getSuitColor(card.suit)}`}>
    <div className="text-xs absolute top-1 left-1">{card.rank}</div>
    <div className="text-3xl">{getSuitSymbol(card.suit)}</div>
    <div className="text-xs absolute bottom-1 right-1 rotate-180">{card.rank}</div>
  </div>
</div>
```

**Features:**
- White card background
- Suit symbols (♥♦♣♠)
- Rank in all 4 corners
- Red for hearts/diamonds
- Black for clubs/spades
- Yellow border when top card
- Ring effect for emphasis

### PlayerArea Component
```jsx
<div className="flex flex-col items-center p-4 transition-all"
     onDragOver={handleDragOver}
     onDrop={(e) => { /* handle drop */ }}>
  
  <div className={`bg-gray-800 rounded-lg p-4 border-4 text-center 
                  min-w-[200px] transition-all cursor-pointer hover:scale-105 ${
    isCurrentTurn ? 'border-green-400 shadow-lg shadow-green-400/50' :
    isYou ? 'border-blue-400 shadow-lg shadow-blue-400/50' :
    drawnCard && selectedTarget === player.position ? 'border-yellow-400 shadow-lg shadow-yellow-400/50 ring-4 ring-yellow-300' :
    drawnCard && waitingForPlacement ? 'border-gray-600 hover:border-yellow-300 hover:shadow-lg' :
    'border-gray-600'
  }`}
    onClick={() => { if (drawnCard && waitingForPlacement) handlePlaceCard(player.position); }}>
    
    {/* Player Name */}
    <p className={`font-bold text-sm ${isYou ? 'text-blue-400' : 'text-white'}`}>
      {player.user?.username || 'Unknown'}
      {isYou && ' (You)'}
    </p>

    {isCurrentTurn && <p className="text-xs text-green-400 font-bold">🎯 Playing</p>}

    {/* Top Card */}
    <div className="mt-2">
      {shouldHideCard ? (
        <div className="w-20 h-28 bg-gray-700/50 rounded-lg flex items-center justify-center text-gray-400 text-xs border-2 border-dashed border-gray-500">
          (placing)
        </div>
      ) : player.top_card ? (
        <PlayingCard card={player.top_card} isTopCard={true} size="normal" />
      ) : (
        <div className="w-20 h-28 bg-gray-700 rounded-lg flex items-center justify-center text-gray-500 text-xs border-2 border-gray-600">
          No Cards
        </div>
      )}
    </div>

    {/* Card Count */}
    <p className="text-xs text-gray-400 mt-2">📚 {player.card_count} cards</p>
    
  </div>
  
</div>
```

**Features:**
- Dark background (gray-800)
- Dynamic border color based on state
- Hover effect (scale up)
- Drag-and-drop support
- Shows current player indicator
- Shows "You" label for your pile
- Shows top card of player
- Shows card count
- Responsive to drawn card state

---

## 🎨 Color Scheme

### Table Colors
```css
/* Table Background */
bg-gradient-to-br from-green-800 via-green-700 to-green-900
/* Looks like: Green poker table felt */

/* Table Border */
border-8 border-yellow-900
/* Looks like: Gold/wooden table edge */
```

### Player Area Colors
```css
/* Default Border */
border-gray-600
/* Shows: Inactive player */

/* Your Pile */
border-blue-400
shadow-lg shadow-blue-400/50
/* Shows: Blue glow, this is you */

/* Current Player's Turn */
border-green-400
shadow-lg shadow-green-400/50
/* Shows: Green glow, someone playing */

/* Selected Target */
border-yellow-400
shadow-lg shadow-yellow-400/50
ring-4 ring-yellow-300
/* Shows: Yellow glow with ring, about to place here */

/* Hovering over clickable pile */
border-gray-600 hover:border-yellow-300 hover:shadow-lg
/* Shows: Yellow highlight on hover */
```

### Card Colors
```css
/* Drawn Card Border */
border-4 border-yellow-400
/* Shows: Important, revealed card */

/* Drawn Card Background */
bg-yellow-900/40
/* Shows: Highlighted area */

/* Suit Colors */
Hearts/Diamonds: text-red-500
Clubs/Spades: text-gray-900
```

---

## 🎬 Animations

### Deck Button
```css
/* When disabled */
opacity-60 cursor-not-allowed

/* When enabled and hovered */
hover:scale-105        /* Grows 5% */
hover:-translate-y-2   /* Moves up 8px */
cursor-pointer

/* TAP Indicator */
animate-bounce         /* Bounces continuously */
```

### Drawn Card
```css
/* Container */
animate-pulse          /* Fades in/out */
cursor-move            /* Indicates draggable */

/* Text */
animate-bounce         /* Bounces continuously */

/* Overall */
shadow-2xl            /* Large shadow */
transition-all        /* Smooth transitions */
```

### Player Areas
```css
/* Hover */
hover:scale-105       /* Grows on hover */
transition-all        /* Smooth change */

/* Glow effects */
shadow-lg shadow-green-400/50    /* Green glow */
shadow-lg shadow-blue-400/50     /* Blue glow */
shadow-lg shadow-yellow-400/50   /* Yellow glow */
```

---

## 📱 Responsive Behavior

### Desktop (1200px+)
```
minHeight: 600px
minWidth: 600px
gap-8 (opponents spacing)
top-8 bottom-8 (padding)
text-sm (normal size)
```
Result: Full-size circular table

### Tablet (768px - 1200px)
```
Scales down proportionally
maxWidth: 100% keeps it bounded
Still square and circular
Still fully functional
```
Result: Medium circular table

### Mobile (< 768px)
```
gap-2 (tight spacing)
top-4 bottom-4 (less padding)
text-xs (smaller text)
Still responsive
Still circular
```
Result: Compact circular table

---

## 🔄 Interaction Flow

### Step 1: Deck Click
```
User clicks DeckCard button
  ↓
onClick={handleDrawCard}
  ↓
setWaitingForPlacement(true)
setDrawnCard(response.card)
  ↓
PlayingCard component appears above deck
  ↓
Instructions show: "⬇️ Tap a player to place card ⬇️"
```

### Step 2: Player Click
```
User clicks PlayerArea component
  ↓
onClick={() => handlePlaceCard(player.position)}
  ↓
POST /api/v1/games/{id}/place
  ↓
Backend validates placement
  ↓
Card is placed
  ↓
onUpdate() fetches new game state
  ↓
Component re-renders
  ↓
drawnCard cleared
selectedTarget cleared
```

### Step 3: Turn Change
```
Backend updates current_player_position
  ↓
fetchGameState() called
  ↓
game.current_player_position changed
  ↓
isYourTurn recalculated
  ↓
UI updates
  ↓
New player's turn begins
```

---

## ✨ Key Features Explained

### 1. Circular Table
- Shape: Perfect circle (aspectRatio: 1)
- Size: Minimum 600px, responsive to screen
- Border: Thick yellow-900 (looks like wood/gold edge)
- Background: Green gradient (poker table felt)
- Texture: Subtle dot pattern overlay (cloth texture)

### 2. Center Deck
- Position: Absolute center (flex centering)
- Appearance: Blue card deck back symbol
- Count: Shows number of remaining cards
- Interactive: Clickable when enabled
- Feedback: "TAP" indicator bounces when available

### 3. Revealed Card
- Position: Above deck, in center
- Size: Large (bigger than player cards)
- Appearance: White card with suit and rank
- Animation: Pulse effect (fades in/out)
- Color: Suit color-coded (red/black)
- Border: Thick yellow (gold border)

### 4. Player Piles
- Position: You at bottom, opponents at top
- Appearance: Dark background with border
- Content: Player name + top card + card count
- Color: Blue for you, varies for others
- Interactive: Click to place drawn card
- Feedback: Color changes on selection

### 5. Status Indicators
- Current player: Green glow border
- Your pile: Blue glow border
- Selected target: Yellow glow + ring
- Available target: Hover highlight

---

## 🎮 Complete User Journey

### 1. Game Starts
```
See: Circular table with players positioned
See: Deck in center, ready to click
See: Your pile at bottom (with blue border)
See: Opponent piles at top
```

### 2. Your Turn
```
See: "YOUR TURN!" indicator
See: Green glow around your name
Click: Deck in center
See: Card appears above deck (animated)
See: All player piles light up yellow
See: Instructions: "Tap a player to place card"
```

### 3. Place Card
```
Choose: Click opponent or your own pile
See: Your selected target glows gold
See: Card animates to selected pile
See: Pile updates with new card
```

### 4. Turn Passes
```
See: Next player gets green glow
See: Your pile goes back to blue
See: Watch opponent draw card
See: Cycle repeats
```

---

## 📊 State Management

### Component State
```javascript
const [drawnCard, setDrawnCard] = useState(null);
const [waitingForPlacement, setWaitingForPlacement] = useState(false);
const [selectedTarget, setSelectedTarget] = useState(null);
const isPlacingCardRef = useRef(false);
const drawnCardRef = useRef(null);
```

### Conditional Rendering
```javascript
// Show deck button only when not already drawn
disabled={!isYourTurn || drawnCard !== null}

// Show drawn card only when exists
{drawnCard && <DrawCard />}

// Show player interactive only when card drawn
isClickable = isYourTurn && !isYou && waitingForPlacement && drawnCard

// Hide your top card when placing
shouldHideCard = isYou && (drawnCard || isPlacingCardRef.current)
```

---

## 🎊 Summary

The circular game table is implemented with:

✅ **Perfect circular shape** - Green felt poker table  
✅ **Centered deck** - Blue button with card count  
✅ **Revealed card** - Large, animated, yellow-bordered  
✅ **Natural positioning** - You at bottom, opponents at top  
✅ **Interactive placement** - Click any pile to place card  
✅ **Visual feedback** - Color changes, glows, highlights  
✅ **Responsive design** - Works on all screen sizes  
✅ **Smooth animations** - Pulse, bounce, scale effects  
✅ **Clear instructions** - On-screen guidance  
✅ **Turn indicators** - Know when it's your turn  

**It's a complete, beautiful, functional game table!** 🎮

