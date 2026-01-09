#!/bin/bash
# Complete iOS E2E Test Setup and Execution Script
# Works for both local development and GitHub CI
#
# Usage:
#   ./scripts/run-ios-e2e.sh              # Run with defaults
#   ./scripts/run-ios-e2e.sh --clean      # Clean build
#   ./scripts/run-ios-e2e.sh --no-build   # Skip iOS build (if already built)
#   ./scripts/run-ios-e2e.sh --flow home  # Run specific test flow

set -e  # Exit on error

# ============================================
# Configuration
# ============================================
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
E2E_DIR="$PROJECT_ROOT/e2e-tests"
IOS_DIR="$PROJECT_ROOT/src/ios/PlatePilot"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
IOS_SIMULATOR_NAME="${IOS_SIMULATOR_NAME:-iPhone 15 Pro}"
IOS_VERSION="${IOS_VERSION:-17.5}"  # Adjust based on available runtime
XCODE_SCHEME="PlatePilot"
APP_BUNDLE_ID="com.platepilot.app"

# Flags
CLEAN_BUILD=false
SKIP_BUILD=false
SKIP_BACKEND=false
TEST_FLOW=""
DEBUG_MODE=false

# ============================================
# Helper Functions
# ============================================
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo ""
    echo -e "${GREEN}===================================================${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}===================================================${NC}"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --clean)
                CLEAN_BUILD=true
                shift
                ;;
            --no-build)
                SKIP_BUILD=true
                shift
                ;;
            --no-backend)
                SKIP_BACKEND=true
                shift
                ;;
            --flow)
                TEST_FLOW="$2"
                shift 2
                ;;
            --debug)
                DEBUG_MODE=true
                shift
                ;;
            --simulator)
                IOS_SIMULATOR_NAME="$2"
                shift 2
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

show_help() {
    cat << EOF
iOS E2E Test Runner

Usage: $0 [OPTIONS]

Options:
    --clean             Clean build (remove derived data)
    --no-build          Skip iOS app build (use existing build)
    --no-backend        Skip backend setup (assumes already running)
    --flow <name>       Run specific test flow (e.g., home, recipes)
    --debug             Run Maestro in debug mode
    --simulator <name>  Use specific simulator (default: iPhone 15 Pro)
    --help              Show this help message

Examples:
    $0                              # Full setup and run all tests
    $0 --clean                      # Clean build and run all tests
    $0 --no-build --flow home       # Skip build, run only home tests
    $0 --debug --flow recipes       # Debug mode with specific flow

EOF
}

# ============================================
# Dependency Checks
# ============================================
check_dependencies() {
    log_step "Checking Dependencies"

    # Check if running on macOS
    if [[ "$OSTYPE" != "darwin"* ]]; then
        log_error "This script requires macOS to run iOS simulators"
        exit 1
    fi

    # Check Xcode
    if ! command -v xcodebuild &> /dev/null; then
        log_error "Xcode is not installed"
        exit 1
    fi
    log_info "Xcode: $(xcodebuild -version | head -n 1)"

    # Check XcodeGen
    if ! command -v xcodegen &> /dev/null; then
        log_warning "XcodeGen not found. Installing..."
        brew install xcodegen
    fi
    log_info "XcodeGen: $(xcodegen --version)"

    # Check Maestro
    if ! command -v maestro &> /dev/null; then
        log_warning "Maestro not found. Installing..."
        curl -Ls "https://get.maestro.mobile.dev" | bash
        export PATH="$HOME/.maestro/bin:$PATH"

        if ! command -v maestro &> /dev/null; then
            log_error "Failed to install Maestro"
            exit 1
        fi
    fi
    log_info "Maestro: $(maestro --version)"

    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    # Check if Docker is running
    if ! docker info &> /dev/null; then
        log_error "Docker is not running. Please start Docker Desktop."
        exit 1
    fi
    log_info "Docker: $(docker --version)"

    log_success "All dependencies are available"
}

# ============================================
# iOS Project Setup
# ============================================
generate_xcode_project() {
    log_step "Generating Xcode Project"

    cd "$IOS_DIR"

    if [ ! -f "project.yml" ]; then
        log_error "project.yml not found in $IOS_DIR"
        exit 1
    fi

    log_info "Running xcodegen..."
    xcodegen generate

    if [ ! -f "PlatePilot.xcodeproj/project.pbxproj" ]; then
        log_error "Failed to generate Xcode project"
        exit 1
    fi

    log_success "Xcode project generated"
}

