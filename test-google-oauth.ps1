# Google OAuth Configuration Test Script
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Google OAuth Configuration Test" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

$BASE_URL = "http://localhost:8080"
$GOOGLE_CONSOLE_URL = "https://console.cloud.google.com/apis/credentials?project=fasiolas-auth"

# Step 1: Check .env file
Write-Host "Step 1: Checking .env file..." -ForegroundColor Cyan
if (Test-Path ".env") {
    Write-Host "[OK] .env file found" -ForegroundColor Green
    $envContent = Get-Content ".env" -Raw

    if ($envContent -match "GOOGLE_CLIENT_ID=(.+)") {
        $clientId = $matches[1].Trim()
        if ($clientId -ne "default-client-id" -and $clientId -ne "") {
            Write-Host "[OK] GOOGLE_CLIENT_ID is configured" -ForegroundColor Green
        } else {
            Write-Host "[FAIL] GOOGLE_CLIENT_ID not properly configured" -ForegroundColor Red
        }
    }

    if ($envContent -match "GOOGLE_CLIENT_SECRET=(.+)") {
        $clientSecret = $matches[1].Trim()
        if ($clientSecret -ne "default-secret" -and $clientSecret -ne "") {
            Write-Host "[OK] GOOGLE_CLIENT_SECRET is configured" -ForegroundColor Green
        } else {
            Write-Host "[FAIL] GOOGLE_CLIENT_SECRET not properly configured" -ForegroundColor Red
        }
    }

    if ($envContent -match "GOOGLE_REDIRECT_URL=(.+)") {
        $redirectUrl = $matches[1].Trim()
        Write-Host "[OK] GOOGLE_REDIRECT_URL: $redirectUrl" -ForegroundColor Green
    }
} else {
    Write-Host "[FAIL] .env file not found" -ForegroundColor Red
}
Write-Host ""

# Step 2: Check server
Write-Host "Step 2: Checking if server is running..." -ForegroundColor Cyan
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/health" -Method Get -ErrorAction Stop
    Write-Host "[OK] Server is running!" -ForegroundColor Green
} catch {
    Write-Host "[FAIL] Server is not running" -ForegroundColor Red
    Write-Host "Start server with: docker-compose up -d" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Step 3: Test OAuth endpoint
Write-Host "Step 3: Testing Google OAuth endpoint..." -ForegroundColor Cyan
try {
    $oauth = Invoke-RestMethod -Uri "$BASE_URL/api/v1/auth/google" -Method Get -ErrorAction Stop
    Write-Host "[OK] OAuth endpoint responding!" -ForegroundColor Green
    Write-Host ""
    Write-Host "OAuth URL:" -ForegroundColor Yellow
    Write-Host $oauth.url -ForegroundColor Blue
    Write-Host ""
    Write-Host "Next Steps:" -ForegroundColor Cyan
    Write-Host "1. Copy the URL above" -ForegroundColor White
    Write-Host "2. Open it in your browser" -ForegroundColor White
    Write-Host "3. Sign in with Google" -ForegroundColor White
    Write-Host ""
} catch {
    Write-Host "[FAIL] OAuth endpoint error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Step 4: Show configuration summary
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Configuration Summary" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Client ID: 845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com" -ForegroundColor White
Write-Host "Project: fasiolas-auth" -ForegroundColor White
Write-Host "Redirect URI: http://localhost:8080/api/v1/auth/google/callback" -ForegroundColor White
Write-Host ""
Write-Host "Google Cloud Console:" -ForegroundColor Yellow
Write-Host $GOOGLE_CONSOLE_URL -ForegroundColor Blue
Write-Host ""
Write-Host "Verify in Google Console:" -ForegroundColor Yellow
Write-Host "- Authorized redirect URIs includes: http://localhost:8080/api/v1/auth/google/callback" -ForegroundColor White
Write-Host "- OAuth consent screen is configured" -ForegroundColor White
Write-Host "- Test users are added (if app in Testing mode)" -ForegroundColor White
Write-Host ""
Write-Host "See GOOGLE_OAUTH_CHECKLIST.md for detailed instructions" -ForegroundColor Cyan

