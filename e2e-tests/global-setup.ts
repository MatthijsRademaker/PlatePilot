import type { FullConfig } from "@playwright/test";

const BASE_URL = process.env.E2E_BASE_URL || "http://localhost:9001";
const API_URL = process.env.E2E_API_URL || "http://localhost:8081";
const MAX_RETRIES = 60;
const RETRY_INTERVAL_MS = 2000;

/**
 * Global setup for E2E tests
 * Waits for both the frontend and API to be ready before running tests
 */
async function globalSetup(_config: FullConfig): Promise<void> {
  console.log("\n===========================================");
  console.log("PlatePilot E2E Test Setup");
  console.log("===========================================\n");
  console.log(`Frontend URL: ${BASE_URL}`);
  console.log(`API URL: ${API_URL}`);
  console.log("");

  // Wait for API to be healthy
  await waitForService(`${API_URL}/health`, "API (BFF)");

  // Wait for API to be fully ready (all dependencies connected)
  await waitForService(`${API_URL}/ready`, "API Ready Check");

  // Wait for frontend to be accessible
  await waitForService(BASE_URL, "Frontend");

  // Verify we have seeded data
  await verifySeededData(`${API_URL}/v1/recipe/all?pageIndex=1&pageSize=5`);

  console.log("\n===========================================");
  console.log("Setup complete - starting tests");
  console.log("===========================================\n");
}

async function waitForService(url: string, name: string): Promise<void> {
  console.log(`Waiting for ${name}...`);

  for (let i = 0; i < MAX_RETRIES; i++) {
    try {
      const response = await fetch(url, {
        method: "GET",
      });

      if (response.ok) {
        console.log(`  ✓ ${name} is ready`);
        return;
      }
    } catch {
      // Service not ready yet
    }

    await sleep(RETRY_INTERVAL_MS);
  }

  throw new Error(
    `${name} did not become ready at ${url} after ${(MAX_RETRIES * RETRY_INTERVAL_MS) / 1000}s`,
  );
}

async function verifySeededData(url: string): Promise<void> {
  console.log("Verifying seeded data...");

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: { Accept: "application/json" },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch recipes: ${response.status}`);
    }

    const recipes = await response.json();

    if (!Array.isArray(recipes) || recipes.length === 0) {
      throw new Error("No recipes found - seeding may have failed");
    }

    console.log(`  ✓ Found ${recipes.length} seeded recipes`);
  } catch (error) {
    throw new Error(`Failed to verify seeded data: ${error}`);
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export default globalSetup;
