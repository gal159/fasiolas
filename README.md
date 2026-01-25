# 🎴 Fasiolas Card Game - REST API

A comprehensive REST API for the Lithuanian card game "Fasiolas" built with Go, featuring OAuth2 authentication, PostgreSQL database, and external API integration.

## ⚡ Quick Start

**Paleisti visą sistemą:**
```powershell
.\start-all.ps1
```

**Jei matote seną frontend versiją:**
```powershell
.\rebuild-frontend.ps1
```

**Po to testuokite:**
1. Atidarykite Incognito langą (`Ctrl+Shift+N`)
2. Eikite į `http://localhost:3000/login`
3. Paspauskite "Google Login"

📖 **Daugiau informacijos:** [GREITAS_START.md](GREITAS_START.md) | [KAIP_PALEISTI.md](KAIP_PALEISTI.md)

---

## 📋 Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Game Rules](#game-rules)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Authentication](#authentication)
- [Role-Based Access Control](#role-based-access-control)
- [External API Integration](#external-api-integration)
- [Project Structure](#project-structure)
- [Development](#development)
- [Testing](#testing)

## ✨ Features

- **Full CRUD REST API** for card game management
- **OAuth2 Authentication** with 3 providers (Google, GitHub, Discord)
- **JWT Token-based** session management
- **Role-Based Access Control** (Admin, Player, Spectator)
- **PostgreSQL Database** with migrations
- **External API Integration** (Deck of Cards API)
- **Real-time game state** management
- **Comprehensive game rules** implementation
- **Action logging** and game history
- **Docker support** for easy deployment

## 🛠 Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL 15
- **Authentication**: OAuth2 (golang.org/x/oauth2)
- **JWT**: golang-jwt/jwt
- **Database Migrations**: golang-migrate
- **External API**: Deck of Cards API
- **Containerization**: Docker & Docker Compose

## 🎮 Game Rules

### Fasiolas - Two-Phase Card Game

**Players**: 2-8 players using a standard 52-card deck

### Phase 1: Accumulation (+1 Phase)

**Goal**: Accumulate cards on other players' piles following the +1 rank rule.

**Rules**:
1. Each player starts with 1 face-up card
2. Player with the lowest card starts
3. **Before drawing**, must place top card if possible (+1 rank on another player)
4. **Cannot draw** if a valid placement exists (this is cheating!)
5. Draw 1 card when no placement is possible
6. If drawn card fits (+1), must place it
7. If it doesn't fit, add to own pile
8. **Cheating penalty**: All players with >1 card give 1 card to cheater
9. Wrong cheat accusation = lose ability to call cheats
10. Phase ends when deck is empty

### Phase 2: Trick-Taking Phase

**Goal**: Get rid of all cards. Last player with cards **loses**.

**Rules**:
1. **Trump suit** = last card drawn in Phase 1 (not ♠)
2. Player with **9♠** (nine of spades) starts
3. **Playing cards**:
   - Empty table: any card
   - Same suit + higher rank, OR
   - Trump (beats anything except spades if trump isn't spades)
4. **Cannot play**: Take oldest card from table, add to bottom of pile, **skip turn**
5. Round ends when all played → last player starts next round
6. Game ends when only 1 player has cards (they lose)

## 📦 Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker & Docker Compose (optional)
- Git

## 🚀 Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd cardGame
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Install development tools

```bash
# Install golang-migrate for database migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install swag for API documentation generation
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. Setup PostgreSQL

#### Option A: Using Docker

```bash
docker-compose up -d postgres
```

#### Option B: Local PostgreSQL

Create a database:
```sql
CREATE DATABASE fasiolas_game;
```

### 5. Run database migrations

```bash
# Using Makefile
make migrate-up

# Or directly
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

## ⚙️ Configuration

### 1. Create environment file

```bash
cp .env.example .env
```

### 2. Configure OAuth2 providers

Edit `.env` and add your OAuth credentials:

#### Google OAuth2
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a project and enable Google+ API
3. Create OAuth 2.0 credentials
4. Add authorized redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
5. Copy Client ID and Client Secret to `.env`

#### GitHub OAuth2
1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Create a new OAuth App
3. Set Authorization callback URL: `http://localhost:8080/api/v1/auth/github/callback`
4. Copy Client ID and Client Secret to `.env`

#### Discord OAuth2
1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Add redirect URI in OAuth2 settings: `http://localhost:8080/api/v1/auth/discord/callback`
4. Copy Client ID and Client Secret to `.env`

### 3. Configure JWT Secret

```env
JWT_SECRET=your-very-secure-random-secret-key-here
```

**Generate a secure secret:**
```bash
openssl rand -base64 32
```

## 🏃 Running the Application

### Option 1: Using Go directly

```bash
# Run the server
go run cmd/server/main.go

# Or build and run
make build
./bin/server
```

### Option 2: Using Docker Compose

```bash
# Start all services (PostgreSQL + App)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Option 3: Using Makefile

```bash
# Run application
make run

# Build binary
make build

# Run tests
make test

# Generate Swagger docs
make swagger
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### Get OAuth URL
```http
GET /api/v1/auth/{provider}
```
Providers: `google`, `github`, `discord`

**Response:**
```json
{
  "url": "https://accounts.google.com/o/oauth2/auth?..."
}
```

#### OAuth Callback
```http
GET /api/v1/auth/{provider}/callback?code={code}&state={state}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "John Doe",
    "role": "player"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Get Profile
```http
GET /api/v1/auth/profile
Authorization: Bearer {token}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Authorization: Bearer {token}
```

### Game Endpoints

All game endpoints require authentication and player/admin role.

#### Create Game
```http
POST /api/v1/games
Authorization: Bearer {token}
Content-Type: application/json

{
  "max_players": 4
}
```

#### Join Game
```http
POST /api/v1/games/join
Authorization: Bearer {token}
Content-Type: application/json

{
  "room_code": "ABC123"
}
```

#### Start Game
```http
POST /api/v1/games/{id}/start
Authorization: Bearer {token}
```

#### Get Game State
```http
GET /api/v1/games/{id}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "game": {
    "id": 1,
    "room_code": "ABC123",
    "state": "phase1",
    "phase": 1,
    "current_player_position": 0,
    "deck_cards": [...],
    "table_cards": []
  },
  "players": [
    {
      "id": 1,
      "position": 0,
      "card_count": 3,
      "top_card": {"suit": "hearts", "rank": "A", "value": 14},
      "user": {"username": "Player 1"}
    }
  ],
  "recent_actions": [...]
}
```

#### Place Card
```http
POST /api/v1/games/{id}/place
Authorization: Bearer {token}
Content-Type: application/json

{
  "target_player_position": 1
}
```

#### Draw Card
```http
POST /api/v1/games/{id}/draw
Authorization: Bearer {token}
```

#### Call Cheat
```http
POST /api/v1/games/{id}/cheat?cheater_id=2
Authorization: Bearer {token}
```

#### List Games
```http
GET /api/v1/games?state=waiting&limit=20&offset=0
Authorization: Bearer {token}
```

### External API Endpoints

#### Get Card Image
```http
GET /api/v1/external/card-image?suit=hearts&rank=A
Authorization: Bearer {token}
```

#### Get Game Statistics
```http
GET /api/v1/external/stats
Authorization: Bearer {token}
```

## 🗄️ Database Schema

### Users Table
```sql
- id: SERIAL PRIMARY KEY
- email: VARCHAR(255) UNIQUE NOT NULL
- username: VARCHAR(100) UNIQUE NOT NULL
- oauth_provider: VARCHAR(50) NOT NULL
- oauth_id: VARCHAR(255) NOT NULL
- role: VARCHAR(20) NOT NULL (admin, player, spectator)
- avatar_url: VARCHAR(500)
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
- last_login: TIMESTAMP
```

### Games Table
```sql
- id: SERIAL PRIMARY KEY
- room_code: VARCHAR(10) UNIQUE NOT NULL
- state: VARCHAR(20) (waiting, phase1, phase2, finished, cancelled)
- phase: INTEGER (1, 2)
- current_player_position: INTEGER
- trump_suit: VARCHAR(10)
- deck_cards: JSONB
- table_cards: JSONB
- created_by: INTEGER REFERENCES users(id)
- max_players: INTEGER
- created_at, updated_at, started_at, finished_at: TIMESTAMP
- winner_id: INTEGER REFERENCES users(id)
```

### Game Players Table
```sql
- id: SERIAL PRIMARY KEY
- game_id: INTEGER REFERENCES games(id)
- user_id: INTEGER REFERENCES users(id)
- position: INTEGER
- cards: JSONB
- top_card: JSONB
- card_count: INTEGER
- status: VARCHAR(20) (active, eliminated, winner, left)
- can_call_cheat: BOOLEAN
- joined_at: TIMESTAMP
```

### Game Actions Table
```sql
- id: SERIAL PRIMARY KEY
- game_id: INTEGER REFERENCES games(id)
- user_id: INTEGER REFERENCES users(id)
- action_type: VARCHAR(50)
- action_data: JSONB
- phase: INTEGER
- timestamp: TIMESTAMP
```

## 🔐 Authentication

### OAuth2 Flow

1. Client requests auth URL for provider
2. User is redirected to OAuth provider
3. User authorizes the application
4. Provider redirects back with authorization code
5. Server exchanges code for access token
6. Server retrieves user info from provider
7. Server creates/updates user in database
8. Server generates JWT token
9. Client uses JWT for subsequent requests

### JWT Token Structure

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "username": "John Doe",
  "role": "player",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Using JWT Tokens

Include the token in the Authorization header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 👥 Role-Based Access Control

### Roles

1. **Admin** (`admin`)
   - Full access to all endpoints
   - Can manage all games
   - Can delete games
   - Can modify user roles

2. **Player** (`player`)
   - Can create and join games
   - Can play in games
   - Can view game state
   - Default role for new users

3. **Spectator** (`spectator`)
   - Can view games
   - Cannot create or join games
   - Read-only access

### Role Middleware

Protected endpoints automatically check user roles:
```go
games.Use(middleware.RequireRole("player", "admin"))
```

## 🌐 External API Integration

### Deck of Cards API

The application integrates with the [Deck of Cards API](https://deckofcardsapi.com/) to provide:

- Card image URLs
- SVG/PNG card representations
- Deck shuffling (alternative implementation)

**Example Usage:**
```bash
curl -H "Authorization: Bearer {token}" \
  "http://localhost:8080/api/v1/external/card-image?suit=hearts&rank=A"
```

**Response:**
```json
{
  "code": "AH",
  "image": "https://deckofcardsapi.com/static/img/AH.png",
  "images": {
    "svg": "https://deckofcardsapi.com/static/img/AH.svg",
    "png": "https://deckofcardsapi.com/static/img/AH.png"
  },
  "value": "ACE",
  "suit": "HEARTS"
}
```

## 📁 Project Structure

```
cardGame/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   └── models.go            # Data models
│   ├── repository/
│   │   ├── db.go                # Database connection
│   │   ├── user_repository.go
│   │   ├── game_repository.go
│   │   ├── game_player_repository.go
│   │   └── game_action_repository.go
│   ├── service/
│   │   ├── auth_service.go      # OAuth2 & JWT logic
│   │   ├── game_service.go      # Game business logic
│   │   └── external_api_service.go
│   ├── handler/
│   │   ├── auth_handler.go      # HTTP handlers for auth
│   │   └── game_handler.go      # HTTP handlers for games
│   ├── middleware/
│   │   └── middleware.go        # Auth, RBAC, CORS middleware
│   └── game/
│       └── engine.go            # Core game rules engine
├── pkg/
│   ├── cards/
│   │   └── cards.go             # Card deck implementation
│   └── utils/
│       └── utils.go             # Utility functions
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000002_create_games_table.up.sql
│   ├── 000003_create_game_players_table.up.sql
│   └── 000004_create_game_actions_table.up.sql
├── .env.example                 # Environment variables template
├── .gitignore
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker image definition
├── Makefile                     # Build automation
├── go.mod                       # Go module dependencies
├── go.sum
└── README.md                    # This file
```

## 🧪 Testing

### Run Unit Tests
```bash
make test
```

### Manual API Testing

Use tools like:
- **Postman**: Import API endpoints
- **curl**: Command-line testing
- **Insomnia**: REST client

### Example Test Flow

1. **Get OAuth URL**
```bash
curl http://localhost:8080/api/v1/auth/google
```

2. **Complete OAuth flow in browser**

3. **Use returned token**
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get profile
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/profile

# Create game
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"max_players": 4}' \
  http://localhost:8080/api/v1/games
```

## 🐛 Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker-compose ps

# View PostgreSQL logs
docker-compose logs postgres

# Recreate database
docker-compose down -v
docker-compose up -d postgres
make migrate-up
```

### OAuth Issues
- Verify OAuth credentials in `.env`
- Check redirect URLs match exactly
- Ensure OAuth apps are enabled in provider consoles

### Port Already in Use
```bash
# Change port in .env
PORT=8081

# Or kill process using port 8080 (Windows)
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

## 📝 License

This project is licensed under the MIT License.

## 👤 Author

Created as a final project for Server-Side Web Development course.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📞 Support

For issues and questions:
- Create an issue on GitHub
- Contact: support@fasiolas.com

---

**Happy Gaming! 🎴**

