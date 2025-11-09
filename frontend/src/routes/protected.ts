import type { RouteRecordRaw } from 'vue-router';

const protectedRoutes: Array<RouteRecordRaw> = [
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../pages/DashboardPage.vue'),
    meta: { requiresAuth: true },
  },
];

export default protectedRoutes;
