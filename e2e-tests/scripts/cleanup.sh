#!/bin/bash
# Cleanup script for E2E tests
# Stops backend services and cleans up test data

set -e  # Exit on error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "🧹 Cleaning up PlatePilot E2E test environment..."

# Navigate to project root
cd "$PROJECT_ROOT"

# Stop Docker Compose services
echo "📦 Stopping Docker Compose services..."
docker compose down

# Optional: Remove volumes to clean database
if [ "$1" == "--volumes" ] || [ "$1" == "-v" ]; then
    echo "🗑️  Removing Docker volumes..."
    docker compose down -v
    echo "✅ Volumes removed"
fi

echo "✅ Cleanup complete!"
echo ""
echo "To restart services, run: ./scripts/setup-backend.sh"
echo ""
