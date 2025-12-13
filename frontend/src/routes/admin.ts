import type { RouteRecordRaw } from 'vue-router';
import SidebarLayout from '@/components/SidebarLayout.vue';

const adminRoutes: Array<RouteRecordRaw> = [
  {
    path: '/admin',
    component: SidebarLayout,
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/pages/protected/DashboardPage.vue'),
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('@/pages/protected/UserManagementPage.vue'),
      },
      {
        path: 'approval',
        name: 'Approval',
        component: () => import('@/pages/protected/Approval/admin/ApprovalList.vue'),
      },
      {
        path: 'news',
        name: 'News',
        component: () => import('@/pages/protected/News/NewsList.vue'),
      },
      {
        path: 'companies',
        name: 'Companies',
        component: () => import('@/pages/protected/DashboardPage.vue'),
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('@/pages/protected/DashboardPage.vue'),
      },
      {
        path: 'screening',
        name: 'screening',
        component: () => import('@/pages/protected/Screening/screening.vue'),
      },
      {
        path: 'interview',
        name: 'Interview',
        component: () => import('@/pages/protected/DashboardPage.vue'),
      },
      {
        path: 'consider',
        name: 'Consider',
        component: () => import('@/pages/protected/DashboardPage.vue'),
      },
      {
        path: 'assistance',
        name: 'AdminAssistance',
        component: () => import('@/pages/protected/DashboardPage.vue'),
      },
    ],
    meta: { requiresAuth: true, role: 'admin' },
  },
];

export default adminRoutes;
