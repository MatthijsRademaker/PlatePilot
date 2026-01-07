import type { RouteRecordRaw } from 'vue-router';
import { homeRoutes } from '@features/home/routes';
import { recipeRoutes } from '@features/recipe/routes';
import { mealplanRoutes } from '@features/mealplan/routes';
import { searchRoutes } from '@features/search/routes';
import { authRoutes } from '@features/auth/routes';

const routes: RouteRecordRaw[] = [
  // Auth routes (no layout, guest only)
  ...authRoutes,

  // Protected routes with main layout
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
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
