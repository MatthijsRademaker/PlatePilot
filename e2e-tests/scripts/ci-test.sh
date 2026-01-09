#!/bin/bash
# Simplified wrapper for CI environments
# This script is optimized for GitHub Actions but works locally too

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Detect CI environment
IS_CI=${CI:-false}

if [ "$IS_CI" = "true" ]; then
    echo "Running in CI environment"

    # CI-specific settings
    export IOS_SIMULATOR_NAME="${IOS_SIMULATOR_NAME:-iPhone 15 Pro}"
    export IOS_VERSION="${IOS_VERSION:-17.5}"

    # Run the main script with appropriate flags
    "$SCRIPT_DIR/run-ios-e2e.sh" "$@"
else
    echo "Running locally"

    # Local development - might want to skip backend if already running
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "Backend is already running, skipping backend setup"
        "$SCRIPT_DIR/run-ios-e2e.sh" --no-backend "$@"
    else
        "$SCRIPT_DIR/run-ios-e2e.sh" "$@"
    fi
fi
