package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"zitadel/app"
	"zitadel/domain"
	"zitadel/zitauth"
)

func prettyPrint(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func main() {
	// 1. CREATE CLIENT FROM TOKEN FILE
	tokenPath := "./secrets/admin.pat"
	var client domain.Auth
	client, err := zitauth.NewClientFromToken("http://localhost:8080", tokenPath)
	if err != nil {
		fmt.Println("Error reading token file. Is ZITADEL fully started and the 'secrets' volume mounted correctly?")
		fmt.Printf("Details: %v\n", err)
		os.Exit(1)
	}

	// Generate unique suffix
	suffix := time.Now().Format("150405")

	// 2. SETUP ORGANIZATION WITH PROJECT AND APPLICATION
	orgName := fmt.Sprintf("TestOrg-%s", suffix)
	projectName := "Main Project"
	appName := "Web Application"
	redirectURIs := []string{
		"http://localhost:3000/auth/callback",
		"http://localhost:8080/auth/callback",
	}

	setup, err := app.SetupOrgWithApp(client, orgName, projectName, appName, redirectURIs)
	if err != nil {
		fmt.Printf("Error during setup: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== Setup Complete ===")
	fmt.Println("\nOrganization:")
	prettyPrint(setup.Organization)
	fmt.Println("\nProject:")
	prettyPrint(setup.Project)
	fmt.Println("\nApplication:")
	prettyPrint(setup.Application)

	// 3. CREATE USER WITH PROJECT ACCESS
	fmt.Println("\n=== Creating User with Project Access ===")
	username := fmt.Sprintf("john.doe-%s", suffix)
	email := fmt.Sprintf("john.doe-%s@example.com", suffix)
	user, err := app.CreateUser(
		client,
		setup.Organization.ID,
		setup.Project.ID,
		username,
		email,
		"John",
		"Doe",
		"Password123!",
	)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}
	prettyPrint(user)

}

func lists(client domain.Auth) {
	// 4. LIST ORGANIZATIONS
	fmt.Println("\n=== Listing Organizations ===")
	orgList, err := client.ListOrganizations()
	if err != nil {
		fmt.Printf("Error listing organizations: %v\n", err)
		os.Exit(1)
	}
	prettyPrint(orgList.Organizations)

	// 5. LIST USERS
	fmt.Println("\n=== Listing Users ===")
	userList, err := client.ListUsers()
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		os.Exit(1)
	}
	prettyPrint(userList.Users)
}
