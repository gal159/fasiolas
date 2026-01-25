# Start development servers (without Docker)

Write-Host "=== Fasiolas Card Game - Dev Mode Start ===" -ForegroundColor Cyan
Write-Host ""

$projectPath = "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
Set-Location $projectPath

# 1. Stop any existing processes
Write-Host "1. Stopping existing processes..." -ForegroundColor Yellow
Get-Process -Name "server" -ErrorAction SilentlyContinue | Stop-Process -Force
Get-Process -Name "node" -ErrorAction SilentlyContinue | Stop-Process -Force
Start-Sleep -Seconds 2
Write-Host "   ✓ Stopped" -ForegroundColor Green
Write-Host ""

# 2. Start Backend Server
Write-Host "2. Starting Backend Server..." -ForegroundColor Yellow
if (Test-Path ".\server.exe") {
    Start-Process -FilePath ".\server.exe" -WindowStyle Normal
    Write-Host "   ✓ Backend started on http://localhost:8080" -ForegroundColor Green
} else {
    Write-Host "   ✗ server.exe not found. Run: go build -o server.exe .\cmd\server" -ForegroundColor Red
}
Write-Host ""

# 3. Start Frontend Dev Server
Write-Host "3. Starting Frontend Dev Server..." -ForegroundColor Yellow
Start-Sleep -Seconds 2
Set-Location "$projectPath\frontend"
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PWD'; Write-Host 'Frontend Dev Server' -ForegroundColor Cyan; npm run dev"
Set-Location $projectPath
Write-Host "   ✓ Frontend dev server starting..." -ForegroundColor Green
Write-Host ""

Start-Sleep -Seconds 3

Write-Host "=== SERVERS STARTED ===" -ForegroundColor Green
Write-Host ""
Write-Host "Services:" -ForegroundColor Cyan
Write-Host "  Frontend: http://localhost:3000" -ForegroundColor White
Write-Host "  Backend:  http://localhost:8080" -ForegroundColor White
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  1. Wait 5 seconds for frontend to compile" -ForegroundColor White
Write-Host "  2. Open browser: http://localhost:3000" -ForegroundColor White
Write-Host "  3. Press F12 to open Developer Console" -ForegroundColor White
Write-Host "  4. Go to game room - you should see debug logs!" -ForegroundColor White
Write-Host ""

