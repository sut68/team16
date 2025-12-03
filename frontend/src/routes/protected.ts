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
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('../pages/protected/UserManagementPage.vue'),
      },
      {
        path: 'approval',
        name: 'Approval',
        component: () => import('../pages/protected/Approval/ApprovalList.vue'),
      },
    ],
    meta: { requiresAuth: true },
  },
];

export default protectedRoutes;
