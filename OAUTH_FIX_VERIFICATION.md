# 🧪 OAuth State Parameter Fix - Verification Plan

## ✅ What Was Fixed

### **Fix #1: State Storage in Memory**
- Added `sync.Mutex` for thread-safe access
- State is NOW stored in `h.stateStore` when user requests OAuth URL
- Debug log: `[OAuth] State generated and stored in memory: {state}`

### **Fix #2: State Validation with Fallback**
- Strategy 1: Check in-memory store (primary)
- Strategy 2: Check cookie (fallback)
- Debug logs show which strategy validates the state

### **Fix #3: Proper Error Handling**
- Now validates state parameter before processing code
- Returns `"invalid state parameter"` if validation fails
- Clear error messages in logs

### **Fix #4: Thread Safety**
- Added `sync.Mutex` to protect concurrent access to state store
- Multiple requests won't corrupt state data

---

## 🔍 Verification Steps

### Step 1: Start Docker Containers
```powershell
cd "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
docker compose up -d
```

Wait for all containers to start:
- PostgreSQL on port 5432
- Backend on port 8080
- Frontend on port 3000

### Step 2: Check Logs for Startup Errors
```powershell
docker logs fasiolas_app -f
```

Expected startup output:
```
=== Fasiolas Card Game Server Starting ===
✓ Config loaded
✓ Database connection successful
✓ Server listening on :8080
```

### Step 3: Test OAuth URL Generation
```powershell
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google"
$response | ConvertTo-Json

# Expected output:
# {
#   "url": "https://accounts.google.com/o/oauth2/auth?..."
# }
```

Check backend logs for:
```
[OAuth] State generated and stored in memory: {state_value}
```

### Step 4: Manual Callback Test (Without Real Google)
```powershell
# Simulate callback with WRONG state
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google/callback?state=WRONG_STATE&code=TEST_CODE" -ErrorAction SilentlyContinue

# Expected output:
# {
#   "error": "invalid state parameter"
# }
```

Check logs for:
```
[OAuth] Callback received - Provider: google, State: WRONG_STATE, Code: TEST_CODE
[OAuth] ❌ VALIDATION FAILED - Invalid state parameter
```

### Step 5: Test Full OAuth Flow (Real Google)
1. Open `http://localhost:3000/login`
2. Click "Google Login" button
3. Sign in with your Google account
4. Watch backend logs in real-time

Expected log sequence:
```
[OAuth] State generated and stored in memory: {state}
[OAuth] Callback received - Provider: google, State: {state}, Code: {code}
[OAuth] ✅ State validated from memory store
[OAuth] State validated successfully, exchanging code for token
[OAuth] ✅ Token generated successfully, redirecting to frontend
```

Expected frontend behavior:
- Redirected to Google login
- After signing in, redirected back to `http://localhost:3000/auth/callback?token=JWT_TOKEN`
- AuthCallback page shows "Completing login..."
- Redirected to dashboard
- ✅ LOGGED IN

---

## 📋 Troubleshooting

### Issue: Still Getting "invalid state parameter"
**Check these in order:**

1. **Is backend logging state generation?**
   ```powershell
   docker logs fasiolas_app | grep "State generated"
   ```
   If NO: State not being stored, fix GetAuthURL

2. **Is backend logging callback receipt?**
   ```powershell
   docker logs fasiolas_app | grep "Callback received"
   ```
   If NO: Callback not reaching backend

3. **Is state validation succeeding?**
   ```powershell
   docker logs fasiolas_app | grep "State validated"
   ```
   If showing "VALIDATION FAILED": State mismatch, check strategies

4. **Check cookie in browser**
   - Open DevTools (F12)
   - Go to Application → Cookies
   - Look for `oauth_state` cookie
   - Should exist for `localhost` domain

---

## 🧬 Code Changes Summary

### File: `internal/handler/auth_handler.go`

#### Change 1: Added Mutex
```go
type AuthHandler struct {
    authService *service.AuthService
    stateStore  map[string]bool
    stateMutex  sync.Mutex      // ← NEW
}
```

