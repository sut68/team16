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
      {
        path: 'news',
        name: 'News',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'companies',
        name: 'Companies',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'screening',
        name: 'Screening',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'interview',
        name: 'Interview',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'consider',
        name: 'Consider',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'assistance',
        name: 'Assistance',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
    ],
    meta: { requiresAuth: true },
  },
];

export default protectedRoutes;
