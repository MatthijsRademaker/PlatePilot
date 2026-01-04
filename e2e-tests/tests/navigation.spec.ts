import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
  test('can navigate through all main pages', async ({ page }) => {
    // Start at home
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'PlatePilot' })).toBeVisible();

    // Navigate to recipes
    await page.getByText('Browse Recipes').click();
    await expect(page).toHaveURL(/.*recipes/);
    await expect(page.getByRole('heading', { name: 'Recipes' })).toBeVisible();

    // Navigate back to home (using browser back or nav)
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'PlatePilot' })).toBeVisible();

    // Navigate to meal plan
    await page.getByText('Meal Plan').click();
    await expect(page).toHaveURL(/.*mealplan/);

    // Navigate to search
    await page.goto('/');
    await page.getByText('Search').click();
    await expect(page).toHaveURL(/.*search/);
  });

  test('handles 404 for unknown routes', async ({ page }) => {
    await page.goto('/unknown-page-that-does-not-exist');

    // Should show error page or redirect
    // The exact behavior depends on your 404 page implementation
    await expect(page.getByText(/not found|error|404/i)).toBeVisible({
      timeout: 5000,
    });
  });

  test('maintains state when using browser back button', async ({ page }) => {
    // Navigate to recipes
    await page.goto('/recipes');
    await expect(page.getByRole('heading', { name: 'Recipes' })).toBeVisible({
      timeout: 10000,
    });

    // Wait for recipes to load
    await page.waitForTimeout(1000);

    // Go to home
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'PlatePilot' })).toBeVisible();

    // Go back
    await page.goBack();

    // Should be at recipes page
    await expect(page).toHaveURL(/.*recipes/);
  });
});

test.describe('Responsive Layout', () => {
  test('displays correctly on mobile viewport', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    await page.goto('/');

    // Main elements should still be visible
    await expect(page.getByRole('heading', { name: 'PlatePilot' })).toBeVisible();
    await expect(page.getByText('Browse Recipes')).toBeVisible();
  });

  test('displays correctly on tablet viewport', async ({ page }) => {
    // Set tablet viewport
    await page.setViewportSize({ width: 768, height: 1024 });

    await page.goto('/');

    // Main elements should still be visible
    await expect(page.getByRole('heading', { name: 'PlatePilot' })).toBeVisible();
    await expect(page.getByText('Browse Recipes')).toBeVisible();
  });

  test('displays correctly on desktop viewport', async ({ page }) => {
    // Set desktop viewport
    await page.setViewportSize({ width: 1920, height: 1080 });

    await page.goto('/');

    // Main elements should still be visible
    await expect(page.getByRole('heading', { name: 'PlatePilot' })).toBeVisible();
    await expect(page.getByText('Browse Recipes')).toBeVisible();
  });
});
