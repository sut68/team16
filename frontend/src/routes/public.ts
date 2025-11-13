import type { RouteRecordRaw } from 'vue-router';

const publicRoutes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../pages/public/HomePage.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/public/LoginPage.vue'),
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../pages/public/RegisterPage.vue'),
  },
];

export default publicRoutes;