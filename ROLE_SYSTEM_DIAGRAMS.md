# Role-Based Access Control - Visual Diagrams

## 1. User Role Hierarchy

```
┌─────────────────────────────────────────────────────┐
│                    ALL USERS                         │
│              (Authenticated)                         │
└────────────┬────────────────────┬────────────────────┘
             │                    │
     ┌───────▼────────┐   ┌───────▼────────┐
     │     PLAYER     │   │   SPECTATOR    │
     │                │   │                │
     │ • Create game  │   │ • View games   │
     │ • Join game    │   │ • Watch play   │
     │ • Play cards   │   │ ✗ Play cards   │
     │ • View games   │   │ ✗ Join games   │
     │ ✗ Manage users │   │ ✗ Manage users │
     └────────────────┘   └────────────────┘
             │
     ┌───────▼────────┐
     │     ADMIN      │
     │                │
     │ • All player   │
     │   permissions  │
     │ • Manage users │
     │ • View stats   │
     └────────────────┘
```

## 2. API Endpoint Access Matrix

```
╔════════════════════════════════════════════════════════════════╗
║ Endpoint                          │ Player │ Spectator │ Admin ║
╠════════════════════════════════════════════════════════════════╣
║ GET /api/v1/games                 │   ✅   │    ✅     │  ✅   ║
║ GET /api/v1/games/:id             │   ✅   │    ✅     │  ✅   ║
║ GET /api/v1/games/:id/spectate     │   ❌   │    ✅     │  ✅   ║
║────────────────────────────────────────────────────────────────║
║ POST /api/v1/games                │   ✅   │    ❌     │  ✅   ║
║ POST /api/v1/games/join           │   ✅   │    ❌     │  ✅   ║
║ POST /api/v1/games/:id/start      │   ✅   │    ❌     │  ✅   ║
║ POST /api/v1/games/:id/draw       │   ✅   │    ❌     │  ✅   ║
║ POST /api/v1/games/:id/place      │   ✅   │    ❌     │  ✅   ║
║ POST /api/v1/games/:id/skip       │   ✅   │    ❌     │  ✅   ║
║ POST /api/v1/games/:id/cheat      │   ✅   │    ❌     │  ✅   ║
║────────────────────────────────────────────────────────────────║
║ GET /api/v1/admin/users           │   ❌   │    ❌     │  ✅   ║
║ GET /api/v1/admin/users/:id       │   ❌   │    ❌     │  ✅   ║
║ PUT /api/v1/admin/users/:id/role  │   ❌   │    ❌     │  ✅   ║
║ GET /api/v1/admin/stats           │   ❌   │    ❌     │  ✅   ║
╚════════════════════════════════════════════════════════════════╝
```

## 3. Request Flow - Authentication & Authorization

```
┌─────────────────────────────────────┐
│  User Makes Request                 │
│  GET /api/v1/games/:id/place       │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│  1. AuthMiddleware                  │
│  - Validate JWT token               │
│  - Extract user_id, role            │
│  - Check token not expired          │
│  ❌ If invalid → 401 Unauthorized   │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│  2. RequireRole Middleware          │
│  - Check if role in allowed list    │
│  - RequireRole("player", "admin")   │
│  ❌ If spectator → 403 Forbidden    │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│  3. Handler + Service               │
│  - Additional business logic        │
│  - Verify user in game              │
│  - Check turn order                 │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│  4. Return Response                 │
│  ✅ 200 OK - Place card             │
│  OR                                 │
│  ❌ 400/403 - Error message         │
└─────────────────────────────────────┘
```

## 4. Frontend Component Logic

```
GameBoard Component
│
├─ Read user role from context
│  const isSpectator = user?.role === 'spectator'
│
├─ Disable UI Elements
│  ├─ handlePlaceCard
│  │  └─ if (isSpectator) { setError(...); return; }
│  │
│  ├─ handleDrawCard
│  │  └─ if (isSpectator) { setError(...); return; }
│  │
│  ├─ Deck Click Button
│  │  └─ disabled={isSpectator}
│  │
│  └─ Card Drag/Drop
│     └─ draggable={!isSpectator && isYourTurn}
│
└─ Show Spectator UI
   ├─ Hide action panel
   ├─ Hide hand info
   └─ Show "Spectator Mode" message
```

## 5. Role Transition States

```
              User Signup
              │
              ▼
         (Default Role)
              │
         PLAYER ◄──────────────┐
              │                │
              │                │ (Admin promotes)
              │                │
      ┌───────┴────────┐       │
      │                │       │
      ▼ (Admin changes)│       │
   SPECTATOR           │       │
      │                │       │
      └────────┬───────┘       │
               │               │
               └───────►ADMIN◄─┘
                        │
                   (Only admins here)
                   (Default: player)

Flow:
- New user → PLAYER (default)
- Admin can change to SPECTATOR
- Admin can change to ADMIN
- Users can be downgraded back to PLAYER or SPECTATOR
```

## 6. Database Role Storage

