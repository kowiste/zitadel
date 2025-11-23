package main

import (
	"encoding/json"
	"fmt"
	"os"

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

	// 2. CREATE ORGANIZATION
	orgName := "TestOrg"
	fmt.Println("\n=== Creating Organization ===")
	org, err := client.CreateOrganization(orgName)
	if err != nil {
		fmt.Printf("Error creating organization: %v\n", err)
		os.Exit(1)
	}
	prettyPrint(org)
	fmt.Printf("\nCreated org with ID: %s\n", org.ID)

	// 3. CREATE USER IN THE ORGANIZATION
	fmt.Println("\n=== Creating User in Organization ===")
	user, err := client.CreateUser(
		"john.doe",
		"john.doe@example.com",
		"John",
		"Doe",
		"Password123!",
	)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	} else {
		prettyPrint(user)
	}

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
