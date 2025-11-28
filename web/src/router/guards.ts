import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

/**
 * Authentication guard for protected routes
 * Checks if user is authenticated before allowing access
 */
export async function authGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  // Allow access to public routes
  if (!requiresAuth) {
    return next()
  }

  // Check if user is authenticated
  if (!authStore.isAuthenticated) {
    // Try to load user from storage first
    await authStore.loadUserFromStorage()
  }

  // If authenticated, allow access
  if (authStore.isAuthenticated) {
    return next()
  }

  // Not authenticated - save intended destination and redirect to login
  sessionStorage.setItem('returnUrl', to.fullPath)
  return next({
    name: 'login',
    query: { redirect: to.fullPath }
  })
}
