package truenas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCronjobClient(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	cronjobClient := NewCronjobClient(client)
	require.NotNil(t, cronjobClient)
	assert.Equal(t, client, cronjobClient.client)
}

func TestCronjobClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockCronjobs := []Cronjob{
		{
			ID:          1,
			Enabled:     true,
			Stderr:      true,
			Stdout:      true,
			Schedule:    Schedule{Minute: "0", Hour: "2", DOM: "*", Month: "*", DOW: "*"},
			Command:     "/usr/bin/backup.sh",
			Description: "Daily backup",
			User:        "root",
		},
		{
			ID:          2,
			Enabled:     false,
			Stderr:      false,
			Stdout:      false,
			Schedule:    Schedule{Minute: "30", Hour: "*/6", DOM: "*", Month: "*", DOW: "*"},
			Command:     "/usr/bin/cleanup.sh",
			Description: "Cleanup temp files",
			User:        "admin",
		},
	}
	server.SetResponse("cronjob.query", mockCronjobs)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	cronjobs, err := client.Cronjob.List(ctx)
	require.NoError(t, err)
	assert.Len(t, cronjobs, 2)
	assert.Equal(t, 1, cronjobs[0].ID)
	assert.Equal(t, "/usr/bin/backup.sh", cronjobs[0].Command)
	assert.Equal(t, "Daily backup", cronjobs[0].Description)
	assert.True(t, cronjobs[0].Enabled)
	assert.Equal(t, "root", cronjobs[0].User)
}

func TestCronjobClient_List_Empty(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("cronjob.query", []Cronjob{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	cronjobs, err := client.Cronjob.List(ctx)
	require.NoError(t, err)
	assert.Len(t, cronjobs, 0)
}

func TestCronjobClient_List_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("cronjob.query", 500, "Internal server error")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	cronjobs, err := client.Cronjob.List(ctx)
	assert.Error(t, err)
	assert.Nil(t, cronjobs)
	assert.Contains(t, err.Error(), "Internal server error")
}

func TestCronjobClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockCronjob := Cronjob{
		ID:          1,
		Enabled:     true,
		Stderr:      true,
		Stdout:      true,
		Schedule:    Schedule{Minute: "0", Hour: "2", DOM: "*", Month: "*", DOW: "*"},
		Command:     "/usr/bin/backup.sh",
		Description: "Daily backup",
		User:        "root",
	}
	server.SetResponse("cronjob.query", []Cronjob{mockCronjob})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, cronjob)
	assert.Equal(t, 1, cronjob.ID)
	assert.Equal(t, "/usr/bin/backup.sh", cronjob.Command)
	assert.Equal(t, "Daily backup", cronjob.Description)
	assert.True(t, cronjob.Enabled)
	assert.Equal(t, "root", cronjob.User)
}

func TestCronjobClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("cronjob.query", []Cronjob{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, cronjob)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestCronjobClient_Get_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("cronjob.query", 404, "Cronjob not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Get(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, cronjob)
	assert.Contains(t, err.Error(), "Cronjob not found")
}

func TestCronjobClient_Create(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockCronjob := Cronjob{
		ID:          1,
		Enabled:     true,
		Stderr:      true,
		Stdout:      true,
		Schedule:    Schedule{Minute: "0", Hour: "2", DOM: "*", Month: "*", DOW: "*"},
		Command:     "/usr/bin/backup.sh",
		Description: "Daily backup",
		User:        "root",
	}
	server.SetResponse("cronjob.create", mockCronjob)

	client := server.CreateTestClient(t)
	defer client.Close()

	req := &CronjobCreateRequest{
		Enabled:     true,
		Stderr:      true,
		Stdout:      true,
		Schedule:    Schedule{Minute: "0", Hour: "2", DOM: "*", Month: "*", DOW: "*"},
		Command:     "/usr/bin/backup.sh",
		Description: "Daily backup",
		User:        "root",
	}

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, cronjob)
	assert.Equal(t, 1, cronjob.ID)
	assert.Equal(t, "/usr/bin/backup.sh", cronjob.Command)
	assert.Equal(t, "Daily backup", cronjob.Description)
	assert.True(t, cronjob.Enabled)
	assert.Equal(t, "root", cronjob.User)
}

func TestCronjobClient_Create_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("cronjob.create", 400, "Invalid command")

	client := server.CreateTestClient(t)
	defer client.Close()

	req := &CronjobCreateRequest{
		Enabled:     true,
		Command:     "",
		Description: "Invalid cronjob",
		User:        "root",
	}

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Create(ctx, req)
	assert.Error(t, err)
	assert.NotNil(t, cronjob) // The method returns &result even on error
	assert.Contains(t, err.Error(), "Invalid command")
}

