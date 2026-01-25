# Load environment variables from .env
$envFile = "$(Get-Location)\.env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#=]+)\s*=\s*(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            Set-Item -Path "env:$name" -Value $value
            Write-Host "Set $name"
        }
    }
    Write-Host "Environment variables loaded from .env"
    Write-Host ""
}

# Start the server
Write-Host "Starting Fasiolas Card Game Server..."
Write-Host "API will be available at: http://localhost:8080"
Write-Host "Press Ctrl+C to stop the server"
Write-Host ""

go run cmd/server/main.go
