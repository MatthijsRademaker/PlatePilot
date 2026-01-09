# PlatePilot iOS E2E Tests

End-to-end tests for the PlatePilot iOS app using [Maestro](https://maestro.mobile.dev/).

## Overview

This directory contains automated UI tests for the native iOS app. Tests are written in Maestro's YAML format and can be run locally or in CI/CD.

## Prerequisites

### 1. Install Maestro

```bash
# macOS/Linux
curl -Ls "https://get.maestro.mobile.dev" | bash

# Add to PATH (if not already)
export PATH="$PATH:$HOME/.maestro/bin"

# Verify installation
maestro --version
```

### 2. iOS Simulator Setup

```bash
# List available simulators
xcrun simctl list devices

# Boot a simulator (iOS 26+ required)
xcrun simctl boot "iPhone 15 Pro"

# Or use Xcode UI: Xcode -> Window -> Devices and Simulators
```

### 3. Build the iOS App

```bash
cd src/ios/PlatePilot
xcodegen generate
xcodebuild -project PlatePilot.xcodeproj \
           -scheme PlatePilot \
           -sdk iphonesimulator \
           -configuration Debug \
           -derivedDataPath ./build
```

Alternatively, build from Xcode (Cmd+B).

### 4. Start Backend Services

```bash
# From project root
docker compose up -d

# Verify BFF is running
curl http://localhost:8080/health
# Should return: OK
```

## Running Tests

### Quick Start (One Command)

The easiest way to run all tests:

```bash
cd e2e-tests
./scripts/run-ios-e2e.sh
```

This single command will:
1. Check all dependencies
2. Generate Xcode project
3. Build the iOS app
4. Start backend services
5. Boot iOS Simulator
6. Install the app
7. Run all Maestro tests
8. Clean up automatically

### Using Make (Recommended)

```bash
cd e2e-tests

# Run all tests
make test

# Run specific test flow
make test-flow name=home
make test-flow name=recipes

# Run in debug mode
make test-debug

# Clean build and test
make test-clean

# Quick test (skip build if possible)
make test-quick

# Show all available commands
make help
```

### Advanced Usage

```bash
# Clean build
./scripts/run-ios-e2e.sh --clean

# Skip build (use existing)
./scripts/run-ios-e2e.sh --no-build

# Skip backend setup
./scripts/run-ios-e2e.sh --no-backend

# Run specific flow
./scripts/run-ios-e2e.sh --flow home

# Debug mode
./scripts/run-ios-e2e.sh --debug

# Custom simulator
./scripts/run-ios-e2e.sh --simulator "iPhone 14 Pro"

# Combine options
./scripts/run-ios-e2e.sh --no-backend --flow recipes --debug
```

### Manual Maestro Commands

If you prefer to run Maestro directly:

```bash
# Prerequisites: backend running, app built and installed
maestro test flows/                    # All tests
maestro test flows/home.yaml          # Specific test
maestro test --debug flows/home.yaml  # Debug mode
```

## Test Structure

```
e2e-tests/
├── flows/                  # Maestro test flows (YAML)
│   ├── home.yaml           # Home screen tests
│   ├── recipes.yaml        # Recipe browsing tests
│   ├── recipe-detail.yaml  # Recipe detail tests
│   ├── meal-plan.yaml      # Meal planning tests
│   ├── navigation.yaml     # Tab navigation tests
│   └── auth.yaml           # Authentication tests
├── fixtures/               # Test data
│   ├── test-users.json     # Test user credentials
│   └── test-recipes.json   # Sample recipe data
├── scripts/                # Helper scripts
│   ├── setup-backend.sh    # Start backend services
│   ├── seed-data.sh        # Seed test data
│   └── cleanup.sh          # Clean up test data
└── .maestro/               # Maestro configuration
    └── config.yaml         # Global test configuration
```

## Writing Tests

Maestro tests are YAML files with a simple, readable syntax:

```yaml
appId: com.platepilot.app
---
# Test: Home screen loads correctly
- launchApp
- assertVisible: "Today's Plan"
- assertVisible: "Daily Calories"

# Tap on first recipe card
- tapOn: "RecipeCard"
- assertVisible: "Ingredients"
- assertVisible: "Directions"
```

### Common Commands

```yaml
# Launch app
- launchApp

# Wait for element
- waitForElement: "ButtonText"

# Tap on element
- tapOn: "ButtonText"

# Input text
- tapOn: "Username field"
- inputText: "testuser@example.com"

# Scroll
- scroll

# Assert element is visible
- assertVisible: "Expected text"

# Take screenshot
- takeScreenshot: test-step-1
```

## Test Data

### Backend Seeding

Tests expect certain data to be present in the backend:

```bash
# Seed test data
./scripts/seed-data.sh
```

This will:
1. Create test user accounts
2. Seed sample recipes
3. Create test meal plans

### Test User Credentials

See `fixtures/test-users.json` for test account credentials:

```json
{
  "testUser": {
    "email": "test@platepilot.com",
    "password": "Test123!"
  }
}
```

## CI/CD Integration

### GitHub Actions

Tests run automatically on PR and push to main:

```yaml
# .github/workflows/e2e-ios.yml
- name: Run Maestro tests
  run: |
    cd e2e-tests
    maestro test flows/
```

See `.github/workflows/e2e-ios.yml` for full configuration.

## Troubleshooting

### App doesn't launch

**Issue**: Maestro can't find or launch the app

**Solution**:
```bash
# Ensure app is built for simulator
cd src/ios/PlatePilot
xcodebuild -project PlatePilot.xcodeproj \
           -scheme PlatePilot \
           -sdk iphonesimulator \
           -configuration Debug

# Verify app is installed in simulator
xcrun simctl listapps booted | grep platepilot
```

### Backend connection fails

**Issue**: App can't connect to localhost:8080

**Solution**:
```bash
# Verify backend is running
docker compose ps

# Check BFF health
curl http://localhost:8080/health

# iOS Simulator uses host machine's localhost, so this should work
```

### Tests fail intermittently

**Issue**: Timing issues, elements not found

**Solution**:
- Add explicit waits: `- waitForElement: "ElementId"`
- Increase timeout in `.maestro/config.yaml`
- Use `assertVisible` with retry

### Element not found

**Issue**: Maestro can't find UI elements

**Solution**:
- Add accessibility identifiers to SwiftUI views:
  ```swift
  Text("Recipe Title")
    .accessibilityIdentifier("RecipeTitle")
  ```
- Use text matching instead of IDs:
  ```yaml
  - tapOn: "Recipe Title"  # Finds by text
  ```

## Best Practices

1. **Use Accessibility Identifiers**: Add `.accessibilityIdentifier()` to SwiftUI views for reliable element selection

2. **Keep Tests Independent**: Each test should set up its own data and not rely on other tests

3. **Use Descriptive Names**: Name test files and flows clearly (e.g., `recipe-detail.yaml`, not `test1.yaml`)

4. **Add Comments**: Use YAML comments to explain complex test steps

5. **Test Happy Path First**: Focus on critical user journeys before edge cases

6. **Clean Up After Tests**: Reset app state or use fresh test data for each run

## Resources

- [Maestro Documentation](https://maestro.mobile.dev/)
- [Maestro Examples](https://github.com/mobile-dev-inc/maestro/tree/main/maestro-test)
- [SwiftUI Accessibility](https://developer.apple.com/documentation/swiftui/view-accessibility)

## Notes

- Maestro works best with accessibility identifiers. Update SwiftUI views to add them where needed.
- Tests run against the **Debug** build of the app (not Release).
- Backend must be running on `localhost:8080` for tests to pass.
- Tests can run in parallel by specifying multiple simulator devices.