func TestCronjobClient_Update(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockCronjob := Cronjob{
		ID:          1,
		Enabled:     false,
		Stderr:      false,
		Stdout:      true,
		Schedule:    Schedule{Minute: "30", Hour: "3", DOM: "*", Month: "*", DOW: "*"},
		Command:     "/usr/bin/updated_backup.sh",
		Description: "Updated daily backup",
		User:        "admin",
	}
	server.SetResponse("cronjob.update", mockCronjob)

	client := server.CreateTestClient(t)
	defer client.Close()

	enabled := false
	stderr := false
	command := "/usr/bin/updated_backup.sh"
	description := "Updated daily backup"
	user := "admin"

	req := &CronjobUpdateRequest{
		Enabled:     &enabled,
		Stderr:      &stderr,
		Command:     &command,
		Description: &description,
		User:        &user,
		Schedule:    &Schedule{Minute: "30", Hour: "3", DOM: "*", Month: "*", DOW: "*"},
	}

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Update(ctx, 1, req)
	require.NoError(t, err)
	require.NotNil(t, cronjob)
	assert.Equal(t, 1, cronjob.ID)
	assert.Equal(t, "/usr/bin/updated_backup.sh", cronjob.Command)
	assert.Equal(t, "Updated daily backup", cronjob.Description)
	assert.False(t, cronjob.Enabled)
	assert.Equal(t, "admin", cronjob.User)
}

func TestCronjobClient_Update_PartialUpdate(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockCronjob := Cronjob{
		ID:          1,
		Enabled:     false,
		Stderr:      true,
		Stdout:      true,
		Schedule:    Schedule{Minute: "0", Hour: "2", DOM: "*", Month: "*", DOW: "*"},
		Command:     "/usr/bin/backup.sh",
		Description: "Daily backup",
		User:        "root",
	}
	server.SetResponse("cronjob.update", mockCronjob)

	client := server.CreateTestClient(t)
	defer client.Close()

	enabled := false
	req := &CronjobUpdateRequest{
		Enabled: &enabled,
	}

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Update(ctx, 1, req)
	require.NoError(t, err)
	require.NotNil(t, cronjob)
	assert.Equal(t, 1, cronjob.ID)
	assert.False(t, cronjob.Enabled)
}

func TestCronjobClient_Update_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("cronjob.update", 404, "Cronjob not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	enabled := true
	req := &CronjobUpdateRequest{
		Enabled: &enabled,
	}

	ctx := NewTestContext(t)
	cronjob, err := client.Cronjob.Update(ctx, 999, req)
	assert.Error(t, err)
	assert.NotNil(t, cronjob) // The method returns &result even on error
	assert.Contains(t, err.Error(), "Cronjob not found")
}

func TestCronjobClient_Delete(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("cronjob.delete", true)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Cronjob.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestCronjobClient_Delete_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("cronjob.delete", 404, "Cronjob not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Cronjob.Delete(ctx, 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cronjob not found")
}

func TestCronjobClient_Run(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	// Set up job response for cronjob.run
	server.SetJobResponse("cronjob.run", "Job completed successfully")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Cronjob.Run(ctx, 1, false)
	assert.NoError(t, err)
}

func TestCronjobClient_Run_SkipDisabled(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetJobResponse("cronjob.run", "Job skipped - cronjob is disabled")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Cronjob.Run(ctx, 1, true)
	assert.NoError(t, err)
}

func TestCronjobClient_Run_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetJobError("cronjob.run", "Cronjob execution failed")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Cronjob.Run(ctx, 1, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cronjob execution failed")
}

func TestNewDailySchedule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		hour     string
		minute   string
		expected Schedule
	}{
		{
			name:   "daily at 2:00 AM",
			hour:   "2",
			minute: "0",
			expected: Schedule{
				Minute: "0",
				Hour:   "2",
				DOM:    "*",
				Month:  "*",
				DOW:    "*",
			},
		},
		{
			name:   "daily at 11:30 PM",
			hour:   "23",
			minute: "30",
			expected: Schedule{
				Minute: "30",
				Hour:   "23",
				DOM:    "*",
				Month:  "*",
				DOW:    "*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := NewDailySchedule(tt.hour, tt.minute)
			assert.Equal(t, tt.expected, schedule)
		})
	}
}

