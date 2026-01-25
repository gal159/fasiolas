#!/usr/bin/env powershell
# Complete restart script with database migrations
# Run from the cardGame directory

param(
    [switch]$NoWait = $false,
    [switch]$Verbose = $false
)

$ErrorActionPreference = 'Stop'

function Write-Step {
    param([string]$message, [int]$step)
    Write-Host "[$step] $message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$message)
    Write-Host "✓ $message" -ForegroundColor Green
}

function Write-Error-Custom {
    param([string]$message)
    Write-Host "✗ $message" -ForegroundColor Red
}

function Write-Info {
    param([string]$message)
    Write-Host "  $message" -ForegroundColor Gray
}

try {
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Fasiolas Card Game - Complete Restart" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan

    # Step 1: Stop Docker
    Write-Step "Stopping Docker containers..." 1
    docker compose down -v --remove-orphans 2>&1 | Out-Null
    Write-Success "Containers stopped"

    # Wait
    Write-Info "Waiting 3 seconds for cleanup..."
    Start-Sleep -Seconds 3

    # Step 2: Start Docker
    Write-Step "Starting Docker containers..." 2
    docker compose up -d 2>&1 | Out-Null
    Write-Success "Containers started"

    # Wait for PostgreSQL
    Write-Info "Waiting for PostgreSQL to be ready..."
    $postgresReady = $false
    for ($i = 1; $i -le 30; $i++) {
        $output = docker compose exec -T postgres pg_isready -U postgres 2>&1
        if ($output -like "*accepting connections*") {
            $postgresReady = $true
            Write-Success "PostgreSQL is ready"
            break
        }
        Write-Info "Attempt $i/30..."
        Start-Sleep -Seconds 1
    }

    if (-not $postgresReady) {
        throw "PostgreSQL failed to start after 30 seconds"
    }

    # Step 3: Apply migrations
    Write-Step "Applying database migrations..." 3

    $migrations = @(
        "migrations/000001_create_users_table.up.sql",
        "migrations/000002_create_games_table.up.sql",
        "migrations/000003_create_game_players_table.up.sql",
        "migrations/000004_create_game_actions_table.up.sql"
    )

    $migrationCount = 0
    $env:PGPASSWORD = "123456"

    foreach ($migration in $migrations) {
        $migrationCount++

        if (-not (Test-Path $migration)) {
            Write-Error-Custom "Migration file not found: $migration"
            continue
        }

        Write-Info "[$migrationCount/4] Applying $migration..."

        try {
            $sqlContent = Get-Content $migration -Raw
            $sqlContent | docker compose exec -T postgres psql -U postgres -d fasiolas_game 2>&1 | Out-Null
            Write-Success "Migration $migrationCount/4 applied"
        }
        catch {
            Write-Error-Custom "Failed to apply migration: $migration"
            Write-Error-Custom $_.Exception.Message
        }
    }

    # Step 4: Verify tables
    Write-Step "Verifying database tables..." 4

    try {
        $query = "SELECT table_name FROM information_schema.tables WHERE table_schema='public' ORDER BY table_name;"
        $tableList = $query | docker compose exec -T postgres psql -U postgres -d fasiolas_game 2>&1
        Write-Host $tableList
        Write-Success "Database verification complete"
    }
    catch {
        Write-Error-Custom "Failed to verify tables: $($_.Exception.Message)"
    }

    # Clear password
    Remove-Item env:PGPASSWORD -ErrorAction SilentlyContinue

    # Step 5: Check containers
    Write-Step "Checking container status..." 5
    docker ps
    Write-Success "Container status displayed"

    # Final message
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Restart Complete! ✓" -ForegroundColor Green
    Write-Host "========================================`n" -ForegroundColor Cyan

    Write-Host "Access your application:" -ForegroundColor Yellow
    Write-Host "  Frontend: http://localhost:3000" -ForegroundColor Cyan
    Write-Host "  Backend:  http://localhost:8080" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "  1. Open http://localhost:3000 in your browser" -ForegroundColor Cyan
    Write-Host "  2. Click 'Login with Google'" -ForegroundColor Cyan
    Write-Host "  3. You should now be able to login successfully!" -ForegroundColor Cyan
    Write-Host ""

}
catch {
    Write-Host "`nError during restart:" -ForegroundColor Red
    Write-Error-Custom $_.Exception.Message
    Remove-Item env:PGPASSWORD -ErrorAction SilentlyContinue
    exit 1
}
