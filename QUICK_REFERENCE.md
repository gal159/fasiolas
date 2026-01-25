# ⚡ Quick Commands Reference Card

Save this for quick reference!

---

## 🚀 START EVERYTHING (3 commands)

```powershell
# 1. Start Database
docker-compose up -d postgres

# 2. Run Migrations (first time only)
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up

# 3. Start Application
go run cmd/server/main.go
```

---

## 🛑 STOP EVERYTHING (2 commands)

```powershell
# 1. Stop Application: Press Ctrl + C

# 2. Stop Database
docker-compose down
```

---

## 🧪 QUICK TESTS

```powershell
# Health Check
curl http://localhost:8080/health

# Get OAuth URL (Google)
curl http://localhost:8080/api/v1/auth/google

# Get OAuth URL (GitHub)
curl http://localhost:8080/api/v1/auth/github
```

---

## 🔍 CHECK STATUS

```powershell
# Check Docker containers
docker-compose ps

# Check Docker logs
docker-compose logs postgres

# Check if port 8080 is in use
netstat -ano | findstr :8080

# Check if port 5432 is in use
netstat -ano | findstr :5432
```

---

## 🗄️ DATABASE COMMANDS

```powershell
# Connect to database
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game

# Inside psql:
\dt              # List tables
\d users         # Describe users table
SELECT * FROM users;
\q               # Exit

# View tables from PowerShell
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "\dt"

# Reset database (WARNING: Deletes all data!)
docker-compose down -v
docker-compose up -d postgres
# Wait 10 seconds then run migrations
```

---

## 🔧 TROUBLESHOOTING

### Port Already in Use

```powershell
# Find what's using the port
netstat -ano | findstr :8080

# Kill process (replace <PID>)
taskkill /PID <PID> /F
```

### Docker Not Running

```powershell
# Open Docker Desktop and wait for it to start
# Then retry: docker-compose up -d postgres
```

### Migration Errors

```powershell
# Check migration version
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" version

# Rollback all
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" down

# Reapply
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

### Clear Go Cache

```powershell
go clean -modcache
go mod download
go mod tidy
```

---

## 📝 COMMON TASKS

### Create New User Admin (Direct DB)

```powershell
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "UPDATE users SET role='admin' WHERE email='your@email.com';"
```

### View Recent Game Actions

```powershell
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "SELECT * FROM game_actions ORDER BY timestamp DESC LIMIT 10;"
```

### Count Active Games

```powershell
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "SELECT state, COUNT(*) FROM games GROUP BY state;"
```

---

## 🏗️ BUILD & DEPLOY

```powershell
# Build binary
go build -o bin/server.exe cmd/server/main.go

# Run binary
.\bin\server.exe

# Build Docker image
docker build -t fasiolas-game:latest .

# Run everything with Docker
docker-compose up --build -d
```

---

## 📦 DEPENDENCY MANAGEMENT

```powershell
# Install new package
go get github.com/some/package

# Update all dependencies
go get -u ./...
go mod tidy

# List dependencies
go list -m all

# Verify dependencies
go mod verify
```

---

## 🧹 CLEANUP

```powershell
# Remove Docker containers and volumes
docker-compose down -v

# Remove Docker images
docker rmi fasiolas_app fasiolas_postgres

# Clean Go build cache
go clean -cache

# Remove binary
Remove-Item -Path .\bin\server.exe
```

---

## 💡 DEVELOPMENT WORKFLOW

**Daily Startup:**
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker-compose up -d postgres
go run cmd/server/main.go
```

**Daily Shutdown:**
```powershell
Ctrl + C              # Stop Go app
docker-compose down   # Stop database
```

**Full Reset:**
```powershell
docker-compose down -v
docker-compose up -d postgres
# Wait 10 seconds
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
go run cmd/server/main.go
```

---

## 🎯 PROJECT PATHS

```
Project Root: C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame

Key Files:
- .env                    # Configuration
- docker-compose.yml      # Docker setup
- cmd/server/main.go      # Application entry
- migrations/             # Database schemas
- internal/               # Business logic
- README.md              # Full documentation
- SETUP_GUIDE.md         # Detailed setup
- THIS FILE              # Quick reference
```

---

## 📞 HELP

```powershell
# Go help
go help

# Docker help
docker-compose --help

# Migration help
migrate --help

# View project README
cat README.md
```

---

**Print this file and keep it next to your keyboard! 📌**

