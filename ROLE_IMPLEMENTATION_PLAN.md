# Role-Based Access Control Implementation Plan

## Overview
Implement three user roles:
- **ADMIN**: Can manage users, view all games, monitor system
- **PLAYER**: Can create and join games, play (default for new users)
- **SPECTATOR**: Can view ongoing games without playing (cannot join as player)

## Implementation Steps

### 1. Backend Changes

#### 1.1 Models (✓ Already Complete)
- `RoleAdmin`, `RolePlayer`, `RoleSpectator` constants defined in models.go
- User model already has `Role` field

#### 1.2 Middleware (✓ Already Complete)
- `RequireRole()` middleware exists for role checks
- `GetUserRole()` helper exists

#### 1.3 Auth Service (✓ Already Complete)
- New users default to `RolePlayer`
- JWT includes role claim

#### 1.4 New Backend Features Needed
- [ ] New repository methods for spectators
- [ ] Game visibility logic (spectators can view, not join)
- [ ] New API endpoints for spectator features
- [ ] Admin endpoints for user role management

#### 1.5 Game Service Enhancements
- [ ] Allow spectators to view games without joining
- [ ] Block spectators from player actions (place, draw, etc.)
- [ ] Allow spectators to see game state

#### 1.6 Route Protection
- [ ] Add spectator-only routes
- [ ] Update game routes to support spectator viewing
- [ ] Add admin routes for user management

### 2. Frontend Changes

#### 2.1 Navigation & Auth
- [ ] Display user role in profile
- [ ] Show role-specific menu options
- [ ] Block UI actions based on role

#### 2.2 Game Board
- [ ] Disable game actions for spectators
- [ ] Show spectator-only view (cannot drag/drop cards)
- [ ] Display spectator badge

#### 2.3 Game List
- [ ] Show all games to spectators
- [ ] Allow spectators to "view" games
- [ ] Prevent spectators from joining as player

#### 2.4 Admin Panel (Optional Phase 2)
- [ ] User management interface
- [ ] Role assignment
- [ ] User statistics

### 3. Database
- [ ] Ensure role column exists in users table ✓
- [ ] Migration if needed

### 4. Testing
- [ ] Test player can create/join games
- [ ] Test spectator can view games
- [ ] Test spectator cannot play
- [ ] Test admin can manage users
- [ ] Test role-based API access

## Priority Order
1. **SPECTATOR** - Implement view-only game access
2. **PLAYER** - Ensure player actions blocked for spectators
3. **ADMIN** - Add admin endpoints for user management (Phase 2)

## Key Files to Modify

### Backend
- `internal/handler/game_handler.go` - Add spectator viewing
- `internal/service/game_service.go` - Add spectator logic
- `cmd/server/main.go` - Add new routes

### Frontend
- `frontend/src/components/GameBoard.jsx` - Disable actions for spectators
- `frontend/src/components/GameList.jsx` - Show spectator options
- `frontend/src/pages/GamePage.jsx` - Route to spectator view
- `frontend/src/context/AuthContext.jsx` - Store role in context

