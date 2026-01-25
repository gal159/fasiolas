# Restart Docker containers and verify
Write-Host "=== Restarting Docker Containers ===" -ForegroundColor Cyan

# Stop containers
Write-Host "Stopping containers..." -ForegroundColor Yellow
docker compose down

# Start containers with build
Write-Host "Starting containers..." -ForegroundColor Yellow
docker compose up -d --build

# Wait for containers to start
Write-Host "Waiting for containers to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check status
Write-Host "`n=== Container Status ===" -ForegroundColor Cyan
docker compose ps

# Check backend logs
Write-Host "`n=== Backend Logs (last 20 lines) ===" -ForegroundColor Cyan
docker compose logs app --tail 20

Write-Host "`n=== Ready to Test ===" -ForegroundColor Green
Write-Host "1. Go to http://localhost:3000/login" -ForegroundColor White
Write-Host "2. Create a NEW game (important!)" -ForegroundColor White
Write-Host "3. Add 2 players and start" -ForegroundColor White
Write-Host "4. Click 'Draw Card' button" -ForegroundColor White
Write-Host "5. Card should appear in the MIDDLE only" -ForegroundColor White
Write-Host "6. Drag or click to place it" -ForegroundColor White

