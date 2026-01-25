# Paleisti visą sistemą iš naujo

Write-Host "=== Fasiolas Card Game - Pilnas Paleidimas ===" -ForegroundColor Cyan
Write-Host ""

$projectPath = "c:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
Set-Location $projectPath

# 1. Sustabdyti viską
Write-Host "1. Sustabdome visus konteinerius..." -ForegroundColor Yellow
docker compose down
Write-Host "   OK Sustabdyta" -ForegroundColor Green
Write-Host ""

# 2. Paleisti iš naujo
Write-Host "2. Paleidžiame konteinerius..." -ForegroundColor Yellow
docker compose up -d
Write-Host "   OK Paleista" -ForegroundColor Green
Write-Host ""

# 3. Laukti kol DB pasiruoš
Write-Host "3. Laukiame kol DB pasiruoš (10 sek)..." -ForegroundColor Yellow
Start-Sleep -Seconds 10
Write-Host "   OK Baigta laukti" -ForegroundColor Green
Write-Host ""

# 4. Nukopijuoti migracijas
Write-Host "4. Kopijuojame migracijas..." -ForegroundColor Yellow
docker compose cp migrations/. postgres:/migrations
Write-Host "   OK Nukopijuota" -ForegroundColor Green
Write-Host ""

# 5. Pritaikyti migracijas
Write-Host "5. Pritaikome migracijas..." -ForegroundColor Yellow
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000001_create_users_table.up.sql | Out-Null
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000002_create_games_table.up.sql | Out-Null
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000003_create_game_players_table.up.sql | Out-Null
docker compose exec postgres psql -U postgres -d fasiolas_game -f /migrations/000004_create_game_actions_table.up.sql | Out-Null
Write-Host "   OK Migracijos pritaikytos" -ForegroundColor Green
Write-Host ""

# 6. Perkrauti backend
Write-Host "6. Perkrauname backend..." -ForegroundColor Yellow
docker compose restart app
Start-Sleep -Seconds 3
Write-Host "   OK Backend perkrautas" -ForegroundColor Green
Write-Host ""

# 7. Patikrinti lentelių sukūrimą
Write-Host "7. Tikriname DB lenteles..." -ForegroundColor Yellow
Write-Host ""
docker compose exec postgres psql -U postgres -d fasiolas_game -c "\dt"
Write-Host ""
Write-Host "   OK DB lenteles sukurtos" -ForegroundColor Green
Write-Host ""

# 8. Patikrinti konteinerius
Write-Host "8. Konteinerių būsena:" -ForegroundColor Yellow
Write-Host ""
docker compose ps
Write-Host ""

Write-Host "=== VISKAS PARUOŠTA ===" -ForegroundColor Green
Write-Host ""
Write-Host "Kaip testuoti:" -ForegroundColor Cyan
Write-Host "1. Uždarykite visus naršykles langus" -ForegroundColor White
Write-Host "2. Atidarykite naują INCOGNITO langą (Ctrl+Shift+N)" -ForegroundColor White
Write-Host "3. Eikite į: http://localhost:3000/login" -ForegroundColor White
Write-Host "4. Paspauskite Google Login" -ForegroundColor White
Write-Host "5. Prisijunkite su Google paskyra" -ForegroundColor White
Write-Host "6. Turėtumėte būti nukreipti į dashboard" -ForegroundColor White
Write-Host ""
Write-Host "Jei vis dar matote seną versiją:" -ForegroundColor Yellow
Write-Host "- Perkompiliuokite frontend: .\rebuild-frontend.ps1" -ForegroundColor White
Write-Host ""

