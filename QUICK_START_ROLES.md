# Role System - Quick Start Guide

## 🚀 Get Started in 5 Minutes

### Step 1: Understand the Roles
```
PLAYER (default)  - Creates and plays games
SPECTATOR         - Watches games, cannot play
ADMIN             - Manages users + plays games
```

### Step 2: No Code Changes Needed!
The system is already implemented. Just rebuild and restart:

```bash
cd cardGame
go build ./cmd/server
./bin/server  # or server.exe on Windows
```

### Step 3: Test the Roles

**Test 1: View Profile to See Your Role**
```bash
# Get your JWT token from login, then:
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"

# Response includes: { "role": "player", ... }
```

**Test 2: Change Someone's Role (Admin Only)**
```bash
# As admin, change user 5 to spectator:
curl -X PUT http://localhost:8080/api/v1/admin/users/5/role \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'
```

**Test 3: Spectator Tries to Play (Should Fail)**
```bash
# As spectator, try to draw a card:
curl -X POST http://localhost:8080/api/v1/games/1/draw \
  -H "Authorization: Bearer SPECTATOR_TOKEN"

# Response: 403 Forbidden - insufficient permissions
```

**Test 4: Spectator Can View (Should Work)**
```bash
# As spectator, view a game:
curl http://localhost:8080/api/v1/games/1 \
  -H "Authorization: Bearer SPECTATOR_TOKEN"

# Response: 200 OK - full game state
```

---

## 📱 Frontend Testing

### For Players
1. Login as player
2. See "Create Game" button
3. Click it
4. Start playing
5. ✅ All game actions available

### For Spectators  
1. Login as spectator (or admin changes you to spectator)
2. See game list
3. Click on a game
4. See "👁️ Spectator Mode" message
5. See action buttons DISABLED
6. ✅ Can view, cannot play

### For Admins
1. Login as admin
2. Do everything players can do
3. Access admin features (when implemented)
4. ✅ Full access

---

## 🔧 Admin Commands Reference

### View All Users
```bash
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### View Specific User
```bash
curl http://localhost:8080/api/v1/admin/users/5 \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### Change User Role
```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/5/role \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "spectator"}'
```

Valid roles: `"player"`, `"spectator"`, `"admin"`

### View System Stats
```bash
curl http://localhost:8080/api/v1/admin/stats \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## 🆘 Common Issues & Quick Fixes

### "403 Forbidden - insufficient permissions"
**Problem**: User trying to play but not authorized
**Fix**: Check user role is "player" or "admin"
```bash
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer TOKEN"
# Should show: "role": "player"
```

### Spectator still seeing play buttons
**Problem**: Frontend cache
**Fix**: Clear browser cache
- Ctrl+Shift+Delete (or Cmd+Shift+Delete on Mac)
- Select "All time"
- Check "Cookies and other site data"
- Click Clear data

### Admin endpoints not working
**Problem**: User role is not "admin"
**Fix**: Verify in database
```sql
SELECT id, username, role FROM users WHERE id = 1;
-- Should show: role = 'admin'
```

### Need to make first admin
```sql
-- If no admins exist, promote first user:
UPDATE users SET role = 'admin' WHERE id = 1;
```

---

## 🎯 Typical Workflows

### Workflow 1: Tournament Setup
```
1. Admin creates accounts for 4 players
2. Admin creates accounts for 10 spectators
3. Players join games
4. Spectators watch games
5. Admin monitors everything
```

### Workflow 2: Convert Player to Spectator
```
1. Admin changes user role to spectator
2. User sees read-only game view
3. User cannot play
4. Admin can change back to player later
```

### Workflow 3: Add New Admin
```
1. Current admin creates new user (default = player)
2. Current admin changes role to admin
3. New admin can now manage other users
```

---

## 📚 Full Documentation

For detailed information, see these files:

| File | Purpose |
|------|---------|
| ROLE_SYSTEM_INDEX.md | Navigation hub (START HERE) |
| ROLE_SYSTEM_QUICK_REF.md | Quick reference table |
| ROLE_BASED_ACCESS_CONTROL.md | Complete API reference |
| ROLE_SYSTEM_TROUBLESHOOTING.md | FAQs & troubleshooting |
| ROLE_SYSTEM_DIAGRAMS.md | Visual diagrams |
| COMPLETION_SUMMARY.md | Summary of work done |

---

## ✅ Verification Checklist

Before deploying, verify:

- [ ] Go code compiles: `go build ./cmd/server`
- [ ] Server starts: `./bin/server`
- [ ] Can login: Test OAuth
- [ ] Can view profile: GET /api/v1/auth/profile
- [ ] Can create game: POST /api/v1/games (as player)
- [ ] Cannot create game: POST /api/v1/games (as spectator) → 403
- [ ] Can change role: PUT /api/v1/admin/users/:id/role (as admin)
- [ ] Spectator UI disabled: Login as spectator, view game

---

## 🎓 What's New vs What Changed

### What's New (Added)
```
✅ SPECTATOR role - view-only access
✅ Admin endpoints - user management
✅ GetGameStateSpectator() - spectator viewing
✅ Admin routes - protected by middleware
```

### What Changed (Enhanced)
```
✅ GameBoard.jsx - now checks for spectators, disables UI
✅ GetGameState() - accepts role parameter
✅ Route protection - updated to allow spectators to view
```

### What's Unchanged (Preserved)
```
✅ Player functionality - 100% same
✅ Game rules - 100% same
✅ Database schema - 100% compatible
✅ OAuth - 100% same
✅ All existing APIs - backward compatible
```

---

## 💡 Key Points

1. **No Migrations Needed** - Role column already exists
2. **No Database Changes** - Everything works as-is
3. **Backward Compatible** - Existing code unaffected
4. **Default Role** - New users are "player"
5. **Easy to Use** - Just change role via API

---

## 🚀 Deploy Checklist

- [ ] Code changes reviewed
- [ ] Go binary rebuilt
- [ ] Server restarted
- [ ] Token format verified (includes role)
- [ ] Routes registered (check startup logs)
- [ ] Test as player (can play)
- [ ] Test as spectator (cannot play)
- [ ] Test as admin (can manage users)
- [ ] Documentation available to team
- [ ] Users notified of new features

---

## 📞 Getting Help

**Quick Question?** → Check ROLE_SYSTEM_QUICK_REF.md  
**Technical Details?** → Check ROLE_BASED_ACCESS_CONTROL.md  
**Something Broken?** → Check ROLE_SYSTEM_TROUBLESHOOTING.md  
**Visual Learner?** → Check ROLE_SYSTEM_DIAGRAMS.md  
**Understand Changes?** → Check COMPLETION_SUMMARY.md  

---

## 🎉 You're Ready!

The role system is implemented, tested, and documented. 

**Next steps:**
1. Rebuild the server
2. Restart it
3. Test the roles
4. Deploy with confidence!

All documentation is included in the project. Share ROLE_SYSTEM_INDEX.md with your team to get them started.

