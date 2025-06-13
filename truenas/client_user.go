package truenas

import (
	"context"
	"fmt"
)

// UserClient provides methods for user management
type UserClient struct {
	client *Client
}

// NewUserClient creates a new user client
func NewUserClient(client *Client) *UserClient {
	return &UserClient{client: client}
}

// User represents a system user
type User struct {
	ID               int            `json:"id"`
	UID              int            `json:"uid"`
	Username         string         `json:"username"`
	UnixHash         string         `json:"unixhash"`
	SMBHash          string         `json:"smbhash"`
	Group            Group          `json:"group"`
	Home             string         `json:"home"`
	Shell            string         `json:"shell"`
	FullName         string         `json:"full_name"`
	Builtin          bool           `json:"builtin"`
	Email            string         `json:"email"`
	PasswordDisabled bool           `json:"password_disabled"`
	Locked           bool           `json:"locked"`
	Sudo             bool           `json:"sudo"`
	SudoNoPasswd     bool           `json:"sudo_nopasswd"`
	SudoCommands     []string       `json:"sudo_commands"`
	MicrosoftAccount bool           `json:"microsoft_account"`
	Attributes       map[string]any `json:"attributes"`
	Groups           []int          `json:"groups"`
	SSHPubKey        string         `json:"sshpubkey"`
	LocalGroups      []Group        `json:"local_groups"`
	SMB              bool           `json:"smb"`
	HomeMode         string         `json:"home_mode"`
}

// UserCreateRequest represents parameters for user.create
type UserCreateRequest struct {
	UID              int            `json:"uid,omitempty"`
	Username         string         `json:"username"`
	Group            int            `json:"group,omitempty"`
	GroupCreate      *bool          `json:"group_create,omitempty"`
	Home             string         `json:"home,omitempty"`
	HomeMode         string         `json:"home_mode,omitempty"`
	Shell            string         `json:"shell,omitempty"`
	FullName         string         `json:"full_name,omitempty"`
	Email            string         `json:"email,omitempty"`
	Password         string         `json:"password,omitempty"`
	PasswordDisabled *bool          `json:"password_disabled,omitempty"`
	Locked           *bool          `json:"locked,omitempty"`
	MicrosoftAccount *bool          `json:"microsoft_account,omitempty"`
	SMB              *bool          `json:"smb,omitempty"`
	Sudo             *bool          `json:"sudo,omitempty"`
	SudoNoPasswd     *bool          `json:"sudo_nopasswd,omitempty"`
	SudoCommands     []string       `json:"sudo_commands,omitempty"`
	SSHPubKey        string         `json:"sshpubkey,omitempty"`
	Groups           []int          `json:"groups,omitempty"`
	Attributes       map[string]any `json:"attributes,omitempty"`
}

// UserUpdateRequest represents parameters for user.update
type UserUpdateRequest struct {
	UID              int            `json:"uid,omitempty"`
	Username         string         `json:"username,omitempty"`
	Group            int            `json:"group,omitempty"`
	Home             string         `json:"home,omitempty"`
	HomeMode         string         `json:"home_mode,omitempty"`
	Shell            string         `json:"shell,omitempty"`
	FullName         string         `json:"full_name,omitempty"`
	Email            string         `json:"email,omitempty"`
	Password         string         `json:"password,omitempty"`
	PasswordDisabled *bool          `json:"password_disabled,omitempty"`
	Locked           *bool          `json:"locked,omitempty"`
	MicrosoftAccount *bool          `json:"microsoft_account,omitempty"`
	SMB              *bool          `json:"smb,omitempty"`
	Sudo             *bool          `json:"sudo,omitempty"`
	SudoNoPasswd     *bool          `json:"sudo_nopasswd,omitempty"`
	SudoCommands     []string       `json:"sudo_commands,omitempty"`
	SSHPubKey        string         `json:"sshpubkey,omitempty"`
	Groups           []int          `json:"groups,omitempty"`
	Attributes       map[string]any `json:"attributes,omitempty"`
}

// UserDeleteRequest represents parameters for user.delete
type UserDeleteRequest struct {
	DeleteGroup *bool `json:"delete_group,omitempty"`
}

// UserGetRequest represents parameters for user.get_user_obj
type UserGetRequest struct {
	Username string `json:"username,omitempty"`
	UID      int    `json:"uid,omitempty"`
}

