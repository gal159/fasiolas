# 🔧 Developer Workflow Guide

## Getting Started

### Prerequisites
- Go 1.21+
- Docker Desktop (for PostgreSQL)
- Git
- PowerShell or bash terminal

### Initial Setup (One-time)

```powershell
# 1. Clone and navigate to project
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"

# 2. Install Go dependencies
go mod download
go mod tidy

# 3. Start PostgreSQL with Docker
docker-compose up -d postgres

# 4. Wait for database to be ready (check with: docker-compose logs postgres)
# Look for "database system is ready to accept connections"

# 5. Install migration tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 6. Run migrations
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up

# 7. Build the application
go build -o bin/server.exe cmd/server/main.go

# 8. Start the application
go run cmd/server/main.go
```

---

## Daily Development Workflow

### Start Services
```powershell
# Start database only (app development mode)
docker-compose up -d postgres

# Wait for database to be ready
docker-compose logs postgres  # Look for "ready to accept connections"

# Start application in development mode with hot reload (optional, requires: go get github.com/cosmtrek/air)
air
# OR manual rebuild
go run cmd/server/main.go
```

### Testing Workflow
```powershell
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/game/...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestCardComparison ./pkg/cards/...
```

### Code Quality
```powershell
# Format code
go fmt ./...

# Lint code (requires: go get github.com/golangci/golangci-lint)
golangci-lint run ./...

# Check for vet issues
go vet ./...

# Get code coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Database Operations
```powershell
# Connect to database
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game

# Within psql:
# \dt              - List tables
# \d users         - Describe users table
# SELECT * FROM users;
# \q               - Exit

# Quick queries from PowerShell:
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "SELECT COUNT(*) FROM games;"

# Reset database (DELETE ALL DATA!)
docker-compose down -v
docker-compose up -d postgres
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

### Stop Services
```powershell
# Stop application: Press Ctrl+C

# Stop database
docker-compose down

# Stop database and remove data
docker-compose down -v
```

---

## Project Structure Reference

```
cardGame/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── game/
│   │   └── engine.go              # Game rules & logic
│   ├── handler/
│   │   ├── auth_handler.go        # OAuth & auth endpoints
│   │   └── game_handler.go        # Game endpoints
│   ├── middleware/
│   │   └── middleware.go          # Auth, CORS, logging middleware
│   ├── models/
│   │   └── models.go              # Data structures
│   ├── repository/
│   │   ├── db.go                  # Database connection
│   │   ├── user_repository.go     # User data access
│   │   ├── game_repository.go     # Game data access
│   │   ├── game_player_repository.go
│   │   └── game_action_repository.go
│   └── service/
│       ├── auth_service.go        # Authentication logic
│       ├── game_service.go        # Game business logic
│       └── external_api_service.go # Deck API integration
├── pkg/
│   ├── cards/
│   │   └── cards.go              # Card & deck structures
│   └── utils/
│       └── utils.go              # Utility functions
├── migrations/
│   ├── 000001_create_users_table.*sql
│   ├── 000002_create_games_table.*sql
│   ├── 000003_create_game_players_table.*sql
│   └── 000004_create_game_actions_table.*sql
├── .env                           # Environment variables (create this)
├── docker-compose.yml             # Docker services
├── Dockerfile                     # Container image
├── go.mod, go.sum                # Dependency management
└── README.md, SETUP_GUIDE.md      # Documentation
```

---

## Adding a New Feature

### Example: Adding a "LeaveGame" endpoint

#### 1. Add to Models (internal/models/models.go)
```go
const (
    ActionLeave ActionType = "leave"
)
```

#### 2. Add to Service (internal/service/game_service.go)
```go
func (s *GameService) LeaveGame(gameID, userID int) error {
    // Implementation
}
```

#### 3. Add to Handler (internal/handler/game_handler.go)
```go
// LeaveGame godoc
// @Summary Leave a game
// @Tags games
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/games/{id}/leave [post]
func (h *GameHandler) LeaveGame(c *gin.Context) {
    // Implementation
}
```

