# 🚀 Quick Start Guide - Fasiolas Card Game API

This guide will help you get the Fasiolas Card Game API up and running quickly.

## ⚡ Quick Setup (5 minutes)

### Step 1: Start PostgreSQL Database

```powershell
# Start PostgreSQL using Docker Compose
docker-compose up -d postgres

# Wait for database to be ready (check logs)
docker-compose logs postgres
```

### Step 2: Run Database Migrations

```powershell
# Apply all migrations
$env:DB_CONNECTION="postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable"
migrate -path migrations -database $env:DB_CONNECTION up
```

Or use the simpler command if you have Make installed:
```powershell
make migrate-up
```

### Step 3: Configure OAuth2 (Optional for Testing)

For local testing, the app will run without OAuth credentials, but authentication endpoints won't work. To enable OAuth:

1. **Google OAuth2**
   - Visit: https://console.cloud.google.com/
   - Create project → APIs & Services → Credentials
   - Create OAuth 2.0 Client ID
   - Add redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
   - Update `.env` with Client ID and Secret

2. **GitHub OAuth2**
   - Visit: https://github.com/settings/developers
   - New OAuth App
   - Callback URL: `http://localhost:8080/api/v1/auth/github/callback`
   - Update `.env` with Client ID and Secret

3. **Discord OAuth2**
   - Visit: https://discord.com/developers/applications
   - New Application → OAuth2
   - Add redirect: `http://localhost:8080/api/v1/auth/discord/callback`
   - Update `.env` with Client ID and Secret

### Step 4: Run the Application

```powershell
# Run directly
go run cmd/server/main.go

# Or build first
go build -o bin/server.exe cmd/server/main.go
.\bin\server.exe
```

You should see:
```
✓ Database connected successfully
✓ Server starting on :8080
✓ Environment: development
✓ API Documentation: http://localhost:8080/api/v1
```

### Step 5: Test the API

Open a new PowerShell window and test the health endpoint:

```powershell
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "fasiolas-card-game"
}
```

## 🎮 Testing Game Flow

### 1. Authentication (with OAuth configured)

```powershell
# Get Google OAuth URL
curl http://localhost:8080/api/v1/auth/google
```

Copy the URL from response and open in browser to complete OAuth flow.

### 2. Create a Game

```powershell
$token = "YOUR_JWT_TOKEN_HERE"

curl -X POST http://localhost:8080/api/v1/games `
  -H "Authorization: Bearer $token" `
  -H "Content-Type: application/json" `
  -d '{\"max_players\": 4}'
```

### 3. Join a Game

```powershell
curl -X POST http://localhost:8080/api/v1/games/join `
  -H "Authorization: Bearer $token" `
  -H "Content-Type: application/json" `
  -d '{\"room_code\": \"ABC123\"}'
```

### 4. Get Game State

```powershell
curl http://localhost:8080/api/v1/games/1 `
  -H "Authorization: Bearer $token"
```

### 5. Start Game

```powershell
curl -X POST http://localhost:8080/api/v1/games/1/start `
  -H "Authorization: Bearer $token"
```

### 6. Draw Card (Phase 1)

```powershell
curl -X POST http://localhost:8080/api/v1/games/1/draw `
  -H "Authorization: Bearer $token"
```

### 7. Place Card

```powershell
curl -X POST http://localhost:8080/api/v1/games/1/place `
  -H "Authorization: Bearer $token" `
  -H "Content-Type: application/json" `
  -d '{\"target_player_position\": 1}'
```

## 🐛 Troubleshooting

### Database Connection Failed

```powershell
# Check if PostgreSQL is running
docker-compose ps

# View PostgreSQL logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Port 8080 Already in Use

```powershell
# Find process using port 8080
netstat -ano | findstr :8080

# Kill the process (replace PID with actual process ID)
taskkill /PID <PID> /F

# Or change port in .env
# PORT=8081
```

### Migration Errors

```powershell
# Check migration status
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" version

# Rollback migrations
make migrate-down

# Reapply migrations
make migrate-up
```

### Module Not Found Errors

```powershell
# Clear Go module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy up
go mod tidy
```

## 📊 Database Management

### View Database Tables

```powershell
# Connect to PostgreSQL
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game

# List tables
\dt

# View users
SELECT * FROM users;

# View games
SELECT * FROM games;

# Exit
\q
```

### Reset Database

```powershell
# Stop and remove containers with volumes
docker-compose down -v

# Start fresh
docker-compose up -d postgres

# Wait a few seconds, then run migrations
timeout /t 10
make migrate-up
```

## 🔧 Development Commands

```powershell
# Run tests
go test ./...

# Build for production
go build -ldflags="-s -w" -o bin/server.exe cmd/server/main.go

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Generate dependency graph
go mod graph

# Update dependencies
go get -u ./...
go mod tidy
```

## 🐳 Docker Deployment

### Build and Run with Docker Compose

```powershell
# Build and start all services
docker-compose up --build -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Manual Docker Build

```powershell
# Build image
docker build -t fasiolas-game:latest .

# Run container
docker run -d `
  --name fasiolas-app `
  -p 8080:8080 `
  --env-file .env `
  fasiolas-game:latest
```

## 📈 Performance Testing

### Load Testing with Apache Bench (if installed)

```powershell
# Test health endpoint
ab -n 1000 -c 10 http://localhost:8080/health
```

### Memory Profiling

```powershell
# Run with profiling
go run cmd/server/main.go -cpuprofile=cpu.prof -memprofile=mem.prof

# Analyze with pprof
go tool pprof cpu.prof
go tool pprof mem.prof
```

## 🔒 Security Notes

1. **JWT Secret**: Change `JWT_SECRET` in `.env` before production
   ```powershell
   # Generate secure secret (requires OpenSSL)
   openssl rand -base64 32
   ```

2. **Database Password**: Use strong passwords in production

3. **OAuth Credentials**: Never commit `.env` file to version control

4. **HTTPS**: Use HTTPS in production (configure reverse proxy like nginx)

## 📚 Additional Resources

- **API Documentation**: http://localhost:8080/api/v1
- **Database Migrations**: `./migrations/`
- **Game Rules**: See README.md "Game Rules" section
- **Source Code**: Explore `./internal/` and `./pkg/` directories

## 🎯 Next Steps

1. ✅ Complete OAuth2 configuration
2. ✅ Test all API endpoints
3. ✅ Create multiple test users
4. ✅ Play a full game through the API
5. ✅ Explore external API integrations
6. ✅ Review game logs in database

## 💡 Tips

- Use Postman or Insomnia for easier API testing
- Enable detailed logging by setting `ENV=development` in `.env`
- Check `game_actions` table to see game history
- Use Docker for consistent development environment
- Keep PostgreSQL running in background for development

---

**Need Help?** Check the main README.md or create an issue on GitHub.

