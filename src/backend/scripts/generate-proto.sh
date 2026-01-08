#!/bin/bash

# Proto generation script for PlatePilot
# Generates Go code from protobuf definitions

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PROTO_DIR="$PROJECT_ROOT/api/proto"

echo "Generating protobuf code..."

# Create output directories
mkdir -p "$PROJECT_ROOT/internal/recipe/pb"
mkdir -p "$PROJECT_ROOT/internal/mealplanner/pb"
mkdir -p "$PROJECT_ROOT/internal/shoppinglist/pb"

# Generate Recipe API protos
echo "  - Generating recipe.proto..."
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/recipe/v1/recipe.proto"

# Move generated files to correct location
mv "$PROJECT_ROOT/recipe/v1/"*.go "$PROJECT_ROOT/internal/recipe/pb/" 2>/dev/null || true
rm -rf "$PROJECT_ROOT/recipe" 2>/dev/null || true

# Generate MealPlanner API protos
echo "  - Generating mealplanner.proto..."
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/mealplanner/v1/mealplanner.proto"

# Move generated files to correct location
mv "$PROJECT_ROOT/mealplanner/v1/"*.go "$PROJECT_ROOT/internal/mealplanner/pb/" 2>/dev/null || true
rm -rf "$PROJECT_ROOT/mealplanner" 2>/dev/null || true

# Generate ShoppingList API protos
echo "  - Generating shoppinglist.proto..."
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$PROJECT_ROOT" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$PROJECT_ROOT" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/shoppinglist/v1/shoppinglist.proto"

# Move generated files to correct location
mv "$PROJECT_ROOT/shoppinglist/v1/"*.go "$PROJECT_ROOT/internal/shoppinglist/pb/" 2>/dev/null || true
rm -rf "$PROJECT_ROOT/shoppinglist" 2>/dev/null || true

echo "Done! Generated files:"
find "$PROJECT_ROOT/internal" -name "*.pb.go" -o -name "*_grpc.pb.go" | head -20

echo ""
echo "Proto generation complete!"
