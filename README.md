# 🎴 Fasiolas Card Game

A multiplayer Lithuanian card game with OAuth authentication and role-based access control.

---

## 🚀 Quick Start (5 Minutes)

### Prerequisites
- Docker Desktop installed and running
- Ports 3000, 8080, 5432 available

### Launch Application

```bash
# Navigate to project folder
cd cardGame

# Start all services (frontend, backend, database)
docker compose up

# Wait for these messages:
# ✓ postgres: ready
# ✓ app: Server starting on :8080
# ✓ frontend: webpack compiled successfully
```

**Access the application:**
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/health

---

## 👤 Create First Admin User

After starting the application, create an admin account:

```bash
# Open new terminal and run:
docker compose exec db psql -U postgres -d fasiolas_game -c "INSERT INTO users (email, username, role, created_at, updated_at) VALUES ('your-email@gmail.com', 'Admin', 'admin', NOW(), NOW()) ON CONFLICT DO NOTHING;"
```

Replace `your-email@gmail.com` with your actual email.

---

## 📋 User Roles

| Role | Create Games | Join Games | Play | Watch Games | Admin Panel |
|------|:------------:|:----------:|:----:|:-----------:|:-----------:|
| **👨‍💼 ADMIN** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **🎮 PLAYER** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **👁️ SPECTATOR** | ❌ | ❌ | ❌ | ✅ | ❌ |

**Default role for new users**: PLAYER

---

## 🎮 How to Play

### 1. Login
- Visit http://localhost:3000
- Click "Google Login" (or GitHub/Discord)
- Authorize the application

### 2. Create Game (Players & Admins)
- Click **"Create Game"** button
- Select number of players (2-8)
- Click **"Create"**
- Share room code with friends

### 3. Join Game (Players & Admins)
- Enter room code in **"Join by Room Code"** section
- Click **"Join Game"**
- Wait for other players

### 4. Start Game
- When 2+ players joined, click **"START GAME"**
- Game begins!

### 5. Gameplay
- **Phase 1**: Place cards on other players (+1 rank rule)
- **Phase 2**: Get rid of all your cards
- **Last player with cards loses**

### 6. Watch Game (All Roles)
- Spectators see **"Watch"** button (purple)
- Players/Admins see **"Join"** button (blue)
- Click to view game in real-time

---

## 🛠️ Common Commands

```bash
# Start application (normal startup after restart)
docker compose up

# Start with rebuild (after code changes)
docker compose up --build

# Stop application
docker compose down

# Stop and remove volumes (fresh start)
docker compose down -v

# View logs
docker compose logs -f

# View specific service logs
docker compose logs -f app
docker compose logs -f frontend
docker compose logs -f db
```

---

## 🔧 Database Commands

```bash
# Access database shell
docker compose exec db psql -U postgres -d fasiolas_game

# List all users
docker compose exec db psql -U postgres -d fasiolas_game -c "SELECT id, username, email, role FROM users;"

# Change user role
docker compose exec db psql -U postgres -d fasiolas_game -c "UPDATE users SET role='admin' WHERE email='user@example.com';"

# Delete user
docker compose exec db psql -U postgres -d fasiolas_game -c "DELETE FROM users WHERE email='user@example.com';"
```

---

## 👨‍💼 Admin Panel

**Access**: http://localhost:3000/admin (admin role required)

### Features:
- View all registered users
- Change user roles (ADMIN/PLAYER/SPECTATOR)
- Delete users (cannot delete yourself)
- View system statistics
- Monitor user count by role

### How to Use:
1. Login as admin
2. Click "Admin Panel" in navigation
3. See user list with roles
4. Click role dropdown to change
5. Click "Delete" to remove user

---

## 🎯 Project Structure

```
cardGame/
├── frontend/                    # React application
│   ├── src/
│   │   ├── pages/
│   │   │   ├── Dashboard.jsx   # Game lobby (role-aware)
│   │   │   ├── Game.jsx        # Game page (spectator mode)
│   │   │   ├── AdminPanel.jsx  # User management
│   │   │   └── Login.jsx       # OAuth login
│   │   ├── components/
│   │   │   ├── GameCard.jsx    # Game card (watch/join)
│   │   │   ├── GameBoard.jsx   # Game board
│   │   │   └── Navigation.jsx  # Top nav bar
│   │   └── App.jsx
│   └── package.json
│
├── internal/                    # Go backend
│   ├── handler/
│   │   ├── auth_handler.go     # OAuth & JWT
│   │   ├── game_handler.go     # Game logic
│   │   └── admin_handler.go    # Admin endpoints
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── game_service.go
│   │   └── admin_service.go
│   ├── middleware/
│   │   └── middleware.go       # Role-based access
│   ├── models/
│   │   └── models.go           # User, Game, Player
│   └── repository/
│       └── repository.go       # Database queries
│
├── migrations/                  # Database migrations
├── docker-compose.yml          # Service orchestration
├── Dockerfile                  # Backend container
└── README.md                   # This file
```

---

## 🔐 Security Features

