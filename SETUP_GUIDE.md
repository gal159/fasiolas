# 🚀 Step-by-Step Setup and Launch Instructions

## Prerequisites Check

Before starting, ensure you have installed:
- ✅ Docker Desktop (download from https://www.docker.com/products/docker-desktop/)
- ✅ Go 1.21+ (check with: `go version`)
- ✅ Git

---

## 📋 OPTION 1: Using Docker (Recommended - Easiest)

This option automatically sets up PostgreSQL database in a container.

### Step 1: Start Docker Desktop
1. Open **Docker Desktop** application on Windows
2. Wait until Docker Desktop shows "Engine running" (green indicator)

### Step 2: Start PostgreSQL Database with Docker Compose

Open PowerShell in your project directory and run:

```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker-compose up -d postgres
```

**What this does:**
- Downloads PostgreSQL 15 image (first time only)
- Creates and starts PostgreSQL container
- Database will be accessible on `localhost:5432`

**Verify it's running:**
```powershell
docker-compose ps
```

You should see:
```
NAME                  COMMAND                  SERVICE    STATUS
fasiolas_postgres     "docker-entrypoint.s…"   postgres   Up
```

### Step 3: Run Database Migrations

Install the migration tool (one-time setup):
```powershell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations to create tables:
```powershell
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

**Expected output:**
```
1/u create_users_table (123.456ms)
2/u create_games_table (89.123ms)
3/u create_game_players_table (67.890ms)
4/u create_game_actions_table (45.678ms)
```

### Step 4: Configure OAuth2 Providers (Optional but Recommended)

Edit the `.env` file and add your OAuth credentials:

**For Google OAuth:**
1. Go to https://console.cloud.google.com/
2. Create a new project or select existing
3. Enable "Google+ API"
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client ID"
5. Add authorized redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
6. Copy Client ID and Client Secret to `.env` file

**For GitHub OAuth:**
1. Go to https://github.com/settings/developers
2. Click "New OAuth App"
3. Set Authorization callback URL: `http://localhost:8080/api/v1/auth/github/callback`
4. Copy Client ID and Client Secret to `.env` file

**For Discord OAuth:**
1. Go to https://discord.com/developers/applications
2. Create "New Application"
3. Go to OAuth2 → Add redirect: `http://localhost:8080/api/v1/auth/discord/callback`
4. Copy Client ID and Client Secret to `.env` file

> **Note:** You need at least ONE OAuth provider configured to test authentication!

### Step 5: Run the Application

```powershell
go run cmd/server/main.go
```

**Expected output:**
```
✓ Database connected successfully
✓ Server starting on :8080
✓ Environment: development
✓ API Documentation: http://localhost:8080/api/v1
[GIN-debug] Listening and serving HTTP on :8080
```

### Step 6: Test the API

Open browser or use curl:

**Health Check:**
```powershell
curl http://localhost:8080/health
```

**Get Google OAuth URL:**
```powershell
curl http://localhost:8080/api/v1/auth/google
```

### Step 7: Stop Everything When Done

**Stop the application:**
- Press `Ctrl + C` in the terminal where Go is running

**Stop the database:**
```powershell
docker-compose down
```

**Stop and remove all data (full reset):**
```powershell
docker-compose down -v
```

---

## 📋 OPTION 2: Using Local PostgreSQL (Without Docker)

If you have PostgreSQL already installed locally.

### Step 1: Start PostgreSQL Service

**On Windows (if installed as service):**
```powershell
# Check if running
Get-Service -Name postgresql*

# Start if not running
Start-Service postgresql-x64-15
```

**Or use pgAdmin:**
1. Open pgAdmin 4
2. Start the PostgreSQL server

### Step 2: Create Database

Connect to PostgreSQL (using psql or pgAdmin) and run:

```sql
CREATE DATABASE fasiolas_game;
```

**Using psql command line:**
```powershell
psql -U postgres
```
Then in psql:
```sql
CREATE DATABASE fasiolas_game;
\q
```

### Step 3: Update .env File

Make sure your `.env` has correct PostgreSQL connection details:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=fasiolas_game
DB_SSL_MODE=disable
```

### Step 4: Run Migrations

```powershell
migrate -path migrations -database "postgresql://postgres:your_password@localhost:5432/fasiolas_game?sslmode=disable" up
```

### Step 5: Run the Application

```powershell
go run cmd/server/main.go
```

---

## 🔧 Troubleshooting Common Issues

### Issue 1: "Cannot connect to Docker daemon"

**Solution:**
1. Open Docker Desktop
2. Wait for it to fully start (shows "Engine running")
3. Try the command again

### Issue 2: "Port 5432 already in use"

**Solution:**
You have another PostgreSQL instance running.

**Option A:** Stop local PostgreSQL
```powershell
Stop-Service postgresql-x64-15
```

**Option B:** Change Docker port in `docker-compose.yml`:
```yaml
ports:
  - "5433:5432"  # Use 5433 instead
```
Then update `.env`:
```env
DB_PORT=5433
```

### Issue 3: "Port 8080 already in use"

**Solution:**
Change the port in `.env`:
```env
PORT=8081
```

Or kill the process using port 8080:
```powershell
# Find process using port 8080
netstat -ano | findstr :8080

# Kill it (replace <PID> with actual process ID)
taskkill /PID <PID> /F
```

### Issue 4: Migration fails - "relation already exists"

**Solution:**
Reset the database:
```powershell
# Drop and recreate database
psql -U postgres -c "DROP DATABASE fasiolas_game;"
psql -U postgres -c "CREATE DATABASE fasiolas_game;"

# Run migrations again
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

### Issue 5: "OAuth error" when testing authentication

**Solution:**
1. Verify OAuth credentials in `.env` are correct
2. Check redirect URLs match exactly (including http vs https)
3. Ensure OAuth app is enabled in provider console

### Issue 6: "go: command not found" or "migrate: command not found"

**Solution:**
1. Ensure Go is installed: Download from https://go.dev/dl/
2. Add Go to PATH:
   - Search "Environment Variables" in Windows
   - Add `C:\Go\bin` and `%USERPROFILE%\go\bin` to PATH
3. Restart PowerShell

---

## 🧪 Quick Test Workflow

Once everything is running, test the complete flow:

### 1. Health Check
```powershell
curl http://localhost:8080/health
```

### 2. Get OAuth URL (example with Google)
```powershell
curl http://localhost:8080/api/v1/auth/google
```

### 3. Complete OAuth in Browser
- Copy the URL from step 2
- Open in browser
- Log in with your Google account
- You'll get redirected with a token

### 4. Use the Token
```powershell
# Save your token
$TOKEN = "your-jwt-token-here"

# Get your profile
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/auth/profile

# Create a game
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{\"max_players\": 4}' http://localhost:8080/api/v1/games

# List games
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/games
```

---

## 📊 Verify Database Tables Were Created

**Using Docker:**
```powershell
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "\dt"
```

**Using local psql:**
```powershell
psql -U postgres -d fasiolas_game -c "\dt"
```

**Expected output:**
```
              List of relations
 Schema |      Name       | Type  |  Owner   
--------+-----------------+-------+----------
 public | game_actions    | table | postgres
 public | game_players    | table | postgres
 public | games           | table | postgres
 public | schema_migrations | table | postgres
 public | users           | table | postgres
```

---

## 🎮 Ready to Play!

Your Fasiolas card game API is now running! 

**Base URL:** http://localhost:8080/api/v1

**Next Steps:**
1. Configure at least one OAuth provider
2. Test authentication flow
3. Create a game
4. Invite friends with the room code
5. Start playing!

---

## 📱 Development Tips

### Run with Auto-Restart (using Air)

Install Air for hot reload:
```powershell
go install github.com/cosmtrek/air@latest
```

Create `.air.toml` and run:
```powershell
air
```

### View Docker Logs
```powershell
docker-compose logs -f postgres
```

### View All Running Containers
```powershell
docker ps
```

### Build Production Binary
```powershell
go build -o bin/server.exe cmd/server/main.go
./bin/server.exe
```

---

## 🛑 Shutdown Checklist

When you're done for the day:

1. **Stop the Go application:** Press `Ctrl + C`
2. **Stop Docker containers:** `docker-compose down`
3. **Close Docker Desktop** (optional, saves resources)

To completely reset everything (including database data):
```powershell
docker-compose down -v
```

---

## ✅ Success Indicators

You'll know everything is working when:
- ✅ Docker shows "fasiolas_postgres" container is running
- ✅ Go server shows "Database connected successfully"
- ✅ Health check returns `{"status":"ok"}`
- ✅ OAuth URL endpoints return authorization URLs
- ✅ No error messages in the console

---

**Need Help?** 
- Check the logs carefully for error messages
- Verify all environment variables in `.env`
- Make sure Docker is running before starting containers
- Ensure no port conflicts (5432, 8080)

**Good Luck! 🎴**

