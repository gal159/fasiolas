# Fasiolas Card Game - Startup Script
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  FASIOLAS CARD GAME - STARTING UP" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Check Docker
Write-Host "[1/5] Checking Docker..." -ForegroundColor Yellow
$dockerProcess = Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue
if ($dockerProcess) {
    Write-Host "  [OK] Docker Desktop is running" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Docker Desktop is not running" -ForegroundColor Red
    Write-Host "  Please start Docker Desktop and try again" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Step 2: Stop existing containers
Write-Host "[2/5] Stopping existing containers..." -ForegroundColor Yellow
docker compose down 2>&1 | Out-Null
Write-Host "  [OK] Cleaned up existing containers" -ForegroundColor Green
Write-Host ""

# Step 3: Start containers
Write-Host "[3/5] Starting containers..." -ForegroundColor Yellow
docker compose up -d
Write-Host "  [WAIT] Waiting for PostgreSQL to be healthy..." -ForegroundColor Yellow

# Wait for PostgreSQL to be ready
$dbReady = $false
$maxAttempts = 60
for ($i = 1; $i -le $maxAttempts; $i++) {
    $result = docker exec fasiolas_postgres pg_isready -U postgres 2>&1
    if ($result -eq "accepting connections") {
        Write-Host "  [OK] PostgreSQL is ready!" -ForegroundColor Green
        $dbReady = $true
        break
    }
    Write-Host "  [WAIT] Waiting... (attempt $i/$maxAttempts)" -ForegroundColor Yellow
    Start-Sleep -Seconds 1
}

if (-not $dbReady) {
    Write-Host "  [FAIL] PostgreSQL failed to start after $maxAttempts seconds" -ForegroundColor Red
    Write-Host "  Check Docker logs: docker logs fasiolas_postgres" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Step 3.5: Run database migrations
Write-Host "[3.5/5] Running database migrations..." -ForegroundColor Yellow
$migrationFile = "all_migrations.sql"
if (Test-Path $migrationFile) {
    Write-Host "  [INFO] Copying migrations to container..." -ForegroundColor Gray
    docker cp $migrationFile fasiolas_postgres:/tmp/migrations.sql 2>&1 | Out-Null
    Write-Host "  [INFO] Running migrations..." -ForegroundColor Gray
    docker exec fasiolas_postgres psql -U postgres -d fasiolas_game -f /tmp/migrations.sql 2>&1 | Out-Null
    Write-Host "  [OK] Database migrations completed" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Migration file not found: $migrationFile" -ForegroundColor Red
}
Write-Host ""

# Wait a bit more to ensure backend is ready
Write-Host "  [WAIT] Waiting for backend to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 15
Write-Host ""

# Step 4: Check status
Write-Host "[4/6] Checking container status..." -ForegroundColor Yellow
$containers = docker ps --format "{{.Names}}" 2>&1
if ($containers -like "*fasiolas_app*") {
    Write-Host "  [OK] Backend API:    RUNNING" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Backend API:    NOT RUNNING" -ForegroundColor Red
}

if ($containers -like "*fasiolas_postgres*") {
    Write-Host "  [OK] PostgreSQL:     RUNNING" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] PostgreSQL:     NOT RUNNING" -ForegroundColor Red
}

if ($containers -like "*fasiolas_frontend*") {
    Write-Host "  [OK] Frontend:       RUNNING" -ForegroundColor Green
} else {
    Write-Host "  [WARN] Frontend:       NOT RUNNING (optional)" -ForegroundColor Yellow
}
Write-Host ""

# Step 5: Test endpoints
Write-Host "[5/6] Testing endpoints..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -TimeoutSec 5 -ErrorAction Stop
    Write-Host "  [OK] Health Check:   PASS" -ForegroundColor Green
    Write-Host "    Status: $($health.status)" -ForegroundColor Gray
} catch {
    Write-Host "  [FAIL] Health Check:   FAIL" -ForegroundColor Red
    Write-Host "    Waiting for server to start..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    try {
        $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -TimeoutSec 5 -ErrorAction Stop
        Write-Host "  [OK] Health Check:   PASS (retry)" -ForegroundColor Green
    } catch {
        Write-Host "  [FAIL] Server not responding. Check logs: docker logs fasiolas_app" -ForegroundColor Red
    }
}

try {
    $auth = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/google" -TimeoutSec 5 -ErrorAction Stop
    Write-Host "  [OK] OAuth Endpoint: PASS" -ForegroundColor Green
} catch {
    Write-Host "  [FAIL] OAuth Endpoint: FAIL" -ForegroundColor Red
}
Write-Host ""

# Display status
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  APPLICATION STATUS" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Backend API:      http://localhost:8080" -ForegroundColor White
Write-Host "Health Check:     http://localhost:8080/health" -ForegroundColor White
Write-Host "API Base:         http://localhost:8080/api/v1" -ForegroundColor White
Write-Host "Frontend:         http://localhost:3000" -ForegroundColor White
Write-Host "PostgreSQL:       localhost:5432" -ForegroundColor White
Write-Host ""

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  NEXT STEPS" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Test OAuth Login:" -ForegroundColor Yellow
Write-Host "  ./test-google-oauth.ps1" -ForegroundColor White
Write-Host ""
Write-Host "Test Complete API:" -ForegroundColor Yellow
Write-Host "  ./test-api-complete.ps1" -ForegroundColor White
Write-Host ""
Write-Host "View Logs:" -ForegroundColor Yellow
Write-Host "  docker logs fasiolas_app -f" -ForegroundColor White
Write-Host ""
Write-Host "Stop Application:" -ForegroundColor Yellow
Write-Host "  docker compose down" -ForegroundColor White
Write-Host ""
Write-Host "[OK] Startup complete!" -ForegroundColor Green
Write-Host ""

