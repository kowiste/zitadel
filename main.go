package main

import (
	"fmt"
	"os"

	"zitadel/zitauth"
)

func main() {
	// 1. CREATE CLIENT FROM TOKEN FILE
	tokenPath := "./secrets/admin.pat"
	client, err := zitauth.NewClientFromToken("http://localhost:8080", tokenPath)
	if err != nil {
		fmt.Println("Error reading token file. Is ZITADEL fully started and the 'secrets' volume mounted correctly?")
		fmt.Printf("Details: %v\n", err)
		os.Exit(1)
	}

	// 2. CREATE ORGANIZATION
	orgName := "t2esgt" //+ time.Now().Format("150405") // Use timestamp to ensure uniqueness
	fmt.Println("\n=== Creating Organization ===")
	createResp, err := client.CreateOrganization(orgName)
	if err != nil {
		fmt.Printf("Error creating organization: %v\n", err)
		os.Exit(1)
	}
	if err := zitauth.PrettyPrintJSON(createResp); err != nil {
		fmt.Println("Response:", createResp)
	}

	// 3. LIST ORGANIZATIONS
	fmt.Println("\n=== Listing Organizations ===")
	listResp, err := client.ListOrganizations()
	if err != nil {
		fmt.Printf("Error listing organizations: %v\n", err)
		os.Exit(1)
	}
	if err := zitauth.PrettyPrintJSON(listResp); err != nil {
		fmt.Println("Response:", listResp)
	}
}
