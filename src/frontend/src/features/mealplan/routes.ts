import type { RouteRecordRaw } from 'vue-router';

export const mealplanRoutes: RouteRecordRaw[] = [
  {
    path: 'mealplan',
    name: 'mealplan',
    component: () => import('./pages/MealplanPage.vue'),
  },
];