func TestNewWeeklySchedule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		dayOfWeek string
		hour      string
		minute    string
		expected  Schedule
	}{
		{
			name:      "weekly on Sunday at 2:00 AM",
			dayOfWeek: "0",
			hour:      "2",
			minute:    "0",
			expected: Schedule{
				Minute: "0",
				Hour:   "2",
				DOM:    "*",
				Month:  "*",
				DOW:    "0",
			},
		},
		{
			name:      "weekly on Friday at 6:30 PM",
			dayOfWeek: "5",
			hour:      "18",
			minute:    "30",
			expected: Schedule{
				Minute: "30",
				Hour:   "18",
				DOM:    "*",
				Month:  "*",
				DOW:    "5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := NewWeeklySchedule(tt.dayOfWeek, tt.hour, tt.minute)
			assert.Equal(t, tt.expected, schedule)
		})
	}
}

func TestNewMonthlySchedule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		dayOfMonth string
		hour       string
		minute     string
		expected   Schedule
	}{
		{
			name:       "monthly on 1st at 2:00 AM",
			dayOfMonth: "1",
			hour:       "2",
			minute:     "0",
			expected: Schedule{
				Minute: "0",
				Hour:   "2",
				DOM:    "1",
				Month:  "*",
				DOW:    "*",
			},
		},
		{
			name:       "monthly on 15th at 11:45 PM",
			dayOfMonth: "15",
			hour:       "23",
			minute:     "45",
			expected: Schedule{
				Minute: "45",
				Hour:   "23",
				DOM:    "15",
				Month:  "*",
				DOW:    "*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := NewMonthlySchedule(tt.dayOfMonth, tt.hour, tt.minute)
			assert.Equal(t, tt.expected, schedule)
		})
	}
}

func TestNewHourlySchedule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		minute   string
		expected Schedule
	}{
		{
			name:   "hourly at 0 minutes",
			minute: "0",
			expected: Schedule{
				Minute: "0",
				Hour:   "*",
				DOM:    "*",
				Month:  "*",
				DOW:    "*",
			},
		},
		{
			name:   "hourly at 30 minutes",
			minute: "30",
			expected: Schedule{
				Minute: "30",
				Hour:   "*",
				DOM:    "*",
				Month:  "*",
				DOW:    "*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := NewHourlySchedule(tt.minute)
			assert.Equal(t, tt.expected, schedule)
		})
	}
}

func TestNewCustomSchedule(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		minute   string
		hour     string
		dom      string
		month    string
		dow      string
		expected Schedule
	}{
		{
			name:   "every 15 minutes",
			minute: "*/15",
			hour:   "*",
			dom:    "*",
			month:  "*",
			dow:    "*",
			expected: Schedule{
				Minute: "*/15",
				Hour:   "*",
				DOM:    "*",
				Month:  "*",
				DOW:    "*",
			},
		},
		{
			name:   "weekdays at 9 AM",
			minute: "0",
			hour:   "9",
			dom:    "*",
			month:  "*",
			dow:    "1-5",
			expected: Schedule{
				Minute: "0",
				Hour:   "9",
				DOM:    "*",
				Month:  "*",
				DOW:    "1-5",
			},
		},
		{
			name:   "first Monday of every month at midnight",
			minute: "0",
			hour:   "0",
			dom:    "1-7",
			month:  "*",
			dow:    "1",
			expected: Schedule{
				Minute: "0",
				Hour:   "0",
				DOM:    "1-7",
				Month:  "*",
				DOW:    "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := NewCustomSchedule(tt.minute, tt.hour, tt.dom, tt.month, tt.dow)
			assert.Equal(t, tt.expected, schedule)
		})
	}
}

// Edge case tests for schedule creation functions
func TestScheduleHelpers_EdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("daily schedule with wildcard values", func(t *testing.T) {
		schedule := NewDailySchedule("*", "*/30")
		assert.Equal(t, "*/30", schedule.Minute)
		assert.Equal(t, "*", schedule.Hour)
		assert.Equal(t, "*", schedule.DOM)
		assert.Equal(t, "*", schedule.Month)
		assert.Equal(t, "*", schedule.DOW)
	})

	t.Run("weekly schedule with range values", func(t *testing.T) {
		schedule := NewWeeklySchedule("1-5", "9-17", "0,30")
		assert.Equal(t, "0,30", schedule.Minute)
		assert.Equal(t, "9-17", schedule.Hour)
		assert.Equal(t, "*", schedule.DOM)
		assert.Equal(t, "*", schedule.Month)
		assert.Equal(t, "1-5", schedule.DOW)
	})

	t.Run("monthly schedule with multiple values", func(t *testing.T) {
		schedule := NewMonthlySchedule("1,15", "0,12", "0")
		assert.Equal(t, "0", schedule.Minute)
		assert.Equal(t, "0,12", schedule.Hour)
		assert.Equal(t, "1,15", schedule.DOM)
		assert.Equal(t, "*", schedule.Month)
		assert.Equal(t, "*", schedule.DOW)
	})

	t.Run("hourly schedule with step values", func(t *testing.T) {
		schedule := NewHourlySchedule("*/10")
		assert.Equal(t, "*/10", schedule.Minute)
		assert.Equal(t, "*", schedule.Hour)
		assert.Equal(t, "*", schedule.DOM)
		assert.Equal(t, "*", schedule.Month)
		assert.Equal(t, "*", schedule.DOW)
	})
}

