// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import { jwtDecode } from 'jwt-decode';
import publicRoutes from './public';
import adminRoutes from './admin';
import studentRoutes from './student';
import { getToken } from '../services/api/httpsTest'; // your existing helper (may return token from storage for back-compat)

// Routes
const routes: Array<RouteRecordRaw> = [
  ...publicRoutes,
  ...adminRoutes,
  ...studentRoutes,
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});


// Helper
function decodeTokenToUser(token?: string | null): { id?: number | string | null; role?: string | null } | null {
  if (!token) return null;
  try {
    const decoded: any = jwtDecode(token as string);
    return {
      id: decoded.user_id ?? decoded.userId ?? decoded.sub ?? null,
      role: decoded.role ?? decoded.role_name ?? null,
    };
  } catch (e) {
    console.warn('[router] invalid token for decode', e);
    return null;
  }
}

// Router guard
router.beforeEach((to, _from, next) => {
  // 1) Try sessionStorage
  let currentUser: any = null;
  const raw = sessionStorage.getItem('currentUser');
  if (raw) {
    try {
      currentUser = JSON.parse(raw);
    } catch (e) {
      console.warn('[router] invalid currentUser in sessionStorage, clearing', e);
      sessionStorage.removeItem('currentUser');
      currentUser = null;
    }
  }

  // Fallback: decode token (back-compat)
  if (!currentUser) {
    const token = getToken();
    const maybeUser = decodeTokenToUser(token);
    if (maybeUser) {
      currentUser = maybeUser;
      try {
        sessionStorage.setItem('currentUser', JSON.stringify(currentUser));
      } catch (e) {
        // ignore storage errors
      }
    }
  }

  const loggedIn = !!currentUser;
  const roleRaw = currentUser?.role ?? null;
  const role = roleRaw ? String(roleRaw).toLowerCase() : null;
  const isStudent = role === 'user' || role === 'student';

  const requiresAuth = to.matched.some((rec) => !!rec.meta?.requiresAuth);
  const guestOnly = to.matched.some((rec) => !!rec.meta?.guest);

  const requiredRole = to.matched.reduce<string | undefined>((acc, r) => {
    if (acc) return acc;
    if (r.meta?.role) return String(r.meta.role);
    if (r.meta?.requireAdmin) return 'admin';
    return acc;
  }, undefined)?.toLowerCase();

  if (to.path === '/' && loggedIn) {
    if (role === 'admin') return next('/admin');
    if (isStudent) return next('/dashboard');
  }

  if (requiresAuth) {
    if (!loggedIn) {
      return next('/');
    }

    if (requiredRole) {
      if (requiredRole === 'user' && (role === 'user' || role === 'student')) {
        return next();
      }
      if (role !== requiredRole) {
        if (role === 'admin') return next('/admin');
        if (isStudent) return next('/dashboard');
        return next('/');
      }
    }
  } else if (guestOnly && loggedIn) {
    // guest-only pages (login/register)
    if (role === 'admin') return next('/admin');
    if (isStudent) return next('/dashboard');
  }

  // default: allow
  return next();
});

export default router;
