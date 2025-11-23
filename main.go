package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Define the known hardcoded client credentials, matching the docker-compose setup
const (
	clientID     = "go-client-id"
	clientSecret = "my-super-secure-secret-2025"
	zitadelURL   = "http://localhost:8080"
)

// TokenResponse structure to unmarshal the OAuth token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// getAccessToken performs the OAuth 2.0 Client Credentials Grant flow.
func getAccessToken() (string, error) {
	fmt.Println("--- 1. Performing Client Credentials Flow to get Access Token ---")

	tokenEndpoint := zitadelURL + "/oauth/v2/token"

	// Form data for the Client Credentials Grant
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	// Request the necessary scope for management API access
	formData.Set("scope", "openid profile email urn:zitadel:iam:org:project:id:zitadel:aud")

	// Set a timeout for the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}
	
	fmt.Printf("Successfully obtained Access Token (Expires in %d seconds).\n", tokenResponse.ExpiresIn)
	return tokenResponse.AccessToken, nil
}

// createOrg uses the retrieved token to call the ZITADEL Management API.
func createOrg(token string) {
	fmt.Println("\n--- 2. Calling Management API to Create Organization ---")
	
	orgsEndpoint := zitadelURL + "/management/v1/orgs"
	
	// Ensure the organization name is unique
	orgName := "Client-Creds-Org-" + time.Now().Format("150405")
	jsonBody := []byte(fmt.Sprintf(`{"name": "%s"}`, orgName))

	req, err := http.NewRequest("POST", orgsEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	// Set the Authorization header with the newly obtained token
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print result
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("----------------------------------------")
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response:", string(body))
	fmt.Println("----------------------------------------")
}

func main() {
	// Step 1: Get the temporary Access Token
	token, err := getAccessToken()
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		os.Exit(1) 
	}

	// Step 2: Use the token to create the Organization
	createOrg(token)
}