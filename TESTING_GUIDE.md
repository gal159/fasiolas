# 🧪 OAuth State Parameter Fix - Testing Guide

## Quick Start: Test the Fix Immediately

### Step 1: Restart Docker Containers
```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose down -v
docker compose up -d
```

Wait for containers to start:
```
✔ Container fasiolas_postgres      Healthy
✔ Container fasiolas_app           Started
✔ Container fasiolas_frontend      Started
```

### Step 2: Test 1 - OAuth URL Generation
```powershell
# Request OAuth URL
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google"
$response.url

# Expected: URL starting with https://accounts.google.com/o/oauth2/auth
```

**Check Backend Logs:**
```powershell
docker compose logs app | Select-String "State generated"

# Expected Output:
# [OAuth] State generated and stored in memory: {random_32_byte_string}
```

---

### Step 3: Test 2 - Browser OAuth Flow (Real Test)

1. **Open browser:**
   ```
   http://localhost:3000/login
   ```

2. **Click "Google Login" button**

3. **Sign in with your Google account**
   - Enter email
   - Enter password
   - Grant permission if prompted

4. **Check what happens:**
   - ✅ You should be redirected to dashboard
   - ✅ You should see your profile info
   - ✅ You should be logged in

5. **Check Backend Logs:**
   ```powershell
   docker compose logs app --tail=50
   ```

   **Expected Log Sequence:**
   ```
   [OAuth] State generated and stored in memory: {state}
   [OAuth] Callback received - Provider: google, State: {state}, Code: {code}
   [OAuth] ✅ State validated from memory store
   [OAuth] State validated successfully, exchanging code for token
   [OAuth] ✅ Token generated successfully, redirecting to frontend
   ```

---

### Step 4: Test 3 - Error Handling (Invalid State)
```powershell
# Try callback with WRONG state
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google/callback?state=WRONG_STATE&code=FAKE_CODE" -ErrorAction SilentlyContinue

$response | ConvertTo-Json

# Expected Output:
# {
#   "error": "invalid state parameter"
# }
```

**Check Backend Logs:**
```powershell
docker compose logs app --tail=10

# Expected:
# [OAuth] Callback received - Provider: google, State: WRONG_STATE, Code: FAKE_CODE
# [OAuth] ❌ VALIDATION FAILED - Invalid state parameter
```

---

### Step 5: Test 4 - Error Handling (Missing State)
```powershell
# Try callback without state
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google/callback?code=FAKE_CODE" -ErrorAction SilentlyContinue

$response | ConvertTo-Json

# Expected Output:
# {
#   "error": "missing state parameter"
# }
```

---

### Step 6: Test 5 - Error Handling (Missing Code)
```powershell
# Try callback without code
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google/callback?state=SOME_STATE" -ErrorAction SilentlyContinue

$response | ConvertTo-Json

# Expected Output:
# {
#   "error": "missing authorization code"
# }
```

---

## 📊 Test Results Summary

Create this checklist as you test:

- [ ] **Test 1:** OAuth URL generated successfully
  - Expected: URL with state parameter
  - Log shows: `[OAuth] State generated and stored in memory: ...`
  
- [ ] **Test 2:** Real Google login works
  - Expected: Logged in on dashboard
  - Log shows: `[OAuth] ✅ State validated from memory store`
  
- [ ] **Test 3:** Invalid state rejected
  - Expected: `{"error":"invalid state parameter"}`
  - Log shows: `[OAuth] ❌ VALIDATION FAILED`
  
- [ ] **Test 4:** Missing state rejected
  - Expected: `{"error":"missing state parameter"}`
  - Log shows: `[OAuth] ❌ Missing state parameter`
  
- [ ] **Test 5:** Missing code rejected
  - Expected: `{"error":"missing authorization code"}`
  - Log shows: `[OAuth] ❌ Missing authorization code`

---

## 🔍 Live Log Monitoring

To see real-time logs while testing:

```powershell
# Terminal 1: Watch backend logs
docker compose logs app -f

# Terminal 2: Test the API
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google"
```

---

## 🎯 Expected Log Output Flow

When a user completes real Google OAuth, you should see:

```
[OAuth] Callback received - Provider: google, State: abc123def456..., Code: 4/0ASc3gC1Tn...
[OAuth] ✅ State validated from memory store
[OAuth] State validated successfully, exchanging code for token
2026/01/18 17:45:12 User created/updated: id=1, email=user@gmail.com
[OAuth] ✅ Token generated successfully, redirecting to frontend
```

