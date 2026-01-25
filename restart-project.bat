@echo off
REM Complete Restart Script for Fasiolas Card Game with Database Migrations

echo ====================================
echo Restarting Fasiolas Card Game
echo ====================================
echo.

REM Step 1: Stop Docker
echo Step 1: Stopping Docker containers...
docker compose down -v --remove-orphans
timeout /t 2 /nobreak
echo ✓ Docker containers stopped
echo.

REM Step 2: Build server
echo Step 2: Building Go server...
go build -o server.exe cmd/server/main.go
if errorlevel 1 (
    echo ✗ Failed to build Go server
    exit /b 1
)
echo ✓ Go server built successfully
echo.

REM Step 3: Start Docker
echo Step 3: Starting Docker containers...
docker compose up -d
echo ✓ Docker containers started
echo.

REM Step 4: Wait for PostgreSQL
echo Step 4: Waiting for PostgreSQL to be ready...
setlocal enabledelayedexpansion
set /a attempt=0
set /a maxAttempts=30

:wait_loop
set /a attempt=!attempt!+1
if !attempt! gtr !maxAttempts! (
    echo ✗ PostgreSQL failed to start
    goto error
)

docker compose exec -T postgres pg_isready -U postgres >nul 2>&1
if errorlevel 1 (
    echo   Waiting... (attempt !attempt!/!maxAttempts!)
    timeout /t 1 /nobreak
    goto wait_loop
)

echo ✓ PostgreSQL is ready
echo.

REM Step 5: Apply migrations
echo Step 5: Applying database migrations...

for %%F in (
    "migrations/000001_create_users_table.up.sql"
    "migrations/000002_create_games_table.up.sql"
    "migrations/000003_create_game_players_table.up.sql"
    "migrations/000004_create_game_actions_table.up.sql"
) do (
    if exist %%F (
        echo   Running: %%F
        type %%F | docker compose exec -T postgres psql -U postgres -d fasiolas_game >nul 2>&1
        if errorlevel 1 (
            echo   ✗ Failed: %%F
        ) else (
            echo   ✓ Applied: %%F
        )
    ) else (
        echo   ✗ Not found: %%F
    )
)

echo ✓ All migrations applied
echo.

REM Step 6: Verify tables
echo Step 6: Verifying database tables...
docker compose exec -T postgres psql -U postgres -d fasiolas_game -c "SELECT table_name FROM information_schema.tables WHERE table_schema='public' ORDER BY table_name;"
echo.

REM Step 7: Check containers
echo Step 7: Checking container status...
docker ps --all
echo.

REM Step 8: Done
echo ====================================
echo Project restart complete!
echo ====================================
echo.
echo Access the application at:
echo   Frontend: http://localhost:3000
echo   Backend:  http://localhost:8080
echo.
echo Next steps:
echo   1. Open http://localhost:3000 in your browser
echo   2. Click 'Login with Google'
echo   3. You should now be able to login without database errors
echo.
pause
exit /b 0

:error
echo.
echo Check Docker logs: docker logs fasiolas_postgres
pause
exit /b 1
