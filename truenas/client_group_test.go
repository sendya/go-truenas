package truenas

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGroupClient(t *testing.T) {
	t.Parallel()
	client := &Client{}
	groupClient := NewGroupClient(client)

	assert.NotNil(t, groupClient)
	assert.Equal(t, client, groupClient.client)
}

func TestNewGroupClient_NilClient(t *testing.T) {
	t.Parallel()
	groupClient := NewGroupClient(nil)

	assert.NotNil(t, groupClient)
	assert.Nil(t, groupClient.client)
}

func TestGroupClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockGroups := []Group{
		{ID: 1, GID: 1000, Name: "testgroup1", Builtin: false, Local: true},
		{ID: 2, GID: 1001, Name: "testgroup2", Builtin: true, Local: true},
	}
	server.SetResponse("group.query", mockGroups)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	groups, err := client.Group.List(ctx)
	require.NoError(t, err)
	assert.Len(t, groups, 2)
	assert.Equal(t, "testgroup1", groups[0].Name)
	assert.Equal(t, "testgroup2", groups[1].Name)
}

func TestGroupClient_List_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.query", 500, "Internal server error")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	groups, err := client.Group.List(ctx)
	assert.Error(t, err)
	assert.Nil(t, groups)
	assert.Contains(t, err.Error(), "Internal server error")
}

func TestGroupClient_ListWithDSCache(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockGroups := []Group{
		{ID: 1, GID: 1000, Name: "testgroup1", Builtin: false, Local: true},
		{ID: 2, GID: 1001, Name: "adgroup", Builtin: false, Local: false},
	}
	server.SetResponse("group.query", mockGroups)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	groups, err := client.Group.ListWithDSCache(ctx)
	require.NoError(t, err)
	assert.Len(t, groups, 2)
	assert.Equal(t, "testgroup1", groups[0].Name)
	assert.Equal(t, "adgroup", groups[1].Name)
	assert.True(t, groups[0].Local)
	assert.False(t, groups[1].Local)
}

func TestGroupClient_ListWithDSCache_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.query", 403, "Permission denied")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	groups, err := client.Group.ListWithDSCache(ctx)
	assert.Error(t, err)
	assert.Nil(t, groups)
	assert.Contains(t, err.Error(), "Permission denied")
}

func TestGroupClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockGroup := Group{
		ID:      1,
		GID:     1000,
		Name:    "testgroup",
		Builtin: false,
		Sudo:    true,
		Users:   []int{1001, 1002},
		Local:   true,
	}
	server.SetResponse("group.query", []Group{mockGroup})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	group, err := client.Group.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, group)
	assert.Equal(t, "testgroup", group.Name)
	assert.Equal(t, 1000, group.GID)
	assert.True(t, group.Sudo)
	assert.Len(t, group.Users, 2)
}

func TestGroupClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("group.query", []Group{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	group, err := client.Group.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, group)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestGroupClient_Get_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.query", 404, "Group not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	group, err := client.Group.Get(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, group)
	assert.Contains(t, err.Error(), "Group not found")
}

func TestGroupClient_GetByName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		groupName     string
		mockResponse  []Group
		expectedGroup *Group
		expectError   bool
	}{
		{
			name:      "successful get by name",
			groupName: "testgroup",
			mockResponse: []Group{
				{ID: 1, GID: 1000, Name: "testgroup", Builtin: false, Local: true},
			},
			expectedGroup: &Group{ID: 1, GID: 1000, Name: "testgroup", Builtin: false, Local: true},
			expectError:   false,
		},
		{
			name:          "group not found",
			groupName:     "nonexistent",
			mockResponse:  []Group{},
			expectedGroup: nil,
			expectError:   true,
		},
		{
			name:      "multiple groups with same name - returns first",
			groupName: "duplicate",
			mockResponse: []Group{
				{ID: 1, GID: 1000, Name: "duplicate", Builtin: false, Local: true},
				{ID: 2, GID: 1001, Name: "duplicate", Builtin: false, Local: false},
			},
			expectedGroup: &Group{ID: 1, GID: 1000, Name: "duplicate", Builtin: false, Local: true},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.Close()

			server.SetResponse("group.query", tt.mockResponse)

			client := server.CreateTestClient(t)
			defer client.Close()

			ctx := NewTestContext(t)
			group, err := client.Group.GetByName(ctx, tt.groupName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, group)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedGroup, group)
			}
		})
	}
}