---

## ⚠️ Troubleshooting

### Problem: Backend container not starting
```powershell
docker compose logs app
```
**Solution:** Check logs, ensure PostgreSQL is healthy

### Problem: "Connection refused" error
```powershell
docker ps
```
**Solution:** Ensure all containers are running. If not, rebuild:
```powershell
docker compose down -v
docker compose build --no-cache
docker compose up -d
```

### Problem: OAuth URL returns 404
**Solution:** Backend is not running. Check logs and restart

### Problem: Still getting "invalid state parameter" error
**Cause:** Multiple possible reasons

**Debug Steps:**
```powershell
# 1. Check if state is being stored
docker compose logs app | Select-String "State generated"

# 2. Check if callback is being received
docker compose logs app | Select-String "Callback received"

# 3. Check if state validation is working
docker compose logs app | Select-String "State validated"

# 4. Check browser cookies
# Press F12 → Application → Cookies → localhost
# Look for: oauth_state cookie
```

---

## 📈 Success Indicators

✅ **You'll know it's working when:**

1. **OAuth URL Generation**
   - Backend returns a valid OAuth URL
   - Log shows state stored in memory

2. **Google Authentication**
   - Google login works without errors
   - You can sign in with your Google account

3. **Callback Processing**
   - Backend validates state successfully
   - Token is generated
   - You're redirected to dashboard

4. **Error Handling**
   - Invalid states are rejected immediately
   - Clear error messages returned
   - No silent failures

5. **Security**
   - State changes with each login attempt
   - Can't reuse old state values
   - CSRF protection is active

---

## 🧬 Technical Verification

### Check Code Implementation

**1. State Storage in Memory:**
```powershell
# Look for this line in backend logs
"[OAuth] State generated and stored in memory:"

# If missing, the memory storage isn't working
```

**2. State Validation:**
```powershell
# Look for this line after callback
"[OAuth] ✅ State validated from"

# Shows which strategy worked (memory or cookie)
```

**3. Error Handling:**
```powershell
# When testing invalid state
"[OAuth] ❌ VALIDATION FAILED"

# Shows proper error handling
```

---

## 🚀 Production Checklist

Before deploying to production:

- [ ] All 5 test scenarios pass
- [ ] No errors in backend logs
- [ ] No race condition issues (test concurrent logins)
- [ ] State is cleaned up after use (prevent memory leak)
- [ ] Database user creation works
- [ ] JWT tokens are generated correctly
- [ ] Frontend receives token and stores it
- [ ] User profile is fetched successfully
- [ ] Dashboard loads after login

---

## 📊 Performance Notes

The fix includes:
- ✅ **Memory Store**: O(1) lookup time for state validation
- ✅ **Mutex**: Minimal lock contention (microseconds)
- ✅ **One-time Use**: Automatic cleanup after validation
- ✅ **No Database Calls**: State validation is in-memory (fast)

**Expected Performance:**
- OAuth URL generation: < 1ms
- State validation: < 1μs
- Token exchange: < 500ms (includes Google API call)

---

## 🎓 Understanding the Fix

The fix works in 3 phases:

### Phase 1: State Generation (GetAuthURL)
```
User clicks "Google Login"
↓
Backend generates random 32-byte token
↓
State stored in: h.stateStore[state] = true
State stored in: cookie "oauth_state"
↓
OAuth URL returned with state parameter
```

### Phase 2: User Authentication (Google)
```
Frontend redirects to Google
↓
User signs in with Google
↓
Google verifies authorization
↓
Google redirects back with code + state
```

### Phase 3: State Validation (HandleCallback)
```
Backend receives callback with code + state
↓
Check 1: Is state in memory store? 
  - YES: ✅ Validation succeeds
  - NO: Check Strategy 2
↓
Check 2: Is state in cookie?
  - YES: ✅ Validation succeeds
  - NO: ❌ Validation fails, return error
↓
If validated: Exchange code for token
If not: Return error immediately
```

---

## 📞 Support

If tests fail:

1. **Check logs:** `docker compose logs app`
2. **Restart containers:** `docker compose down -v && docker compose up -d`
3. **Verify compilation:** The binary was built successfully
4. **Check environment:** All OAuth credentials set in .env


