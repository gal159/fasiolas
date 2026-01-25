# Google OAuth Configuration Checklist

## ✅ Current Configuration Status

### 1. Credentials Found
- **Client ID**: `845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com`
- **Client Secret**: `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`
- **Project ID**: `fasiolas-auth`

### 2. Environment Variables (.env)
✅ GOOGLE_CLIENT_ID is set correctly
✅ GOOGLE_CLIENT_SECRET is set correctly
✅ GOOGLE_REDIRECT_URL is set to: `http://localhost:8080/api/v1/auth/google/callback`

### 3. Code Configuration (internal/service/auth_service.go)
✅ OAuth2 config using Google endpoint
✅ Scopes configured:
   - `https://www.googleapis.com/auth/userinfo.email`
   - `https://www.googleapis.com/auth/userinfo.profile`

### 4. API Routes (cmd/server/main.go)
✅ Route: `GET /api/v1/auth/google` - Get OAuth URL
✅ Route: `GET /api/v1/auth/google/callback` - Handle callback

## 🔧 Required Google Cloud Console Configuration

### Step 1: Access Google Cloud Console
1. Go to: https://console.cloud.google.com/
2. Select project: **fasiolas-auth**

### Step 2: Configure OAuth Consent Screen
1. Navigate to: **APIs & Services** → **OAuth consent screen**
2. Ensure the following settings:
   - **User Type**: External (for testing with any Google account)
   - **App name**: Fasiolas Card Game
   - **User support email**: Your email
   - **Developer contact information**: Your email
   - **Scopes**: Add these scopes:
     - `.../auth/userinfo.email`
     - `.../auth/userinfo.profile`
     - `openid`

### Step 3: Configure Authorized Redirect URIs
1. Navigate to: **APIs & Services** → **Credentials**
2. Click on the OAuth 2.0 Client ID: `845779433699...`
3. Under **Authorized redirect URIs**, add:
   ```
   http://localhost:8080/api/v1/auth/google/callback
   ```
4. Click **Save**

### Step 4: Configure Authorized JavaScript Origins (Optional but recommended)
Add these origins:
```
http://localhost:8080
http://localhost:3000
```

### Step 5: Enable Required APIs
Ensure these APIs are enabled:
1. Navigate to: **APIs & Services** → **Library**
2. Search and enable:
   - **Google+ API** (or People API)
   - **OAuth2 API**

### Step 6: Test User Configuration (If app is in Testing mode)
1. Navigate to: **OAuth consent screen** → **Test users**
2. Add your Google account email as a test user
3. Or publish the app if ready

## 🧪 Testing Steps

### Test 1: Get OAuth URL
```bash
curl http://localhost:8080/api/v1/auth/google
```

Expected response:
```json
{
  "url": "https://accounts.google.com/o/oauth2/auth?client_id=845779433699...&redirect_uri=http://localhost:8080/api/v1/auth/google/callback&response_type=code&scope=https://www.googleapis.com/auth/userinfo.email+https://www.googleapis.com/auth/userinfo.profile&state=..."
}
```

### Test 2: Complete OAuth Flow
1. Copy the URL from Test 1 response
2. Open it in your browser
3. Sign in with Google
4. Authorize the app
5. You should be redirected to: `http://localhost:8080/api/v1/auth/google/callback?code=...&state=...`
6. Check the response - should contain user info and JWT token

### Test 3: Use JWT Token
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/v1/auth/profile
```

Expected response:
```json
{
  "id": 1,
  "email": "your@email.com",
  "username": "Your Name",
  "role": "player"
}
```

## 🚀 Quick Start Commands

### Start the backend server:
```bash
# Using Docker Compose
docker-compose up -d

# Or build and run manually
go build -o bin/server cmd/server/main.go
./bin/server
```

### Check server health:
```bash
curl http://localhost:8080/health
```

### Test OAuth flow with Postman:
1. Import: `Fasiolas-API.postman_collection.json`
2. Use "Get Google OAuth URL" request
3. Open the returned URL in browser
4. Complete authentication
5. Save the returned JWT token in Postman variables

## ❗ Common Issues and Solutions

### Issue 1: "redirect_uri_mismatch" error
**Solution**: 
- Verify the redirect URI in Google Console matches exactly: `http://localhost:8080/api/v1/auth/google/callback`
- No trailing slashes
- Check for http vs https

### Issue 2: "Access blocked: This app's request is invalid"
**Solution**:
- Enable Google+ API or People API in Google Cloud Console
- Verify OAuth consent screen is properly configured
- Add your email as a test user if app is in testing mode

### Issue 3: "invalid_client" error
**Solution**:
- Double-check CLIENT_ID and CLIENT_SECRET in .env file
- Ensure credentials are not expired
- Verify you're using the correct credentials from Google Console

### Issue 4: "state parameter mismatch"
**Solution**:
- Clear browser cookies
- The state parameter is stored in cookies for CSRF protection
- Try the flow in an incognito window

### Issue 5: "unsupported_grant_type" or token exchange fails
**Solution**:
- Ensure the authorization code is used immediately (expires quickly)
- Check that the oauth2 library is properly configured
- Verify network connectivity to Google's token endpoint

## 📝 Verification Checklist

Before testing, verify:

- [ ] Google Cloud project exists (fasiolas-auth)
- [ ] OAuth 2.0 Client ID is created
- [ ] Redirect URI is added in Google Console
- [ ] OAuth consent screen is configured
- [ ] Required APIs are enabled
- [ ] .env file has correct CLIENT_ID and CLIENT_SECRET
- [ ] Server is running on port 8080
- [ ] Database is running and connected
- [ ] Health check endpoint responds: http://localhost:8080/health

## 🔒 Security Notes

1. **Never commit credentials to Git**
   - .env file should be in .gitignore
   - Use environment variables in production

2. **Change JWT_SECRET in production**
   - Current secret is for development only

3. **Use HTTPS in production**
   - Update redirect URLs to use https://
   - Configure SSL certificates

4. **Restrict CORS origins**
   - Update CORS_ALLOWED_ORIGINS for production domains

5. **Review OAuth scopes**
   - Only request minimum required permissions
   - Current scopes: email, profile (appropriate for this app)

