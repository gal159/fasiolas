# 🚀 LAUNCH GUIDE - Fasiolas Card Game

**Status:** Docker is not running on your system  
**Date:** January 15, 2026

---

## ⚠️ ISSUE: Docker Desktop Not Running

### What's the Problem?
Docker Desktop needs to be running to start the application with `docker-compose up`.

---

## ✅ SOLUTION 1: Start Docker Desktop (Recommended)

### Step 1: Open Docker Desktop
1. **Windows Start Menu** → Search for "Docker Desktop"
2. Click to open **Docker Desktop**
3. Wait for it to start (you'll see a green indicator)
4. This takes 1-2 minutes

### Step 2: Verify Docker is Running
```powershell
docker --version
docker-compose --version
```

Should output version numbers if working.

### Step 3: Launch the Project
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker-compose up --build
```

Wait for all services to start (2-3 minutes).

### Step 4: Access the Application
- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **Health Check:** curl http://localhost:8080/health

---

## ✅ SOLUTION 2: Local Development (No Docker)

If Docker won't start, run backend and frontend separately:

### Prerequisites
- Go 1.21+ installed
- Node.js 18+ installed
- PostgreSQL running locally (or install)

### Terminal 1 - Start Database (if not running)
```powershell
# Option A: Install PostgreSQL locally
# Download from: https://www.postgresql.org/download/

# Option B: Use Docker just for postgres
docker run -d --name fasiolas_postgres `
  -e POSTGRES_USER=postgres `
  -e POSTGRES_PASSWORD=postgres `
  -e POSTGRES_DB=fasiolas_game `
  -p 5432:5432 `
  postgres:15-alpine
```

### Terminal 1 - Run Migrations
```powershell
# Install migrate tool (one-time)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/fasiolas_game?sslmode=disable" up
```

### Terminal 1 - Start Backend
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
go run cmd/server/main.go
```

Expected output:
```
✓ Database connected successfully
✓ Server starting on :8080
✓ Environment: development
```

### Terminal 2 - Install Frontend Dependencies
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame\frontend"
npm install
```

This takes 2-3 minutes. Wait for it to complete.

### Terminal 2 - Start Frontend
```powershell
npm start
```

Expected output:
```
On Your Network:  http://192.168.x.x:3000
Compiled successfully!
```

### Terminal 3 - Test Backend (Optional)
```powershell
curl http://localhost:8080/health
```

Should return: `{"status":"ok","service":"fasiolas-card-game"}`

### Access Application
- **Frontend:** http://localhost:3000
- **Backend:** http://localhost:8080

---

## ✅ SOLUTION 3: Install Docker Desktop

If Docker isn't installed:

1. **Download:** https://www.docker.com/products/docker-desktop/
2. **Install:** Run the installer and follow prompts
3. **Restart:** Restart your computer
4. **Launch:** Open Docker Desktop from Start Menu
5. **Wait:** Give it 2-3 minutes to fully start
6. **Verify:** Run `docker --version`

Then follow **Solution 1** above.

---

## 🆘 TROUBLESHOOTING

### "npm command not found"
1. Install Node.js: https://nodejs.org/ (LTS version)
2. Restart PowerShell
3. Verify: `node --version`

### "go command not found"
1. Install Go: https://golang.org/doc/install
2. Restart PowerShell
3. Verify: `go version`

### "PostgreSQL connection refused"
1. Option A: Install PostgreSQL locally
2. Option B: Run Docker for just postgres: `docker run -d --name postgres postgres:15-alpine`
3. Update connection string if needed

### "Port 3000 already in use"
```powershell
# Find what's using port 3000
netstat -ano | findstr :3000

# Kill the process (replace PID)
taskkill /PID <PID> /F

# Or use different port
$env:PORT=3001; npm start
```

### "Port 8080 already in use"
```powershell
# Find what's using port 8080
netstat -ano | findstr :8080

# Kill the process (replace PID)
taskkill /PID <PID> /F
```

### "Port 5432 already in use"
PostgreSQL already running. If you don't need it:
```powershell
# Stop PostgreSQL service or remove Docker container
docker stop fasiolas_postgres
```

---

## 📋 QUICK CHECKLIST

Before launching, verify:

- [ ] Docker Desktop installed AND running (for docker-compose method)
  OR
- [ ] Go 1.21+ installed: `go version`
- [ ] Node.js 18+ installed: `node --version`
- [ ] PostgreSQL running (local or Docker)
- [ ] Ports available (3000, 8080, 5432)

---

## 🎯 RECOMMENDED PATH

### Best for Quick Testing
**Solution 1: Docker Desktop Method**
- Simplest setup
- No additional installations
- Everything isolated in containers
- Just one command: `docker-compose up --build`

### Best for Development
**Solution 2: Local Development**
- Faster code changes
- No Docker overhead
- Easier debugging
- Good for learning

---

## 📞 COMMON COMMANDS

### After Services Start

**Health Check Backend**
```powershell
curl http://localhost:8080/health
```

**Health Check Database**
```powershell
# If using Docker postgres
docker exec -it fasiolas_postgres psql -U postgres -d fasiolas_game -c "SELECT 1;"
```

**Stop Everything (Docker)**
```powershell
docker-compose down
```

**Stop Backend (Local)**
```powershell
# Press Ctrl+C in terminal running go run
```

**Stop Frontend (Local)**
```powershell
# Press Ctrl+C in terminal running npm start
```

---

## ✅ YOU'RE READY WHEN

✅ You see "npm start" output showing localhost:3000  
✅ You see "go run" output showing "Server starting on :8080"  
✅ http://localhost:3000 opens in browser showing login page  
✅ http://localhost:8080/health returns JSON

Then you can:
1. Click OAuth login button
2. Create a game
3. Join a game
4. Play Fasiolas!

---

## 📚 DETAILED GUIDES

For complete information, read:
- `FULL_STACK_SETUP.md` - Complete setup guide
- `FRONTEND_GUIDE.md` - Frontend setup
- `DEVELOPER_GUIDE.md` - Development workflow
- `QUICK_REFERENCE.md` - Common commands

---

## 🎉 Next Steps

1. **Choose a solution above** (Docker or Local)
2. **Follow the steps** for your chosen method
3. **Wait for services to start**
4. **Open http://localhost:3000**
5. **Click OAuth button to login**
6. **Create/join game and play!**

---

**Questions? Check the documentation files listed above.**

**Good luck! 🚀**

