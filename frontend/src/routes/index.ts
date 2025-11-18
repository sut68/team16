import { createRouter, createWebHistory } from 'vue-router';
import publicRoutes from './public';
import protectedRoutes from './protected';
import { getToken } from '../services/api';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [...publicRoutes, ...protectedRoutes],
});

router.beforeEach((to, _from, next) => {
  const loggedIn = getToken();

  if (to.meta.requiresAuth && !loggedIn) {
    next('/login');
  } else {
    next();
  }
});

export default router;