func TestGroupClient_GetByName_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.query", 500, "Database error")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	group, err := client.Group.GetByName(ctx, "testgroup")
	assert.Error(t, err)
	assert.Nil(t, group)
	assert.Contains(t, err.Error(), "Database error")
}

func TestGroupClient_GetByGID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		gid           int
		mockResponse  []Group
		expectedGroup *Group
		expectError   bool
	}{
		{
			name: "successful get by GID",
			gid:  1000,
			mockResponse: []Group{
				{ID: 1, GID: 1000, Name: "testgroup", Builtin: false, Local: true},
			},
			expectedGroup: &Group{ID: 1, GID: 1000, Name: "testgroup", Builtin: false, Local: true},
			expectError:   false,
		},
		{
			name:          "GID not found",
			gid:           9999,
			mockResponse:  []Group{},
			expectedGroup: nil,
			expectError:   true,
		},
		{
			name: "zero GID",
			gid:  0,
			mockResponse: []Group{
				{ID: 1, GID: 0, Name: "root", Builtin: true, Local: true},
			},
			expectedGroup: &Group{ID: 1, GID: 0, Name: "root", Builtin: true, Local: true},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.Close()

			server.SetResponse("group.query", tt.mockResponse)

			client := server.CreateTestClient(t)
			defer client.Close()

			ctx := NewTestContext(t)
			group, err := client.Group.GetByGID(ctx, tt.gid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, group)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedGroup, group)
			}
		})
	}
}

func TestGroupClient_GetByGID_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.query", 403, "Access denied")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	group, err := client.Group.GetByGID(ctx, 1000)
	assert.Error(t, err)
	assert.Nil(t, group)
	assert.Contains(t, err.Error(), "Access denied")
}

func TestGroupClient_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		request     *GroupCreateRequest
		mockGroup   Group
		expectError bool
	}{
		{
			name: "basic group creation",
			request: &GroupCreateRequest{
				Name: "newgroup",
				GID:  1500,
				Smb:  true,
			},
			mockGroup: Group{
				ID:      10,
				GID:     1500,
				Name:    "newgroup",
				Builtin: false,
				Smb:     true,
				Local:   true,
			},
			expectError: false,
		},
		{
			name: "group with sudo permissions",
			request: &GroupCreateRequest{
				Name:         "sudogroup",
				GID:          1501,
				Sudo:         true,
				SudoNoPasswd: true,
				SudoCommands: []string{"/bin/ls", "/bin/cat"},
				Users:        []int{1001, 1002},
			},
			mockGroup: Group{
				ID:           11,
				GID:          1501,
				Name:         "sudogroup",
				Builtin:      false,
				Sudo:         true,
				SudoNoPasswd: true,
				SudoCommands: []string{"/bin/ls", "/bin/cat"},
				Users:        []int{1001, 1002},
				Local:        true,
			},
			expectError: false,
		},
		{
			name: "group with duplicate GID allowed",
			request: &GroupCreateRequest{
				Name:              "dupgroup",
				GID:               1000,
				AllowDuplicateGID: true,
			},
			mockGroup: Group{
				ID:      12,
				GID:     1000,
				Name:    "dupgroup",
				Builtin: false,
				Local:   true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.Close()

			server.SetResponse("group.create", tt.mockGroup)

			client := server.CreateTestClient(t)
			defer client.Close()

			ctx := NewTestContext(t)
			group, err := client.Group.Create(ctx, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, group)
			} else {
				require.NoError(t, err)
				require.NotNil(t, group)
				assert.Equal(t, tt.mockGroup.Name, group.Name)
				assert.Equal(t, tt.mockGroup.GID, group.GID)
				assert.Equal(t, tt.mockGroup.Sudo, group.Sudo)
				assert.Equal(t, tt.mockGroup.SudoNoPasswd, group.SudoNoPasswd)
				assert.Equal(t, tt.mockGroup.SudoCommands, group.SudoCommands)
				assert.Equal(t, tt.mockGroup.Users, group.Users)
			}
		})
	}
}

func TestGroupClient_Create_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.create", 422, "Group name already exists")

	client := server.CreateTestClient(t)
	defer client.Close()

	req := &GroupCreateRequest{
		Name: "existinggroup",
		GID:  1000,
	}

	ctx := NewTestContext(t)
	group, err := client.Group.Create(ctx, req)
	assert.Error(t, err)
	assert.NotNil(t, group) // The method returns &result even on error
	assert.Contains(t, err.Error(), "Group name already exists")
}