#### 4. Add Route (cmd/server/main.go)
```go
protected.POST("/games/:id/leave", gameHandler.LeaveGame)
```

#### 5. Write Tests
```go
func TestLeaveGame(t *testing.T) {
    // Test implementation
}
```

#### 6. Test Manually
```bash
curl -X POST http://localhost:8080/api/v1/games/1/leave \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Debugging Tips

### Enable Debug Logging
Add to your code:
```go
import "log"

log.Printf("DEBUG: %+v", variable)
```

### Database Query Debugging
Enable PostgreSQL query logging in docker-compose.yml:
```yaml
environment:
  - POSTGRES_INITDB_ARGS=-c log_statement=all
```

### Gin Debug Mode
```go
import "github.com/gin-gonic/gin"

gin.SetMode(gin.DebugMode)  // For verbose logging
```

### Inspect HTTP Requests
Use Postman or add middleware:
```go
router.Use(gin.Logger())
```

### Common Issues

**Issue:** "database not ready"
```
Solution: 
docker-compose logs postgres
# Wait for "ready to accept connections"
```

**Issue:** "port already in use"
```powershell
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

**Issue:** "migration failed"
```powershell
# Check current version
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" version

# Rollback one step
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" down
```

---

## Writing Tests

### Unit Test Template
```go
package game

import (
    "testing"
    "cardgame/pkg/cards"
)

func TestCardIsOnePlus(t *testing.T) {
    card1 := cards.NewCard(cards.Hearts, cards.Seven)
    card2 := cards.NewCard(cards.Hearts, cards.Eight)
    
    if !card2.IsOnePlus(card1) {
        t.Errorf("Expected 8 to be +1 of 7")
    }
}
```

### Integration Test Template
```go
func TestGameFlow(t *testing.T) {
    // Setup
    db := setupTestDB()
    defer db.Close()
    
    gameService := service.NewGameService(...)
    
    // Test
    game, err := gameService.CreateGame(1, 4)
    if err != nil {
        t.Fatalf("CreateGame failed: %v", err)
    }
    
    // Verify
    if game.State != models.StateWaiting {
        t.Errorf("Expected state waiting, got %v", game.State)
    }
}
```

---

## Deployment to Staging/Production

### Build Docker Image
```powershell
# Development
docker build -t fasiolas-game:dev .

# Production
docker build --build-arg ENV=production -t fasiolas-game:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=postgres \
  -e DB_PASSWORD=<secure-password> \
  -e GOOGLE_CLIENT_ID=<your-id> \
  -e GOOGLE_CLIENT_SECRET=<your-secret> \
  fasiolas-game:latest
```

### Docker Compose Production
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Environment Variables for Production
```bash
# Create .env.prod
SERVER_ENV=production
JWT_SECRET=<very-strong-random-secret>
GOOGLE_CLIENT_SECRET=<secure-value>
GITHUB_CLIENT_SECRET=<secure-value>
DB_PASSWORD=<very-secure-password>
```

---

## Git Workflow

```powershell
# Create feature branch
git checkout -b feature/add-leaderboard

# Make changes and test
go test ./...
go fmt ./...

# Commit with clear message
git add .
git commit -m "feat: add leaderboard endpoint"

# Push and create PR
git push origin feature/add-leaderboard
```

---

## Performance Profiling

```go
import _ "net/http/pprof"

// In main.go, add:
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// Then visit:
// http://localhost:6060/debug/pprof/
```

---

## Useful Commands Reference

```powershell
# Build
go build ./cmd/server

# Run
go run cmd/server/main.go

# Test
go test ./...

# Format
go fmt ./...

# Dependencies
go mod tidy
go mod download
go list -m all

# Generate docs (if using godoc)
godoc -http=:6060

# Check for security issues (requires gosec)
gosec ./...
```

---

## Notes

- Always run `go fmt` before committing
- Write tests for new features
- Update README if adding features
- Keep .env file with local development values
- Never commit .env to Git (use .gitignore)
- Use meaningful commit messages
- Test both success and error cases