// SetRootPasswordRequest represents parameters for user.set_root_password
type SetRootPasswordRequest struct {
	Password string                  `json:"password"`
	Options  *SetRootPasswordOptions `json:"options,omitempty"`
}

// SetRootPasswordOptions represents options for setting root password
type SetRootPasswordOptions struct {
	EC2 *SetRootPasswordEC2Options `json:"ec2,omitempty"`
}

// SetRootPasswordEC2Options represents EC2-specific options
type SetRootPasswordEC2Options struct {
	Enabled bool `json:"enabled"`
}

// List returns all users
func (u *UserClient) List(ctx context.Context) ([]User, error) {
	var result []User
	err := u.client.Call(ctx, "user.query", []any{}, &result)
	return result, err
}

// ListWithDSCache returns all users including directory service users
func (u *UserClient) ListWithDSCache(ctx context.Context) ([]User, error) {
	var result []User
	options := map[string]any{
		"extra": map[string]any{
			"search_dscache": true,
		},
	}
	err := u.client.Call(ctx, "user.query", []any{[]any{}, options}, &result)
	return result, err
}

// Get returns a specific user by ID
func (u *UserClient) Get(ctx context.Context, id int) (*User, error) {
	var result []User
	err := u.client.Call(ctx, "user.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("user", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// GetByUsername returns a specific user by username
func (u *UserClient) GetByUsername(ctx context.Context, username string) (*User, error) {
	var result []User
	err := u.client.Call(ctx, "user.query", []any{[]any{[]any{"username", "=", username}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("user", fmt.Sprintf("username %s", username))
	}
	return &result[0], nil
}

// Create creates a new user
func (u *UserClient) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	var result User
	err := u.client.Call(ctx, "user.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing user
func (u *UserClient) Update(ctx context.Context, id int, req *UserUpdateRequest) (*User, error) {
	var result User
	err := u.client.Call(ctx, "user.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes a user
func (u *UserClient) Delete(ctx context.Context, id int, req *UserDeleteRequest) error {
	params := []any{id}
	if req != nil {
		params = append(params, *req)
	}
	return u.client.Call(ctx, "user.delete", params, nil)
}

// GetNextUID returns the next available UID
func (u *UserClient) GetNextUID(ctx context.Context) (int, error) {
	var result int
	err := u.client.Call(ctx, "user.get_next_uid", []any{}, &result)
	return result, err
}

// GetUserObj returns user information from struct passwd
func (u *UserClient) GetUserObj(ctx context.Context, req UserGetRequest) (map[string]any, error) {
	var result map[string]any
	err := u.client.Call(ctx, "user.get_user_obj", []any{req}, &result)
	return result, err
}

// HasRootPassword checks if root user has a valid password set
func (u *UserClient) HasRootPassword(ctx context.Context) (bool, error) {
	var result bool
	err := u.client.Call(ctx, "user.has_root_password", []any{}, &result)
	return result, err
}

// SetRootPassword sets the root user password
func (u *UserClient) SetRootPassword(ctx context.Context, req SetRootPasswordRequest) error {
	params := []any{req.Password}
	if req.Options != nil {
		params = append(params, req.Options)
	}
	return u.client.Call(ctx, "user.set_root_password", params, nil)
}

// SetRootPasswordSimple sets the root password with just a password string
func (u *UserClient) SetRootPasswordSimple(ctx context.Context, password string) error {
	return u.client.Call(ctx, "user.set_root_password", []any{password}, nil)
}

// GetShellChoices returns available shell choices
func (u *UserClient) GetShellChoices(ctx context.Context, userID *int) (map[string]string, error) {
	var result map[string]string
	params := []any{}
	if userID != nil {
		params = append(params, *userID)
	}
	err := u.client.Call(ctx, "user.shell_choices", params, &result)
	return result, err
}

// SetAttribute sets a user attribute
func (u *UserClient) SetAttribute(ctx context.Context, id int, key string, value any) error {
	return u.client.Call(ctx, "user.set_attribute", []any{id, key, value}, nil)
}

// PopAttribute removes a user attribute
func (u *UserClient) PopAttribute(ctx context.Context, id int, key string) error {
	return u.client.Call(ctx, "user.pop_attribute", []any{id, key}, nil)
}
