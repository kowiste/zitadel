package domain

type Auth interface {
	CreateOrganization(orgName string) (*Organization, error)
	ListOrganizations() (*OrganizationList, error)
	CreateUser(orgID, username, email, firstName, lastName, password string) (*User, error)
	ListUsers() (*UserList, error)
	DeleteUser(userID string) error
	CreateProject(orgID, projectName string) (*Project, error)
	CreateOIDCWebApplication(orgID, projectID, appName string, redirectURIs []string) (*Application, error)
	GrantUserToProject(orgID, projectID, userID string) error
	GrantOrgToProject(orgID, projectID, grantedOrgID string) error
}
