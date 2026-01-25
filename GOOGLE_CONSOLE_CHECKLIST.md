# 🔍 Google OAuth Configuration Checklist

## ⚠️ "Invalid State Parameter" Error - Google Console Checklist

If you're still getting `{"error":"invalid state parameter"}`, let's verify your Google Cloud Console configuration.

---

## ✅ Step-by-Step Google Console Setup

### 1. Go to Google Cloud Console
```
https://console.cloud.google.com/apis/credentials?project=fasiolas-auth
```

### 2. Check OAuth 2.0 Client IDs

Click on your OAuth Client ID (should show your Client ID: `845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com`)

#### ✅ Required Settings:

**Application type:** Web application

**Authorized JavaScript origins:**
```
http://localhost:3000
http://localhost:8080
```

**Authorized redirect URIs:** (THIS IS CRITICAL!)
```
http://localhost:8080/api/v1/auth/google/callback
```

⚠️ **MUST BE EXACT** - No trailing slash, exact path

---

### 3. OAuth Consent Screen Configuration

Go to: OAuth consent screen (left sidebar)

#### ✅ Required Settings:

**User Type:** External

**App name:** Fasiolas Card Game (or your choice)

**User support email:** Your email

**Developer contact information:** Your email

**Scopes:** 
- `.../auth/userinfo.email`
- `.../auth/userinfo.profile`  
- `openid`

**Test users:** ADD YOUR GOOGLE ACCOUNT EMAIL HERE

⚠️ **CRITICAL:** Your Google account MUST be added as a test user!

---

### 4. Enable Required APIs

Go to: Library (left sidebar)

Search and Enable:
- ✅ **Google+ API** (or People API)
- ✅ **Google OAuth2 API**

---

## 🔍 Common Google Console Issues

### Issue #1: Redirect URI Mismatch
**Symptom:** redirect_uri_mismatch error or state parameter fails
**Fix:** 
1. Go to Credentials
2. Edit your OAuth Client
3. Add EXACTLY: `http://localhost:8080/api/v1/auth/google/callback`
4. Click SAVE
5. Wait 5 minutes for propagation

### Issue #2: Test User Not Added
**Symptom:** "Access blocked: This app's request is invalid"
**Fix:**
1. Go to OAuth consent screen
2. Scroll to "Test users"
3. Click "ADD USERS"
4. Add your Google account email
5. Save

### Issue #3: Wrong Consent Screen Type
**Symptom:** Can't add test users or access denied
**Fix:**
1. OAuth consent screen should be "External"
2. Publishing status: "Testing"
3. Add your email as test user

### Issue #4: APIs Not Enabled
**Symptom:** User info fetch fails after OAuth
**Fix:**
1. Go to Library
2. Search "Google+ API" or "People API"
3. Click ENABLE
4. Search "OAuth2 API"  
5. Click ENABLE

---

## 🧪 Quick Test Checklist

Run through this checklist:

- [ ] Redirect URI is **EXACTLY**: `http://localhost:8080/api/v1/auth/google/callback`
- [ ] No extra slashes, no typos
- [ ] JavaScript origins include `http://localhost:3000`
- [ ] OAuth consent screen is "External"
- [ ] Your email is added as a test user
- [ ] Google+ API (or People API) is enabled
- [ ] Client ID matches your `.env` file
- [ ] Client Secret matches your `.env` file

---

## 📝 Your Current Configuration

From your `.env` file:
```
GOOGLE_CLIENT_ID=845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

**Verify these match your Google Console!**

---

## 🔧 How to Verify Settings

### Check Redirect URI:
1. Go to https://console.cloud.google.com/apis/credentials
2. Click your OAuth 2.0 Client ID
3. Look under "Authorized redirect URIs"
4. Should see: `http://localhost:8080/api/v1/auth/google/callback`

### Check Test Users:
1. Go to OAuth consent screen
2. Scroll down to "Test users"
3. Your Google account email should be listed

### Check APIs:
1. Go to Library
2. Search "People API" - should show "MANAGE" (means enabled)
3. Search "OAuth2" - should show "MANAGE" (means enabled)

---

## 🚨 Most Common Mistake

**THE #1 CAUSE OF "Invalid State Parameter":**

The redirect URI in Google Console doesn't match EXACTLY what your backend is using.

**Must be:** `http://localhost:8080/api/v1/auth/google/callback`

**Common wrong versions:**
- ❌ `http://localhost:8080/api/v1/auth/google/callback/` (extra slash)
- ❌ `http://localhost:8080/auth/google/callback` (missing /api/v1)
- ❌ `https://localhost:8080/api/v1/auth/google/callback` (https instead of http)
- ❌ `http://127.0.0.1:8080/api/v1/auth/google/callback` (127.0.0.1 instead of localhost)

---

## 🔄 After Making Changes

1. Save changes in Google Console
2. **Wait 5 minutes** for Google to propagate changes
3. Clear browser cookies
4. Try OAuth flow again

---

## 📞 Debug Commands

### Test OAuth URL Generation:
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google"
```

Should return URL with your Client ID

### Check Backend Logs:
```powershell
docker logs fasiolas_app -f
```

Look for state validation messages

### Check Google Console Project:
```
Project ID: fasiolas-auth
Project Number: 845779433699
```

---

## ✅ Final Verification Steps

### Step 1: Google Console
- [ ] Open https://console.cloud.google.com/apis/credentials?project=fasiolas-auth
- [ ] Click your OAuth Client ID
- [ ] **Take a screenshot** of "Authorized redirect URIs"
- [ ] Verify it shows: `http://localhost:8080/api/v1/auth/google/callback`

### Step 2: Test Users
- [ ] Go to OAuth consent screen
- [ ] Scroll to "Test users"
- [ ] **Take a screenshot**
- [ ] Verify your Google email is listed

### Step 3: Try Again
- [ ] Go to http://localhost:3000/login
- [ ] Click "Google Login"
- [ ] Should work now!

---

## 🆘 Still Not Working?

### Share These Screenshots:
1. OAuth Client settings (Authorized redirect URIs section)
2. OAuth consent screen (Test users section)
3. Browser error (if any)
4. Backend logs showing state parameter

### Alternative: Recreate OAuth Client
If all else fails:
1. Create a NEW OAuth 2.0 Client ID
2. Copy new Client ID and Secret to `.env`
3. Add redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
4. Add yourself as test user
5. Restart backend: `docker compose restart fasiolas_app`
6. Try again

---

**Next Action:** Verify your Google Console settings match this checklist!


