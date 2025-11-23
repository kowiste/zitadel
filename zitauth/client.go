package zitauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
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

func (c *Client) CreateOrganization(orgName string) (string, error) {
	url := fmt.Sprintf("%s/management/v1/orgs", c.baseURL)
	jsonBody := []byte(fmt.Sprintf(`{"name": "%s"}`, orgName))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func (c *Client) ListOrganizations() (string, error) {
	url := fmt.Sprintf("%s/admin/v1/orgs/_search", c.baseURL)
	jsonBody := []byte(`{"queries": []}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func PrettyPrintJSON(jsonStr string) error {
	var obj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	fmt.Println(string(pretty))
	return nil
}
