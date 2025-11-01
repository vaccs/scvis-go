package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// GitHubClient represents a GitHub API client
type GitHubClient struct {
	baseURL    string
	token      string
	owner      string
	repo       string
	httpClient *http.Client
}

// IssueRequest represents a GitHub issue creation request
type IssueRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Labels    []string `json:"labels,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
}

// IssueResponse represents a GitHub issue creation response
type IssueResponse struct {
	ID        int    `json:"id"`
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	HTMLURL   string `json:"html_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(owner, repo string) *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		// For development/testing, you might want to log this
		// but don't fail if no token is provided
	}

	return &GitHubClient{
		baseURL: "https://api.github.com",
		token:   token,
		owner:   owner,
		repo:    repo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateIssue creates a new issue in the GitHub repository
func (c *GitHubClient) CreateIssue(issue *IssueRequest) (*IssueResponse, error) {
	if c.token == "" {
		return nil, fmt.Errorf("GitHub token not configured")
	}

	// Prepare the request body
	jsonData, err := json.Marshal(issue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal issue request: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/repos/%s/%s/issues", c.baseURL, c.owner, c.repo)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "SCVis-Go/1.0")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var issueResponse IssueResponse
	if err := json.Unmarshal(body, &issueResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &issueResponse, nil
}

// GetRepository returns the owner and repository name
func (c *GitHubClient) GetRepository() (string, string) {
	return c.owner, c.repo
}

// IsConfigured returns true if the GitHub client is properly configured
func (c *GitHubClient) IsConfigured() bool {
	return c.token != ""
}

// TestConnection tests the GitHub API connection
func (c *GitHubClient) TestConnection() error {
	if c.token == "" {
		return fmt.Errorf("GitHub token not configured")
	}

	// Test with a simple repository info request
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, c.owner, c.repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create test request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "SCVis-Go/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API connection test failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
