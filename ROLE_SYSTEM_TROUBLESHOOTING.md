# Role System - Troubleshooting & FAQs

## Troubleshooting

### Issue: Spectator can still play cards

**Problem**: Frontend allowed spectator to click buttons or drag cards

**Solutions**:
1. Check browser console for errors
2. Verify user role in token: `curl /api/v1/auth/profile`
3. Clear browser cache and localStorage
4. Verify GameBoard.jsx has role checks

```javascript
// Should see this check
const isSpectator = currentUser?.role === 'spectator';
if (isSpectator) {
  setError('Spectators cannot play...');
  return;
}
```

### Issue: Admin endpoints return 403 Forbidden

**Problem**: Admin user can't access admin endpoints

**Solutions**:
1. Verify user role is actually "admin": `curl /api/v1/auth/profile`
2. Check JWT token is in Authorization header: `Authorization: Bearer TOKEN`
3. Verify token hasn't expired
4. Check role in JWT matches what database says

```bash
# Debug: Get profile as admin
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Should return: { "role": "admin", ... }
```

### Issue: Player can't join or create games

**Problem**: "insufficient permissions" error

**Solutions**:
1. Verify role is "player" not "spectator": `curl /api/v1/auth/profile`
2. If spectator, admin must change role: `PUT /api/v1/admin/users/:id/role`
3. Check middleware is correctly set up in main.go

```bash
# Change spectator back to player (admin only)
curl -X PUT http://localhost:8080/api/v1/admin/users/5/role \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "player"}'
```

### Issue: Spectator endpoint returns 404

**Problem**: GET /api/v1/games/:id/spectate not found

**Solutions**:
1. Verify route is registered in main.go
2. Restart server after code changes
3. Check logs for route initialization

```go
// Should be in main.go protected routes
games.GET("/:id/spectate", gameHandler.GetGameStateSpectator)
```

### Issue: Token claims don't include role

**Problem**: Role not in JWT token

**Solutions**:
1. Verify GenerateJWT includes role: `pkg/utils/jwt.go`
2. Check role is fetched from database during login
3. Recreate token (logout and login again)

```bash
# Decode JWT to see claims
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN" | jq .
```

### Issue: Spectator can see other players' hands

**Problem**: Hidden information visible

**Solutions**:
1. This is by design - spectators can see all game info
2. If you want hidden hands, implement on backend:
   - Create GetGameStateForSpectator with hidden card data
   - Frontend to obfuscate sensitive data

---

## FAQs

### Q1: What's the default role for new users?

**A:** `player` - New users can create and join games immediately.

### Q2: How do I make someone a spectator?

**A:** Use admin endpoint:
```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/5/role \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'
```

### Q3: Can a spectator become a player again?

**A:** Yes, admin can change role back to player.

### Q4: What if no admins exist?

**A:** 
- First user created gets promoted (currently just player)
- Need to manually update database:
  ```sql
  UPDATE users SET role = 'admin' WHERE id = 1;
  ```

### Q5: Can spectators chat during games?

**A:** Not implemented yet - on Phase 2 roadmap.

### Q6: Can spectators see other players' cards?

**A:** Yes, spectators see full game state including all cards.

### Q7: Is there a limit on number of spectators per game?

**A:** No, unlimited spectators can view one game.

### Q8: Can a spectator switch to player mid-game?

**A:** They can switch role, but won't be in the current game. They'd need to join a new game as player.

### Q9: What happens if admin removes their own admin role?

**A:** They become a player - be careful!

### Q10: How are roles stored in database?

**A:** In `users.role` column as string: 'player', 'spectator', or 'admin'

### Q11: Can I have multiple roles?

**A:** No, each user has exactly one role at a time.

### Q12: Are roles stored in JWT token?

**A:** Yes, JWT includes role claim. Used for quick checks without database lookup.

---

## Common Scenarios

### Scenario 1: Tournament Setup

```
Goal: Have viewers watch tournament games

Steps:
1. Create admin account
2. Create player accounts for 4 competitors
3. Create spectator accounts for viewers
4. Players join games, spectators watch
5. Admin monitors everything
```

### Scenario 2: Convert Player to Spectator Temporarily

```
Goal: Player becomes spectator for one game, then plays again

Steps:
1. Admin changes user role to spectator
2. User views games, can't play
3. Admin changes back to player
4. User can join/create games again
```