# ============================================
# iOS Build
# ============================================
build_ios_app() {
    log_step "Building iOS App"

    cd "$IOS_DIR"

    # Clean if requested
    if [ "$CLEAN_BUILD" = true ]; then
        log_info "Cleaning build directory..."
        rm -rf ./build
        rm -rf ~/Library/Developer/Xcode/DerivedData/PlatePilot-*
    fi

    log_info "Building $XCODE_SCHEME for iOS Simulator..."
    log_info "This may take a few minutes..."

    # Build for testing (includes building the app)
    xcodebuild \
        -project PlatePilot.xcodeproj \
        -scheme "$XCODE_SCHEME" \
        -sdk iphonesimulator \
        -configuration Debug \
        -derivedDataPath ./build \
        build-for-testing \
        | grep -E "^\*\*|error:|warning:|succeeded|failed" || true

    # Check if build succeeded
    if [ ${PIPESTATUS[0]} -ne 0 ]; then
        log_error "iOS build failed"
        exit 1
    fi

    # Find the built app
    APP_PATH=$(find ./build/Build/Products/Debug-iphonesimulator -name "*.app" -type d | head -n 1)
    if [ -z "$APP_PATH" ]; then
        log_error "Could not find built .app file"
        exit 1
    fi

    log_info "App built at: $APP_PATH"
    log_success "iOS app built successfully"

    # Export for later use
    export APP_PATH
}

# ============================================
# Backend Setup
# ============================================
setup_backend() {
    log_step "Setting Up Backend Services"

    cd "$PROJECT_ROOT"

    # Check if docker-compose.yml exists
    if [ ! -f "docker-compose.yml" ]; then
        log_error "docker-compose.yml not found in project root"
        exit 1
    fi

    # Stop any existing services
    log_info "Stopping existing services..."
    docker compose down 2>/dev/null || true

    # Start services
    log_info "Starting Docker Compose services..."
    docker compose up -d

    # Wait for services
    log_info "Waiting for services to be ready..."
    sleep 5

    # Check PostgreSQL
    log_info "Checking PostgreSQL..."
    local retries=0
    local max_retries=30
    until docker compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; do
        retries=$((retries + 1))
        if [ $retries -ge $max_retries ]; then
            log_error "PostgreSQL failed to start"
            docker compose logs postgres
            exit 1
        fi
        sleep 2
    done
    log_success "PostgreSQL is ready"

    # Check RabbitMQ
    log_info "Checking RabbitMQ..."
    retries=0
    until docker compose exec -T rabbitmq rabbitmq-diagnostics -q ping > /dev/null 2>&1; do
        retries=$((retries + 1))
        if [ $retries -ge $max_retries ]; then
            log_error "RabbitMQ failed to start"
            docker compose logs rabbitmq
            exit 1
        fi
        sleep 2
    done
    log_success "RabbitMQ is ready"

    # Check Mobile BFF
    log_info "Checking Mobile BFF..."
    retries=0
    until curl -f http://localhost:8080/health > /dev/null 2>&1; do
        retries=$((retries + 1))
        if [ $retries -ge $max_retries ]; then
            log_error "Mobile BFF failed to start"
            docker compose logs mobile-bff
            exit 1
        fi
        sleep 2
    done
    log_success "Mobile BFF is ready"

    # Seed data
    log_info "Seeding test data..."
    "$E2E_DIR/scripts/seed-data.sh"

    log_success "Backend services are ready"
}

# ============================================
# iOS Simulator Setup
# ============================================
setup_simulator() {
    log_step "Setting Up iOS Simulator"

    # List available simulators
    log_info "Available simulators:"
    xcrun simctl list devices | grep -E "iPhone|iPad" | grep -v "unavailable"

    # Find or create simulator
    log_info "Looking for simulator: $IOS_SIMULATOR_NAME"

    # Get the device UDID
    DEVICE_UDID=$(xcrun simctl list devices available | grep "$IOS_SIMULATOR_NAME" | grep -v "unavailable" | head -n 1 | grep -o '[A-F0-9]\{8\}-[A-F0-9]\{4\}-[A-F0-9]\{4\}-[A-F0-9]\{4\}-[A-F0-9]\{12\}')

    if [ -z "$DEVICE_UDID" ]; then
        log_warning "Simulator not found, creating new one..."

        # Try to create simulator (this might fail if runtime not installed)
        DEVICE_UDID=$(xcrun simctl create "PlatePilot Test" "com.apple.CoreSimulator.SimDeviceType.${IOS_SIMULATOR_NAME// /-}" "com.apple.CoreSimulator.SimRuntime.iOS-${IOS_VERSION//./-}" 2>&1) || true

        if [ -z "$DEVICE_UDID" ] || [[ "$DEVICE_UDID" == *"Error"* ]]; then
            log_error "Failed to create simulator. Available runtimes:"
            xcrun simctl list runtimes
            exit 1
        fi
    fi

    log_info "Using simulator: $DEVICE_UDID"

    # Boot simulator
    log_info "Booting simulator..."
    xcrun simctl boot "$DEVICE_UDID" 2>/dev/null || true
    xcrun simctl bootstatus "$DEVICE_UDID" -b

    # Wait a bit for simulator to fully boot
    sleep 3

    log_success "Simulator is ready"

    # Export for later use
    export DEVICE_UDID
}

