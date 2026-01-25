#!/usr/bin/env pwsh

Write-Host "=== FASIOLAS OAUTH DEBUG SCRIPT ===" -ForegroundColor Cyan
Write-Host ""

# Check Docker
Write-Host "1. Checking Docker..." -ForegroundColor Yellow
try {
    docker --version
    Write-Host "   ✓ Docker is installed" -ForegroundColor Green
} catch {
    Write-Host "   ✗ Docker is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

# Check Docker containers
Write-Host ""
Write-Host "2. Checking Docker containers..." -ForegroundColor Yellow
$containers = docker ps --filter "name=fasiolas" --format "{{.Names}}"
if ($containers) {
    Write-Host "   ✓ Found containers:" -ForegroundColor Green
    docker ps --filter "name=fasiolas" --format "   - {{.Names}} ({{.Status}})"
} else {
    Write-Host "   ✗ No fasiolas containers running" -ForegroundColor Red
    Write-Host "   Starting containers..." -ForegroundColor Yellow
    Set-Location "C:\Users\ciuta\Desktop\viko\interneto svetainių serverio dalies kūrimas\cardGame"
    docker-compose up -d
    Start-Sleep -Seconds 10
}

# Check backend health
Write-Host ""
Write-Host "3. Checking backend health..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 5
    Write-Host "   ✓ Backend is responding" -ForegroundColor Green
    Write-Host "   Response: $($response.Content)" -ForegroundColor Gray
} catch {
    Write-Host "   ✗ Backend is not responding: $($_.Exception.Message)" -ForegroundColor Red
}

# Check OAuth endpoint
Write-Host ""
Write-Host "4. Checking OAuth endpoint..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/google" -Method Get -TimeoutSec 5
    $content = $response.Content | ConvertFrom-Json
    Write-Host "   ✓ OAuth endpoint is responding" -ForegroundColor Green
    Write-Host "   Status: $($response.StatusCode)" -ForegroundColor Gray
    if ($content.url) {
        Write-Host "   ✓ OAuth URL received: $($content.url.Substring(0, 80))..." -ForegroundColor Green
    } else {
        Write-Host "   ✗ No URL in response" -ForegroundColor Red
        Write-Host "   Response: $($response.Content)" -ForegroundColor Gray
    }
} catch {
    Write-Host "   ✗ OAuth endpoint error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "   Response body: $responseBody" -ForegroundColor Gray
    }
}

# Check frontend
Write-Host ""
Write-Host "5. Checking frontend..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:3000" -Method Get -TimeoutSec 5
    Write-Host "   ✓ Frontend is responding" -ForegroundColor Green
} catch {
    Write-Host "   ✗ Frontend is not responding: $($_.Exception.Message)" -ForegroundColor Red
}

# Check logs
Write-Host ""
Write-Host "6. Recent backend logs:" -ForegroundColor Yellow
docker logs fasiolas_app --tail 20 2>&1 | ForEach-Object { Write-Host "   $_" -ForegroundColor Gray }

Write-Host ""
Write-Host "=== DEBUG COMPLETE ===" -ForegroundColor Cyan

