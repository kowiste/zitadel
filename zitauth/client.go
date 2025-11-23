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

func (c *Client) CreateUser(username, email, firstName, lastName, password string) (*domain.User, error) {
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
		SetBody(body).
		Post("/management/v1/users/human/_import")

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	var user domain.User
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
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
