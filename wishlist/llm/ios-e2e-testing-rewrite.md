# iOS E2E Testing Rewrite

## Context

PlatePilot has migrated from a Vue.js/Quasar web frontend to a native iOS SwiftUI application. The existing E2E tests in `e2e-tests/` are written using Playwright and target the deprecated web frontend. These tests need to be rewritten to test the native iOS app.

## Current State

### Existing E2E Test Structure
- **Location**: `e2e-tests/`
- **Framework**: Playwright (browser automation)
- **Target**: Vue.js/Quasar web app (DEPRECATED)
- **Test Files**:
  - `tests/home.spec.ts` - Home page tests
  - `tests/recipes.spec.ts` - Recipe browsing tests
  - `tests/navigation.spec.ts` - Navigation tests
  - `tests/api.spec.ts` - API integration tests

### iOS App Structure
- **Location**: `src/ios/PlatePilot/`
- **Framework**: SwiftUI (iOS 26+)
- **Features**:
  - Home dashboard with daily plan
  - Recipe browsing and detail views
  - Meal planning (week view)
  - Search interface
  - Authentication flow
  - Shopping list
  - Insights/analytics

## Goal

Rewrite the E2E test suite to test the native iOS app using **XCUITest** (Apple's native UI testing framework).

## Requirements

### Test Framework
- **Primary**: XCUITest (native iOS UI testing)
- **Language**: Swift (to match the app codebase)

### Test Coverage

The new test suite should cover the same scenarios as the existing Playwright tests, adapted for the iOS app:

#### 1. Home Screen Tests
- Home dashboard loads correctly
- Today's plan card displays meal information
- Daily calorie tracker shows current progress
- Recipe suggestions carousel is functional
- Navigation to recipe detail from home

#### 2. Recipe Tests
- Recipe list displays cards correctly
- Recipe search/filtering works
- Recipe detail view shows all information (ingredients, directions, nutrition)
- Recipe card tap navigates to detail view
- Recipe creation flow (if implemented)

#### 3. Navigation Tests
- Tab bar navigation between Home, Recipes, Meal Plan, Search, Insights
- Back navigation works correctly
- Deep linking to specific recipes (if supported)
- Sheet/modal presentation and dismissal

#### 4. Meal Planning Tests
- Week view displays correctly
- Meal slots can be filled with recipes
- Drag-and-drop functionality (if implemented)
- Meal plan persistence

#### 5. Authentication Tests
- Sign in flow
- Sign up flow
- Token persistence (Keychain)
- Logged out state handling

#### 6. API Integration Tests
- Backend connectivity (Mobile BFF at `http://localhost:8080/v1`)
- Recipe fetching from API
- Error handling for network failures
- Loading states

### Test Infrastructure

#### Backend Setup
- Tests should start/stop the backend services (Docker Compose)
- Seed test data before running tests
- Clean up test data after tests

#### iOS Simulator Setup
- Tests should run on iOS 26+ simulator
- App should be built and installed before tests
- App should connect to local backend (localhost tunneling)

#### CI/CD Integration
- Tests should run in GitHub Actions
- macOS runner required for iOS simulator
- Xcode installed and configured

## Implementation Approach

### Phase 1: Setup and Infrastructure
1. Create new test target in `src/ios/PlatePilot/` called `PlatePilotUITests`
2. Configure `project.yml` (XcodeGen) to include UI test target
3. Set up XCUITest framework and dependencies
4. Create test helper utilities (app launch, element finders, etc.)

### Phase 2: Port Existing Tests
1. Rewrite home screen tests using XCUITest
2. Rewrite recipe browsing tests
3. Rewrite navigation tests
4. Add meal planning tests
5. Add authentication tests

### Phase 3: CI/CD Integration
1. Update GitHub Actions workflow to run iOS UI tests
2. Configure macOS runner with Xcode
3. Set up backend services in CI
4. Configure test data seeding

### Phase 4: Documentation
1. Update `e2e-tests/README.md` or create `src/ios/PlatePilot/UITests/README.md`
2. Document how to run tests locally
3. Document CI/CD test execution
4. Add troubleshooting guide

## Example Test Structure

```swift
// PlatePilotUITests/HomeScreenTests.swift
import XCTest

final class HomeScreenTests: XCTestCase {
    var app: XCUIApplication!

    override func setUpWithError() throws {
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["--uitesting"]
        app.launch()
    }

    func testHomeScreenLoads() throws {
        // Verify home screen elements
        XCTAssertTrue(app.staticTexts["Today's Plan"].exists)
        XCTAssertTrue(app.otherElements["TodayPlanCard"].exists)
        XCTAssertTrue(app.otherElements["DailyCalorieTracker"].exists)
    }

    func testNavigateToRecipeDetail() throws {
        // Tap on a recipe card
        let firstRecipeCard = app.otherElements["RecipeCard"].firstMatch
        XCTAssertTrue(firstRecipeCard.waitForExistence(timeout: 5))
        firstRecipeCard.tap()

        // Verify recipe detail screen
        XCTAssertTrue(app.navigationBars["Recipe"].exists)
        XCTAssertTrue(app.staticTexts["Ingredients"].exists)
        XCTAssertTrue(app.staticTexts["Directions"].exists)
    }
}
```

## Acceptance Criteria

- [ ] All existing test scenarios are ported to XCUITest
- [ ] Tests run successfully on iOS 26+ simulator
- [ ] Tests can be run locally with a single command
- [ ] Tests are integrated into CI/CD pipeline
- [ ] Documentation is updated with new testing approach
- [ ] Backend services are automatically started/stopped for tests
- [ ] Test data is seeded and cleaned up properly
- [ ] Legacy Playwright tests in `e2e-tests/` are removed or clearly marked as deprecated

## Related Files

### To Update
- `e2e-tests/README.md` - Update to reference iOS tests or mark as deprecated
- `.github/workflows/ci.yml` - Add iOS UI test job
- `src/ios/PlatePilot/project.yml` - Add UI test target

### To Create
- `src/ios/PlatePilot/PlatePilotUITests/` - New UI test target directory
- UI test files for each feature (Home, Recipes, MealPlan, etc.)
- Test helper utilities and base classes

### To Deprecate/Remove
- `e2e-tests/tests/*.spec.ts` - Playwright test files (after migration)
- `e2e-tests/fixtures/` - May need to adapt for iOS tests
- `e2e-tests/playwright.config.ts` - No longer needed

## Estimated Effort

- **Phase 1 (Setup)**: 2-3 days
- **Phase 2 (Port Tests)**: 5-7 days
- **Phase 3 (CI/CD)**: 2-3 days
- **Phase 4 (Documentation)**: 1 day

**Total**: ~10-14 days

## Notes

- XCUITest is the native and most reliable option for iOS UI testing
- Consider using Page Object pattern for maintainability
- Accessibility identifiers should be added to SwiftUI views for easier element selection
- Network stubbing/mocking may be needed for some tests (consider using URLProtocol)
- Simulator networking with localhost should work, but may need special configuration

## Priority

**Medium-High** - While the app is functional, automated E2E tests are crucial for:
- Regression prevention
- Confidence in releases
- CI/CD pipeline quality gates
- Onboarding new developers

The lack of iOS E2E tests is a technical debt item that should be addressed soon.
