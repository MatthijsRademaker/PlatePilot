import { test, expect } from '@playwright/test';
import { TEST_RECIPES, RECIPE_COUNT } from '../fixtures/test-data';

test.describe('Recipe List Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/#/recipes');
  });

  test('displays recipes heading', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'Recipes' })).toBeVisible();
  });

  test('loads and displays seeded recipes', async ({ page }) => {
    // Wait for recipes to load
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.name)).toBeVisible({
      timeout: 10000,
    });

    // Check that multiple recipes are displayed
    await expect(page.getByText(TEST_RECIPES.chickenTikkaMasala.name)).toBeVisible();
    await expect(page.getByText(TEST_RECIPES.beefTacos.name)).toBeVisible();
  });

  test('displays recipe cuisine information', async ({ page }) => {
    // Wait for page to load
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.name)).toBeVisible({
      timeout: 10000,
    });

    // Cuisines should be visible
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.cuisine)).toBeVisible();
  });

  test('navigates to recipe detail when clicking a recipe', async ({ page }) => {
    // Wait for recipes to load
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.name)).toBeVisible({
      timeout: 10000,
    });

    // Click on the first recipe
    await page.getByText(TEST_RECIPES.spaghettiCarbonara.name).click();

    // Should navigate to recipe detail page
    await expect(page).toHaveURL(new RegExp(`.*#/recipes/${TEST_RECIPES.spaghettiCarbonara.id}`));
  });

  test('has refresh button', async ({ page }) => {
    await expect(page.getByRole('button', { name: /refresh/i })).toBeVisible();
  });

  test('can refresh recipe list', async ({ page }) => {
    // Wait for initial load
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.name)).toBeVisible({
      timeout: 10000,
    });

    // Click refresh
    await page.getByRole('button', { name: /refresh/i }).click();

    // Recipes should still be visible after refresh
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.name)).toBeVisible({
      timeout: 10000,
    });
  });
});

test.describe('Recipe Detail Page', () => {
  test('displays recipe details', async ({ page }) => {
    // Navigate directly to a known recipe
    await page.goto(`/#/recipes/${TEST_RECIPES.spaghettiCarbonara.id}`);

    // Check for recipe name
    await expect(
      page.getByRole('heading', { name: TEST_RECIPES.spaghettiCarbonara.name })
    ).toBeVisible({ timeout: 10000 });

    // Check for description
    await expect(
      page.getByText('Classic Italian pasta dish with eggs, cheese, and pancetta')
    ).toBeVisible();
  });

  test('displays ingredients section', async ({ page }) => {
    await page.goto(`/#/recipes/${TEST_RECIPES.spaghettiCarbonara.id}`);

    // Wait for page to load
    await expect(
      page.getByRole('heading', { name: TEST_RECIPES.spaghettiCarbonara.name })
    ).toBeVisible({ timeout: 10000 });

    // Check for ingredients heading
    await expect(page.getByText('Ingredients')).toBeVisible();

    // Check for specific ingredients
    await expect(page.getByText('Spaghetti')).toBeVisible();
    await expect(page.getByText('Eggs')).toBeVisible();
    await expect(page.getByText('Pancetta')).toBeVisible();
  });

  test('displays instructions section', async ({ page }) => {
    await page.goto(`/#/recipes/${TEST_RECIPES.spaghettiCarbonara.id}`);

    // Wait for page to load
    await expect(
      page.getByRole('heading', { name: TEST_RECIPES.spaghettiCarbonara.name })
    ).toBeVisible({ timeout: 10000 });

    // Check for instructions heading
    await expect(page.getByText('Instructions')).toBeVisible();

    // Check for specific instruction steps
    await expect(page.getByText(/Bring a large pot of salted water to boil/)).toBeVisible();
    await expect(page.getByText(/Cook spaghetti until al dente/)).toBeVisible();
  });

  test('displays cuisine information', async ({ page }) => {
    await page.goto(`/#/recipes/${TEST_RECIPES.spaghettiCarbonara.id}`);

    // Wait for page to load
    await expect(
      page.getByRole('heading', { name: TEST_RECIPES.spaghettiCarbonara.name })
    ).toBeVisible({ timeout: 10000 });

    // Check for cuisine
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.cuisine)).toBeVisible();
  });

  test('has back button that navigates to list', async ({ page }) => {
    // Start at recipes list
    await page.goto('/#/recipes');
    await expect(page.getByText(TEST_RECIPES.spaghettiCarbonara.name)).toBeVisible({
      timeout: 10000,
    });

    // Navigate to detail
    await page.getByText(TEST_RECIPES.spaghettiCarbonara.name).click();
    await expect(
      page.getByRole('heading', { name: TEST_RECIPES.spaghettiCarbonara.name })
    ).toBeVisible({ timeout: 10000 });

    // Click back button
    await page.getByRole('button', { name: /back/i }).click();

    // Should be back at list
    await expect(page).toHaveURL(/.*#\/recipes$/);
  });

  test('handles non-existent recipe', async ({ page }) => {
    await page.goto('/#/recipes/non-existent-recipe-id');

    // Should show error state - look for specific error text
    await expect(page.getByText('HTTP error: Not Found')).toBeVisible({ timeout: 10000 });
  });
});
