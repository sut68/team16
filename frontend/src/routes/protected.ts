import type { RouteRecordRaw } from 'vue-router';

const protectedRoutes: Array<RouteRecordRaw> = [
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../pages/test.vue'),
    meta: { requiresAuth: true },
  },
];

export default protectedRoutes;
