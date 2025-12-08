import type { RouteRecordRaw } from 'vue-router';
import StudentSidebar from '../components/StudentSidebar.vue';

const studentRoutes: Array<RouteRecordRaw> = [
  {
    path: '/dashboard',
    component: StudentSidebar,
    children: [
      {
        path: '',
        name: 'StudentDashboard',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'apply-scholarship',
        name: 'ApplyScholarship',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'track-status',
        name: 'TrackStatus',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
      {
        path: 'assistance',
        name: 'StudentAssistance',
        component: () => import('../pages/protected/DashboardPage.vue'),
      },
    ],
    meta: { requiresAuth: true, role: 'user' },
  },
];

export default studentRoutes;
