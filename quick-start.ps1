# Greitas paleidimas - jei konteineriai jau sukurti

Write-Host "=== Greitas Paleidimas ===" -ForegroundColor Cyan
Write-Host ""

$projectPath = "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
Set-Location $projectPath

Write-Host "Paleidžiame konteinerius..." -ForegroundColor Yellow
docker compose up -d

Write-Host ""
Write-Host "Laukiame 5 sekundes..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

Write-Host ""
Write-Host "Būsena:" -ForegroundColor Yellow
docker compose ps

Write-Host ""
Write-Host "=== PALEISTA ===" -ForegroundColor Green
Write-Host ""
Write-Host "Testuokite: http://localhost:3000/login" -ForegroundColor Cyan
Write-Host ""

