package domain

type Auth interface {
	CreateOrganization(orgName string) (*Organization, error)
	ListOrganizations() (*OrganizationList, error)
	CreateUser(username, email, firstName, lastName, password string) (*User, error)
	ListUsers() (*UserList, error)
	DeleteUser(userID string) error
}
