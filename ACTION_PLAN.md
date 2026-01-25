# ⚡ IMMEDIATE ACTION - Get Your Project Running

## The Problem (FIXED NOW)
Your backend was crashing because it tried to connect to the database too quickly before PostgreSQL was ready.

## The Solution (APPLIED)
- ✅ Added automatic retry logic to database connections
- ✅ Improved PostgreSQL startup checks
- ✅ Backend now waits patiently for database

## DO THIS RIGHT NOW

### Step 1: Stop Everything
```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose down -v
```

Wait for it to complete.

### Step 2: Rebuild with New Code
```powershell
docker compose up -d --build
```

This will rebuild the backend with the new database retry logic.

### Step 3: Wait 20-30 Seconds
Let the containers fully initialize. This is important!

### Step 4: Verify Backend is Running
```powershell
curl http://localhost:8080/health
```

**You should see:**
```json
{"status":"ok","service":"fasiolas-card-game"}
```

If you see this, ✅ **DATABASE CONNECTION IS WORKING**

### Step 5: Open the Frontend
Open in your browser:
```
http://localhost:3000/login
```

You should see the login page with three OAuth buttons.

### Step 6: Test Google Login
Click the **"Google Login"** button and complete authentication.

---

## If Something Goes Wrong

### Issue: curl returns nothing / Backend not responding

**Solution:**
```powershell
# Check container status
docker compose ps

# View backend logs
docker logs fasiolas_app --tail 100

# If logs show it's still connecting, wait 10 more seconds
```

### Issue: Database connection timeout

**Solution:**
```powershell
# Check if PostgreSQL is running
docker logs fasiolas_postgres --tail 50

# If PostgreSQL logs show errors, restart everything
docker compose down -v
docker compose up -d --build
# Wait 30 seconds
```

### Issue: Frontend loads but OAuth doesn't work

**Solution:**
```powershell
# Make sure backend is running
curl http://localhost:8080/api/v1/auth/google

# Check browser console for errors
# Open DevTools (F12) and check Console tab
```

---

## Quick Commands Reference

```powershell
# View all logs
docker compose logs -f

# View just backend logs
docker compose logs fasiolas_app -f

# View just database logs
docker compose logs fasiolas_postgres -f

# View just frontend logs
docker compose logs fasiolas_frontend -f

# Stop everything
docker compose down

# Stop and remove all data (fresh start)
docker compose down -v

# Rebuild and restart
docker compose down -v
docker compose up -d --build

# Check health
curl http://localhost:8080/health
```

---

## Expected Timeline

| Time | Event | Status |
|------|-------|--------|
| 0s | You run `docker compose up -d --build` | ⏳ Building images |
| 5s | PostgreSQL container starts | 🔄 Initializing database |
| 10-15s | Backend container starts | 🔄 Connecting to database |
| 15-20s | Backend connects to database | ✅ Connection successful |
| 20-25s | All services running | ✅ Ready |
| 25-30s | You can start testing | 🎮 Go play! |

---

## The Fix in Plain English

**Before:**
- Backend starts
- Backend immediately tries to connect to database
- Database isn't ready yet
- Connection fails
- Backend crashes ❌

**After:**
- Backend starts
- Backend tries to connect to database
- If not ready, backend waits 1 second
- Backend tries again (repeats up to 30 times)
- Database becomes ready
- Connection succeeds
- Backend starts serving requests ✅

---

## Your Next Steps

1. **Run:** `docker compose down -v`
2. **Run:** `docker compose up -d --build`
3. **Wait:** 20-30 seconds
4. **Test:** `curl http://localhost:8080/health`
5. **Open:** http://localhost:3000/login
6. **Click:** "Google Login"
7. **Play:** Enjoy Fasiolas! 🎮

---

## Estimated Time: 3 minutes

- 1 minute: Run commands
- 2 minutes: Wait for containers
- Then: Your app is live!

---

**Status: ✅ Ready to Launch**
**Next Action: Follow the steps above**
**Expected Result: Working application**

🚀 **Let's go!**

