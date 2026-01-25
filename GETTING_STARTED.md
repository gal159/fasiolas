# 🎯 GETTING STARTED - 5 Minute Quick Start

**Read this if you want to test the application RIGHT NOW!**

---

## ⚡ The Fastest Way to Get Running

### 1️⃣ Prerequisites (2 min)
Ensure you have:
- ✅ Docker Desktop installed and running
- ✅ Go 1.21+ installed (`go version`)

### 2️⃣ Start Database (1 min)

Open PowerShell in the project folder:
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"

# Start PostgreSQL in Docker
docker-compose up -d postgres

# Wait for database to be ready (watch logs)
docker-compose logs postgres
# Look for: "database system is ready to accept connections"
```

### 3️⃣ Setup Database (1 min)

Install migrations tool (first time only):
```powershell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations:
```powershell
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

### 4️⃣ Start Application (30 sec)

```powershell
go run cmd/server/main.go
```

You should see:
```
✓ Database connected successfully
✓ Server starting on :8080
✓ Environment: development
```

### 5️⃣ Test It! (30 sec)

Open a new terminal window:

```powershell
# Test health check
curl http://localhost:8080/health

# You should get:
# {"status":"ok","service":"fasiolas-card-game"}
```

✅ **SUCCESS!** Server is running!

---

## 🎮 Quick Game Test

### Get Google OAuth URL

```powershell
curl http://localhost:8080/api/v1/auth/google
```

Expected response:
```json
{
  "url": "https://accounts.google.com/o/oauth2/v2/auth?..."
}
```

> **Note:** If you see an error about `GOOGLE_CLIENT_ID`, that's expected. You need to configure OAuth credentials (see section below).

---

## 🔐 (Optional) Setup OAuth for Full Testing

### Google OAuth Setup

1. Go to https://console.cloud.google.com/
2. Create a new project
3. Enable "Google+ API"
4. Go to Credentials → Create OAuth 2.0 Client ID (Web application)
5. Add authorized redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
6. Copy your Client ID and Client Secret

### Update .env File

Edit `.env` in your project root:

```env
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here
```

### Restart Application

```powershell
# Stop current app: Ctrl+C
# Restart:
go run cmd/server/main.go
```

Now you can test OAuth flow! ✅

---

## 📚 Important Files to Know

| File | Purpose |
|------|---------|
| `README.md` | Project overview and features |
| `SETUP_GUIDE.md` | Detailed setup instructions |
| `API_TESTING.md` | How to test all API endpoints |
| `GAME_RULES.md` | Game rules explained |
| `DEVELOPER_GUIDE.md` | Development workflow |
| `IMPLEMENTATION_STATUS.md` | What's done, what's next |
| `COMPLETION_CHECKLIST.md` | Full project checklist |

---

## 🆘 Common Issues

### "Could not find docker"
→ Install Docker Desktop https://www.docker.com/products/docker-desktop/

### "Port 5432 already in use"
→ Another PostgreSQL is running. Stop it or use different port in docker-compose.yml

### "Database connection refused"
→ PostgreSQL not ready yet. Check logs: `docker-compose logs postgres`

### "Migration failed"
→ Database already has tables. Run: `docker-compose down -v` to reset

### "GOOGLE_CLIENT_ID not configured"
→ That's okay for testing! Not all OAuth tests require configured credentials.

---

## 🎯 What to Test

Once running, test these:

```powershell
# 1. Health check
curl http://localhost:8080/health

# 2. Get OAuth URL
curl http://localhost:8080/api/v1/auth/google

# 3. Check database
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "SELECT COUNT(*) FROM users;"
```

---

## 📖 Next Steps

After confirming the app runs:

1. **Read `API_TESTING.md`** - Learn how to test all endpoints
2. **Check `GAME_RULES.md`** - Understand how Fasiolas works
3. **Review `IMPLEMENTATION_STATUS.md`** - See what's implemented
4. **Follow `DEVELOPER_GUIDE.md`** - Start developing

---

## 🚀 Stop & Cleanup

### Stop Application
```powershell
# In the terminal running go run: Ctrl+C
```

### Stop Database
```powershell
docker-compose down
```

### Full Reset (Delete all data)
```powershell
docker-compose down -v
docker-compose up -d postgres
# Then re-run migrations
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

---

## ✅ You're Ready!

The application is now:
- ✅ Built and running
- ✅ Connected to database
- ✅ Ready for API testing
- ✅ Ready for development

**Next:** See `API_TESTING.md` for detailed endpoint testing!

---

**Questions?** Check the appropriate documentation file above. Each has detailed explanations!

