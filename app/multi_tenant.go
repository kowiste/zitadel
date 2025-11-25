package app

import (
	"encoding/json"
	"fmt"
	"os"
	"zitadel/domain"
)

// PlatformSetup holds the shared platform application details
type PlatformSetup struct {
	PlatformOrg *domain.Organization
	Project     *domain.Project
	SharedApp   *domain.Application
}

// SetupPlatformApp creates a platform organization with a shared application
// This should be run ONCE to get the shared Client ID that all tenants will use
func SetupPlatformApp(
	client domain.Auth,
	platformName string,
	projectName string,
	appName string,
	redirectURIs []string,
) (*PlatformSetup, error) {

	// 1. Create Platform Organization
	fmt.Printf("\n=== Creating Platform Organization: %s ===\n", platformName)
	platformOrg, err := client.CreateOrganization(platformName)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform org: %w", err)
	}
	fmt.Printf("✓ Platform Org ID: %s\n", platformOrg.ID)

	// 2. Create Platform Project
	fmt.Printf("\n=== Creating Platform Project: %s ===\n", projectName)
	project, err := client.CreateProject(platformOrg.ID, projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	fmt.Printf("✓ Platform Project ID: %s\n", project.ID)

	// 3. Create Shared Application
	fmt.Printf("\n=== Creating Shared Application: %s ===\n", appName)
	app, err := client.CreateOIDCWebApplication(
		platformOrg.ID,
		project.ID,
		appName,
		redirectURIs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}
	fmt.Printf("✓ Application created\n")
	fmt.Printf("✓ Shared Client ID: %s\n", app.ClientID)
	fmt.Printf("✓ Client Secret: %s\n", app.ClientSecret)

	fmt.Println("\n" + repeat("=", 60))
	fmt.Println("IMPORTANT: Save these credentials!")
	fmt.Println(repeat("=", 60))
	fmt.Printf("Shared Client ID: %s\n", app.ClientID)
	fmt.Printf("Client Secret: %s\n", app.ClientSecret)
	fmt.Println("\nHardcode the Client ID in your frontend application.")
	fmt.Println("All tenants will use this same Client ID.")
	fmt.Println(repeat("=", 60))

	return &PlatformSetup{
		PlatformOrg: platformOrg,
		Project:     project,
		SharedApp:   app,
	}, nil
}

// TenantSetup holds the tenant organization details
type TenantSetup struct {
	Organization *domain.Organization
}

// LoadBrandingConfig loads the default branding configuration from theme/branding.json
func LoadBrandingConfig() (*domain.BrandingConfig, error) {
	data, err := os.ReadFile("./theme/branding.json")
	if err != nil {
		// If file doesn't exist, return default config
		fmt.Println("⚠️  No branding.json found, using default config")
		return domain.DefaultBrandingConfig(), nil
	}

	var config domain.BrandingConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse branding.json: %w", err)
	}

	return &config, nil
}

// AddTenant creates a new tenant organization and grants access to the platform
// Run this when a new customer/tenant registers
func AddTenant(
	client domain.Auth,
	platform *PlatformSetup,
	tenantName string,
) (*TenantSetup, error) {

	// 1. Create Tenant Organization
	fmt.Printf("\n=== Creating Tenant Organization: %s ===\n", tenantName)
	tenantOrg, err := client.CreateOrganization(tenantName)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant org: %w", err)
	}
	fmt.Printf("✓ Tenant Org ID: %s\n", tenantOrg.ID)
	fmt.Printf("✓ Primary Domain: %s\n", tenantOrg.PrimaryDomain)

	// 2. Grant Project Access to Tenant
	fmt.Println("\nGranting project access to tenant...")
	err = client.GrantOrgToProject(
		platform.PlatformOrg.ID,
		platform.Project.ID,
		tenantOrg.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to grant project: %w", err)
	}
	fmt.Println("✓ Project granted to tenant")

	// 3. Apply Default Branding
	fmt.Println("\nApplying default branding...")
	brandingConfig, err := LoadBrandingConfig()
	if err != nil {
		fmt.Printf("⚠️  Failed to load branding config: %v\n", err)
	} else {
		// Set branding colors and settings
		err = client.SetOrgBranding(tenantOrg.ID, brandingConfig)
		if err != nil {
			fmt.Printf("⚠️  Failed to set branding: %v\n", err)
		} else {
			fmt.Println("✓ Branding colors and settings applied")
		}

		// Upload logo if it exists
		logoPath := "./theme/logo.svg"
		if _, err := os.Stat(logoPath); err == nil {
			err = client.UploadOrgLogo(tenantOrg.ID, logoPath)
			if err != nil {
				fmt.Printf("⚠️  Failed to upload logo: %v\n", err)
			} else {
				fmt.Println("✓ Logo uploaded (both banner logo and icon)")
			}
		} else {
			fmt.Println("⚠️  No logo.svg found in theme directory")
		}

		// Activate branding
		err = client.ActivateBranding(tenantOrg.ID)
		if err != nil {
			fmt.Printf("⚠️  Failed to activate branding: %v\n", err)
		} else {
			fmt.Println("✓ Branding activated")
		}
	}

	fmt.Println("\n" + repeat("-", 60))
	fmt.Printf("Tenant '%s' can now use the shared Client ID\n", tenantName)
	fmt.Printf("Organization ID: %s\n", tenantOrg.ID)
	fmt.Printf("Primary Domain: %s\n", tenantOrg.PrimaryDomain)
	fmt.Println(repeat("-", 60))

	return &TenantSetup{
		Organization: tenantOrg,
	}, nil
}

