export interface UserProfile {
  sub: string
  name?: string
  email?: string
  email_verified?: boolean
  preferred_username?: string
  given_name?: string
  family_name?: string
  [key: string]: any
}
