package truenas

import (
	"context"
)

// AuthClient provides methods for authentication and user management
type AuthClient struct {
	client *Client
}

// NewAuthClient creates a new auth client
func NewAuthClient(client *Client) *AuthClient {
	return &AuthClient{client: client}
}

// LoginRequest represents parameters for auth.login
type LoginRequest struct {
	Username string
	Password string
}

// LoginWithAPIKeyRequest represents parameters for auth.login_with_api_key
type LoginWithAPIKeyRequest struct {
	APIKey string
}

// GenerateTokenRequest represents parameters for auth.generate_token
type GenerateTokenRequest struct {
	TTL        int `json:"ttl,omitempty"`
	Attributes any `json:"attributes,omitempty"`
}

// TokenResponse represents the response from auth.generate_token
type TokenResponse struct {
	Token string `json:"token"`
}

// CheckPasswordRequest represents parameters for auth.check_password
type CheckPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login authenticates with username and password
func (a *AuthClient) Login(ctx context.Context, username, password string) (bool, error) {
	var result bool
	err := a.client.Call(ctx, "auth.login", []any{username, password}, &result)
	return result, err
}

// LoginWithAPIKey authenticates with an API key
func (a *AuthClient) LoginWithAPIKey(ctx context.Context, apiKey string) (bool, error) {
	var result bool
	err := a.client.Call(ctx, "auth.login_with_api_key", []any{apiKey}, &result)
	return result, err
}

// Logout ends the current session
func (a *AuthClient) Logout(ctx context.Context) error {
	return a.client.Call(ctx, "auth.logout", []any{}, nil)
}

// CheckPassword validates a user's password
func (a *AuthClient) CheckPassword(ctx context.Context, username, password string) (bool, error) {
	var result bool
	err := a.client.Call(ctx, "auth.check_password", []any{username, password}, &result)
	return result, err
}

// GenerateToken creates a new authentication token
func (a *AuthClient) GenerateToken(ctx context.Context, req GenerateTokenRequest) (*TokenResponse, error) {
	var result TokenResponse
	params := []any{}
	if req.TTL > 0 || req.Attributes != nil {
		tokenParams := map[string]any{}
		if req.TTL > 0 {
			tokenParams["ttl"] = req.TTL
		}
		if req.Attributes != nil {
			tokenParams["attributes"] = req.Attributes
		}
		params = append(params, tokenParams)
	}
	err := a.client.Call(ctx, "auth.generate_token", params, &result)
	return &result, err
}
