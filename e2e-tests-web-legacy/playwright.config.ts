import { defineConfig, devices } from '@playwright/test';

/**
 * PlatePilot E2E Test Configuration
 *
 * This config is designed to work with the docker-compose.e2e.yml stack:
 * - Frontend runs on port 9001
 * - API (BFF) runs on port 8081
 *
 * Usage:
 *   bun run test:e2e        # Run all E2E tests
 *   bun run test:e2e:ui     # Run with Playwright UI
 *   bun run test:e2e:debug  # Run in debug mode
 */

// Environment configuration
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:9001';
const API_URL = process.env.E2E_API_URL || 'http://localhost:8081';

export default defineConfig({
  testDir: './tests',

  // Run tests in files in parallel
  fullyParallel: true,

  // Fail the build on CI if you accidentally left test.only in the source code
  forbidOnly: !!process.env.CI,

  // Retry on CI only
  retries: process.env.CI ? 2 : 0,

  // Opt out of parallel tests on CI for more stable results
  workers: process.env.CI ? 1 : undefined,

  // Reporter configuration
  reporter: [
    ['html', { open: 'never' }],
    ['list'],
    ...(process.env.CI ? [['github' as const]] : []),
  ],

  // Timeout configuration
  timeout: 30000,
  expect: {
    timeout: 10000,
  },

  // Global setup - waits for the stack to be ready
  globalSetup: './global-setup.ts',

  // Shared settings for all projects
  use: {
    // Base URL for navigation
    baseURL: BASE_URL,

    // Collect trace on first retry
    trace: 'on-first-retry',

    // Take screenshot on failure
    screenshot: 'only-on-failure',

    // Record video on first retry
    video: 'on-first-retry',

    // Extra HTTP headers for API requests
    extraHTTPHeaders: {
      'Accept': 'application/json',
    },
  },

  // Configure projects for major browsers
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    // Uncomment to test on additional browsers
    // {
    //   name: 'firefox',
    //   use: { ...devices['Desktop Firefox'] },
    // },
    // {
    //   name: 'webkit',
    //   use: { ...devices['Desktop Safari'] },
    // },
    // Mobile viewport testing
    // {
    //   name: 'Mobile Chrome',
    //   use: { ...devices['Pixel 5'] },
    // },
  ],

  // Output directory for test artifacts
  outputDir: './test-results',
});

// Export for use in tests
export { BASE_URL, API_URL };