### Scenario 3: Multi-Admin System

```
Goal: Have multiple admins managing users

Steps:
1. Create first admin manually in database
2. First admin creates other admin accounts
3. Admins can now manage users independently
4. Each admin has full permissions
```

---

## Performance Notes

### Database Queries
- Role check per request: Very fast (string comparison)
- User lookup for role changes: Indexed on user.id
- List users: Uses pagination (limit/offset)

### JWT Token
- Includes role claim: No performance impact
- Checked on every protected request: ~1ms per check
- No additional database calls after authentication

### Frontend
- Role checks: Zero-cost (string comparison)
- Disabling UI: No re-renders
- Spectator mode: Same performance as player mode

---

## Security Considerations

### 1. JWT Token Signing
- Token signed with secret key
- Cannot be forged
- Expiration enforced

### 2. Role Escalation Prevention
- Can't grant yourself admin role
- Only existing admins can create admins
- Middleware validates role on every request

### 3. Database Role Column
- Cannot be NULL
- Default value: 'player'
- Validated against known roles

### 4. API Endpoint Protection
- Double-checked with middleware
- Service layer can also validate
- Business logic enforces role rules

### 5. Frontend Validation
- UI doesn't show spectator play options
- API calls blocked for spectators
- Backend enforces regardless of frontend

### 6. Logging
- Role change logged (future)
- Admin actions auditable (future)
- Failed role checks can be logged

---

## Future Enhancements

### Phase 2 - Medium Priority
- [ ] Spectator chat/comments
- [ ] Admin UI panel
- [ ] Role history/audit log
- [ ] Invite system for players

### Phase 3 - Lower Priority
- [ ] Role-based game creation rules
- [ ] Custom roles
- [ ] Permission flags per role
- [ ] Spectator notifications
- [ ] Game recordings for spectators

---

## Testing Commands

### Quick Test Script

```bash
#!/bin/bash

# Get tokens
PLAYER_TOKEN=$(curl -s ... | jq -r .token)
SPECTATOR_TOKEN=$(curl -s ... | jq -r .token)
ADMIN_TOKEN=$(curl -s ... | jq -r .token)

# Test 1: All can view games
echo "Test 1: View games (all roles)"
curl http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Test 2: Only players can create
echo "Test 2: Create game (player only)"
curl -X POST http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer $SPECTATOR_TOKEN" \
  -d '{"max_players":2}'
# Should fail with 403

# Test 3: Only admins can change roles
echo "Test 3: Change role (admin only)"
curl -X PUT http://localhost:8080/api/v1/admin/users/2/role \
  -H "Authorization: Bearer $PLAYER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}'
# Should fail with 403

# Test 4: Admin can change roles
echo "Test 4: Admin changes role"
curl -X PUT http://localhost:8080/api/v1/admin/users/2/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"spectator"}'
# Should succeed with 200
```

---

## Debug Mode

### Enable Verbose Logging

In main.go, add before router setup:
```go
if cfg.Server.Env == "development" {
  log.SetFlags(log.LstdFlags | log.Lshortfile)
}
```

### Check JWT Claims

Add debug endpoint:
```go
router.GET("/debug/token", func(c *gin.Context) {
  claims, _ := authService.ValidateToken(token)
  c.JSON(200, claims)
})
```

### Monitor Role Checks

Add logging to RequireRole:
```go
func RequireRole(roles ...string) gin.HandlerFunc {
  return func(c *gin.Context) {
    userRole, _ := c.Get("role")
    log.Printf("[RBAC] Checking role %v against %v\n", userRole, roles)
    // ... rest of code
  }
}
```

---

## Rollback Plan

If something goes wrong:

### Quick Rollback
```bash
# 1. Stop server
docker-compose down

# 2. Revert code changes
git checkout HEAD -- .

# 3. Rebuild
go build ./cmd/server

# 4. Restart
docker-compose up
```

### Database Rollback
If roles are corrupted:
```sql
-- Reset all users to player
UPDATE users SET role = 'player';

-- Or reset specific user
UPDATE users SET role = 'player' WHERE id = 5;
```

### Make First Admin
```sql
-- Promote first user to admin
UPDATE users SET role = 'admin' WHERE id = 1;
```

