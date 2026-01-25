# 🚀 Complete Fasiolas Setup (Backend + Frontend)

This guide shows how to run the complete Fasiolas Card Game application with both backend API and React frontend.

---

## 📋 Prerequisites

- Docker Desktop (running)
- Go 1.21+
- Node.js 18+ & npm
- Git

---

## ⚡ Option 1: Quick Start (Everything in Docker)

### 1. Start All Services

```bash
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"

# Start database, backend, and frontend
docker-compose up --build
```

### 2. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### 3. Stop All Services

```bash
docker-compose down
```

---

## 🏃 Option 2: Development Setup (Local + Docker)

### 1. Start Database with Docker

```bash
docker-compose up -d postgres
docker-compose logs postgres  # Wait for "ready to accept connections"
```

### 2. Run Backend (Go) - Terminal 1

```bash
# Run migrations
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up

# Start backend
go run cmd/server/main.go
```

You should see:
```
✓ Database connected successfully
✓ Server starting on :8080
```

### 3. Run Frontend (React) - Terminal 2

```bash
cd frontend
npm install  # Only needed first time
npm start
```

Frontend opens at: http://localhost:3000

### 4. Stop Everything

```bash
# Terminal 1: Ctrl+C (backend)
# Terminal 2: Ctrl+C (frontend)
docker-compose down  # Stop database
```

---

## 📁 Project Structure

```
cardGame/
├── cmd/server/              ✅ Backend (Go)
├── internal/                ✅ Backend code
├── pkg/                     ✅ Backend utilities
├── migrations/              ✅ Database migrations
├── frontend/                🆕 Frontend (React)
│   ├── src/
│   │   ├── pages/          # Login, Dashboard, Game
│   │   ├── components/     # Navigation, GameCard, GameBoard
│   │   └── App.jsx         # Main app
│   ├── public/             # HTML entry point
│   ├── package.json        # React dependencies
│   ├── Dockerfile          # Container for frontend
│   └── tailwind.config.js  # Styling config
├── docker-compose.yml      # 🆕 Now includes frontend!
└── FRONTEND_GUIDE.md       # 🆕 Frontend documentation
```

---

## 🎯 Application Flow

```
1. User opens http://localhost:3000
   ↓
2. Frontend loads (React)
   ↓
3. User clicks login (Google/GitHub/Discord)
   ↓
4. Frontend calls: GET /api/v1/auth/{provider}
   ↓
5. Backend (Go) returns OAuth URL
   ↓
6. Browser redirects to OAuth provider
   ↓
7. User authenticates
   ↓
8. Frontend receives auth token
   ↓
9. Frontend calls: GET /api/v1/auth/profile
   ↓
10. Stores token in localStorage
    ↓
11. Shows Dashboard with games
    ↓
12. User can create/join games
    ↓
13. Frontend polls: GET /api/v1/games/{id}
    ↓
14. Displays game state in real-time
```

---

## 🔐 OAuth Configuration Required

Before first login, configure OAuth credentials:

### .env File

Create or update `.env` in project root:

```env
GOOGLE_CLIENT_ID=your-google-id
GOOGLE_CLIENT_SECRET=your-google-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

GITHUB_CLIENT_ID=your-github-id
GITHUB_CLIENT_SECRET=your-github-secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/auth/github/callback

DISCORD_CLIENT_ID=your-discord-id
DISCORD_CLIENT_SECRET=your-discord-secret
DISCORD_REDIRECT_URL=http://localhost:8080/api/v1/auth/discord/callback
```

Or get them from OAuth providers:
- Google: https://console.cloud.google.com/
- GitHub: https://github.com/settings/developers
- Discord: https://discord.com/developers/applications

---

## 🌐 Frontend Features

### ✅ Authentication
- OAuth2 login (3 providers)
- Automatic token management
- Protected routes
- User profile display

### ✅ Game Management
- Create new games
- Browse available games
- Join games with room codes
- List all players

### ✅ Game Play
- Real-time game state
- Place cards
- Draw cards
- View current player
- Display table cards (Phase 2)

### ✅ UI/UX
- Dark theme (Tailwind CSS)
- Responsive design
- Real-time updates (2-second polling)
- Error handling
- Copy room code button

---

## 🛠️ Development Workflow

### Making Changes to Backend

1. Edit code in `internal/`, `cmd/`, or `pkg/`
2. Backend reloads automatically (with `go run`)
3. Frontend auto-refreshes via proxy

### Making Changes to Frontend

1. Edit code in `frontend/src/`
2. Frontend auto-reloads (React dev server)
3. Changes visible instantly at http://localhost:3000

### Adding API Endpoints

1. Add handler in `internal/handler/`
2. Add route in `cmd/server/main.go`
3. Call from frontend using `axios.get('/api/v1/...')`

---

## 🐳 Docker Commands

### Build Images
```bash
docker-compose build --no-cache
```