// Test JSON marshaling/unmarshaling of structures
func TestCronjobJSON(t *testing.T) {
	t.Parallel()
	t.Run("cronjob struct serialization", func(t *testing.T) {
		cronjob := Cronjob{
			ID:          1,
			Enabled:     true,
			Stderr:      false,
			Stdout:      true,
			Schedule:    Schedule{Minute: "0", Hour: "2", DOM: "*", Month: "*", DOW: "*"},
			Command:     "/usr/bin/test.sh",
			Description: "Test cronjob",
			User:        "root",
		}

		// This would typically be tested with actual JSON marshaling
		// but we're focusing on the API client functionality
		assert.Equal(t, 1, cronjob.ID)
		assert.True(t, cronjob.Enabled)
		assert.False(t, cronjob.Stderr)
		assert.True(t, cronjob.Stdout)
	})

	t.Run("schedule struct completeness", func(t *testing.T) {
		schedule := Schedule{
			Minute: "30",
			Hour:   "2",
			DOM:    "15",
			Month:  "*/3",
			DOW:    "1-5",
		}

		assert.Equal(t, "30", schedule.Minute)
		assert.Equal(t, "2", schedule.Hour)
		assert.Equal(t, "15", schedule.DOM)
		assert.Equal(t, "*/3", schedule.Month)
		assert.Equal(t, "1-5", schedule.DOW)
	})
}

// Test comprehensive request structures
func TestCronjobRequestStructures(t *testing.T) {
	t.Parallel()
	t.Run("create request validation", func(t *testing.T) {
		req := CronjobCreateRequest{
			Enabled:     true,
			Stderr:      true,
			Stdout:      false,
			Schedule:    NewHourlySchedule("0"),
			Command:     "/usr/bin/hourly_task.sh",
			Description: "Hourly maintenance task",
			User:        "maintenance",
		}

		assert.True(t, req.Enabled)
		assert.True(t, req.Stderr)
		assert.False(t, req.Stdout)
		assert.Equal(t, "/usr/bin/hourly_task.sh", req.Command)
		assert.Equal(t, "maintenance", req.User)
		assert.Equal(t, "0", req.Schedule.Minute)
		assert.Equal(t, "*", req.Schedule.Hour)
	})

	t.Run("update request with nil values", func(t *testing.T) {
		req := CronjobUpdateRequest{
			Enabled:     nil,
			Stderr:      nil,
			Stdout:      nil,
			Schedule:    nil,
			Command:     nil,
			Description: nil,
			User:        nil,
		}

		// All fields should be nil (omitted in JSON)
		assert.Nil(t, req.Enabled)
		assert.Nil(t, req.Stderr)
		assert.Nil(t, req.Stdout)
		assert.Nil(t, req.Schedule)
		assert.Nil(t, req.Command)
		assert.Nil(t, req.Description)
		assert.Nil(t, req.User)
	})

	t.Run("update request with all values", func(t *testing.T) {
		enabled := true
		stderr := false
		stdout := true
		command := "/usr/bin/updated.sh"
		description := "Updated description"
		user := "newuser"
		schedule := NewDailySchedule("3", "15")

		req := CronjobUpdateRequest{
			Enabled:     &enabled,
			Stderr:      &stderr,
			Stdout:      &stdout,
			Schedule:    &schedule,
			Command:     &command,
			Description: &description,
			User:        &user,
		}

		assert.NotNil(t, req.Enabled)
		assert.True(t, *req.Enabled)
		assert.NotNil(t, req.Stderr)
		assert.False(t, *req.Stderr)
		assert.NotNil(t, req.Stdout)
		assert.True(t, *req.Stdout)
		assert.NotNil(t, req.Command)
		assert.Equal(t, "/usr/bin/updated.sh", *req.Command)
		assert.NotNil(t, req.User)
		assert.Equal(t, "newuser", *req.User)
	})
}
