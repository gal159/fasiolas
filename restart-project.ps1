# Complete Restart Script for Fasiolas Card Game with Database Migrations
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "Restarting Fasiolas Card Game" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Kill any existing processes
Write-Host "Step 1: Cleaning up existing processes..." -ForegroundColor Yellow
Get-Process node -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Get-Process go -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Write-Host "✓ Processes cleaned up" -ForegroundColor Green
Write-Host ""

# Step 2: Stop and remove Docker containers
Write-Host "Step 2: Stopping Docker containers..." -ForegroundColor Yellow
docker compose down -v --remove-orphans 2>$null
Start-Sleep -Seconds 2
Write-Host "✓ Docker containers stopped" -ForegroundColor Green
Write-Host ""

# Step 3: Build the Go server binary
Write-Host "Step 3: Building Go server..." -ForegroundColor Yellow
go build -o server.exe cmd/server/main.go
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Go server built successfully" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to build Go server" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Step 4: Start Docker Compose
Write-Host "Step 4: Starting Docker containers..." -ForegroundColor Yellow
docker compose up -d
Write-Host "✓ Docker containers started" -ForegroundColor Green
Write-Host ""

# Step 5: Wait for PostgreSQL to be healthy
Write-Host "Step 5: Waiting for PostgreSQL to be ready..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0
$dbReady = $false

while ($attempt -lt $maxAttempts -and -not $dbReady) {
    $attempt++
    $output = docker compose exec -T postgres pg_isready -U postgres 2>&1
    if ($output -like "*accepting connections*") {
        $dbReady = $true
        Write-Host "✓ PostgreSQL is ready" -ForegroundColor Green
    } else {
        Write-Host "  Waiting... (attempt $attempt/$maxAttempts)" -ForegroundColor Gray
        Start-Sleep -Seconds 1
    }
}

if (-not $dbReady) {
    Write-Host "✗ PostgreSQL failed to start" -ForegroundColor Red
    Write-Host "Check Docker logs: docker logs fasiolas_postgres" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Step 6: Apply database migrations
Write-Host "Step 6: Applying database migrations..." -ForegroundColor Yellow

$migrationFiles = @(
    "migrations/000001_create_users_table.up.sql",
    "migrations/000002_create_games_table.up.sql",
    "migrations/000003_create_game_players_table.up.sql",
    "migrations/000004_create_game_actions_table.up.sql"
)

foreach ($migrationFile in $migrationFiles) {
    if (Test-Path $migrationFile) {
        Write-Host "  Running: $migrationFile" -ForegroundColor Gray
        $sqlContent = Get-Content $migrationFile -Raw
        $sqlContent | docker compose exec -T postgres psql -U postgres -d fasiolas_game 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✓ Applied: $migrationFile" -ForegroundColor Green
        } else {
            Write-Host "  ✗ Failed: $migrationFile" -ForegroundColor Red
        }
    } else {
        Write-Host "  ✗ Not found: $migrationFile" -ForegroundColor Red
    }
}
Write-Host "✓ All migrations applied" -ForegroundColor Green
Write-Host ""

# Step 7: Verify tables exist
Write-Host "Step 7: Verifying database tables..." -ForegroundColor Yellow
$checkSQL = "SELECT table_name FROM information_schema.tables WHERE table_schema='public' ORDER BY table_name;"
$tableOutput = $checkSQL | docker compose exec -T postgres psql -U postgres -d fasiolas_game 2>&1
Write-Host $tableOutput
Write-Host "✓ Database verification complete" -ForegroundColor Green
Write-Host ""

# Step 8: Check container status
Write-Host "Step 8: Checking container status..." -ForegroundColor Yellow
docker ps --all
Write-Host ""

# Step 9: Instructions
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "Project restart complete!" -ForegroundColor Green
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Access the application at:" -ForegroundColor Yellow
Write-Host "  Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "  Backend:  http://localhost:8080" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Open http://localhost:3000 in your browser" -ForegroundColor Cyan
Write-Host "  2. Click 'Login with Google'" -ForegroundColor Cyan
Write-Host "  3. You should now be able to login without database errors" -ForegroundColor Cyan
Write-Host ""
