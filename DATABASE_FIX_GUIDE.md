# 🎮 Fasiolas Card Game - Database Setup & Launch Guide

## Problem You're Experiencing
When you log in with Google, you get this error:
```
{"error":"failed to get user: pq: relation \"users\" does not exist"}
```

**Why this happens:** The database migrations haven't been applied, so the `users` table doesn't exist in PostgreSQL.

---

## ✅ Solution - Complete Fix

### Option A: Automated Fix (Recommended)

Open PowerShell in your project directory and run this command:

```powershell
powershell -ExecutionPolicy Bypass -File .\restart-complete.ps1
```

This script will:
1. Stop all Docker containers
2. Clean up old volumes
3. Start fresh Docker containers
4. Wait for PostgreSQL to be ready
5. Apply all 4 database migrations
6. Verify tables were created
7. Display the status

**Then skip to Step 6 (Test the Application) below.**

---

### Option B: Manual Step-by-Step Fix

If the automated script doesn't work, follow these steps manually:

#### Step 1: Stop Everything
```powershell
docker compose down -v --remove-orphans
```
Wait 3-5 seconds.

#### Step 2: Start Docker
```powershell
docker compose up -d
```
Wait about 10 seconds for PostgreSQL to start.

#### Step 3: Copy Migration Files to Container

Copy all migration files to the database container:

```powershell
docker compose cp migrations/000001_create_users_table.up.sql postgres:/tmp/
docker compose cp migrations/000002_create_games_table.up.sql postgres:/tmp/
docker compose cp migrations/000003_create_game_players_table.up.sql postgres:/tmp/
docker compose cp migrations/000004_create_game_actions_table.up.sql postgres:/tmp/
```

#### Step 4: Apply Migrations

Run the migrations in order:

```powershell
# Set the password environment variable
$env:PGPASSWORD = "123456"

# Apply each migration
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000001_create_users_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000002_create_games_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000003_create_game_players_table.up.sql
docker compose exec -T postgres psql -U postgres -d fasiolas_game -f /tmp/000004_create_game_actions_table.up.sql

# Clear the password variable
Remove-Item env:PGPASSWORD
```

If you see `CREATE TABLE` or `CREATE INDEX` in the output, they worked! ✓

#### Step 5: Verify Tables Were Created

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

#### Step 6: Test the Application

1. Open http://localhost:3000 in your browser
2. Click "Login with Google"
3. Complete the Google OAuth login
4. You should now be logged in successfully! ✅

---

## 🔧 Troubleshooting

### Problem: "connection refused" or "PostgreSQL not ready"
**Solution:** Wait another 10-15 seconds and try again. PostgreSQL takes time to start.

### Problem: "table already exists"
**Solution:** This means the migrations already worked! Just verify all 4 tables exist using Step 5 above.

### Problem: Containers won't start
**Solution:** 
1. Make sure Docker Desktop is running
2. Try: `docker system prune -a --volumes` (careful - this deletes everything)
3. Then: `docker compose up -d` again

### Problem: "Docker not found"
**Solution:**
1. Install Docker Desktop from https://www.docker.com/products/docker-desktop
2. Restart your computer
3. Try again

---

## 📋 Verification Checklist

After following the fix, verify:

- [ ] `docker ps` shows:
  - [ ] `fasiolas_postgres` - running
  - [ ] `fasiolas_app` - running
  - [ ] `fasiolas_frontend` - running

- [ ] Database has 4 tables: `users`, `games`, `game_players`, `game_actions`

- [ ] Can access http://localhost:3000 in browser

- [ ] Can click "Login with Google" without errors

- [ ] Can complete Google login successfully

---

## 🚀 Next Steps (After Database is Fixed)

Once the database is working:

1. **Test Google Login**
   - Go to http://localhost:3000
   - Click "Login with Google"
   - You should be logged in as your Google account

2. **Start Building the Game**
   - Create a game
   - Join a game
   - Test game mechanics

3. **Check the Backend**
   - API docs: http://localhost:8080
   - Health check: http://localhost:8080/health

---

## 📞 Emergency Reset

If nothing works and you need to start completely fresh:

```powershell
# Stop everything
docker compose down -v --remove-orphans

# Remove all Docker data (WARNING: This deletes everything!)
docker system prune -a --volumes

# Start fresh
docker compose up -d

# Then apply migrations again (Step 4 above)
```

---

## 📚 Additional Resources

- **Docker Logs (to see what's happening):**
  ```powershell
  docker logs fasiolas_postgres  # Database logs
  docker logs fasiolas_app       # Backend logs
  docker logs fasiolas_frontend  # Frontend logs
  ```

- **Check Database Directly:**
  ```powershell
  docker compose exec -T postgres psql -U postgres -d fasiolas_game
  # Then in psql: SELECT * FROM users;
  ```

- **Rebuild Everything:**
  ```powershell
  docker compose build --no-cache
  docker compose up -d
  ```

---

**Good luck! 🎮 Once this is set up, the game should work smoothly.**
