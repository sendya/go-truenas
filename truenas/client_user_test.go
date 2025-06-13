package truenas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockUsers := []User{
		{ID: 1, UID: 1000, Username: "testuser1", FullName: "Test User 1"},
		{ID: 2, UID: 1001, Username: "testuser2", FullName: "Test User 2"},
	}
	server.SetResponse("user.query", mockUsers)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	users, err := client.User.List(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "testuser1", users[0].Username)
	assert.Equal(t, "testuser2", users[1].Username)
}

func TestUserClient_ListWithDSCache(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockUsers := []User{
		{ID: 1, UID: 1000, Username: "testuser1", FullName: "Test User 1"},
	}
	server.SetResponse("user.query", mockUsers)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	users, err := client.User.ListWithDSCache(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "testuser1", users[0].Username)
}

func TestUserClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockUser := User{ID: 1, UID: 1000, Username: "testuser", FullName: "Test User"}
	server.SetResponse("user.query", []User{mockUser})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	user, err := client.User.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, 1000, user.UID)
}

func TestUserClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.query", []User{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	user, err := client.User.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "not found")
}

func TestUserClient_GetByUsername(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockUser := User{ID: 1, UID: 1000, Username: "testuser", FullName: "Test User"}
	server.SetResponse("user.query", []User{mockUser})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	user, err := client.User.GetByUsername(ctx, "testuser")
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, 1, user.ID)
}

func TestUserClient_Create(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockUser := User{ID: 1, UID: 1000, Username: "newuser", FullName: "New User"}
	server.SetResponse("user.create", mockUser)

	client := server.CreateTestClient(t)
	defer client.Close()

	req := &UserCreateRequest{
		Username:         "newuser",
		FullName:         "New User",
		Password:         "password123",
		Group:            1000,
		Home:             "/home/newuser",
		Shell:            "/bin/bash",
		SSHPubKey:        "ssh-rsa AAAA...",
		Groups:           []int{1001, 1002},
		Email:            "newuser@example.com",
		Locked:           Ptr(false),
		Sudo:             Ptr(false),
		MicrosoftAccount: Ptr(false),
		SMB:              Ptr(true),
	}

	ctx := NewTestContext(t)
	user, err := client.User.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, 1000, user.UID)
}

func TestUserClient_Update(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockUser := User{ID: 1, UID: 1000, Username: "testuser", FullName: "Updated User"}
	server.SetResponse("user.update", mockUser)

	client := server.CreateTestClient(t)
	defer client.Close()

	req := &UserUpdateRequest{
		FullName: "Updated User",
		Email:    "updated@example.com",
		Locked:   Ptr(false),
	}

	ctx := NewTestContext(t)
	user, err := client.User.Update(ctx, 1, req)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Updated User", user.FullName)
}

func TestUserClient_Delete(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.delete", 1)

	client := server.CreateTestClient(t)
	defer client.Close()

	req := &UserDeleteRequest{
		DeleteGroup: Ptr(true),
	}

	ctx := NewTestContext(t)
	err := client.User.Delete(ctx, 1, req)
	assert.NoError(t, err)
}

func TestUserClient_GetNextUID(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.get_next_uid", 1001)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	uid, err := client.User.GetNextUID(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1001, uid)
}

func TestUserClient_GetUserObj(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockObj := map[string]any{
		"pw_name": "testuser",
		"pw_uid":  1000,
		"pw_gid":  1000,
	}
	server.SetResponse("user.get_user_obj", mockObj)

	client := server.CreateTestClient(t)
	defer client.Close()

	req := UserGetRequest{
		Username: "testuser",
	}

	ctx := NewTestContext(t)
	obj, err := client.User.GetUserObj(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "testuser", obj["pw_name"])
	assert.Equal(t, float64(1000), obj["pw_uid"])
}

func TestUserClient_HasRootPassword(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.has_root_password", true)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	hasPassword, err := client.User.HasRootPassword(ctx)
	require.NoError(t, err)
	assert.True(t, hasPassword)
}

func TestUserClient_SetRootPassword(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.set_root_password", true)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	req := SetRootPasswordRequest{
		Password: "newpassword",
		Options: &SetRootPasswordOptions{
			EC2: &SetRootPasswordEC2Options{Enabled: false},
		},
	}
	err := client.User.SetRootPassword(ctx, req)
	assert.NoError(t, err)
}

func TestUserClient_GetShellChoices(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockChoices := map[string]string{
		"/bin/bash": "Bourne Again SHell",
		"/bin/sh":   "Bourne SHell",
		"/bin/zsh":  "Z SHell",
		"/bin/csh":  "C SHell",
	}
	server.SetResponse("user.shell_choices", mockChoices)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	choices, err := client.User.GetShellChoices(ctx, nil)
	require.NoError(t, err)
	assert.Len(t, choices, 4)
	assert.Equal(t, "Bourne Again SHell", choices["/bin/bash"])
	assert.Equal(t, "Z SHell", choices["/bin/zsh"])
}

func TestUserClient_GetShellChoices_WithUserID(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockChoices := map[string]string{
		"/bin/bash":     "Bourne Again SHell",
		"/bin/sh":       "Bourne SHell",
		"/sbin/nologin": "No Login",
	}
	server.SetResponse("user.shell_choices", mockChoices)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	userID := 1000
	choices, err := client.User.GetShellChoices(ctx, &userID)
	require.NoError(t, err)
	assert.Len(t, choices, 3)
	assert.Equal(t, "Bourne Again SHell", choices["/bin/bash"])
	assert.Equal(t, "No Login", choices["/sbin/nologin"])
}

func TestUserClient_GetShellChoices_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("user.shell_choices", 422, "Invalid user ID")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	userID := -1
	choices, err := client.User.GetShellChoices(ctx, &userID)
	assert.Error(t, err)
	assert.Nil(t, choices)
	assert.Contains(t, err.Error(), "Invalid user ID")
}

func TestUserClient_SetAttribute(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.set_attribute", true)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.User.SetAttribute(ctx, 1000, "custom_field", "custom_value")
	assert.NoError(t, err)
}

func TestUserClient_SetAttribute_ComplexValue(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.set_attribute", true)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	complexValue := map[string]any{
		"setting1": "value1",
		"setting2": 42,
		"setting3": true,
		"nested": map[string]any{
			"key": "nested_value",
		},
	}
	err := client.User.SetAttribute(ctx, 1000, "complex_config", complexValue)
	assert.NoError(t, err)
}

func TestUserClient_SetAttribute_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("user.set_attribute", 422, "User not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.User.SetAttribute(ctx, 99999, "key", "value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User not found")
}

func TestUserClient_PopAttribute(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.pop_attribute", true)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.User.PopAttribute(ctx, 1000, "old_setting")
	assert.NoError(t, err)
}

func TestUserClient_PopAttribute_NonExistentKey(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("user.pop_attribute", true)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.User.PopAttribute(ctx, 1000, "non_existent_key")
	assert.NoError(t, err)
}

func TestUserClient_PopAttribute_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("user.pop_attribute", 422, "User not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.User.PopAttribute(ctx, 99999, "some_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User not found")
}
