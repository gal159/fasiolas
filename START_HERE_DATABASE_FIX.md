# 🚀 START HERE - Fix Database Error

## Your Problem
When you login with Google, you get:
```
{"error":"failed to get user: pq: relation \"users\" does not exist"}
```

## Quick Solution (Choose One)

### ✅ Option 1: Automatic (Easiest)
Open PowerShell in the cardGame folder and run:
```powershell
powershell -ExecutionPolicy Bypass -File .\restart-complete.ps1
```

The script will:
- Stop Docker
- Start fresh Docker containers  
- Apply all 4 database migrations automatically
- Verify everything is working
- Tell you when it's done

Then go to http://localhost:3000 and login! ✅

---

### ✅ Option 2: Manual (If automatic doesn't work)

**Step 1:** Stop Docker
```powershell
docker compose down -v --remove-orphans
```
Wait 3 seconds.

**Step 2:** Start Docker  
```powershell
docker compose up -d
```
Wait 10 seconds.

**Step 3:** Copy and apply migrations
```powershell
docker compose cp migrations/000001_create_users_table.up.sql postgres:/tmp/
docker compose cp migrations/000002_create_games_table.up.sql postgres:/tmp/
docker compose cp migrations/000003_create_game_players_table.up.sql postgres:/tmp/
docker compose cp migrations/000004_create_game_actions_table.up.sql postgres:/tmp/

docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000001_create_users_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000002_create_games_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000003_create_game_players_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000004_create_game_actions_table.up.sql
```

**Step 4:** Verify tables exist
```powershell
docker compose exec -T postgres psql -U postgres -d fasiolas_game -c "SELECT table_name FROM information_schema.tables WHERE table_schema='public' ORDER BY table_name;"
```

Should see: `users`, `games`, `game_players`, `game_actions`

**Step 5:** Test  
Open http://localhost:3000 and login with Google! ✅

---

## 📖 Full Documentation
- See: `DATABASE_FIX_GUIDE.md` for detailed troubleshooting
- See: `QUICK_FIX_DATABASE_ERROR.md` for quick reference

## ⚡ TL;DR
```powershell
docker compose down -v --remove-orphans
docker compose up -d
Start-Sleep -Seconds 10
# Copy and apply each migration file from migrations/ folder
# Then go to http://localhost:3000 and login!
```

---

**That's it! The database error is fixed. 🎮**
