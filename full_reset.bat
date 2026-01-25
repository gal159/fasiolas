@echo off
REM Full reset with migrations
REM This script will reset everything and start fresh

setlocal enabledelayedexpansion

echo.
echo ╔════════════════════════════════════════════╗
echo ║  FASIOLAS CARD GAME - FULL RESET           ║
echo ╚════════════════════════════════════════════╝
echo.

REM Step 1: Stop all processes
echo [1/6] Stopping all processes...
taskkill /F /IM node.exe >nul 2>&1
taskkill /F /IM npm.exe >nul 2>&1
taskkill /IM server.exe >nul 2>&1
timeout /t 2 >nul

echo ✓ Processes stopped

REM Step 2: Run migrations
echo.
echo [2/6] Running database migrations...
echo.

REM Check if PostgreSQL is running
echo Checking PostgreSQL connection...
cd "%~dp0"

REM Create database if it doesn't exist
psql -U postgres -c "CREATE DATABASE IF NOT EXISTS fasiolas_game;" >nul 2>&1

REM Run migrations one by one
psql -U postgres -d fasiolas_game -f migrations\000001_create_users_table.up.sql >nul 2>&1
if !errorlevel! neq 0 goto migration_error

psql -U postgres -d fasiolas_game -f migrations\000002_create_games_table.up.sql >nul 2>&1
if !errorlevel! neq 0 goto migration_error

psql -U postgres -d fasiolas_game -f migrations\000003_create_game_players_table.up.sql >nul 2>&1
if !errorlevel! neq 0 goto migration_error

psql -U postgres -d fasiolas_game -f migrations\000004_create_game_actions_table.up.sql >nul 2>&1
if !errorlevel! neq 0 goto migration_error

echo ✓ Migrations applied successfully
goto migration_done

:migration_error
echo ✗ Migration error - make sure PostgreSQL is running
echo   Database: fasiolas_game
echo   User: postgres
pause
exit /b 1

:migration_done

REM Step 3: Clear frontend cache
echo.
echo [3/6] Clearing frontend cache...
cd /d "%~dp0frontend"
if exist "node_modules\.cache" rmdir /s /q "node_modules\.cache" >nul 2>&1
if exist "node_modules\.react-scripts-cache" rmdir /s /q "node_modules\.react-scripts-cache" >nul 2>&1
echo ✓ Cache cleared

REM Step 4: Start backend
echo.
echo [4/6] Starting backend server...
cd /d "%~dp0"
start "BACKEND" cmd /k "server.exe"
timeout /t 3 >nul
echo ✓ Backend started (http://localhost:8080)

REM Step 5: Start frontend
echo.
echo [5/6] Starting frontend dev server...
cd /d "%~dp0frontend"
start "FRONTEND" cmd /k "npm start"
timeout /t 3 >nul
echo ✓ Frontend starting (http://localhost:3000)

REM Step 6: Done
echo.
echo [6/6] Complete!
echo.
echo ╔════════════════════════════════════════════╗
echo ║  SYSTEM READY                              ║
echo ╚════════════════════════════════════════════╝
echo.
echo Frontend:  http://localhost:3000
echo Backend:   http://localhost:8080
echo.
echo Wait 10-15 seconds for frontend to compile...
echo.
pause

