# Google Cloud Console Configuration Guide

## 🎯 Quick Reference Card

**Your Google Project:** `fasiolas-auth`  
**Client ID:** `845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com`  
**Redirect URI:** `http://localhost:8080/api/v1/auth/google/callback`

---

## 📋 Step-by-Step Configuration

### Step 1: Access Google Cloud Console

1. Open your browser and go to:
   ```
   https://console.cloud.google.com/
   ```

2. Sign in with your Google account

3. Select project: **fasiolas-auth**
   - If you don't see it, click the project dropdown at the top
   - Search for "fasiolas-auth"

---

### Step 2: Navigate to Credentials

1. In the left sidebar, click **"APIs & Services"**
2. Click **"Credentials"**
3. You should see your OAuth 2.0 Client ID:
   ```
   845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com
   ```

**Direct Link:**
```
https://console.cloud.google.com/apis/credentials?project=fasiolas-auth
```

---

### Step 3: Configure OAuth Client ID

1. Click on your OAuth 2.0 Client ID (the one starting with 845779433699...)

2. Scroll down to **"Authorized redirect URIs"**

3. Click **"ADD URI"**

4. Paste this exact URI:
   ```
   http://localhost:8080/api/v1/auth/google/callback
   ```

5. **IMPORTANT:** 
   - ✅ Use `http://` (not `https://`) for localhost
   - ✅ No trailing slash
   - ✅ Port must be `8080`
   - ✅ Path must be `/api/v1/auth/google/callback`

6. Click **"SAVE"** at the bottom

7. ⏰ Wait 1-2 minutes for changes to propagate

---

### Step 4: Configure OAuth Consent Screen

1. In the left sidebar, click **"OAuth consent screen"**

2. If not already configured, select:
   - **User Type:** External
   - Click **"CREATE"**

3. Fill in **App information:**
   ```
   App name:              Fasiolas Card Game
   User support email:    [Your email]
   Developer contact:     [Your email]
   ```

4. Click **"SAVE AND CONTINUE"**

5. On the **Scopes** page:
   - Click **"ADD OR REMOVE SCOPES"**
   - Select these scopes:
     - ✅ `.../auth/userinfo.email`
     - ✅ `.../auth/userinfo.profile`
     - ✅ `openid`
   - Click **"UPDATE"**
   - Click **"SAVE AND CONTINUE"**

6. On the **Test users** page (if app is in Testing mode):
   - Click **"ADD USERS"**
   - Add your Google email address
   - Click **"ADD"**
   - Click **"SAVE AND CONTINUE"**

7. Review and click **"BACK TO DASHBOARD"**

---

### Step 5: Enable Required APIs

1. In the left sidebar, click **"Library"**

2. Search for and enable these APIs:
   - **Google+ API** (or **People API**)
   - **Google Identity Toolkit API**

3. For each API:
   - Click on it
   - Click **"ENABLE"**
   - Wait for activation

**Direct Link:**
```
https://console.cloud.google.com/apis/library?project=fasiolas-auth
```

---

### Step 6: Verify Configuration

Go back to **Credentials** and verify:

```
✅ OAuth 2.0 Client IDs section shows your client
✅ Authorized redirect URIs includes:
   http://localhost:8080/api/v1/auth/google/callback
✅ OAuth consent screen shows "External" type
✅ Required APIs are enabled
✅ Test users are added (if in Testing mode)
```

---

## 🧪 Test Your Configuration

After completing the Google Console setup:

1. **Run the test script:**
   ```powershell
   ./test-google-oauth.ps1
   ```

2. **Copy the OAuth URL from the output**

3. **Open the URL in your browser:**
   - You should see the Google sign-in page
   - Sign in with your Google account
   - Click "Allow" to grant permissions
   - You should be redirected back with a JWT token

4. **Expected redirect:**
   ```
   http://localhost:8080/api/v1/auth/google/callback?code=4/...&state=...
   ```

5. **Expected response (JSON):**
   ```json
   {
     "user": {
       "id": 1,
       "email": "your@email.com",
       "username": "Your Name",
       "role": "player",
       "avatar_url": "https://..."
     },
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
   }
   ```

---

## ⚠️ Troubleshooting

### Error: "redirect_uri_mismatch"

**Problem:** The redirect URI doesn't match what's configured in Google Console.

**Solution:**
1. Go to Google Cloud Console → Credentials
2. Click on your OAuth Client ID
3. Verify the redirect URI is exactly:
   ```
   http://localhost:8080/api/v1/auth/google/callback
   ```
4. No trailing slash, no extra spaces
5. Save and wait 1-2 minutes

---

### Error: "Access blocked: This app's request is invalid"

**Problem:** OAuth consent screen is not properly configured or you're not a test user.

**Solution:**
1. Complete OAuth consent screen setup (Step 4 above)
2. Add your email as a test user
3. Or publish your app (if ready)

---

### Error: "invalid_client"

**Problem:** Client ID or Client Secret is wrong.

**Solution:**
1. Download the latest credentials JSON from Google Console
2. Update `.env` file with correct values
3. Restart the server:
   ```bash
   docker-compose restart
   ```

---

### Error: "Admin policy" or "Access denied"

**Problem:** Your organization has restricted OAuth apps.

**Solution:**
1. Use a personal Google account for testing
2. Or contact your Google Workspace admin
3. Or create a new Google Cloud project with a personal account

---

## 📸 What You Should See

### In Google Cloud Console - Credentials Page:

```
OAuth 2.0 Client IDs
┌─────────────────────────────────────────────────────────┐
│ Name: Web client 1                                      │
│ Client ID: 845779433699-i7hb2rk5hk...                   │
│ Creation date: ...                                      │
│                                                         │
│ Authorized redirect URIs:                               │
│ • http://localhost:8080/api/v1/auth/google/callback    │
│                                                         │
│ [SAVE] [CANCEL]                                         │
└─────────────────────────────────────────────────────────┘
```

### In Google Cloud Console - OAuth Consent Screen:

```
Publishing status: Testing
User type: External

App information:
• App name: Fasiolas Card Game
• User support email: your@email.com

Scopes:
• .../auth/userinfo.email
• .../auth/userinfo.profile  
• openid

Test users:
• your@email.com
```

---

## ✅ Configuration Checklist

Use this checklist to ensure everything is set up:

- [ ] Logged into Google Cloud Console
- [ ] Selected project: fasiolas-auth
- [ ] Found OAuth 2.0 Client ID (845779433699-...)
- [ ] Added redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
- [ ] Saved changes
- [ ] Configured OAuth consent screen (External)
- [ ] Added required scopes (email, profile, openid)
- [ ] Added test user (your email)
- [ ] Enabled Google+ API or People API
- [ ] Waited 1-2 minutes for changes to propagate
- [ ] Ran test script: `./test-google-oauth.ps1`
- [ ] Tested OAuth flow in browser
- [ ] Received JWT token successfully

---

## 📞 Need Help?

**Documentation Files:**
- `GOOGLE_OAUTH_STATUS.md` - Current configuration status
- `GOOGLE_OAUTH_CHECKLIST.md` - Detailed checklist
- `test-google-oauth.ps1` - Automated test script

**Check Server Logs:**
```bash
docker logs fasiolas_app
```

**Test Endpoint Directly:**
```bash
curl http://localhost:8080/api/v1/auth/google
```

**Google OAuth Documentation:**
- https://developers.google.com/identity/protocols/oauth2
- https://support.google.com/cloud/answer/6158849

---

**Last Updated:** January 15, 2026  
**Status:** Ready for configuration

