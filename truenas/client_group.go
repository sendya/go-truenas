package truenas

import (
	"context"
	"fmt"
)

// GroupClient provides methods for group management
type GroupClient struct {
	client *Client
}

// NewGroupClient creates a new group client
func NewGroupClient(client *Client) *GroupClient {
	return &GroupClient{client: client}
}

// Group represents a system group
type Group struct {
	ID           int      `json:"id"`
	GID          int      `json:"gid"`
	Name         string   `json:"name"`
	Builtin      bool     `json:"builtin"`
	Sudo         bool     `json:"sudo"`
	SudoNoPasswd bool     `json:"sudo_nopasswd"`
	SudoCommands []string `json:"sudo_commands"`
	Smb          bool     `json:"smb"`
	Users        []int    `json:"users"`
	Local        bool     `json:"local"`
}

// GroupCreateRequest represents parameters for group.create
type GroupCreateRequest struct {
	GID               int      `json:"gid,omitempty"`
	Name              string   `json:"name"`
	Smb               bool     `json:"smb,omitempty"`
	Sudo              bool     `json:"sudo,omitempty"`
	SudoNoPasswd      bool     `json:"sudo_nopasswd,omitempty"`
	SudoCommands      []string `json:"sudo_commands,omitempty"`
	AllowDuplicateGID bool     `json:"allow_duplicate_gid,omitempty"`
	Users             []int    `json:"users,omitempty"`
}

// GroupUpdateRequest represents parameters for group.update
type GroupUpdateRequest struct {
	GID               int      `json:"gid,omitempty"`
	Name              string   `json:"name,omitempty"`
	Smb               bool     `json:"smb,omitempty"`
	Sudo              bool     `json:"sudo,omitempty"`
	SudoNoPasswd      bool     `json:"sudo_nopasswd,omitempty"`
	SudoCommands      []string `json:"sudo_commands,omitempty"`
	AllowDuplicateGID bool     `json:"allow_duplicate_gid,omitempty"`
	Users             []int    `json:"users,omitempty"`
}

// GroupDeleteRequest represents parameters for group.delete
type GroupDeleteRequest struct {
	DeleteUsers bool `json:"delete_users,omitempty"`
}

// GroupGetRequest represents parameters for group.get_group_obj
type GroupGetRequest struct {
	GroupName string `json:"groupname,omitempty"`
	GID       int    `json:"gid,omitempty"`
}

// List returns all groups
func (g *GroupClient) List(ctx context.Context) ([]Group, error) {
	var result []Group
	err := g.client.Call(ctx, "group.query", []any{}, &result)
	return result, err
}

// ListWithDSCache returns all groups including directory service groups
func (g *GroupClient) ListWithDSCache(ctx context.Context) ([]Group, error) {
	var result []Group
	options := map[string]any{
		"extra": map[string]any{
			"search_dscache": true,
		},
	}
	err := g.client.Call(ctx, "group.query", []any{[]any{}, options}, &result)
	return result, err
}

// Get returns a specific group by ID
func (g *GroupClient) Get(ctx context.Context, id int) (*Group, error) {
	var result []Group
	err := g.client.Call(ctx, "group.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("group", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// GetByName returns a specific group by name
func (g *GroupClient) GetByName(ctx context.Context, name string) (*Group, error) {
	var result []Group
	err := g.client.Call(ctx, "group.query", []any{[]any{[]any{"name", "=", name}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("group", fmt.Sprintf("name %s", name))
	}
	return &result[0], nil
}

// GetByGID returns a specific group by GID
func (g *GroupClient) GetByGID(ctx context.Context, gid int) (*Group, error) {
	var result []Group
	err := g.client.Call(ctx, "group.query", []any{[]any{[]any{"gid", "=", gid}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("group", fmt.Sprintf("GID %d", gid))
	}
	return &result[0], nil
}

// Create creates a new group
func (g *GroupClient) Create(ctx context.Context, req *GroupCreateRequest) (*Group, error) {
	var result Group
	err := g.client.Call(ctx, "group.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing group
func (g *GroupClient) Update(ctx context.Context, id int, req *GroupUpdateRequest) (*Group, error) {
	var result Group
	err := g.client.Call(ctx, "group.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes a group
func (g *GroupClient) Delete(ctx context.Context, id int, req *GroupDeleteRequest) error {
	params := []any{id}
	if req != nil {
		params = append(params, *req)
	}
	return g.client.Call(ctx, "group.delete", params, nil)
}

// GetNextGID returns the next available GID
func (g *GroupClient) GetNextGID(ctx context.Context) (int, error) {
	var result int
	err := g.client.Call(ctx, "group.get_next_gid", []any{}, &result)
	return result, err
}

// GetGroupObj returns group information from struct grp
func (g *GroupClient) GetGroupObj(ctx context.Context, req GroupGetRequest) (map[string]any, error) {
	var result map[string]any
	err := g.client.Call(ctx, "group.get_group_obj", []any{req}, &result)
	return result, err
}
