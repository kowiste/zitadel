package main

import (
	"fmt"
	"os"
	"time"

	"zitadel/app"
	"zitadel/domain"
	"zitadel/zitauth"
)

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

	fmt.Println("\n" + repeat("━", 80))
	fmt.Println("MULTI-TENANT SETUP - ONE CLIENT ID FOR ALL TENANTS")
	fmt.Println(repeat("━", 80))

	// 2. ONE-TIME PLATFORM SETUP
	// This creates ONE shared application that ALL tenants will use
	fmt.Println("\n🏗️  STEP 1: Creating Platform Application (ONE TIME SETUP)")
	fmt.Println("This gives you the shared Client ID to hardcode in your frontend")

	platform, err := app.SetupPlatformApp(
		client,
		fmt.Sprintf("Platform-%s", suffix),
		"Shared Platform Project",
		"Shared Web App",
		[]string{
			"http://localhost:3000/auth/callback",
			"http://localhost:8080/auth/callback",
		},
	)
	if err != nil {
		fmt.Printf("Error during platform setup: %v\n", err)
		os.Exit(1)
	}

	// 3. CREATE MULTIPLE TENANTS
	// This simulates what happens when new customers/tenants register
	fmt.Println("\n\n🏢 STEP 2: Creating Tenant Organizations")
	fmt.Println("This happens when new customers register on your website")

	tenantNames := []string{"Tesla", "Pepsi", "Nike"}
	tenants := make([]*app.TenantSetup, 0)

	for _, tenantName := range tenantNames {
		fullName := fmt.Sprintf("%s-%s", tenantName, suffix)

		tenant, err := app.AddTenant(client, platform, fullName)
		if err != nil {
			fmt.Printf("❌ Error adding tenant %s: %v\n", tenantName, err)
			continue
		}

		tenants = append(tenants, tenant)

		// Create a test user in this tenant
		username := fmt.Sprintf("user@%s.com", tenantName)
		_, err = app.CreateTenantUser(
			client,
			platform,
			tenant,
			username,
			username,
			"Test",
			"User",
			"Password123!",
		)
		if err != nil {
			fmt.Printf("❌ Error creating user in %s: %v\n", tenantName, err)
		}
	}

	// 4. DEMONSTRATE BRANDING CUSTOMIZATION
	fmt.Println("\n\n" + repeat("━", 80))
	fmt.Println("🎨 STEP 3: Custom Branding (Optional)")
	fmt.Println(repeat("━", 80))
	fmt.Println("\nℹ️  Default branding from theme/branding.json was automatically applied to all tenants.")
	fmt.Println("Now let's demonstrate updating branding for specific tenants...")

	// Example: Update Tesla branding with custom colors
	if len(tenants) > 0 {
		fmt.Println("\n=== Customizing Tesla Branding ===")
		teslaColors := &domain.BrandingConfig{
			PrimaryColor:         "#e31937", // Tesla red
			BackgroundColor:      "#000000", // Black
			WarnColor:            "#ff3b5b",
			FontColor:            "#ffffff", // White
			PrimaryColorDark:     "#e31937",
			BackgroundColorDark:  "#000000",
			WarnColorDark:        "#ff3b5b",
			FontColorDark:        "#ffffff",
			HideLoginNameSuffix:  true,
			DisableWatermark:     true,
		}

		err = client.UpdateOrgColors(tenants[0].Organization.ID, teslaColors)
		if err != nil {
			fmt.Printf("⚠️  Error updating Tesla branding: %v\n", err)
		} else {
			fmt.Println("✓ Tesla branding colors updated to red/black theme")
			err = client.ActivateBranding(tenants[0].Organization.ID)
			if err != nil {
				fmt.Printf("⚠️  Error activating Tesla branding: %v\n", err)
			} else {
				fmt.Println("✓ Tesla branding activated")
			}
		}
	}

	fmt.Println("\n💡 You can update branding anytime using:")
	fmt.Println("   • client.UpdateOrgColors(orgID, config) - Update colors")
	fmt.Println("   • client.UploadOrgLogo(orgID, logoPath) - Upload new logo")
	fmt.Println("   • client.ActivateBranding(orgID) - Activate changes")

	// 5. SHOW LOGIN URLS FOR ALL TENANTS
	fmt.Println("\n\n" + repeat("━", 80))
	fmt.Println("🔗 STEP 4: Login URLs for Testing")
	fmt.Println(repeat("━", 80))
	fmt.Println("\nAll tenants use the SAME Client ID, just different org parameter:")
	fmt.Printf("\n📋 Shared Client ID: %s\n", platform.SharedApp.ClientID)

	for i, tenant := range tenants {
		fmt.Printf("\n%d. %s Login URL:\n", i+1, tenant.Organization.Name)
		fmt.Println(repeat("─", 80))
		loginURL := fmt.Sprintf("http://localhost:8080/oauth/v2/authorize?client_id=%s&redirect_uri=http://localhost:3000/auth/callback&response_type=code&scope=openid%%20email%%20profile%%20urn:zitadel:iam:org:domain:primary:%s",
			platform.SharedApp.ClientID,
			tenant.Organization.PrimaryDomain)
		fmt.Println(loginURL)
	}

	// 6. SUMMARY
	fmt.Println("\n\n" + repeat("━", 80))
	fmt.Println("✅ SETUP COMPLETE - SUMMARY")
	fmt.Println(repeat("━", 80))

	fmt.Printf("\n📋 Shared Client ID: %s\n", platform.SharedApp.ClientID)
	fmt.Printf("📋 Client Secret: %s\n", platform.SharedApp.ClientSecret)
	fmt.Printf("\n🏢 Created %d tenant organizations:\n", len(tenants))
	for _, tenant := range tenants {
		fmt.Printf("   • %s (Domain: %s)\n", tenant.Organization.Name, tenant.Organization.PrimaryDomain)
	}
	fmt.Println("\n🎨 Branding:")
	fmt.Println("   • Default branding applied from theme/branding.json")
	fmt.Println("   • ZITADEL watermark removed")
	fmt.Println("   • Login name suffix hidden")
	fmt.Println("   • Custom logo from theme/logo.svg (if exists)")

	fmt.Println("\n" + repeat("━", 80))
	fmt.Println(repeat("━", 80))
	fmt.Println("\n✅ Frontend KNOWS:")
	fmt.Printf("   • Client ID = '%s' (hardcoded)\n", platform.SharedApp.ClientID)
	fmt.Println("   • User provides: company name (e.g., 'tesla')")
	fmt.Println("\n✅ Frontend DOES:")
	fmt.Println("   • Constructs OAuth URL with Client ID + org scope")
	fmt.Println("   • Redirects to ZITADEL")
	fmt.Println("   • NO backend API call needed!")
	fmt.Println("\n✅ After login:")
	fmt.Println("   • User gets token")
	fmt.Println("   • Frontend makes API requests to your backend")
	fmt.Println("   • Backend validates token and serves data")

	fmt.Println("\n" + repeat("━", 80))
	fmt.Println("DATABASE STORAGE (What to save per tenant)")
	fmt.Println(repeat("━", 80))
	fmt.Println("\nFor each tenant, store:")
	fmt.Println("   • tenant_id (your internal ID)")
	fmt.Println("   • tenant_name (e.g., 'tesla')")
	fmt.Println("   • zitadel_org_id (from ZITADEL)")
	fmt.Println("   • zitadel_org_domain (from ZITADEL)")
	fmt.Println("\nDO NOT store Client ID per tenant - it's the same for all!")

	fmt.Println("\n" + repeat("━", 80))
	fmt.Println("NEXT STEPS")
	fmt.Println(repeat("━", 80))
	fmt.Println("\n1. Copy the Client ID above and hardcode it in your frontend")
	fmt.Println("2. Customize default branding (optional):")
	fmt.Println("   → Edit theme/branding.json with your colors")
	fmt.Println("   → Add theme/logo.svg with your logo")
	fmt.Println("   → These will be applied automatically to new tenants")
	fmt.Println("3. When new tenants register:")
	fmt.Println("   → Call app.AddTenant() to create ZITADEL org")
	fmt.Println("   → Default branding is applied automatically")
	fmt.Println("   → Save tenant-to-org mapping in your database")
	fmt.Println("   → Update tenant branding later if needed")
	fmt.Println("4. Frontend login flow:")
	fmt.Println("   → User enters company name")
	fmt.Println("   → Frontend looks up org domain (or uses company name)")
	fmt.Println("   → Redirect to ZITADEL with org scope")

	fmt.Println("\n" + repeat("━", 80))
	fmt.Println("🧪 TEST LOGIN URLS (Copy & Paste in Browser)")
	fmt.Println(repeat("━", 80))
	fmt.Println("\nAll use the SAME Client ID, different org domains:")
	for i, tenant := range tenants {
		fmt.Printf("\n%d. %s:\n", i+1, tenant.Organization.Name)
		loginURL := fmt.Sprintf("http://localhost:8080/oauth/v2/authorize?client_id=%s&redirect_uri=http://localhost:3000/auth/callback&response_type=code&scope=openid%%20email%%20profile%%20urn:zitadel:iam:org:domain:primary:%s",
			platform.SharedApp.ClientID,
			tenant.Organization.PrimaryDomain)
		fmt.Println(loginURL)
	}

	fmt.Println("\n" + repeat("━", 80) + "\n")
}

// Helper function to repeat strings
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
