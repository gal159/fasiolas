@echo off
REM Comprehensive startup script for Fasiolas Card Game

echo.
echo ======================================
echo Fasiolas Card Game - Complete Startup
echo ======================================
echo.

REM Step 1: Kill any existing node/go processes
echo 1. Cleaning up existing processes...
taskkill /F /IM node.exe >nul 2>&1
taskkill /F /IM go.exe >nul 2>&1

REM Step 2: Stop Docker containers
echo 2. Stopping any existing Docker containers...
docker-compose down --remove-orphans >nul 2>&1
timeout /t 2 /nobreak >nul

REM Step 3: Start PostgreSQL
echo 3. Starting PostgreSQL database...
docker-compose up -d postgres
echo    Waiting for database to be ready...
timeout /t 10 /nobreak >nul

REM Step 4: Check database
echo 4. Checking database connection...
set /a attempts=0
:check_db
set /a attempts=%attempts%+1
if %attempts% gtr 30 (
    echo    ERROR: Database failed to start
    exit /b 1
)

docker exec fasiolas_postgres pg_isready -U postgres >nul 2>&1
if errorlevel 1 (
    echo    Attempt %attempts%/30 - Waiting for PostgreSQL...
    timeout /t 1 /nobreak >nul
    goto check_db
)

echo    PostgreSQL is ready!
echo.

REM Step 5: Run database migrations
echo 5. Running database migrations...
docker cp all_migrations.sql fasiolas_postgres:/tmp/migrations.sql
docker exec fasiolas_postgres psql -U postgres -d fasiolas_game -f /tmp/migrations.sql >nul 2>&1
echo    Database migrations completed!
echo.

REM Step 6: Set environment variables and start backend
echo 6. Starting backend server...
setlocal enabledelayedexpansion
set "GOOGLE_CLIENT_ID=845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com"
set "GOOGLE_CLIENT_SECRET=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"
set "DB_HOST=localhost"
set "DB_PORT=5432"
set "DB_USER=postgres"
set "DB_PASSWORD=postgres"
set "DB_NAME=fasiolas_game"
set "SERVER_PORT=8080"
set "SERVER_ENV=development"

start "Fasiolas Backend" cmd /k "go run cmd/server/main.go"
echo    Backend started - waiting for it to be ready...
timeout /t 5 /nobreak >nul
echo.

REM Step 7: Start frontend
echo 7. Starting frontend server...
cd frontend
start "Fasiolas Frontend" cmd /k "npm start"
cd ..
echo    Frontend started - waiting for it to be ready...
timeout /t 10 /nobreak >nul
echo.

echo ======================================
echo ✓ Startup Complete!
echo ======================================
echo.
echo Frontend: http://localhost:3000
echo Backend:  http://localhost:8080
echo Health:   http://localhost:8080/health
echo.
echo Close this window to stop all services
echo ======================================
echo.

pause
