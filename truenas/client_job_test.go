package truenas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockJobs := []Job{
		{
			ID:     1,
			Method: "pool.scrub.scrub",
			State:  "RUNNING",
			Progress: &JobProgress{
				Percent:     25.5,
				Description: "Scrubbing",
			},
		},
		{
			ID:     2,
			Method: "dataset.create",
			State:  "SUCCESS",
		},
	}
	server.SetResponse("core.get_jobs", mockJobs)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	jobs, err := client.Job.List(ctx)
	require.NoError(t, err)
	assert.Len(t, jobs, 2)
	assert.Equal(t, "pool.scrub.scrub", jobs[0].Method)
	assert.Equal(t, "RUNNING", jobs[0].State)
	assert.Equal(t, 25.5, jobs[0].Progress.Percent)
}

func TestJobClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockJob := Job{
		ID:     123,
		Method: "pool.scrub.scrub",
		State:  "SUCCESS",
		Progress: &JobProgress{
			Percent:     100.0,
			Description: "Scrub completed",
		},
	}
	server.SetResponse("core.get_jobs", []Job{mockJob})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	job, err := client.Job.Get(ctx, 123)
	require.NoError(t, err)
	require.NotNil(t, job)
	assert.Equal(t, 123, job.ID)
	assert.Equal(t, "pool.scrub.scrub", job.Method)
	assert.Equal(t, "SUCCESS", job.State)
}

func TestJob_StateCheckers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		state        string
		isCompleted  bool
		isRunning    bool
		isSuccessful bool
		isFailed     bool
	}{
		{"WAITING", false, false, false, false},
		{"RUNNING", false, true, false, false},
		{"SUCCESS", true, false, true, false},
		{"FAILED", true, false, false, true},
		{"ABORTED", true, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.state, func(t *testing.T) {
			job := &Job{State: tt.state}

			assert.Equal(t, tt.isCompleted, job.IsCompleted(), "IsCompleted")
			assert.Equal(t, tt.isRunning, job.IsRunning(), "IsRunning")
			assert.Equal(t, tt.isSuccessful, job.IsSuccessful(), "IsSuccessful")
			assert.Equal(t, tt.isFailed, job.IsFailed(), "IsFailed")
		})
	}
}

func TestJobClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("core.get_jobs", []Job{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	job, err := client.Job.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, job)

	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
	assert.Equal(t, "job", notFoundErr.ResourceType)
}
