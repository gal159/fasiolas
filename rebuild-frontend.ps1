# Perkompiliuoti ir paleisti frontend su nauju build'u

Write-Host "=== Frontend Rebuild Script ===" -ForegroundColor Cyan
Write-Host ""

$projectPath = "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"

Write-Host "1. Sustabdome frontend..." -ForegroundColor Yellow
Set-Location $projectPath
docker compose stop frontend

Write-Host ""
Write-Host "2. Perkompiluojame frontend be cache..." -ForegroundColor Yellow
docker compose build --no-cache frontend

Write-Host ""
Write-Host "3. Paleidžiame frontend..." -ForegroundColor Yellow
docker compose up -d frontend

Write-Host ""
Write-Host "4. Laukiame 5 sekundes..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

Write-Host ""
Write-Host "5. Tikriname konteinerių būseną..." -ForegroundColor Yellow
docker compose ps

Write-Host ""
Write-Host "=== ATLIKTA ===" -ForegroundColor Green
Write-Host ""
Write-Host "Dabar:" -ForegroundColor Cyan
Write-Host "1. Uždarykite visus naršykles langus" -ForegroundColor White
Write-Host "2. Atidarykite naują INCOGNITO langą (Ctrl+Shift+N)" -ForegroundColor White
Write-Host "3. Eikite į: http://localhost:3000/login" -ForegroundColor White
Write-Host "4. Turėtumėte matyti login puslapį su Google Login mygtuku" -ForegroundColor White
Write-Host ""