// CreateTenantUser creates a user in a tenant org and grants access to the platform
func CreateTenantUser(
	client domain.Auth,
	platform *PlatformSetup,
	tenant *TenantSetup,
	username, email, firstName, lastName, password string,
) (*domain.User, error) {

	// 1. Create user in tenant organization
	fmt.Printf("\n=== Creating User: %s ===\n", username)
	user, err := client.CreateUser(
		tenant.Organization.ID,
		username,
		email,
		firstName,
		lastName,
		password,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Printf("✓ User created with ID: %s\n", user.ID)

	// 2. Grant user access to platform project
	fmt.Println("Granting user access to platform project...")
	err = client.GrantUserToProject(
		tenant.Organization.ID,
		platform.Project.ID,
		user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to grant user to project: %w", err)
	}
	fmt.Println("✓ User granted access to platform")

	return user, nil
}

// PrintOAuthURLExample prints example OAuth URLs for frontend integration
func PrintOAuthURLExample(platform *PlatformSetup, tenant *TenantSetup, zitadelURL string) {
	fmt.Println("\n" + repeat("=", 60))
	fmt.Printf("FRONTEND INTEGRATION - Tenant: %s\n", tenant.Organization.Name)
	fmt.Println(repeat("=", 60))

	fmt.Println("\nOption 1: Using Organization ID")
	fmt.Println("--------------------------------")
	fmt.Printf("const clientId = '%s'; // Hardcoded!\n", platform.SharedApp.ClientID)
	fmt.Printf("const orgId = '%s'; // From user input or database\n", tenant.Organization.ID)
	fmt.Println("const redirectUri = 'http://localhost:3000/auth/callback';")
	fmt.Println()
	fmt.Println("const authUrl = `" + zitadelURL + "/oauth/v2/authorize")
	fmt.Println("  ?client_id=${clientId}")
	fmt.Println("  &redirect_uri=${redirectUri}")
	fmt.Println("  &response_type=code")
	fmt.Println("  &scope=openid email profile urn:zitadel:iam:org:id:${orgId}")
	fmt.Println("  &state=YOUR_RANDOM_STATE`;")

	fmt.Println("\nOption 2: Using Organization Domain (Recommended)")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("const clientId = '%s'; // Hardcoded!\n", platform.SharedApp.ClientID)
	fmt.Printf("const orgDomain = '%s'; // From user input\n", tenant.Organization.PrimaryDomain)
	fmt.Println("const redirectUri = 'http://localhost:3000/auth/callback';")
	fmt.Println()
	fmt.Println("const authUrl = `" + zitadelURL + "/oauth/v2/authorize")
	fmt.Println("  ?client_id=${clientId}")
	fmt.Println("  &redirect_uri=${redirectUri}")
	fmt.Println("  &response_type=code")
	fmt.Println("  &scope=openid email profile urn:zitadel:iam:org:domain:primary:${orgDomain}")
	fmt.Println("  &state=YOUR_RANDOM_STATE`;")

	fmt.Println("\n" + repeat("=", 60))
}

// Helper function to repeat strings (Go doesn't have built-in)
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
