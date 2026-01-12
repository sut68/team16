import { createRouter, createWebHistory } from 'vue-router';
import { jwtDecode } from 'jwt-decode';
import publicRoutes from './public';
import adminRoutes from './admin';
import studentRoutes from './student';
import { getToken } from '../services/api';
import { Get } from '../services/api/https';
import { checkProfileCompleteness, PROFILE_REQUIRED_ROUTES } from '../utils/profileValidator';
import type { StudentProfileResponse, FamilyInfo } from '@/interfaces/user';

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

// Cache profile check to avoid repeated API calls
let profileCheckCache: { isComplete: boolean; timestamp: number } | null = null;
const CACHE_DURATION = 60000; // 1 minute

const checkProfile = async (): Promise<boolean> => {
  // Check cache first
  if (profileCheckCache && Date.now() - profileCheckCache.timestamp < CACHE_DURATION) {
    return profileCheckCache.isComplete;
  }

  try {
    const profileRes: any = await Get('/profile/me');
    if (profileRes && profileRes.role === 'student') {
      const student = profileRes.data as StudentProfileResponse;
      const family = profileRes.family as FamilyInfo | null;
      const result = checkProfileCompleteness(student, family);

      // Update cache
      profileCheckCache = {
        isComplete: result.isComplete,
        timestamp: Date.now(),
      };

      return result.isComplete;
    }
    return true; // Admin doesn't need profile check
  } catch (error) {
    console.error('Error checking profile:', error);
    return true; // Allow access on error to avoid blocking
  }
};

// Clear cache when needed (e.g., after profile update)
export const clearProfileCache = () => {
  profileCheckCache = null;
};

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [...publicRoutes, ...adminRoutes, ...studentRoutes],
});

router.beforeEach(async (to, _from, next) => {
  const loggedIn = !!getToken();
  const role = getUserRole();
  const isStudent = role === 'user' || role === 'student';

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const guestOnly = to.matched.some(record => record.meta.guest);
  const requiredRole = to.meta.role as string | undefined;

  // Redirect from root if logged in
  if (to.path === '/' && loggedIn) {
    if (role === 'admin') {
      return next('/admin/news');
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
        // Continue to profile check below
      } else {
        // User is logged in but does not have the required role
        // Redirect to their respective dashboard
        if (role === 'admin') {
          return next('/admin/news');
        }
        if (isStudent) {
          return next('/dashboard');
        }
        // Fallback if role is unknown
        return next('/');
      }
    }

    // Profile completeness check for students on protected routes
    if (isStudent && PROFILE_REQUIRED_ROUTES.includes(to.path)) {
      const isProfileComplete = await checkProfile();
      if (!isProfileComplete) {
        // Redirect to profile page with a flag
        return next('/dashboard/profilestudent?required=true');
      }
    }
  } else if (guestOnly && loggedIn) {
    // Logged-in user tries to access guest-only pages (login/register)
    if (role === 'admin') {
      return next('/admin/news');
    }
    if (isStudent) {
      return next('/dashboard');
    }
  }

  // All other cases
  next();
});

export default router;