### View Logs
```bash
docker-compose logs -f                    # All services
docker-compose logs -f postgres           # Database only
docker-compose logs -f app                # Backend only
docker-compose logs -f frontend           # Frontend only
```

### Stop and Remove
```bash
docker-compose down                       # Stop containers
docker-compose down -v                    # Also remove volumes (data!)
```

### Run Single Service
```bash
docker-compose up -d postgres             # Only database
docker-compose up app                     # Only backend (needs postgres running)
docker-compose up frontend                # Only frontend (needs app running)
```

---

## 🆘 Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
netstat -ano | findstr :3000              # Frontend
netstat -ano | findstr :8080              # Backend
netstat -ano | findstr :5432              # Database

# Kill process (replace PID)
taskkill /PID <PID> /F
```

### Frontend Can't Connect to Backend

1. Verify backend is running: `curl http://localhost:8080/health`
2. Check proxy in `frontend/package.json`: `"proxy": "http://localhost:8080"`
3. Ensure CORS is enabled in backend

### OAuth Not Working

1. Verify credentials in `.env`
2. Check redirect URIs match in OAuth provider settings
3. Ensure backend and frontend are running
4. Check backend logs: `docker-compose logs app`

### Database Connection Error

1. Check PostgreSQL is running: `docker-compose ps`
2. Verify database is ready: `docker-compose logs postgres`
3. Check connection string in `.env`

### Styles Not Loading (Frontend)

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm start
```

---

## 📊 Service Health Checks

### Test Backend

```bash
# Health check
curl http://localhost:8080/health

# Get games
curl http://localhost:8080/api/v1/games

# Get OAuth URL
curl http://localhost:8080/api/v1/auth/google
```

### Test Database

```bash
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "SELECT COUNT(*) FROM users;"
```

### Test Frontend

Visit: http://localhost:3000

Should see login page with 3 OAuth options.

---

## 🚀 Deployment

### Docker Deployment

```bash
# Build all services
docker-compose build

# Run in production mode
docker-compose -f docker-compose.yml up -d

# View logs
docker-compose logs -f
```

### Scaling

```bash
# Run multiple backend instances (load balancing)
docker-compose up -d --scale app=3

# Frontend still single instance
docker-compose up -d frontend
```

### Environment Variables for Production

Update `.env` or `docker-compose.yml`:

```env
SERVER_ENV=production
JWT_SECRET=<very-strong-random-secret>
DB_PASSWORD=<secure-password>
GOOGLE_CLIENT_SECRET=<secure-secret>
GITHUB_CLIENT_SECRET=<secure-secret>
DISCORD_CLIENT_SECRET=<secure-secret>
```

---

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| FRONTEND_GUIDE.md | Frontend setup and usage |
| DEVELOPER_GUIDE.md | Development workflow |
| API_TESTING.md | API endpoint reference |
| GAME_RULES.md | Game rules explained |
| SETUP_GUIDE.md | Backend setup |
| Getting_STARTED.md | Quick start |

---

## 🎓 Learning Resources

### Frontend (React)
- `frontend/src/App.jsx` - Main component routing
- `frontend/src/pages/Login.jsx` - OAuth implementation
- `frontend/src/pages/Dashboard.jsx` - Game list management
- `frontend/src/pages/Game.jsx` - Game play interface

### Backend (Go)
- `cmd/server/main.go` - Server setup
- `internal/handler/` - HTTP handlers
- `internal/service/` - Business logic
- `internal/game/engine.go` - Game rules

### Full Stack Flow
- Login → Dashboard → Game → Play → End Game

---

## ✅ Checklist

Before going live:

- [ ] OAuth credentials configured
- [ ] Database migrations run
- [ ] Backend runs without errors
- [ ] Frontend loads at http://localhost:3000
- [ ] Can login with OAuth
- [ ] Can create a game
- [ ] Can join a game
- [ ] Can place cards
- [ ] Can draw cards
- [ ] Game state updates in real-time

---

## 🎉 You're All Set!

The complete Fasiolas application is ready:

✅ **Backend API** - Production-ready Go REST API
✅ **Frontend UI** - Modern React web application
✅ **Database** - PostgreSQL with migrations
✅ **Authentication** - OAuth2 with 3 providers
✅ **Game Logic** - Both phases fully implemented
✅ **Docker** - Complete containerization

### Next Steps

1. **Start Development**
   ```bash
   # Terminal 1: Start database and backend
   docker-compose up -d postgres
   go run cmd/server/main.go
   
   # Terminal 2: Start frontend
   cd frontend && npm start
   ```

2. **Create a Game**
   - Visit http://localhost:3000
   - Login with OAuth
   - Click "Create Game"
   - Share room code with friends

3. **Play the Game**
   - Others join with room code
   - Start game with minimum 2 players
   - Place cards and draw as needed

---

**Happy Gaming! 🎴 Let's play Fasiolas! 🚀**

