import type { RouteRecordRaw } from 'vue-router';
import { homeRoutes } from '@features/home/routes';
import { recipeRoutes } from '@features/recipe/routes';
import { mealplanRoutes } from '@features/mealplan/routes';
import { searchRoutes } from '@features/search/routes';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      ...homeRoutes,
      ...recipeRoutes,
      ...mealplanRoutes,
      ...searchRoutes,
    ],
  },

  // Always leave this as last one
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
