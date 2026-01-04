import { test, expect } from '@playwright/test';

test.describe('Home Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('displays PlatePilot branding', async ({ page }) => {
    // Check for main heading
    await expect(page.getByRole('heading', { name: 'PlatePilot' })).toBeVisible();

    // Check for tagline
    await expect(
      page.getByText('Your intelligent meal planning companion')
    ).toBeVisible();
  });

  test('displays navigation cards', async ({ page }) => {
    // Check for Browse Recipes card
    await expect(page.getByText('Browse Recipes')).toBeVisible();
    await expect(
      page.getByText('Explore our collection of delicious recipes')
    ).toBeVisible();

    // Check for Meal Plan card
    await expect(page.getByText('Meal Plan')).toBeVisible();
    await expect(page.getByText('Plan your meals for the week')).toBeVisible();

    // Check for Search card
    await expect(page.getByText('Search')).toBeVisible();
    await expect(page.getByText('Find the perfect recipe')).toBeVisible();
  });

  test('navigates to recipes page when clicking Browse Recipes', async ({
    page,
  }) => {
    await page.getByText('Browse Recipes').click();

    await expect(page).toHaveURL(/.*recipes/);
    await expect(page.getByRole('heading', { name: 'Recipes' })).toBeVisible();
  });

  test('navigates to meal plan page when clicking Meal Plan', async ({
    page,
  }) => {
    await page.getByText('Meal Plan').click();

    await expect(page).toHaveURL(/.*mealplan/);
  });

  test('navigates to search page when clicking Search', async ({ page }) => {
    await page.getByText('Search').click();

    await expect(page).toHaveURL(/.*search/);
  });
});
