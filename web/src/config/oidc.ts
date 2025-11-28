import { UserManager, WebStorageStateStore, Log } from 'oidc-client-ts'
import type { OIDCConfig } from '@/types/auth'

// Enable debug logging in development
if (import.meta.env.DEV) {
  Log.setLogger(console)
  Log.setLevel(Log.DEBUG)
}

export const oidcConfig: OIDCConfig = {
  domain: import.meta.env.VITE_ZITADEL_DOMAIN || 'http://localhost:8080',
  clientId: import.meta.env.VITE_ZITADEL_CLIENT_ID || '',
  redirectUri: import.meta.env.VITE_REDIRECT_URI || 'http://localhost:3000/auth/callback',
  postLogoutRedirectUri: import.meta.env.VITE_POST_LOGOUT_URI || 'http://localhost:3000',
}

if (!oidcConfig.clientId) {
  console.warn('VITE_ZITADEL_CLIENT_ID is not set. Please add it to .env.local')
}

/**
 * Creates a UserManager instance with organization-scoped authentication
 * @param organizationDomain - The primary domain of the organization (e.g., "Tesla-150405")
 * @returns Configured UserManager instance
 */
export function createUserManager(organizationDomain: string): UserManager {
  return new UserManager({
    authority: oidcConfig.domain,
    client_id: oidcConfig.clientId,
    redirect_uri: oidcConfig.redirectUri,
    post_logout_redirect_uri: oidcConfig.postLogoutRedirectUri,
    response_type: 'code',
    scope: `openid email profile urn:zitadel:iam:org:domain:primary:${organizationDomain}`,

    // PKCE settings
    pkce: true,

    // Token storage
    userStore: new WebStorageStateStore({ store: window.localStorage }),

    // Silent refresh configuration
    automaticSilentRenew: true,
    accessTokenExpiringNotificationTimeInSeconds: 60,
    silent_redirect_uri: oidcConfig.redirectUri,

    // Additional settings
    loadUserInfo: true,
    monitorSession: true,
  })
}
