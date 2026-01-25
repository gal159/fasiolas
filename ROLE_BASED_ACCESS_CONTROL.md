# Role-Based Access Control (RBAC) Implementation

## Overview

The card game now supports three user roles with different capabilities:

### 1. **PLAYER** (Default)
- **Permissions:**
  - Create new games
  - Join existing games
  - Play cards (draw, place, cheat call, skip turn)
  - View own game state
  - See all games list
  - View ongoing games they're in

- **Restrictions:**
  - Cannot manage other users
  - Cannot access admin endpoints

### 2. **SPECTATOR** (View-Only)
- **Permissions:**
  - View all games and their state
  - Watch ongoing games in real-time
  - See player hands and table cards
  - View game history and actions

- **Restrictions:**
  - **Cannot** create games
  - **Cannot** join games as a player
  - **Cannot** perform any game actions (draw, place, cheat)
  - **Cannot** manage users

### 3. **ADMIN** (Management)
- **Permissions:**
  - All player permissions
  - List all users
  - View user details
  - Change user roles (player ↔ spectator ↔ admin)
  - View system statistics

- **Restrictions:**
  - Admin-only endpoints secured with role middleware

---

## API Endpoints

### Game Endpoints (All Authenticated Users)

#### View Games
```
GET /api/v1/games
- List all games (available to player, spectator, admin)
- Query params: state, limit, offset
```

#### Get Game State
```
GET /api/v1/games/:id
- Get full game state (players, cards, actions)
- Available to: players in game, spectators, admins
```

#### Spectator View
```
GET /api/v1/games/:id/spectate
- View-only game state
- No position requirement - spectators can watch any game
- Available to: spectators, admins
```

### Game Actions (Players & Admins Only)

These endpoints require `player` or `admin` role:

```
POST /api/v1/games
- Create new game

POST /api/v1/games/join
- Join existing game
- Body: { room_code: string }

POST /api/v1/games/:id/start
- Start game

POST /api/v1/games/:id/draw
- Draw card from deck

POST /api/v1/games/:id/place
- Place card on player
- Body: { target_player_position: int }

POST /api/v1/games/:id/skip
- Skip turn

POST /api/v1/games/:id/cheat
- Call out cheating
- Query: cheater_id
```

### Admin Endpoints (Admin Only)

All require `admin` role:

```
GET /api/v1/admin/users
- List all users
- Query params: limit, offset

GET /api/v1/admin/users/:id
- Get user details

PUT /api/v1/admin/users/:id/role
- Update user role
- Body: { role: "player" | "spectator" | "admin" }

GET /api/v1/admin/stats
- Get system statistics (user counts by role, etc.)
```

---

## Database

### Users Table
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  username VARCHAR(255) NOT NULL,
  oauth_provider VARCHAR(50) NOT NULL,
  oauth_id VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL DEFAULT 'player', -- 'player', 'spectator', 'admin'
  avatar_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP
);
```

### Role Values
- `player` - Default role for new users, can play games
- `spectator` - Can watch games without playing
- `admin` - Full system access

---

## JWT Token Claims

The JWT includes the user's role:

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "username": "john_doe",
  "role": "player",
  "exp": 1234567890,
  "iat": 1234567000
}
```

---

## Frontend Implementation

### Role Detection
```javascript
const currentUser = {
  id: 1,
  email: "user@example.com",
  username: "john_doe",
  role: "player" // or "spectator" or "admin"
};

const isSpectator = currentUser.role === 'spectator';
const isPlayer = currentUser.role === 'player' || currentUser.role === 'admin';
const isAdmin = currentUser.role === 'admin';
```

### Conditional UI
- **Spectators** see read-only game view
  - Card drag/drop disabled
  - Draw button disabled
  - Action panel shows spectator message
  - Cannot see own hand (not in game)

- **Players** see full interactive game
  - Can drag cards when their turn
  - Can draw and place cards
  - Action panel shows turn info

- **Admins** have player + admin features
  - Access to admin panel
  - Can manage users

---

## How to Use

### Change User Role (Admin Only)

```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/2/role \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'
```

### View User Stats (Admin Only)

```bash
curl http://localhost:8080/api/v1/admin/stats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Spectator Viewing Game

```bash
curl http://localhost:8080/api/v1/games/1/spectate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Security Notes

1. **Role-Based Middleware**: All protected endpoints check user role
2. **JWT Validation**: Every request validates JWT signature and role claim
3. **Business Logic Protection**: Backend prevents spectators from playing even if they send play requests
4. **Frontend Validation**: UI doesn't show play options to spectators

---

## Workflow Examples

### Example 1: New User Signs Up
1. OAuth login → Default role = `player`
2. User can create/join games immediately
3. Admin can later change role to `spectator` if needed

### Example 2: User Wants to Watch Only
1. User is `player`
2. Admin visits `/api/v1/admin/users/:id/role` endpoint
3. Updates role to `spectator`
4. User can now only view games
5. Frontend hides play buttons automatically

### Example 3: Spectator Becomes Player
1. User is `spectator`
2. Admin updates role to `player`
3. User can now create/join games
4. Frontend shows full game interface

### Example 4: Multi-Role User (Admin)
1. User created with `player` role
2. Admin promotes to `admin` role
3. User can:
   - Play all games (player permission)
   - Access admin panel
   - Manage other users

---

## Testing

### Test Spectator Cannot Play
```bash
# 1. Get token for spectator user
TOKEN=$(curl ... /api/v1/auth/google) # Get JWT

# 2. Try to draw (should fail)
curl -X POST http://localhost:8080/api/v1/games/1/draw \
  -H "Authorization: Bearer $TOKEN" \
  # Response: 403 Forbidden - insufficient permissions
```

### Test Spectator Can View
```bash
# Same spectator token
curl http://localhost:8080/api/v1/games/1 \
  -H "Authorization: Bearer $TOKEN" \
  # Response: 200 OK - game state returned
```

### Test Admin Can Change Role
```bash
# Admin token
ADMIN_TOKEN=$(curl ... /api/v1/auth/google) # Get JWT with admin role

# Change user 5 to spectator
curl -X PUT http://localhost:8080/api/v1/admin/users/5/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'
  # Response: 200 OK - user updated
```

