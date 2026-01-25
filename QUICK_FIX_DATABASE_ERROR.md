# 🚀 Quick Restart Guide - Fix Database Error

## Problem
When logging in with Google, you get: `{"error":"failed to get user: pq: relation \"users\" does not exist"}`

This happens because the database migrations haven't been applied yet.

## Solution - Manual Step by Step

### Step 1: Stop Everything
```powershell
docker compose down -v --remove-orphans
```

Wait 3-5 seconds for containers to fully stop.

---

### Step 2: Start Docker Compose
```powershell
docker compose up -d
```

Wait about 10 seconds for PostgreSQL to be fully healthy.

---

### Step 3: Apply Database Migrations

Run these 4 commands **one by one** in order:

```powershell
# Set the database password
$env:PGPASSWORD = "123456"

# Migration 1: Create users table
Get-Content migrations/000001_create_users_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game

# Migration 2: Create games table
Get-Content migrations/000002_create_games_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game

# Migration 3: Create game_players table
Get-Content migrations/000003_create_game_players_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game

# Migration 4: Create game_actions table
Get-Content migrations/000004_create_game_actions_table.up.sql | docker compose exec -T postgres psql -U postgres -d fasiolas_game

# Clear the password
Remove-Item env:PGPASSWORD
```

**Each command should show "CREATE TABLE" or "CREATE INDEX" output without errors.**

---

### Step 4: Verify Tables Were Created

```powershell
docker compose exec -T postgres psql -U postgres -d fasiolas_game -c "SELECT table_name FROM information_schema.tables WHERE table_schema='public' ORDER BY table_name;"
```

You should see:
```
 table_name
---------------
 game_actions
 game_players
 games
 users
(4 rows)
```

---

### Step 5: Check Container Status

```powershell
docker ps
```

You should see:
- `fasiolas_postgres` - running (healthy)
- `fasiolas_app` - running
- `fasiolas_frontend` - running (or created)

---

### Step 6: Test the Application

1. Open http://localhost:3000 in your browser
2. Click "Login with Google"
3. You should now be able to login successfully without database errors! ✅

---

## Troubleshooting

### If PostgreSQL is not ready
- Wait another 10 seconds and try again
- Check logs: `docker logs fasiolas_postgres`

### If migrations fail with "connection refused"
- PostgreSQL hasn't started yet
- Wait longer and try again

### If migrations fail with "table already exists"
- The migrations have already been applied successfully
- Just check that all 4 tables exist (Step 4)

### If containers won't start
- Make sure Docker Desktop is running
- Try: `docker system prune -a` to clean up and start fresh
- Then restart with `docker compose up -d`

---

## Automated Solution

If you want to automate this entire process, run the PowerShell script:

```powershell
.\restart-project.ps1
```

Or the batch file:

```cmd
restart-project.bat
```

These scripts do all the steps above automatically.

---

## Next Steps

Once everything is running:
1. ✅ Database should have all 4 tables created
2. ✅ Backend should be connected to database
3. ✅ Frontend should be accessible
4. ✅ Google OAuth should work without errors

Start building and testing! 🎮
