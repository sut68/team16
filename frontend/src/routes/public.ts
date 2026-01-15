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
    meta: { guest: true },
  },
  // {
  //   path: '/register',
  //   name: 'Register',
  //   component: () => import('../pages/public/RegisterPage.vue'),
  //   meta: { guest: true },
  // },
];

export default publicRoutes;