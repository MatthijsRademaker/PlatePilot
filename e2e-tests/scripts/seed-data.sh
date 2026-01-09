#!/bin/bash
# Seed test data for E2E tests
# Seeds recipes and other necessary data into the database

set -e  # Exit on error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "🌱 Seeding test data..."

# Navigate to project root
cd "$PROJECT_ROOT"

# Check if backend is running
if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ Error: Backend is not running. Run ./scripts/setup-backend.sh first."
    exit 1
fi

# Check if recipes.json exists
SEED_FILE="$PROJECT_ROOT/data/recipes.json"
if [ ! -f "$SEED_FILE" ]; then
    echo "❌ Error: Seed file not found at $SEED_FILE"
    exit 1
fi

# Seed using the recipe-api seeder
echo "📊 Seeding recipes from $SEED_FILE..."
docker compose exec -T recipe-api /app/recipe-api -seed /data/recipes.json

echo "✅ Test data seeded successfully!"
echo ""
echo "Recipes seeded. You can verify by running:"
echo "  curl http://localhost:8080/v1/recipe/all"
echo ""
