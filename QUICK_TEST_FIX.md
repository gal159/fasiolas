# Quick Start - Test the "Failed to load games" Fix

## What Was Fixed?

When no games exist, the API now returns `[]` (empty array) instead of `null`. The frontend also now has better error handling and detailed logging.

## Step 1: Rebuild Backend

Open PowerShell and run:
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
go build -o server.exe ./cmd/server
.\server.exe
```

**Expected Output:** Server starts on port 8080

## Step 2: Rebuild Frontend

Open another PowerShell tab and run:
```powershell
cd "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame\frontend"
npm install
npm start
```

**Expected Output:** React app starts on http://localhost:3000

## Step 3: Test the Fix

1. Open browser to `http://localhost:3000`
2. Click "Login with Google"
3. Sign in with your Google account
4. You should be redirected to Dashboard
5. **Check Console (F12 → Console tab):**
   - Should see: `📡 Fetching games from /api/v1/games?state=waiting&limit=20`
   - Should see: `✅ Games response received: 200 []`
   - NO red error messages

## Step 4: Verify UI

After successful login, you should see:
- ✅ Dashboard with "Games" header
- ✅ "No games available" message (if no games exist)
- ✅ "Create the first game!" button
- ✅ **NOT** the red "Failed to load games" error

## Troubleshooting

### Still seeing "Failed to load games"?

Check browser console (F12) for:
- **What status code?** 200? 401? 500?
- **What response data?** Check the ✅ log
- **Authorization header?** Should have `Bearer {token}`

Check server logs for:
- Any error messages
- Database connection issues
- Authorization middleware errors

### Database Issues?

Make sure you have PostgreSQL running and migrations applied:
```powershell
# Check if database exists and has tables
psql -h localhost -U postgres -d fasiolas_game -c "SELECT COUNT(*) FROM games;"
```

## Files Changed

- ✅ `internal/handler/game_handler.go` - Backend now returns `[]` instead of `null`
- ✅ `frontend/src/pages/Dashboard.jsx` - Frontend has better error handling & logging

See `FIX_FAILED_TO_LOAD_GAMES.md` for detailed technical explanation.
