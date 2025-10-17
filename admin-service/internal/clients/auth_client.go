package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// AuthClient communicates with the upstream service responsible for admin authentication flows.
type AuthClient interface {
	RegisterAdmin(ctx context.Context, payload interface{}) (*http.Response, error)
	LoginAdmin(ctx context.Context, payload interface{}) (*http.Response, error)
	ForgotPassword(ctx context.Context, payload interface{}) (*http.Response, error)
	ResetPassword(ctx context.Context, payload interface{}) (*http.Response, error)
	LogoutAdmin(ctx context.Context, token string) (*http.Response, error)
}

type authClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAuthClient wires an HTTP client for admin auth endpoints.
func NewAuthClient(baseURL string) AuthClient {
	return &authClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *authClient) RegisterAdmin(ctx context.Context, payload interface{}) (*http.Response, error) {
	return c.post(ctx, "/api/admin/register", payload)
}

func (c *authClient) LoginAdmin(ctx context.Context, payload interface{}) (*http.Response, error) {
	return c.post(ctx, "/api/admin/login", payload)
}

func (c *authClient) ForgotPassword(ctx context.Context, payload interface{}) (*http.Response, error) {
	return c.post(ctx, "/api/admin/forgot-password", payload)
}

func (c *authClient) ResetPassword(ctx context.Context, payload interface{}) (*http.Response, error) {
	return c.post(ctx, "/api/admin/reset-password", payload)
}

func (c *authClient) LogoutAdmin(ctx context.Context, token string) (*http.Response, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return nil, fmt.Errorf("auth client: base URL is not configured")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("auth client: token is required for logout")
	}

	url := strings.TrimRight(c.baseURL, "/") + "/api/admin/logout"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("auth client: create logout request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth client: perform logout: %w", err)
	}

	return resp, nil
}

func (c *authClient) post(ctx context.Context, path string, payload interface{}) (*http.Response, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return nil, fmt.Errorf("auth client: base URL is not configured")
	}

	url := strings.TrimRight(c.baseURL, "/") + path

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("auth client: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("auth client: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth client: perform request: %w", err)
	}

	return resp, nil
}
