# Force restart with cache clear

Write-Host "=== Force Restart ===" -ForegroundColor Cyan

# Stop everything
Write-Host "Stopping processes..." -ForegroundColor Yellow
Get-Process -Name "server" -ErrorAction SilentlyContinue | Stop-Process -Force
Get-Process -Name "node" -ErrorAction SilentlyContinue | Stop-Process -Force
Start-Sleep -Seconds 2

$projectPath = "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"

# Clear frontend cache
Write-Host "Clearing frontend cache..." -ForegroundColor Yellow
Remove-Item "$projectPath\frontend\node_modules\.vite" -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item "$projectPath\frontend\dist" -Recurse -Force -ErrorAction SilentlyContinue

# Start backend
Write-Host "Starting backend..." -ForegroundColor Yellow
Set-Location $projectPath
Start-Process -FilePath ".\server.exe"
Start-Sleep -Seconds 2

# Start frontend
Write-Host "Starting frontend..." -ForegroundColor Yellow
Set-Location "$projectPath\frontend"
Start-Process powershell -ArgumentList "-NoExit", "-NoProfile", "-Command", "cd '$PWD'; npm run dev"

Set-Location $projectPath
Start-Sleep -Seconds 3

Write-Host "`n✓ Servers restarted with clean cache!" -ForegroundColor Green
Write-Host "`nWait 10 seconds, then:" -ForegroundColor Cyan
Write-Host "1. Open NEW incognito window (Ctrl+Shift+N)" -ForegroundColor White
Write-Host "2. Go to: http://localhost:3000" -ForegroundColor White
Write-Host "3. Login and check game room`n" -ForegroundColor White