```
╔════════════════════════════════════════════════╗
║ USERS TABLE                                    ║
╠═══════╦═══════════╦═══════════╦════════════════╣
║ id    ║ username  ║ email     ║ role           ║
╠═══════╬═══════════╬═══════════╬════════════════╣
║ 1     ║ john      ║ john@...  ║ player         ║
║ 2     ║ jane      ║ jane@...  ║ spectator      ║
║ 3     ║ admin_usr ║ admin@... ║ admin          ║
║ 4     ║ bob       ║ bob@...   ║ player         ║
╚═══════╩═══════════╩═══════════╩════════════════╝

Values:
  • 'player'     - Can play games (default)
  • 'spectator'  - Can only view games
  • 'admin'      - Full system access
```

## 7. Middleware Stack

```
Request
  │
  ▼
┌─────────────────────────────────────┐
│ CORSMiddleware                      │ ◄─ Global
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ LoggingMiddleware                   │ ◄─ Global
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ AuthMiddleware                      │ ◄─ Protected routes
│ (Validates JWT & sets user_id,      │
│  email, username, role in context)  │
└──────────────┬──────────────────────┘
               │
               ▼
    ┌──────────┴──────────┐
    │                     │
    ▼ (Public routes)     ▼ (Protected routes)
 No role check       RequireRole Middleware
                     │
      ┌──────────────┼──────────────┐
      │              │              │
      ▼              ▼              ▼
  All roles   "player","admin"  "admin"
     │            │              │
  GET /games    /games/*play  /admin/*
  GET /games/:id
  GET /games/:id/spectate
```

## 8. Feature Comparison Table

```
╔═════════════════════════════════════════════════════════════════════╗
║ Feature                          │ Player │ Spectator │ Admin      ║
╠═════════════════════════════════════════════════════════════════════╣
║ View game list                   │   ✅   │    ✅     │    ✅      ║
║ View game state (in game)        │   ✅   │    ✅     │    ✅      ║
║ View game state (spectate)       │   ❌   │    ✅     │    ✅      ║
║                                  │        │           │            ║
║ Create new game                  │   ✅   │    ❌     │    ✅      ║
║ Join game                        │   ✅   │    ❌     │    ✅      ║
║ Start game                       │   ✅   │    ❌     │    ✅      ║
║                                  │        │           │            ║
║ Draw card                        │   ✅   │    ❌     │    ✅      ║
║ Place card                       │   ✅   │    ❌     │    ✅      ║
║ Skip turn                        │   ✅   │    ❌     │    ✅      ║
║ Call cheat                       │   ✅   │    ❌     │    ✅      ║
║                                  │        │           │            ║
║ List users                       │   ❌   │    ❌     │    ✅      ║
║ View user details                │   ❌   │    ❌     │    ✅      ║
║ Change user role                 │   ❌   │    ❌     │    ✅      ║
║ View system stats                │   ❌   │    ❌     │    ✅      ║
╚═════════════════════════════════════════════════════════════════════╝
```

## 9. Error Response Examples

```
Spectator tries to place card:

Request:
  POST /api/v1/games/1/place
  Authorization: Bearer SPECTATOR_TOKEN

Response:
  403 Forbidden
  {
    "error": "insufficient permissions"
  }

---

Spectator tries to view admin panel:

Request:
  GET /api/v1/admin/users
  Authorization: Bearer SPECTATOR_TOKEN

Response:
  403 Forbidden
  {
    "error": "insufficient permissions"
  }

---

Spectator views game (allowed):

Request:
  GET /api/v1/games/1
  Authorization: Bearer SPECTATOR_TOKEN

Response:
  200 OK
  {
    "game": { ... },
    "players": [ ... ],
    "actions": [ ... ],
    "deck_count": 42
  }
```

## 10. Authentication Flow with Roles

```
┌──────────────────────────────────────────────────────────────┐
│                    OAuth Login Flow                          │
└──────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┴────────────┐
                ▼                        ▼
        Google OAuth              Discord OAuth
            │                          │
            └───────────┬──────────────┘
                        │
                        ▼
            ┌──────────────────────┐
            │ Get User Info        │
            │ (email, name, etc.)  │
            └──────────┬───────────┘
                        │
                        ▼
            ┌──────────────────────┐
            │ User Exists?         │
            └──────┬───────┬───────┘
                   │       │
            No     │       │ Yes
                   ▼       ▼
            ┌────────┐  ┌─────────────┐
            │Create  │  │Load existing│
            │User    │  │user from DB │
            │        │  │             │
            │role:   │  │Keep current │
            │player  │  │role         │
            └────┬───┘  └──────┬──────┘
                 │             │
                 └──────┬──────┘
                        │
                        ▼
            ┌──────────────────────┐
            │ Generate JWT Token   │
            │ With Claims:         │
            │ - user_id            │
            │ - email              │
            │ - username           │
            │ - role (player/spec) │
            │ - exp                │
            └──────────┬───────────┘
                        │
                        ▼
            ┌──────────────────────┐
            │ Return Token to      │
            │ Frontend             │
            │ (Store in localStorage)
            └────────────────────────┘
```
