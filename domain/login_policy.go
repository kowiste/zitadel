package domain

// LoginPolicyConfig holds the login policy configuration for an organization
// This controls authentication methods and user registration settings
type LoginPolicyConfig struct {
	// AllowRegister controls whether the "Register" button appears on the login page
	// false = Register button is HIDDEN (users can only be created by admins)
	// true = Register button is VISIBLE (users can self-register)
	AllowRegister bool `json:"allowRegister"`

	// AllowUsernamePassword enables/disables username+password authentication
	// true = Users can login with username and password
	// false = Username/password login is disabled (must use other methods like SSO)
	AllowUsernamePassword bool `json:"allowUsernamePassword"`

	// AllowExternalIdp enables/disables external identity provider login (Google, GitHub, etc.)
	// true = Users can login via configured external identity providers
	// false = External IdP login is disabled
	AllowExternalIdp bool `json:"allowExternalIdp"`

	// ForceMfa requires multi-factor authentication for ALL users
	// true = All users MUST set up and use MFA (2FA, OTP, etc.)
	// false = MFA is optional (users can choose to enable it)
	ForceMfa bool `json:"forceMfa"`

	// ForceMfaLocalOnly requires MFA only for local password-based logins
	// true = MFA required only for username/password logins (not for SSO/external IdP)
	// false = MFA requirement applies to all login methods if ForceMfa is enabled
	ForceMfaLocalOnly bool `json:"forceMfaLocalOnly"`

	// HidePasswordReset hides the "Forgot Password?" / "Reset Password" link on login page
	// true = Password reset link is HIDDEN (users must contact admin)
	// false = Password reset link is VISIBLE (users can reset their own passwords)
	HidePasswordReset bool `json:"hidePasswordReset"`

	// IgnoreUnknownUsernames prevents user enumeration attacks
	// true = Shows generic error "Invalid credentials" even if username doesn't exist
	// false = Shows specific error "User not found" vs "Wrong password"
	// Recommendation: Set to true for security (prevents attackers from discovering valid usernames)
	IgnoreUnknownUsernames bool `json:"ignoreUnknownUsernames"`

	// AllowDomainDiscovery enables automatic organization discovery based on email domain
	// true = Users can login with just email, ZITADEL auto-detects their organization
	// false = Users must specify organization explicitly (via org parameter)
	AllowDomainDiscovery bool `json:"allowDomainDiscovery"`

	// DisableLoginWithEmail disables logging in with email address
	// true = Users CANNOT use email to login (must use username)
	// false = Users CAN use email address to login
	DisableLoginWithEmail bool `json:"disableLoginWithEmail"`

	// DisableLoginWithPhone disables logging in with phone number
	// true = Users CANNOT use phone to login (most common setting)
	// false = Users CAN use phone number to login (requires phone verification)
	DisableLoginWithPhone bool `json:"disableLoginWithPhone"`

	// PasswordlessType specifies the passwordless authentication method
	// "" (empty) = Passwordless is disabled
	// "allowed" = Passwordless is optional (users can choose)
	// "required" = Passwordless is mandatory (users must use passkeys/WebAuthn)
	PasswordlessType string `json:"passwordlessType,omitempty"`

	// DefaultRedirectUri is the default redirect URL after login if none is specified
	// Leave empty to use ZITADEL's default behavior
	DefaultRedirectUri string `json:"defaultRedirectUri,omitempty"`
}

// DefaultLoginPolicyConfig returns the default login policy configuration
// with registration disabled by default for security
func DefaultLoginPolicyConfig() *LoginPolicyConfig {
	return &LoginPolicyConfig{
		AllowRegister:              false, // DISABLE REGISTRATION BY DEFAULT
		AllowUsernamePassword:      true,
		AllowExternalIdp:           false,
		ForceMfa:                   false,
		ForceMfaLocalOnly:          false,
		HidePasswordReset:          false,
		IgnoreUnknownUsernames:     false,
		AllowDomainDiscovery:       false,
		DisableLoginWithEmail:      false,
		DisableLoginWithPhone:      true,
		PasswordlessType:           "",
		DefaultRedirectUri:         "",
	}
}