- ✅ OAuth 2.0 authentication (Google, GitHub, Discord)
- ✅ JWT token-based sessions
- ✅ Role-based access control (RBAC)
- ✅ Protected API endpoints
- ✅ Admin cannot delete own account
- ✅ Spectators cannot modify games
- ✅ CORS configured
- ✅ Input validation

---

## 🐛 Troubleshooting

### Services won't start
```bash
# Check if ports are available
netstat -ano | findstr :3000
netstat -ano | findstr :8080
netstat -ano | findstr :5432

# Restart Docker Desktop
# Then: docker compose down -v
# Then: docker compose up --build
```

### Can't login with OAuth
- Check `.env` files have correct OAuth credentials
- Verify redirect URLs in OAuth provider settings
- Check browser console for errors
- Ensure backend is running on port 8080

### Database connection errors
```bash
# Full reset
docker compose down -v
docker compose up
```

### Frontend shows blank page
```bash
# Clear browser cache (Ctrl+Shift+Delete)
# Hard refresh (Ctrl+Shift+R)
# Or rebuild: docker compose up --build
```

### "Invalid or expired token" error
```javascript
// Clear localStorage in browser console (F12):
localStorage.clear()
// Then refresh page and login again
```

### Port already in use
```bash
# Find process using port 3000
netstat -ano | findstr :3000

# Kill process (replace PID with actual process ID)
taskkill /PID <PID> /F

# Or change port in docker-compose.yml
```

---

## 📦 Technology Stack

**Frontend:**
- React 18
- React Router 6
- Axios
- Tailwind CSS
- React Icons

**Backend:**
- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- JWT tokens
- OAuth2 library

**Database:**
- PostgreSQL 15

**Deployment:**
- Docker
- Docker Compose

**External APIs:**
- Deck of Cards API
- Trivia API (Open Trivia Database)

---

## 🎮 Game Rules

### Phase 1: Accumulation
- Each player gets 1 face-up card
- Player with lowest card starts
- **Before drawing**: Must place top card if possible (+1 rank on another player)
- **Cannot draw** if valid placement exists
- Draw 1 card when no placement possible
- If drawn card fits (+1), must place it
- Phase ends when deck is empty

### Phase 2: Trick-Taking
- **Goal**: Get rid of all cards (last player loses)
- Player with 9♠ (nine of spades) starts
- Must follow suit or play trump
- Highest card wins trick
- Winner leads next trick

---

## 🚀 Deployment to Production

### Environment Setup
1. Create production `.env` files
2. Set production OAuth URLs
3. Configure production database
4. Enable HTTPS

### Deploy Steps
```bash
# Build for production
docker compose -f docker-compose.prod.yml up --build -d

# Check services
docker compose ps

# View logs
docker compose logs -f
```

### Recommended Setup
- Use Nginx/Caddy as reverse proxy
- Enable HTTPS with Let's Encrypt
- Set up database backups
- Configure monitoring (logs, metrics)
- Use environment secrets management

---

## 📝 API Endpoints

### Public (No Auth)
- `GET /health` - Health check
- `GET /api/v1/auth/:provider` - Get OAuth URL
- `GET /api/v1/auth/:provider/callback` - OAuth callback

### Authenticated (All Roles)
- `GET /api/v1/auth/profile` - Get user profile
- `GET /api/v1/games` - List games
- `GET /api/v1/games/:id/spectate` - Watch game (spectators)

### Players & Admins Only
- `POST /api/v1/games` - Create game
- `POST /api/v1/games/join` - Join game
- `POST /api/v1/games/:id/start` - Start game
- `POST /api/v1/games/:id/place` - Place card
- `POST /api/v1/games/:id/draw` - Draw card
- `POST /api/v1/games/:id/skip` - Skip turn

### Admin Only
- `GET /api/v1/admin/users` - List users
- `GET /api/v1/admin/users/:id` - Get user
- `PUT /api/v1/admin/users/:id/role` - Change role
- `DELETE /api/v1/admin/users/:id` - Delete user
- `GET /api/v1/admin/stats` - System stats

---

## 🧪 Testing

### Manual Testing
1. Create 2+ accounts (use different Google accounts or incognito)
2. Create game with first account
3. Join with second account
4. Start game
5. Test gameplay (place cards, draw)
6. Test role features (admin panel, spectator watching)

### API Testing with curl
```bash
# Health check
curl http://localhost:8080/health

# Get games (requires auth token)
curl http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Create game
curl -X POST http://localhost:8080/api/v1/games \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"max_players": 4}'
```

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

---

## 📄 License

MIT License - See LICENSE file for details

---

## 📞 Support

For issues or questions:
- Check troubleshooting section above
- Review Docker logs: `docker compose logs`
- Check browser console (F12)
- Open GitHub issue

---

## ⚡ Quick Reference

| Task | Command |
|------|---------|
| Start | `docker compose up` |
| Stop | `docker compose down` |
| Rebuild | `docker compose up --build` |
| Fresh Start | `docker compose down -v && docker compose up` |
| View Logs | `docker compose logs -f` |
| Access DB | `docker compose exec db psql -U postgres -d fasiolas_game` |

---

**Ready to play! Start with `docker compose up` and visit http://localhost:3000** 🎮

**Created**: January 27, 2026  
**Version**: 1.0 with RBAC
