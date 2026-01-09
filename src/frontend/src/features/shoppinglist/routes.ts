import type { RouteRecordRaw } from 'vue-router';

export const shoppinglistRoutes: RouteRecordRaw[] = [
  {
    path: 'shopping-lists',
    name: 'shopping-lists',
    component: () => import('./pages/ShoppingListsPage.vue'),
  },
  {
    path: 'shopping-lists/:id',
    name: 'shoppinglist-detail',
    component: () => import('./pages/ShoppingListDetailPage.vue'),
  },
];
