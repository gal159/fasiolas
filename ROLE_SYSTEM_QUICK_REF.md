# Role System - Quick Reference

## Three Roles

| Role | Create Game | Join Game | Play | View | Manage Users |
|------|-------------|-----------|------|------|--------------|
| **PLAYER** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **SPECTATOR** | ❌ | ❌ | ❌ | ✅ | ❌ |
| **ADMIN** | ✅ | ✅ | ✅ | ✅ | ✅ |

## Key Files Modified/Created

### Backend
- ✅ `internal/models/models.go` - Role constants (already existed)
- ✅ `internal/middleware/middleware.go` - Role middleware (already existed)
- ✅ `internal/handler/game_handler.go` - Added spectator endpoint
- ✅ `internal/handler/admin_handler.go` - NEW - Admin endpoints
- ✅ `internal/service/game_service.go` - Added spectator service methods
- ✅ `internal/service/admin_service.go` - NEW - Admin service
- ✅ `cmd/server/main.go` - Added admin routes

### Frontend
- ✅ `frontend/src/components/GameBoard.jsx` - Added role checks, disabled play for spectators

## Role Constants (Go)
```go
const (
    RoleAdmin     UserRole = "admin"
    RolePlayer    UserRole = "player"
    RoleSpectator UserRole = "spectator"
)
```

## How Spectators Work

### Backend
1. `GetGameState()` updated to accept `userRole` parameter
2. New `GetGameStateForSpectator()` allows viewing without being a player
3. New route `/api/v1/games/:id/spectate` for spectator-only access
4. Game actions (place, draw, etc.) blocked by middleware `RequireRole("player", "admin")`

### Frontend
```javascript
const isSpectator = currentUser?.role === 'spectator';

// All these functions check isSpectator and return early:
- handlePlaceCard() - blocked
- handleDrawCard() - blocked
- All UI elements disabled (drag, drop, buttons)
```

## How Admins Work

### New Admin Endpoints
```
GET    /api/v1/admin/users              - List users
GET    /api/v1/admin/users/:id          - Get user
PUT    /api/v1/admin/users/:id/role     - Change role
GET    /api/v1/admin/stats              - System stats
```

### Middleware
All admin routes protected by: `middleware.RequireRole("admin")`

## Testing

### Change User Role
```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/5/role \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'
```

### Test Spectator Block
```javascript
// Frontend - try to place card as spectator
if (isSpectator) {
  setError('Spectators cannot play. Switch to player role to participate.');
  return;
}
```

## Database
- `users.role` column stores: `player`, `spectator`, or `admin`
- New users default to `player` role

## JWT Token
- Includes `role` claim
- Used by middleware to check permissions
- Validated on every protected request

## Next Steps (Phase 2 - Optional)
- [ ] Add frontend admin panel UI
- [ ] Add user search/filter in admin
- [ ] Add role-based game history
- [ ] Add spectator chat (viewers can chat)
- [ ] Add game invites (players invite others)
