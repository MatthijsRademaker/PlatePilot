import { test, expect } from '@playwright/test';
import { TEST_RECIPES, API_ENDPOINTS } from '../fixtures/test-data';

// API URL from environment or default
const API_URL = process.env.E2E_API_URL || 'http://localhost:8081';

test.describe('API Health Checks', () => {
  test('health endpoint returns OK', async ({ request }) => {
    const response = await request.get(`${API_URL}${API_ENDPOINTS.health}`);

    expect(response.ok()).toBeTruthy();
    expect(await response.text()).toBe('OK');
  });

  test('ready endpoint returns OK', async ({ request }) => {
    const response = await request.get(`${API_URL}${API_ENDPOINTS.ready}`);

    expect(response.ok()).toBeTruthy();
  });
});

test.describe('Recipe API', () => {
  test('GET /v1/recipe/all returns seeded recipes', async ({ request }) => {
    const response = await request.get(
      `${API_URL}${API_ENDPOINTS.allRecipes}?pageIndex=1&pageSize=10`
    );

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.items).toBeDefined();
    expect(Array.isArray(data.items)).toBeTruthy();
    expect(data.items.length).toBeGreaterThan(0);
    expect(data.totalCount).toBeGreaterThan(0);
    expect(data.pageIndex).toBe(1);

    // Check for seeded recipe
    const carbonara = data.items.find(
      (r: { name: string }) => r.name === TEST_RECIPES.spaghettiCarbonara.name
    );
    expect(carbonara).toBeDefined();
  });

  test('GET /v1/recipe/:id returns specific recipe', async ({ request }) => {
    const response = await request.get(
      `${API_URL}${API_ENDPOINTS.recipeById(TEST_RECIPES.spaghettiCarbonara.id)}`
    );

    expect(response.ok()).toBeTruthy();

    const recipe = await response.json();
    expect(recipe.name).toBe(TEST_RECIPES.spaghettiCarbonara.name);
    expect(recipe.id).toBe(TEST_RECIPES.spaghettiCarbonara.id);
  });

  test('GET /v1/recipe/:id returns 404 for non-existent recipe', async ({
    request,
  }) => {
    const response = await request.get(
      `${API_URL}${API_ENDPOINTS.recipeById('00000000-0000-0000-0000-000000000000')}`
    );

    expect(response.status()).toBe(404);
  });

  test('GET /v1/recipe/similar returns similar recipes', async ({ request }) => {
    const response = await request.get(
      `${API_URL}${API_ENDPOINTS.similarRecipes(TEST_RECIPES.spaghettiCarbonara.id, 3)}`
    );

    expect(response.ok()).toBeTruthy();

    const recipes = await response.json();
    expect(Array.isArray(recipes)).toBeTruthy();
  });

  test('POST /v1/recipe/create creates a new recipe', async ({ request }) => {
    const newRecipe = {
      name: `E2E Test Recipe ${Date.now()}`,
      description: 'A recipe created during E2E testing',
      prepTime: '10 minutes',
      cookTime: '15 minutes',
      directions: ['Step 1: Test', 'Step 2: Verify'],
      cuisineName: 'Test Cuisine',
      mainIngredientName: 'Test Ingredient',
      ingredientNames: ['Ingredient A', 'Ingredient B'],
    };

    const response = await request.post(`${API_URL}${API_ENDPOINTS.createRecipe}`, {
      data: newRecipe,
    });

    expect(response.status()).toBe(201);

    const created = await response.json();
    expect(created.id).toBeDefined();
    expect(created.name).toBe(newRecipe.name);
  });
});

test.describe('Meal Plan API', () => {
  test('POST /v1/mealplan/suggest returns suggestions', async ({ request }) => {
    const response = await request.post(`${API_URL}${API_ENDPOINTS.mealPlanSuggest}`, {
      data: {
        amount: 3,
        dailyConstraints: [],
        alreadySelectedRecipes: [],
      },
    });

    expect(response.ok()).toBeTruthy();

    const suggestions = await response.json();
    expect(Array.isArray(suggestions)).toBeTruthy();
  });
});
