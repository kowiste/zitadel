package domain

type Auth interface {
	CreateOrganization(orgName string) (*Organization, error)
	GetOrganization(orgID string) (*Organization, error)
	ListOrganizations() (*OrganizationList, error)
	CreateUser(orgID, username, email, firstName, lastName, password string) (*User, error)
	ListUsers() (*UserList, error)
	DeleteUser(userID string) error
	CreateProject(orgID, projectName string) (*Project, error)
	CreateOIDCWebApplication(orgID, projectID, appName string, redirectURIs []string) (*Application, error)
	GrantUserToProject(orgID, projectID, userID string) error
	GrantOrgToProject(orgID, projectID, grantedOrgID string) error

	// Branding methods
	SetOrgBranding(orgID string, config *BrandingConfig) error
	UploadOrgLogo(orgID string, logoPath string) error
	UpdateOrgColors(orgID string, config *BrandingConfig) error
	ActivateBranding(orgID string) error

	// Login Policy methods
	SetOrgLoginPolicy(orgID string, config *LoginPolicyConfig) error
	GetOrgLoginPolicy(orgID string) (*LoginPolicyConfig, error)
}
