# PlatePilot E2E Tests

End-to-end tests for PlatePilot using Playwright.

## Overview

This E2E test suite:
- Spins up an isolated test environment using Docker Compose
- Seeds deterministic test data into the database
- Runs Playwright tests against the full stack (frontend + backend)
- Uses separate ports to avoid conflicts with the development environment

## Prerequisites

- Docker and Docker Compose
- Bun (package manager)
- Playwright browsers (`bunx playwright install chromium`)

## Quick Start

### Run E2E tests (recommended)

```bash
cd e2e-tests

# Run full E2E suite (starts stack, runs tests, stops stack)
bun run e2e

# Run tests and keep stack running (for debugging)
bun run e2e:keep

# Clean volumes and run fresh (useful if data is corrupted)
bun run e2e:clean
```

### Individual commands

```bash
# Start the E2E test stack
bun run stack:up

# Run tests (stack must be running)
bun run test

# Run tests with UI mode
bun run test:ui

# Run tests in debug mode
bun run test:debug

# View test report
bun run report

# View stack logs
bun run stack:logs

# Stop the stack
bun run stack:down

# Clean up (remove volumes)
bun run stack:clean
```

## Test Environment

The E2E environment runs on separate ports to avoid conflicts with development:

| Service | E2E Port | Dev Port |
|---------|----------|----------|
| Frontend | 9001 | 9000 |
| API (BFF) | 8081 | 8080 |
| PostgreSQL | 5433 | 5432 |
| RabbitMQ | 5673 | 5672 |

## Test Data

Test data is defined in `fixtures/seed-data.json` with deterministic UUIDs.
These IDs are exported in `fixtures/test-data.ts` for use in tests.

### Available Test Recipes

| Name | Cuisine | ID |
|------|---------|-----|
| Spaghetti Carbonara | Italian | `e2e00001-0001-0001-0001-000000000001` |
| Chicken Tikka Masala | Indian | `e2e00001-0001-0001-0001-000000000002` |
| Beef Tacos | Mexican | `e2e00001-0001-0001-0001-000000000003` |
| Caesar Salad | American | `e2e00001-0001-0001-0001-000000000004` |
| Pad Thai | Thai | `e2e00001-0001-0001-0001-000000000005` |

## Writing Tests

Tests are organized by feature:

```
tests/
├── home.spec.ts      # Home page tests
├── recipes.spec.ts   # Recipe list and detail tests
├── navigation.spec.ts # Navigation and responsive tests
└── api.spec.ts       # Direct API tests
```

### Using Test Data

```typescript
import { TEST_RECIPES, TEST_CUISINES, API_ENDPOINTS } from '../fixtures/test-data';

test('displays recipe', async ({ page }) => {
  await page.goto(`/recipes/${TEST_RECIPES.spaghettiCarbonara.id}`);
  await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.name)).toBeVisible();
});
```

## CI/CD

For CI environments, use:

```bash
# Start stack and wait for healthy
docker compose -f docker-compose.e2e.yml up -d --wait

# Run tests with CI-specific settings
CI=true bun run test

# Cleanup
docker compose -f docker-compose.e2e.yml down -v
```

## Troubleshooting

### Tests fail with "service not ready"
- Increase `start_period` in healthchecks
- Check `bun run stack:logs` for service errors

### Data not seeding
- Check seeder logs: `docker logs platepilot-seeder-e2e`
- Try a clean start: `bun run e2e:clean`

### Port conflicts
- Stop your development environment first
- Or update ports in `docker-compose.e2e.yml`
