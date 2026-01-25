#!/bin/bash
# Comprehensive startup script for Fasiolas Card Game

echo "======================================"
echo "Fasiolas Card Game - Complete Startup"
echo "======================================"
echo ""

# Check if docker is running
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed or not in PATH"
    exit 1
fi

echo "1. Stopping any existing processes..."
docker-compose down --remove-orphans 2>/dev/null || true
sleep 2

echo "2. Starting PostgreSQL database..."
docker-compose up -d postgres
echo "   Waiting for database to be ready..."
sleep 10

# Check if database is ready
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if docker exec fasiolas_postgres pg_isready -U postgres > /dev/null 2>&1; then
        echo "   ✓ PostgreSQL is ready"
        break
    fi
    attempt=$((attempt + 1))
    echo "   Waiting... (attempt $attempt/$max_attempts)"
    sleep 1
done

if [ $attempt -eq $max_attempts ]; then
    echo "❌ PostgreSQL failed to start after $max_attempts attempts"
    exit 1
fi

echo ""
echo "3. Starting backend server..."
export GOOGLE_CLIENT_ID="845779433699-i7hb2rk5hk080aasm63dud7chruo80ed.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="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"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="fasiolas_game"
export SERVER_PORT="8080"
export SERVER_ENV="development"

go run cmd/server/main.go > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
echo "   Backend started (PID: $BACKEND_PID)"
sleep 5

echo "4. Checking backend health..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "   ✓ Backend is responding"
else
    echo "⚠️  Backend health check failed - checking logs..."
    head -30 /tmp/backend.log
fi

echo ""
echo "5. Starting frontend..."
cd frontend
npm start > /tmp/frontend.log 2>&1 &
FRONTEND_PID=$!
echo "   Frontend started (PID: $FRONTEND_PID)"
sleep 10

echo ""
echo "======================================"
echo "✓ Startup Complete!"
echo "======================================"
echo ""
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8080"
echo "Health:   http://localhost:8080/health"
echo ""
echo "Backend logs:   tail -f /tmp/backend.log"
echo "Frontend logs:  tail -f /tmp/frontend.log"
echo ""
echo "To stop everything: press Ctrl+C"
echo "======================================"

# Wait for all processes
wait
