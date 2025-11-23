package zitauth

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"zitadel/domain"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *resty.Client
}

func NewClient(baseURL, token string) *Client {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetBaseURL(baseURL)
	client.SetHeader("Authorization", "Bearer "+token)
	client.SetHeader("Content-Type", "application/json")

	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: client,
	}
}

func NewClientFromToken(baseURL, tokenPath string) (*Client, error) {
	fmt.Println("Waiting for token file...")
	var tokenBytes []byte
	var err error

	// Simple retry loop to wait for Docker/ZITADEL to finish writing the file
	for i := 0; i < 20; i++ {
		tokenBytes, err = os.ReadFile(tokenPath)
		// Check for error and ensure the file is not empty
		if err == nil && len(strings.TrimSpace(string(tokenBytes))) > 0 {
			break
		}
		time.Sleep(2 * time.Second)
		fmt.Print(".")
	}
	fmt.Println()

	if err != nil {
		return nil, fmt.Errorf("error reading token file: %w", err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	fmt.Println("Token found successfully!")

	return NewClient(baseURL, token), nil
}

func (c *Client) CreateOrganization(orgName string) (*domain.Organization, error) {
	body := map[string]string{"name": orgName}

	resp, err := c.httpClient.R().
		SetBody(body).
		Post("/management/v1/orgs")

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	var org domain.Organization
	if err := json.Unmarshal(resp.Body(), &org); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &org, nil
}

func (c *Client) ListOrganizations() (*domain.OrganizationList, error) {
	body := map[string][]interface{}{"queries": {}}

	resp, err := c.httpClient.R().
		SetBody(body).
		Post("/admin/v1/orgs/_search")

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	var orgList domain.OrganizationList
	if err := json.Unmarshal(resp.Body(), &orgList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &orgList, nil
}

func (c *Client) CreateUser(orgID, username, email, firstName, lastName, password string) (*domain.User, error) {
	body := map[string]interface{}{
		"userName": username,
		"profile": map[string]string{
			"firstName": firstName,
			"lastName":  lastName,
		},
		"email": map[string]interface{}{
			"email":           email,
			"isEmailVerified": true,
		},
		"password": password,
	}

	resp, err := c.httpClient.R().
		SetHeader("x-zitadel-orgid", orgID).
		SetBody(body).
		Post("/management/v1/users/human/_import")

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	// Parse the response which has userId at the top level
	var createResp struct {
		UserID string `json:"userId"`
	}
	if err := json.Unmarshal(resp.Body(), &createResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Return a User with the ID populated
	return &domain.User{
		ID: createResp.UserID,
	}, nil
}

func (c *Client) ListUsers() (*domain.UserList, error) {
	body := map[string][]interface{}{"queries": {}}

	resp, err := c.httpClient.R().
		SetBody(body).
		Post("/management/v1/users/_search")

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	var userList domain.UserList
	if err := json.Unmarshal(resp.Body(), &userList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &userList, nil
}

func (c *Client) DeleteUser(userID string) error {
	_, err := c.httpClient.R().
		Delete(fmt.Sprintf("/management/v1/users/%s", userID))

	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	return nil
}

func (c *Client) CreateProject(orgID, projectName string) (*domain.Project, error) {
	body := map[string]interface{}{
		"name":                  projectName,
		"projectRoleAssertion":  false,
		"projectRoleCheck":      false,
		"hasProjectCheck":       false,
		"privateLabelingSetting": "PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY",
	}

	resp, err := c.httpClient.R().
		SetHeader("x-zitadel-orgid", orgID).
		SetBody(body).
		Post("/management/v1/projects")

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	var project domain.Project
	if err := json.Unmarshal(resp.Body(), &project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &project, nil
}

func (c *Client) CreateOIDCWebApplication(orgID, projectID, appName string, redirectURIs []string) (*domain.Application, error) {
	body := map[string]interface{}{
		"name":                      appName,
		"redirectUris":              redirectURIs,
		"responseTypes":             []string{"OIDC_RESPONSE_TYPE_CODE"},
		"grantTypes":                []string{"OIDC_GRANT_TYPE_AUTHORIZATION_CODE", "OIDC_GRANT_TYPE_REFRESH_TOKEN"},
		"appType":                   "OIDC_APP_TYPE_WEB",
		"authMethodType":            "OIDC_AUTH_METHOD_TYPE_BASIC",
		"postLogoutRedirectUris":    redirectURIs,
		"version":                   "OIDC_VERSION_1_0",
		"devMode":                   false,
		"accessTokenType":           "OIDC_TOKEN_TYPE_BEARER",
		"accessTokenRoleAssertion":  true,
		"idTokenRoleAssertion":      true,
		"idTokenUserinfoAssertion":  true,
		"clockSkew":                 "0s",
		"skipNativeAppSuccessPage":  false,
	}

	resp, err := c.httpClient.R().
		SetHeader("x-zitadel-orgid", orgID).
		SetBody(body).
		Post(fmt.Sprintf("/management/v1/projects/%s/apps/oidc", projectID))

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	var app domain.Application
	if err := json.Unmarshal(resp.Body(), &app); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &app, nil
}

func (c *Client) GrantUserToProject(orgID, projectID, userID string) error {
	body := map[string]interface{}{
		"projectId": projectID,
		"roleKeys":  []string{},
	}

	resp, err := c.httpClient.R().
		SetHeader("x-zitadel-orgid", orgID).
		SetBody(body).
		Post(fmt.Sprintf("/management/v1/users/%s/grants", userID))

	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

func (c *Client) GrantOrgToProject(orgID, projectID, grantedOrgID string) error {
	body := map[string]interface{}{
		"grantedOrgId": grantedOrgID,
		"roleKeys":     []string{},
	}

	resp, err := c.httpClient.R().
		SetHeader("x-zitadel-orgid", orgID).
		SetBody(body).
		Post(fmt.Sprintf("/management/v1/projects/%s/grants", projectID))

	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
