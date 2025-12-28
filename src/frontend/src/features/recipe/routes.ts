import type { RouteRecordRaw } from 'vue-router';

export const recipeRoutes: RouteRecordRaw[] = [
  {
    path: 'recipes',
    name: 'recipes',
    component: () => import('./pages/RecipeListPage.vue'),
  },
  {
    path: 'recipes/:id',
    name: 'recipe-detail',
    component: () => import('./pages/RecipeDetailPage.vue'),
  },
];
