#!/usr/bin/env pwsh

Write-Host @"
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║      🎮 FASIOLAS CARD GAME - QUICK RESTART & TEST 🎮      ║
║                                                            ║
╔════════════════════════════════════════════════════════════╗
"@ -ForegroundColor Cyan

Write-Host "`n📋 This script will:" -ForegroundColor Yellow
Write-Host "   1. Stop all containers"
Write-Host "   2. Rebuild backend with fixes"
Write-Host "   3. Start all containers"
Write-Host "   4. Show you how to test"
Write-Host ""

Read-Host "Press ENTER to continue"

Write-Host "`n⏹️  Stopping containers..." -ForegroundColor Yellow
docker compose down
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to stop containers" -ForegroundColor Red
    exit 1
}

Write-Host "`n🔨 Building and starting containers..." -ForegroundColor Yellow
docker compose up -d --build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to start containers" -ForegroundColor Red
    exit 1
}

Write-Host "`n⏳ Waiting for containers to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 8

Write-Host "`n📊 Container Status:" -ForegroundColor Cyan
docker compose ps

Write-Host "`n📜 Backend Logs (last 15 lines):" -ForegroundColor Cyan
docker compose logs app --tail 15

Write-Host @"

╔════════════════════════════════════════════════════════════╗
║                    ✅ READY TO TEST! ✅                    ║
╚════════════════════════════════════════════════════════════╝

🎯 TESTING INSTRUCTIONS:

1️⃣  Open browser: http://localhost:3000/login

2️⃣  Create a NEW game:
   - Click "Create Game"
   - Max players: 2
   - Click "Create"

3️⃣  Add second player:
   - Open INCOGNITO window
   - Login with different Google account
   - Enter room code and join

4️⃣  Start the game:
   - First player clicks "START GAME"

5️⃣  Test drawing:
   - Wait for your turn (green border)
   - Click "🎴 Draw Card" button

   ✅ EXPECTED:
   - Card appears in MIDDLE (with yellow border)
   - Your pile shows "(placing above)"
   - Deck count decreases

   ❌ NOT EXPECTED:
   - Card in BOTH middle AND your pile
   - "Deck is empty!" message

6️⃣  Test placement:
   - Drag card to a player OR click on player

   ✅ EXPECTED:
   - Card moves to their pile
   - Turn ends automatically
   - Next player's turn

╔════════════════════════════════════════════════════════════╗
║                 🐛 TROUBLESHOOTING 🐛                      ║
╚════════════════════════════════════════════════════════════╝

If card doesn't appear:
- Open browser console (F12)
- Look for errors in Console tab
- Check Network tab for /draw request

If card appears in both places:
- Check console logs
- Look for: "refPlacing: true"
- Try hard refresh: Ctrl+Shift+R

If placement fails:
- Check backend logs: docker compose logs app --tail 50
- Verify target position is sent correctly

╔════════════════════════════════════════════════════════════╗

"@ -ForegroundColor Green

Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

