import type { RouteRecordRaw } from 'vue-router';
import StudentSidebar from '@/components/StudentSidebar.vue';

const studentRoutes: Array<RouteRecordRaw> = [
  {
    path: '/dashboard',
    component: StudentSidebar,
    children: [
      {
        path: '',
        name: 'StudentDashboard',
        component: () => import('@/pages/protected/News/student/newslist.vue'),
      },
      {
        path: 'apply-scholarship',
        name: 'ApplyScholarship',
        component: () => import('@/pages/protected/Approval/student/ScholarshipDetail.vue'),
      },
      {
        path: 'track-status',
        name: 'TrackStatus',
        component: () => import('@/pages/protected/Approval/student/StudentDocument.vue'),
      },
      {
        path: 'assistance',
        name: 'StudentAssistance',
        component: () => import('@/pages/protected/Test/StudentInterviewBookingTest.vue'),
      },
      {
        path: 'profilestudent',
        name: 'ProfileStudent',
        component: () => import('@/pages/protected/User/MyProfile.vue'),
      },
      {
        path: 'schedule',
        name: 'Schedule',
        component: () => import('@/pages/protected/Interview/student/StudentInterviewBooking.vue'),
      },
    ],
    meta: { requiresAuth: true, role: 'user' },
  },
];

export default studentRoutes;