# Quick Testing Guide - OAuth Fix Verification

## ✅ Verification Checklist

### Step 1: Verify Containers Are Running
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose ps
```

**Expected Output:**
```
NAME                 STATUS
fasiolas_postgres    Up (healthy)
fasiolas_app         Up 
fasiolas_frontend    Up
```

### Step 2: Test Backend API
```powershell
$response = curl -s http://localhost:8080/api/v1/auth/google | ConvertFrom-Json
Write-Host "OAuth URL: $($response.url.Substring(0, 100))..."
```

**Expected:** Returns JSON with a long Google OAuth URL

### Step 3: Test Frontend Connection
1. Open browser and go to: `http://localhost:3000/login`
2. You should see the login page with three buttons:
   - Google Login (blue)
   - GitHub Login (gray)
   - Discord Login (purple)

### Step 4: Test OAuth Flow
1. Click the **"Google Login"** button
2. **Expected Behavior:**
   - ✅ No error message
   - ✅ Browser redirects to Google login page
   - ✅ See Google account selection screen

### Step 5: Monitor Logs
While testing, monitor the logs:
```powershell
docker compose logs -f
```

**Look for:**
- Backend receiving the OAuth request
- Frontend successfully calling `/api/v1/auth/google`
- Redirect URL being generated

## 🔧 If Something Goes Wrong

### Issue: "Failed to start google login: No OAuth URL received"
**Solution:**
1. Check backend health: `curl http://localhost:8080/health`
2. Verify axios config: Check `frontend/src/config.js`
3. Rebuild: `docker compose down && docker compose up -d --build`

### Issue: Frontend can't reach backend
**Solution:**
1. Verify frontend environment: `docker exec fasiolas_frontend sh -c 'echo $REACT_APP_API_URL'`
2. Should show: `http://localhost:8080`
3. If not, rebuild frontend: `docker compose up -d --build`

### Issue: CORS errors in browser console
**Solution:**
1. Check CORS_ALLOWED_ORIGINS in `.env`
2. Ensure `http://localhost:3000` is included
3. Restart backend: `docker compose restart fasiolas_app`

## 📋 Configuration Verification

### Check API_BASE_URL is Set
```powershell
docker exec fasiolas_frontend sh -c 'grep REACT_APP_API_URL /app/.env'
```

### Check Backend OAuth Config
```powershell
docker exec fasiolas_app sh -c 'env | grep GOOGLE'
```

Should show:
```
GOOGLE_CLIENT_ID=845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

## 📊 Full System Test Script

```powershell
# Run this complete test
Write-Host "=== Fasiolas OAuth Fix Verification ===" -ForegroundColor Cyan
Write-Host ""

# Test 1: Backend Health
Write-Host "[1/4] Testing Backend Health..." -ForegroundColor Yellow
try {
    $health = curl -s http://localhost:8080/health | ConvertFrom-Json
    Write-Host "✅ Backend: Healthy" -ForegroundColor Green
} catch {
    Write-Host "❌ Backend: Not responding" -ForegroundColor Red
    exit 1
}

# Test 2: OAuth Endpoint
Write-Host "[2/4] Testing OAuth Endpoint..." -ForegroundColor Yellow
try {
    $oauth = curl -s http://localhost:8080/api/v1/auth/google | ConvertFrom-Json
    if ($oauth.url) {
        Write-Host "✅ OAuth: Endpoint working" -ForegroundColor Green
    } else {
        Write-Host "❌ OAuth: No URL in response" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ OAuth: Endpoint failed" -ForegroundColor Red
}

# Test 3: Frontend Running
Write-Host "[3/4] Testing Frontend..." -ForegroundColor Yellow
try {
    $frontend = Invoke-WebRequest -Uri "http://localhost:3000" -TimeoutSec 2 -ErrorAction Stop
    if ($frontend.StatusCode -eq 200) {
        Write-Host "✅ Frontend: Running" -ForegroundColor Green
    }
} catch {
    Write-Host "❌ Frontend: Not responding" -ForegroundColor Red
}

# Test 4: Container Status
Write-Host "[4/4] Checking Container Status..." -ForegroundColor Yellow
$running = docker compose ps -q | Measure-Object | Select-Object -ExpandProperty Count
if ($running -ge 3) {
    Write-Host "✅ Containers: All running ($running)" -ForegroundColor Green
} else {
    Write-Host "⚠️ Containers: Only $running running" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== All Tests Complete ===" -ForegroundColor Cyan
```

## 🎯 Success Criteria

You'll know the fix is working when:

1. ✅ Frontend loads at `http://localhost:3000`
2. ✅ Login page displays with three OAuth buttons
3. ✅ Click "Google Login" → No error message
4. ✅ Browser redirects to Google authentication
5. ✅ No CORS errors in browser console
6. ✅ Network tab shows successful `/api/v1/auth/google` request
7. ✅ Response includes `"url"` field with Google OAuth URL

---

**Fix Deployed:** January 18, 2026
**Status:** Ready for Testing

