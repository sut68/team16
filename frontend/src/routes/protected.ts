import type { RouteRecordRaw } from 'vue-router';
import SidebarLayout from '../components/SidebarLayout.vue';

const protectedRoutes: Array<RouteRecordRaw> = [
  {
    path: '/dashboard',
    component: SidebarLayout,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
    ],
    meta: { requiresAuth: true },
  },
];

export default protectedRoutes;
