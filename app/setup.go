package app

import (
	"fmt"

	"zitadel/domain"
)

type SetupResult struct {
	Organization *domain.Organization
	Project      *domain.Project
	Application  *domain.Application
}

// SetupOrgWithApp orchestrates the creation of an organization, project, and web application
func SetupOrgWithApp(client domain.Auth, orgName, projectName, appName string, redirectURIs []string) (*SetupResult, error) {
	// 1. Create Organization
	fmt.Printf("\n=== Creating Organization: %s ===\n", orgName)
	org, err := client.CreateOrganization(orgName)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	fmt.Printf("✓ Organization created with ID: %s\n", org.ID)

	// 2. Create Project in the Organization
	fmt.Printf("\n=== Creating Project: %s ===\n", projectName)
	project, err := client.CreateProject(org.ID, projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	fmt.Printf("✓ Project created with ID: %s\n", project.ID)

	// 3. Create OIDC Web Application in the Project
	fmt.Printf("\n=== Creating Application: %s ===\n", appName)
	app, err := client.CreateOIDCWebApplication(org.ID, project.ID, appName, redirectURIs)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}
	fmt.Printf("✓ Application created with ID: %s\n", app.ID)
	fmt.Printf("✓ Client ID: %s\n", app.ClientID)
	fmt.Printf("✓ Client Secret: %s\n", app.ClientSecret)

	return &SetupResult{
		Organization: org,
		Project:      project,
		Application:  app,
	}, nil
}

// CreateUser creates a user and automatically grants them access to the project
func CreateUser(client domain.Auth, orgID, projectID, username, email, firstName, lastName, password string) (*domain.User, error) {
	// 1. Create the user
	fmt.Printf("Creating user: %s\n", username)
	user, err := client.CreateUser(orgID, username, email, firstName, lastName, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Printf("✓ User created with ID: %s\n", user.ID)

	// 2. Grant user access to the project
	fmt.Println("Granting user access to project...")
	err = client.GrantUserToProject(orgID, projectID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to grant user to project: %w", err)
	}
	fmt.Println("✓ User granted access to project")

	return user, nil
}
