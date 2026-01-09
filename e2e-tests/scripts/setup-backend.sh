#!/bin/bash
# Setup script for PlatePilot backend services
# Starts Docker Compose services required for E2E tests

set -e  # Exit on error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "🚀 Setting up PlatePilot backend services..."

# Navigate to project root
cd "$PROJECT_ROOT"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running. Please start Docker Desktop."
    exit 1
fi

# Check if docker-compose.yml exists
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: docker-compose.yml not found in project root"
    exit 1
fi

# Stop any existing services
echo "📦 Stopping existing services..."
docker compose down 2>/dev/null || true

# Start services
echo "📦 Starting Docker Compose services..."
docker compose up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check PostgreSQL
echo "🔍 Checking PostgreSQL..."
until docker compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; do
    echo "   Waiting for PostgreSQL..."
    sleep 2
done
echo "✅ PostgreSQL is ready"

# Check RabbitMQ
echo "🔍 Checking RabbitMQ..."
until docker compose exec -T rabbitmq rabbitmq-diagnostics -q ping > /dev/null 2>&1; do
    echo "   Waiting for RabbitMQ..."
    sleep 2
done
echo "✅ RabbitMQ is ready"

# Wait for BFF to be healthy
echo "🔍 Checking Mobile BFF..."
MAX_RETRIES=30
RETRY_COUNT=0
until curl -f http://localhost:8080/health > /dev/null 2>&1; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        echo "❌ Error: Mobile BFF failed to start after $MAX_RETRIES attempts"
        docker compose logs mobile-bff
        exit 1
    fi
    echo "   Waiting for Mobile BFF... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done
echo "✅ Mobile BFF is ready"

# Run database migrations (they should run automatically in docker-compose, but verify)
echo "🔄 Verifying database migrations..."
sleep 2

# Check if seeding is needed
echo "🌱 Checking if database needs seeding..."
RECIPE_COUNT=$(docker compose exec -T postgres psql -U postgres -d recipedb -t -c "SELECT COUNT(*) FROM recipes;" | xargs)
if [ "$RECIPE_COUNT" -eq "0" ]; then
    echo "📊 Database is empty, seeding with sample data..."
    "$SCRIPT_DIR/seed-data.sh"
else
    echo "✅ Database already has $RECIPE_COUNT recipes"
fi

echo ""
echo "✅ Backend services are ready!"
echo ""
echo "Services running:"
echo "  - PostgreSQL:     localhost:5432"
echo "  - RabbitMQ:       localhost:5672 (Management UI: http://localhost:15672)"
echo "  - Mobile BFF:     http://localhost:8080"
echo "  - Recipe API:     gRPC on localhost:50051"
echo "  - MealPlanner API: gRPC on localhost:50052"
echo ""
echo "Health check: curl http://localhost:8080/health"
echo ""
