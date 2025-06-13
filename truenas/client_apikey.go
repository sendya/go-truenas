package truenas

import (
	"context"
	"fmt"
	"time"
)

// APIKeyClient provides methods for API key management
type APIKeyClient struct {
	client *Client
}

// NewAPIKeyClient creates a new API key client
func NewAPIKeyClient(client *Client) *APIKeyClient {
	return &APIKeyClient{client: client}
}

// APIKey represents an API key
type APIKey struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
}

// APIKeyCreateRequest represents parameters for api_key.create
type APIKeyCreateRequest struct {
	Name string `json:"name"`
}

// APIKeyUpdateRequest represents parameters for api_key.update
type APIKeyUpdateRequest struct {
	Name  *string `json:"name,omitempty"`
	Reset *bool   `json:"reset,omitempty"`
}

// List returns all API keys
func (a *APIKeyClient) List(ctx context.Context) ([]APIKey, error) {
	var result []APIKey
	err := a.client.Call(ctx, "api_key.query", []any{}, &result)
	return result, err
}

// Get returns a specific API key by ID
func (a *APIKeyClient) Get(ctx context.Context, id int) (*APIKey, error) {
	var result []APIKey
	err := a.client.Call(ctx, "api_key.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("api_key", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new API key
func (a *APIKeyClient) Create(ctx context.Context, name string) (*APIKey, error) {
	var result APIKey
	req := APIKeyCreateRequest{Name: name}
	err := a.client.Call(ctx, "api_key.create", []any{req}, &result)
	return &result, err
}

// Update updates an existing API key
func (a *APIKeyClient) Update(ctx context.Context, id int, req *APIKeyUpdateRequest) (*APIKey, error) {
	var result APIKey
	err := a.client.Call(ctx, "api_key.update", []any{id, *req}, &result)
	return &result, err
}

// UpdateName updates the name of an API key
func (a *APIKeyClient) UpdateName(ctx context.Context, id int, name string) (*APIKey, error) {
	req := &APIKeyUpdateRequest{Name: &name}
	return a.Update(ctx, id, req)
}

// Reset regenerates an API key (creates new key value)
func (a *APIKeyClient) Reset(ctx context.Context, id int) (*APIKey, error) {
	reset := true
	req := &APIKeyUpdateRequest{Reset: &reset}
	return a.Update(ctx, id, req)
}

// Delete deletes an API key
func (a *APIKeyClient) Delete(ctx context.Context, id int) error {
	return a.client.Call(ctx, "api_key.delete", []any{id}, nil)
}
