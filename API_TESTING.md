# 🧪 API Testing Guide

This guide helps you test the Fasiolas Card Game REST API using common tools.

## Quick Test - Health Check

```bash
# Test server is running
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "fasiolas-card-game"
}
```

---

## 🔐 Authentication Flow

### Step 1: Get OAuth URL

```bash
# Get Google OAuth URL
curl http://localhost:8080/api/v1/auth/google

# Response:
{
  "url": "https://accounts.google.com/o/oauth2/v2/auth?..."
}
```

### Step 2: Complete OAuth in Browser

1. Copy the URL from above
2. Paste in browser
3. Complete Google login
4. Get redirected with `code` parameter

### Step 3: Handle Callback (Automatic)

The app will exchange the code for a JWT token and return:

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "User Name",
    "role": "player",
    "oauth_provider": "google"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Save the `token` for subsequent requests.

---

## 🎮 Game API Endpoints

All game endpoints require the JWT token in the Authorization header:

```bash
# Example header:
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Create a Game

```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "max_players": 4
  }'

# Response:
{
  "id": 1,
  "room_code": "ABC123",
  "state": "waiting",
  "phase": 1,
  "max_players": 4,
  "created_by": 1,
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Save the `room_code` for other players to join!**

### List Available Games

```bash
curl http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer YOUR_TOKEN"

# Optional filters:
# ?state=waiting    - Only waiting games
# ?limit=20         - Limit results
# ?offset=0         - Pagination offset
```

### Join a Game

```bash
curl -X POST http://localhost:8080/api/v1/games/join \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "room_code": "ABC123"
  }'
```

### Get Game State

```bash
curl http://localhost:8080/api/v1/games/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# Response includes:
# - Game info (state, phase, trump suit)
# - All players with their cards
# - Recent actions
```

### Start a Game

Only the game creator can start:

```bash
curl -X POST http://localhost:8080/api/v1/games/1/start \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Requirements:**
- At least 2 players joined
- Game state is "waiting"

---

## 🃏 Phase 1: Accumulation Moves

### Place a Card (+1 Rule)

```bash
curl -X POST http://localhost:8080/api/v1/games/1/place \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target_player_position": 2
  }'
```

**Rules:**
- Your top card must be exactly +1 rank from target's top card
- Example: If target has 7, you place 8
- Cannot place on self unless no other options available

### Draw a Card

```bash
curl -X POST http://localhost:8080/api/v1/games/1/draw \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**When to draw:**
- No valid placement options (no +1 available)
- Card is automatically added to your pile if it doesn't fit

### Call Out Cheating

```bash
curl -X POST http://localhost:8080/api/v1/games/1/cheat \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cheater_id": 2
  }'
```

**Penalty if correct:**
- All other players (with >1 card) give 1 card to accused cheater
- Cards go to bottom of cheater's pile

**Penalty if wrong:**
- You lose ability to call cheats for rest of game

---

## 🟥 Phase 2: Trick-Taking

After Phase 1 ends (deck is empty), game transitions to Phase 2:

```bash
# Same endpoints work:
# GET  /api/v1/games/1        - Check if phase changed
# POST /api/v1/games/1/place  - Play cards on table
```

**Phase 2 Rules:**
- Trump suit shown as game's `trump_suit` field
- Must play same suit as last card, or trump, or higher rank
- Last player with cards loses

---

## 📊 External API Integration

### Get Card Image

```bash
curl "http://localhost:8080/api/v1/external/card-image?suit=hearts&rank=A" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Response:
{
  "code": "AH",
  "image": "https://deckofcardsapi.com/api/deck/undefined/draw/?cards=AH",
  "images": {
    "svg": "https://deckofcardsapi.com/static/img/AH.svg",
    "png": "https://deckofcardsapi.com/static/img/AH.png"
  },
  "value": "ACE",
  "suit": "HEARTS"
}
```

### Get Game Statistics

```bash
curl http://localhost:8080/api/v1/external/stats \
  -H "Authorization: Bearer YOUR_TOKEN"

# Response:
{
  "total_games": 1000,
  "average_game_time": 15.5,
  "popular_card": "AS"
}
```

---

## Using Postman

Import the included `Fasiolas-API.postman_collection.json`:

1. Open Postman
2. Click "Import"
3. Select the JSON file
4. All endpoints configured with:
   - Correct headers
   - Example payloads
   - Auth token placeholder: `{{TOKEN}}`

**To use:**
1. Get JWT token from OAuth flow
2. Set Postman variable: `TOKEN` = your JWT
3. Run requests from collection

---

## Error Responses

### 401 Unauthorized
```json
{
  "error": "authorization header required"
}
```
**Fix:** Include `Authorization: Bearer YOUR_TOKEN` header

### 400 Bad Request
```json
{
  "error": "invalid game ID"
}
```
**Fix:** Check URL parameters are correct

### 404 Not Found
```json
{
  "error": "game not found"
}
```
**Fix:** Verify game ID exists

---

## Common Issues

### "Game already started"
- Can only join games in "waiting" state
- Create new game or wait for next round

### "Not your turn"
- Wait for game to reach your position
- Game processes actions sequentially

### "Must place card before drawing"
- You have a valid placement (+1)
- Use `/place` endpoint before `/draw`

### "Cannot beat spades with non-spade trump"
- Spades can only be beaten by higher spades
- Trump is for non-spade suits
- Exception: Aces always high

---

## 🧹 Cleanup / Testing Reset

```bash
# Delete game (if you have permissions)
curl -X DELETE http://localhost:8080/api/v1/games/1 \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Get user profile
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"

# Refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📝 Notes

- All timestamps are in UTC (ISO 8601 format)
- Card suits: `hearts`, `diamonds`, `clubs`, `spades`
- Card ranks: `A`, `2`-`10`, `J`, `Q`, `K`
- Responses include pagination headers in list endpoints
- Database transactions ensure game state consistency