func TestGroupClient_Update(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		groupID     int
		request     *GroupUpdateRequest
		mockGroup   Group
		expectError bool
	}{
		{
			name:    "update group name",
			groupID: 1,
			request: &GroupUpdateRequest{
				Name: "updatedgroup",
			},
			mockGroup: Group{
				ID:      1,
				GID:     1000,
				Name:    "updatedgroup",
				Builtin: false,
				Local:   true,
			},
			expectError: false,
		},
		{
			name:    "update sudo settings",
			groupID: 1,
			request: &GroupUpdateRequest{
				Sudo:         true,
				SudoNoPasswd: false,
				SudoCommands: []string{"/usr/bin/vim", "/bin/nano"},
			},
			mockGroup: Group{
				ID:           1,
				GID:          1000,
				Name:         "testgroup",
				Builtin:      false,
				Sudo:         true,
				SudoNoPasswd: false,
				SudoCommands: []string{"/usr/bin/vim", "/bin/nano"},
				Local:        true,
			},
			expectError: false,
		},
		{
			name:    "update user membership",
			groupID: 1,
			request: &GroupUpdateRequest{
				Users: []int{1001, 1002, 1003},
			},
			mockGroup: Group{
				ID:      1,
				GID:     1000,
				Name:    "testgroup",
				Builtin: false,
				Users:   []int{1001, 1002, 1003},
				Local:   true,
			},
			expectError: false,
		},
		{
			name:    "clear user membership",
			groupID: 1,
			request: &GroupUpdateRequest{
				Users: []int{},
			},
			mockGroup: Group{
				ID:      1,
				GID:     1000,
				Name:    "testgroup",
				Builtin: false,
				Users:   []int{},
				Local:   true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.Close()

			server.SetResponse("group.update", tt.mockGroup)

			client := server.CreateTestClient(t)
			defer client.Close()

			ctx := NewTestContext(t)
			group, err := client.Group.Update(ctx, tt.groupID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, group)
			} else {
				require.NoError(t, err)
				require.NotNil(t, group)
				assert.Equal(t, tt.mockGroup.Name, group.Name)
				assert.Equal(t, tt.mockGroup.Sudo, group.Sudo)
				assert.Equal(t, tt.mockGroup.SudoNoPasswd, group.SudoNoPasswd)
				assert.Equal(t, tt.mockGroup.SudoCommands, group.SudoCommands)
				assert.Equal(t, tt.mockGroup.Users, group.Users)
			}
		})
	}
}

func TestGroupClient_Update_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.update", 404, "Group not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	req := &GroupUpdateRequest{
		Name: "newname",
	}

	ctx := NewTestContext(t)
	group, err := client.Group.Update(ctx, 999, req)
	assert.Error(t, err)
	assert.NotNil(t, group) // The method returns &result even on error
	assert.Contains(t, err.Error(), "Group not found")
}