#### Change 2: Store State in Memory
```go
func (h *AuthHandler) GetAuthURL(c *gin.Context) {
    state, _ := utils.GenerateSecureToken(32)
    
    h.stateMutex.Lock()
    h.stateStore[state] = true  // ← NEW: Store in memory
    h.stateMutex.Unlock()
    
    fmt.Printf("[OAuth] State generated and stored in memory: %s\n", state)  // ← NEW: Debug log
    
    c.SetCookie("oauth_state", state, 300, "/", "localhost", false, true)
    // ... rest
}
```

#### Change 3: Validate State with Fallback
```go
func (h *AuthHandler) HandleCallback(c *gin.Context) {
    // ... get params ...
    
    // ← NEW: Strategy 1 - Check memory store
    h.stateMutex.Lock()
    if h.stateStore[state] {
        fmt.Printf("[OAuth] ✅ State validated from memory store\n")
        delete(h.stateStore, state)
        isStateValid = true
    }
    h.stateMutex.Unlock()
    
    // ← NEW: Strategy 2 - Fallback to cookie
    if !isStateValid {
        cookieState, _ := c.Cookie("oauth_state")
        if cookieState == state {
            fmt.Printf("[OAuth] ✅ State validated from cookie\n")
            isStateValid = true
        }
    }
    
    // ← NEW: If still invalid, return error
    if !isStateValid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state parameter"})
        return
    }
    
    // ... continue with callback ...
}
```

---

## ✨ Expected Behavior After Fix

| Scenario | Before | After |
|----------|--------|-------|
| OAuth URL request | ❌ State only in cookie | ✅ State in memory + cookie |
| Callback with valid state | ❌ No validation | ✅ Validated against memory/cookie |
| Callback with invalid state | ❌ Ignored, might error later | ✅ Immediate "invalid state parameter" |
| Log visibility | ❌ Silent failures | ✅ Clear debug messages |
| Thread safety | ❌ Potential race conditions | ✅ Protected with Mutex |

---

## 🚀 Next Steps After Verification

1. ✅ Test local OAuth flow completely
2. ✅ Verify logs show correct state validation
3. ✅ Check that user is created in database
4. ✅ Confirm JWT token is generated
5. ✅ Test multiple users to ensure state store doesn't accumulate
6. ✅ Test error scenarios (wrong state, missing code)

---

## 📊 State Flow Diagram

```
1. User clicks "Google Login"
   ↓
2. Frontend: GET /api/v1/auth/google
   ↓
3. Backend GetAuthURL:
   ├─ Generate random state: "abc123def..."
   ├─ Store in memory: stateStore["abc123def..."] = true
   ├─ Store in cookie: Set-Cookie: oauth_state=abc123def...
   └─ Return OAuth URL with state parameter
   ↓
4. Frontend redirects to Google with state
   ↓
5. User authenticates with Google
   ↓
6. Google redirects to callback: /api/v1/auth/google/callback?state=abc123def...&code=...
   ↓
7. Backend HandleCallback:
   ├─ Get state from URL: "abc123def..."
   ├─ Check memory store: ✅ Found!
   ├─ Delete from memory: stateStore["abc123def..."] deleted
   ├─ Exchange code for token
   ├─ Get user info
   ├─ Create/update user in DB
   ├─ Generate JWT
   └─ Redirect with JWT: /auth/callback?token=JWT_TOKEN
   ↓
8. Frontend AuthCallback:
   ├─ Extract token from URL
   ├─ Store in localStorage
   ├─ Fetch user profile
   └─ Redirect to dashboard
   ↓
9. ✅ User is logged in!
```

---

## 🔐 Security Improvements

✅ **CSRF Protection**: State token prevents cross-site request forgery
✅ **One-time Use**: State deleted after validation (prevents replay attacks)
✅ **Thread Safety**: Mutex prevents concurrent modification issues
✅ **Fallback Strategy**: Cookie backup for edge cases
✅ **Clear Errors**: Invalid states rejected immediately


