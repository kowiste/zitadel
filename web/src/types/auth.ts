import { User } from 'oidc-client-ts'

export interface AuthState {
  user: User | null
  selectedOrganization: string
  isLoading: boolean
  error: string | null
}

export interface OIDCConfig {
  domain: string
  clientId: string
  redirectUri: string
  postLogoutRedirectUri: string
}