# ============================================
# Install App
# ============================================
install_app() {
    log_step "Installing App on Simulator"

    if [ -z "$APP_PATH" ]; then
        log_error "APP_PATH not set. Build the app first."
        exit 1
    fi

    if [ -z "$DEVICE_UDID" ]; then
        log_error "DEVICE_UDID not set. Set up simulator first."
        exit 1
    fi

    # Uninstall old version (if exists)
    log_info "Uninstalling old version (if exists)..."
    xcrun simctl uninstall "$DEVICE_UDID" "$APP_BUNDLE_ID" 2>/dev/null || true

    # Install app
    log_info "Installing app: $APP_PATH"
    xcrun simctl install "$DEVICE_UDID" "$APP_PATH"

    # Verify installation
    if xcrun simctl listapps "$DEVICE_UDID" | grep -q "$APP_BUNDLE_ID"; then
        log_success "App installed successfully"
    else
        log_error "Failed to install app"
        exit 1
    fi
}

# ============================================
# Run Maestro Tests
# ============================================
run_maestro_tests() {
    log_step "Running Maestro E2E Tests"

    cd "$E2E_DIR"

    # Ensure Maestro is in PATH
    export PATH="$HOME/.maestro/bin:$PATH"

    # Build test command
    local maestro_cmd="maestro test"

    if [ "$DEBUG_MODE" = true ]; then
        maestro_cmd="$maestro_cmd --debug"
    fi

    # Determine what to test
    if [ -n "$TEST_FLOW" ]; then
        local flow_file="flows/${TEST_FLOW}.yaml"
        if [ ! -f "$flow_file" ]; then
            log_error "Flow file not found: $flow_file"
            log_info "Available flows:"
            ls -1 flows/*.yaml | xargs -n 1 basename
            exit 1
        fi
        log_info "Running test flow: $TEST_FLOW"
        $maestro_cmd "$flow_file"
    else
        log_info "Running all test flows in flows/"
        $maestro_cmd flows/
    fi

    if [ $? -eq 0 ]; then
        log_success "All tests passed!"
    else
        log_error "Tests failed"
        return 1
    fi
}

# ============================================
# Cleanup
# ============================================
cleanup() {
    log_step "Cleaning Up"

    # Shutdown simulator if we created it
    if [ -n "$DEVICE_UDID" ]; then
        log_info "Shutting down simulator..."
        xcrun simctl shutdown "$DEVICE_UDID" 2>/dev/null || true
    fi

    # Stop backend if we started it
    if [ "$SKIP_BACKEND" = false ]; then
        log_info "Stopping backend services..."
        cd "$PROJECT_ROOT"
        docker compose down 2>/dev/null || true
    fi

    log_success "Cleanup complete"
}

# ============================================
# Main Execution
# ============================================
main() {
    echo -e "${BLUE}"
    cat << "EOF"
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║        PlatePilot iOS E2E Test Runner                ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"

    # Parse arguments
    parse_args "$@"

    # Show configuration
    log_info "Configuration:"
    log_info "  Project Root:      $PROJECT_ROOT"
    log_info "  iOS Directory:     $IOS_DIR"
    log_info "  E2E Directory:     $E2E_DIR"
    log_info "  Simulator:         $IOS_SIMULATOR_NAME"
    log_info "  Clean Build:       $CLEAN_BUILD"
    log_info "  Skip Build:        $SKIP_BUILD"
    log_info "  Skip Backend:      $SKIP_BACKEND"
    log_info "  Test Flow:         ${TEST_FLOW:-all}"
    log_info "  Debug Mode:        $DEBUG_MODE"

    # Set up trap for cleanup
    trap cleanup EXIT

    # Start the process
    START_TIME=$(date +%s)

    # Step 1: Check dependencies
    check_dependencies

    # Step 2: Generate Xcode project
    if [ "$SKIP_BUILD" = false ]; then
        generate_xcode_project
    fi

    # Step 3: Build iOS app
    if [ "$SKIP_BUILD" = false ]; then
        build_ios_app
    else
        log_warning "Skipping iOS build (using existing build)"
        # Try to find existing build
        APP_PATH=$(find "$IOS_DIR/build/Build/Products/Debug-iphonesimulator" -name "*.app" -type d 2>/dev/null | head -n 1)
        if [ -z "$APP_PATH" ]; then
            log_error "No existing build found. Run without --no-build"
            exit 1
        fi
        log_info "Using existing app: $APP_PATH"
        export APP_PATH
    fi

    # Step 4: Setup backend
    if [ "$SKIP_BACKEND" = false ]; then
        setup_backend
    else
        log_warning "Skipping backend setup (assuming already running)"
        # Verify backend is accessible
        if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
            log_error "Backend is not accessible at http://localhost:8080"
            exit 1
        fi
        log_success "Backend is accessible"
    fi

    # Step 5: Setup simulator
    setup_simulator

    # Step 6: Install app
    install_app

    # Step 7: Run tests
    run_maestro_tests

    # Calculate duration
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    MINUTES=$((DURATION / 60))
    SECONDS=$((DURATION % 60))

    # Summary
    echo ""
    log_step "Test Run Complete"
    log_success "Total time: ${MINUTES}m ${SECONDS}s"
    echo ""
}

# Run main function
main "$@"
