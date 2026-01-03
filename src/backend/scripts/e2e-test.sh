#!/bin/bash
# End-to-end integration test script for PlatePilot
# Tests the full flow: BFF -> Recipe API -> Event -> MealPlanner

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BFF_URL="${BFF_URL:-http://localhost:8080}"
MAX_RETRIES=30
RETRY_INTERVAL=2

echo "=============================================="
echo "PlatePilot E2E Integration Tests"
echo "=============================================="
echo ""
echo "BFF URL: $BFF_URL"
echo ""

# Helper functions
log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_error() {
    echo -e "${RED}✗ $1${NC}"
}

log_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Wait for service to be ready
wait_for_service() {
    local url=$1
    local name=$2
    local retries=0

    log_info "Waiting for $name to be ready..."

    while [ $retries -lt $MAX_RETRIES ]; do
        if curl -s -f "$url" > /dev/null 2>&1; then
            log_success "$name is ready"
            return 0
        fi
        retries=$((retries + 1))
        sleep $RETRY_INTERVAL
    done

    log_error "$name did not become ready in time"
    return 1
}

# Test health endpoint
test_health() {
    log_info "Testing health endpoint..."

    response=$(curl -s -w "\n%{http_code}" "$BFF_URL/health")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" == "200" ] && [ "$body" == "OK" ]; then
        log_success "Health check passed"
        return 0
    else
        log_error "Health check failed (HTTP $http_code)"
        return 1
    fi
}

# Test ready endpoint
test_ready() {
    log_info "Testing ready endpoint..."

    response=$(curl -s -w "\n%{http_code}" "$BFF_URL/ready")
    http_code=$(echo "$response" | tail -n1)

    if [ "$http_code" == "200" ]; then
        log_success "Ready check passed"
        return 0
    else
        log_error "Ready check failed (HTTP $http_code)"
        return 1
    fi
}

# Test get all recipes
test_get_all_recipes() {
    log_info "Testing GET /v1/recipe/all..."

    response=$(curl -s -w "\n%{http_code}" "$BFF_URL/v1/recipe/all?pageIndex=1&pageSize=5")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" == "200" ]; then
        # Check if we got a JSON array
        if echo "$body" | jq -e '. | type == "array"' > /dev/null 2>&1; then
            count=$(echo "$body" | jq 'length')
            log_success "Get all recipes passed (found $count recipes)"

            # Store first recipe ID for later tests
            if [ "$count" -gt 0 ]; then
                FIRST_RECIPE_ID=$(echo "$body" | jq -r '.[0].id')
                export FIRST_RECIPE_ID
            fi
            return 0
        else
            log_error "Response is not a JSON array"
            echo "$body" | head -c 500
            return 1
        fi
    else
        log_error "Get all recipes failed (HTTP $http_code)"
        echo "$body" | head -c 500
        return 1
    fi
}

# Test get recipe by ID
test_get_recipe_by_id() {
    if [ -z "$FIRST_RECIPE_ID" ]; then
        log_info "Skipping get by ID test (no recipe ID available)"
        return 0
    fi

    log_info "Testing GET /v1/recipe/{id}..."

    response=$(curl -s -w "\n%{http_code}" "$BFF_URL/v1/recipe/$FIRST_RECIPE_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" == "200" ]; then
        recipe_name=$(echo "$body" | jq -r '.name // empty')
        if [ -n "$recipe_name" ]; then
            log_success "Get recipe by ID passed (name: $recipe_name)"
            return 0
        else
            log_error "Recipe response missing name field"
            return 1
        fi
    else
        log_error "Get recipe by ID failed (HTTP $http_code)"
        return 1
    fi
}

# Test similar recipes
test_similar_recipes() {
    if [ -z "$FIRST_RECIPE_ID" ]; then
        log_info "Skipping similar recipes test (no recipe ID available)"
        return 0
    fi

    log_info "Testing GET /v1/recipe/similar..."

    response=$(curl -s -w "\n%{http_code}" "$BFF_URL/v1/recipe/similar?recipe=$FIRST_RECIPE_ID&amount=3")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" == "200" ]; then
        if echo "$body" | jq -e '. | type == "array"' > /dev/null 2>&1; then
            count=$(echo "$body" | jq 'length')
            log_success "Similar recipes passed (found $count similar)"
            return 0
        else
            log_error "Similar recipes response is not an array"
            return 1
        fi
    else
        log_error "Similar recipes failed (HTTP $http_code)"
        return 1
    fi
}

# Test create recipe
test_create_recipe() {
    log_info "Testing POST /v1/recipe/create..."

    # Create a test recipe
    payload='{
        "name": "E2E Test Recipe",
        "description": "A test recipe created by E2E tests",
        "prepTime": "10 minutes",
        "cookTime": "20 minutes",
        "directions": ["Step 1: Test", "Step 2: Verify"],
        "cuisineName": "Test Cuisine",
        "mainIngredientName": "Test Ingredient",
        "ingredientNames": ["Ingredient 1", "Ingredient 2"]
    }'

    response=$(curl -s -w "\n%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$payload" \
        "$BFF_URL/v1/recipe/create")

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" == "201" ] || [ "$http_code" == "200" ]; then
        recipe_id=$(echo "$body" | jq -r '.id // empty')
        if [ -n "$recipe_id" ]; then
            log_success "Create recipe passed (id: $recipe_id)"
            export CREATED_RECIPE_ID="$recipe_id"
            return 0
        else
            log_error "Created recipe missing ID"
            return 1
        fi
    else
        log_error "Create recipe failed (HTTP $http_code)"
        echo "$body" | head -c 500
        return 1
    fi
}

# Test meal plan suggestion
test_meal_plan_suggest() {
    log_info "Testing POST /v1/mealplan/suggest..."

    payload='{
        "amount": 3,
        "dailyConstraints": [],
        "alreadySelectedRecipes": []
    }'

    response=$(curl -s -w "\n%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$payload" \
        "$BFF_URL/v1/mealplan/suggest")

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" == "200" ]; then
        if echo "$body" | jq -e '. | type == "array"' > /dev/null 2>&1; then
            count=$(echo "$body" | jq 'length')
            log_success "Meal plan suggest passed (got $count suggestions)"
            return 0
        else
            log_error "Meal plan response is not an array"
            return 1
        fi
    else
        log_error "Meal plan suggest failed (HTTP $http_code)"
        echo "$body" | head -c 500
        return 1
    fi
}

# Main test execution
main() {
    local failed=0

    # Wait for BFF to be ready
    if ! wait_for_service "$BFF_URL/health" "Mobile BFF"; then
        log_error "Services are not ready. Make sure docker-compose is running."
        exit 1
    fi

    echo ""
    echo "Running tests..."
    echo "----------------------------------------------"

    # Run tests
    test_health || ((failed++))
    test_ready || ((failed++))
    test_get_all_recipes || ((failed++))
    test_get_recipe_by_id || ((failed++))
    test_similar_recipes || ((failed++))
    test_create_recipe || ((failed++))
    test_meal_plan_suggest || ((failed++))

    echo ""
    echo "----------------------------------------------"

    if [ $failed -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        exit 0
    else
        echo -e "${RED}$failed test(s) failed${NC}"
        exit 1
    fi
}

# Check for jq
if ! command -v jq &> /dev/null; then
    log_error "jq is required for this script. Install with: brew install jq"
    exit 1
fi

main "$@"
