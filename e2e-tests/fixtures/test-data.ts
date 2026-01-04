/**
 * Test data constants for E2E tests
 * These IDs match the seed-data.json file for deterministic testing
 */

export const TEST_RECIPES = {
  spaghettiCarbonara: {
    id: 'e2e00001-0001-0001-0001-000000000001',
    name: 'Spaghetti Carbonara',
    cuisine: 'Italian',
  },
  chickenTikkaMasala: {
    id: 'e2e00001-0001-0001-0001-000000000002',
    name: 'Chicken Tikka Masala',
    cuisine: 'Indian',
  },
  beefTacos: {
    id: 'e2e00001-0001-0001-0001-000000000003',
    name: 'Beef Tacos',
    cuisine: 'Mexican',
  },
  caesarSalad: {
    id: 'e2e00001-0001-0001-0001-000000000004',
    name: 'Caesar Salad',
    cuisine: 'American',
  },
  padThai: {
    id: 'e2e00001-0001-0001-0001-000000000005',
    name: 'Pad Thai',
    cuisine: 'Thai',
  },
} as const;

export const TEST_CUISINES = {
  italian: {
    id: 'e2e00001-2001-0001-0001-000000000001',
    name: 'Italian',
  },
  indian: {
    id: 'e2e00001-2002-0001-0001-000000000001',
    name: 'Indian',
  },
  mexican: {
    id: 'e2e00001-2003-0001-0001-000000000001',
    name: 'Mexican',
  },
  american: {
    id: 'e2e00001-2004-0001-0001-000000000001',
    name: 'American',
  },
  thai: {
    id: 'e2e00001-2005-0001-0001-000000000001',
    name: 'Thai',
  },
} as const;

export const RECIPE_COUNT = 5;

// API endpoints for direct testing
export const API_ENDPOINTS = {
  health: '/health',
  ready: '/ready',
  allRecipes: '/v1/recipe/all',
  recipeById: (id: string) => `/v1/recipe/${id}`,
  similarRecipes: (id: string, amount = 3) =>
    `/v1/recipe/similar?recipe=${id}&amount=${amount}`,
  createRecipe: '/v1/recipe/create',
  mealPlanSuggest: '/v1/mealplan/suggest',
} as const;