func TestGroupClient_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		groupID     int
		request     *GroupDeleteRequest
		expectError bool
	}{
		{
			name:        "delete group without options",
			groupID:     1,
			request:     nil,
			expectError: false,
		},
		{
			name:    "delete group with delete users",
			groupID: 1,
			request: &GroupDeleteRequest{
				DeleteUsers: true,
			},
			expectError: false,
		},
		{
			name:    "delete group without deleting users",
			groupID: 1,
			request: &GroupDeleteRequest{
				DeleteUsers: false,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.Close()

			server.SetResponse("group.delete", true)

			client := server.CreateTestClient(t)
			defer client.Close()

			ctx := NewTestContext(t)
			err := client.Group.Delete(ctx, tt.groupID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGroupClient_Delete_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.delete", 403, "Cannot delete builtin group")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Group.Delete(ctx, 1, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot delete builtin group")
}

func TestGroupClient_GetNextGID(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("group.get_next_gid", 1500)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	gid, err := client.Group.GetNextGID(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1500, gid)
}

func TestGroupClient_GetNextGID_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.get_next_gid", 500, "Unable to calculate next GID")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	gid, err := client.Group.GetNextGID(ctx)
	assert.Error(t, err)
	assert.Equal(t, 0, gid)
	assert.Contains(t, err.Error(), "Unable to calculate next GID")
}

func TestGroupClient_GetGroupObj(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		request      GroupGetRequest
		mockResponse map[string]any
		expectError  bool
	}{
		{
			name: "get group object by name",
			request: GroupGetRequest{
				GroupName: "testgroup",
			},
			mockResponse: map[string]any{
				"gr_name":   "testgroup",
				"gr_gid":    float64(1000),
				"gr_passwd": "x",
				"gr_mem":    []any{"user1", "user2"},
			},
			expectError: false,
		},
		{
			name: "get group object by GID",
			request: GroupGetRequest{
				GID: 1000,
			},
			mockResponse: map[string]any{
				"gr_name":   "testgroup",
				"gr_gid":    float64(1000),
				"gr_passwd": "x",
				"gr_mem":    []any{"user1", "user2"},
			},
			expectError: false,
		},
		{
			name: "get group object with both name and GID",
			request: GroupGetRequest{
				GroupName: "testgroup",
				GID:       1000,
			},
			mockResponse: map[string]any{
				"gr_name":   "testgroup",
				"gr_gid":    float64(1000),
				"gr_passwd": "x",
				"gr_mem":    []any{"user1", "user2"},
			},
			expectError: false,
		},
		{
			name:    "empty request",
			request: GroupGetRequest{},
			mockResponse: map[string]any{
				"error": "No group specified",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.Close()

			server.SetResponse("group.get_group_obj", tt.mockResponse)

			client := server.CreateTestClient(t)
			defer client.Close()

			ctx := NewTestContext(t)
			obj, err := client.Group.GetGroupObj(ctx, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, obj)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockResponse, obj)
			}
		})
	}
}

func TestGroupClient_GetGroupObj_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("group.get_group_obj", 404, "Group not found in system")

	client := server.CreateTestClient(t)
	defer client.Close()

	req := GroupGetRequest{
		GroupName: "nonexistent",
	}

	ctx := NewTestContext(t)
	obj, err := client.Group.GetGroupObj(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, obj)
	assert.Contains(t, err.Error(), "Group not found in system")
}

// Test edge cases and input validation
func TestGroupClient_EdgeCases(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)

	t.Run("Get with negative ID", func(t *testing.T) {
		server.SetResponse("group.query", []Group{})
		group, err := client.Group.Get(ctx, -1)
		assert.Error(t, err)
		assert.Nil(t, group)
		var notFoundErr *NotFoundError
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("GetByName with empty string", func(t *testing.T) {
		server.SetResponse("group.query", []Group{})
		group, err := client.Group.GetByName(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, group)
		var notFoundErr *NotFoundError
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("GetByGID with negative GID", func(t *testing.T) {
		server.SetResponse("group.query", []Group{})
		group, err := client.Group.GetByGID(ctx, -1)
		assert.Error(t, err)
		assert.Nil(t, group)
		var notFoundErr *NotFoundError
		assert.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("Update with zero ID", func(t *testing.T) {
		server.SetError("group.update", 400, "Invalid group ID")
		req := &GroupUpdateRequest{Name: "newname"}
		group, err := client.Group.Update(ctx, 0, req)
		assert.Error(t, err)
		assert.NotNil(t, group) // The method returns &result even on error
	})

	t.Run("Delete with zero ID", func(t *testing.T) {
		server.SetError("group.delete", 400, "Invalid group ID")
		err := client.Group.Delete(ctx, 0, nil)
		assert.Error(t, err)
	})
}

// Test context cancellation
func TestGroupClient_ContextCancellation(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	// Create a context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	t.Run("List with canceled context", func(t *testing.T) {
		groups, err := client.Group.List(ctx)
		assert.Error(t, err)
		assert.Nil(t, groups)
		assert.Contains(t, err.Error(), "context canceled")
	})

	t.Run("Create with canceled context", func(t *testing.T) {
		req := &GroupCreateRequest{Name: "testgroup"}
		group, err := client.Group.Create(ctx, req)
		assert.Error(t, err)
		assert.NotNil(t, group) // The method returns &result even on error
		assert.Contains(t, err.Error(), "context canceled")
	})
}

// Test concurrent access (basic thread safety test)
func TestGroupClient_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("group.query", []Group{
		{ID: 1, GID: 1000, Name: "group1", Builtin: false, Local: true},
	})

	client := server.CreateTestClient(t)
	defer client.Close()

	// Run multiple goroutines concurrently
	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			ctx := NewTestContext(t)
			_, err := client.Group.List(ctx)
			results <- err
		}()
	}

	// Check that all goroutines completed without error
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
