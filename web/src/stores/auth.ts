import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { User, UserManager } from 'oidc-client-ts'
import { createUserManager } from '@/config/oidc'
import type { UserProfile } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const selectedOrganization = ref<string>('')
  const userManager = ref<UserManager | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const isAuthenticated = computed(() => !!user.value && !user.value.expired)
  const accessToken = computed(() => user.value?.access_token || null)
  const userProfile = computed<UserProfile | null>(() => user.value?.profile as UserProfile || null)

  /**
   * Initiates the login flow by redirecting to Zitadel
   * @param organizationDomain - The organization's primary domain
   */
  async function login(organizationDomain: string) {
    try {
      isLoading.value = true
      error.value = null

      // Store organization for callback
      selectedOrganization.value = organizationDomain
      localStorage.setItem('selectedOrganization', organizationDomain)

      // Create UserManager with org scope
      userManager.value = createUserManager(organizationDomain)

      // Redirect to Zitadel login
      await userManager.value.signinRedirect()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed'
      console.error('Login error:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Handles the OAuth callback and exchanges the authorization code for tokens
   * @returns The authenticated user
   */
  async function handleCallback(): Promise<User> {
    try {
      isLoading.value = true
      error.value = null

      // Retrieve stored organization
      const orgDomain = localStorage.getItem('selectedOrganization')
      if (!orgDomain) {
        throw new Error('Organization not found. Please login again.')
      }

      // Create UserManager with same org scope
      userManager.value = createUserManager(orgDomain)

      // Process callback and get user
      const callbackUser = await userManager.value.signinRedirectCallback()
      user.value = callbackUser
      selectedOrganization.value = orgDomain

      console.log('Authentication successful:', {
        name: callbackUser.profile.name,
        email: callbackUser.profile.email,
        organization: orgDomain,
      })

      return callbackUser
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Authentication callback failed'
      console.error('Callback error:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Logs out the user and redirects to Zitadel logout endpoint
   */
  async function logout() {
    try {
      isLoading.value = true

      if (userManager.value) {
        await userManager.value.signoutRedirect()
      }

      // Clear state
      user.value = null
      selectedOrganization.value = ''
      userManager.value = null
      localStorage.removeItem('selectedOrganization')
    } catch (err) {
      console.error('Logout error:', err)
      error.value = err instanceof Error ? err.message : 'Logout failed'
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Loads user from localStorage on app initialization or page refresh
   * @returns The stored user or null
   */
  async function loadUserFromStorage(): Promise<User | null> {
    try {
      const orgDomain = localStorage.getItem('selectedOrganization')
      if (!orgDomain) {
        return null
      }

      // Create UserManager with stored org
      userManager.value = createUserManager(orgDomain)

      // Try to load user from storage
      const storedUser = await userManager.value.getUser()

      if (storedUser && !storedUser.expired) {
        user.value = storedUser
        selectedOrganization.value = orgDomain
        console.log('User loaded from storage:', {
          name: storedUser.profile.name,
          email: storedUser.profile.email,
        })
        return storedUser
      }

      // User expired or not found
      console.log('No valid user found in storage')
      return null
    } catch (err) {
      console.error('Error loading user from storage:', err)
      return null
    }
  }

  /**
   * Clears any error messages
   */
  function clearError() {
    error.value = null
  }

  return {
    // State
    user,
    selectedOrganization,
    isLoading,
    error,

    // Computed
    isAuthenticated,
    accessToken,
    userProfile,

    // Actions
    login,
    handleCallback,
    logout,
    loadUserFromStorage,
    clearError,
  }
})
