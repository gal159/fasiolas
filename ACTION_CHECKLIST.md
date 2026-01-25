# ✅ Database Error Fix - Action Checklist

## Your Current Issue
```
{"error":"failed to get user: pq: relation \"users\" does not exist"}
```
**This appears when you try to login with Google.**

---

## 🎯 IMMEDIATE ACTION (Choose One)

### Method 1: Auto-Run (2 minutes) ⚡
```powershell
# 1. Open PowerShell in your cardGame folder
# 2. Run this ONE command:
powershell -ExecutionPolicy Bypass -File .\restart-complete.ps1

# 3. Wait for it to finish
# 4. Go to http://localhost:3000 and login
```

**✅ DONE!**

---

### Method 2: Manual (5 minutes) 
```powershell
# Step 1: Stop Docker (wait 3 seconds after)
docker compose down -v --remove-orphans

# Step 2: Start Docker (wait 10 seconds after)
docker compose up -d

# Step 3: Copy migration files
docker compose cp migrations/000001_create_users_table.up.sql postgres:/tmp/
docker compose cp migrations/000002_create_games_table.up.sql postgres:/tmp/
docker compose cp migrations/000003_create_game_players_table.up.sql postgres:/tmp/
docker compose cp migrations/000004_create_game_actions_table.up.sql postgres:/tmp/

# Step 4: Apply migrations (run all 4 in order)
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000001_create_users_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000002_create_games_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000003_create_game_players_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000004_create_game_actions_table.up.sql

# Step 5: Verify
docker compose exec -T postgres psql -U postgres -d fasiolas_game -c "SELECT table_name FROM information_schema.tables WHERE table_schema='public' ORDER BY table_name;"

# You should see: game_actions, game_players, games, users
```

**✅ DONE! Go to http://localhost:3000**

---

## ✔️ Verification

After running either method, check:

- [ ] Can you open http://localhost:3000?
- [ ] Can you click "Login with Google"?
- [ ] Can you complete the login without errors?
- [ ] Does it show your Google profile name?

**If YES to all = ✅ Fixed!**

---

## 🚨 If Something Goes Wrong

1. **Error: "connection refused"**
   - PostgreSQL is still starting
   - Wait 20 seconds, try again

2. **Error: "table already exists"**
   - Good news! Migrations already ran
   - Just verify all 4 tables exist (see Verification step 5 above)

3. **Still getting database error**
   - Check: `docker logs fasiolas_postgres` (shows database logs)
   - Check: `docker logs fasiolas_app` (shows backend logs)
   - Check: `docker ps` (shows if containers are running)

4. **Need to completely restart**
   ```powershell
   docker compose down -v --remove-orphans
   docker system prune -a --volumes
   docker compose up -d
   # Then repeat either Method 1 or 2 above
   ```

---

## 📚 More Help

**Read these files for detailed help:**
- `DATABASE_FIX_GUIDE.md` - Full troubleshooting guide
- `START_HERE_DATABASE_FIX.md` - Quick start guide
- `QUICK_FIX_DATABASE_ERROR.md` - Migration steps

---

## ⏱️ Time Estimate

| Method | Time | Difficulty |
|--------|------|-----------|
| Method 1 (Auto) | 2-3 min | ⭐ Easy |
| Method 2 (Manual) | 5-7 min | ⭐⭐ Medium |
| Method 2 + Troubleshooting | 10-15 min | ⭐⭐⭐ Hard |

---

**Status: READY TO FIX ✅**

Pick a method above and let's get your database working!
