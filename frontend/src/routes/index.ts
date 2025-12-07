import { createRouter, createWebHistory } from 'vue-router';
import { jwtDecode } from 'jwt-decode';
import publicRoutes from './public';
import adminRoutes from './admin';
import studentRoutes from './student';
import { getToken } from '../services/api';

// Helper to get user role from JWT
const getUserRole = (): string | null => {
  const token = getToken();
  if (token) {
    try {
      const decodedToken: { role: string } = jwtDecode(token);
      return decodedToken.role;
    } catch (error) {
      console.error('Invalid token:', error);
      return null;
    }
  }
  return null;
};

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [...publicRoutes, ...adminRoutes, ...studentRoutes],
});

router.beforeEach((to, _from, next) => {
  const loggedIn = !!getToken();
  const role = getUserRole();
  const isStudent = role === 'user' || role === 'student';

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const guestOnly = to.matched.some(record => record.meta.guest);
  const requiredRole = to.meta.role as string | undefined;

  // Redirect from root if logged in
  if (to.path === '/' && loggedIn) {
    if (role === 'admin') {
      return next('/admin');
    }
    if (isStudent) {
      return next('/dashboard');
    }
  }
  
  if (requiresAuth) {
    if (!loggedIn) {
      // User is not logged in, redirect to home
      return next('/');
    }
    if (requiredRole && role !== requiredRole) {
      // Special case for student routes, which accept 'user' role
      if (requiredRole === 'user' && isStudent) {
        return next();
      }

      // User is logged in but does not have the required role
      // Redirect to their respective dashboard
      if (role === 'admin') {
        return next('/admin');
      }
      if (isStudent) {
        return next('/dashboard');
      }
      // Fallback if role is unknown
      return next('/');
    }
  } else if (guestOnly && loggedIn) {
    // Logged-in user tries to access guest-only pages (login/register)
    if (role === 'admin') {
      return next('/admin');
    }
    if (isStudent) {
      return next('/dashboard');
    }
  }

  // All other cases
  next();
});

export default router;