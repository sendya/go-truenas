package truenas

import (
	"context"
	"fmt"
	"time"
)

// JobClient provides methods for job management and monitoring
type JobClient struct {
	client *Client
}

// NewJobClient creates a new job client
func NewJobClient(client *Client) *JobClient {
	return &JobClient{client: client}
}

// Job represents a TrueNAS job
type Job struct {
	ID           int              `json:"id"`
	Method       string           `json:"method"`
	Arguments    []any            `json:"arguments"`
	LogsPath     *string          `json:"logs_path"`
	LogsExcerpt  *string          `json:"logs_excerpt"`
	Progress     *JobProgress     `json:"progress"`
	Result       any              `json:"result"`
	Error        *string          `json:"error"`
	Exception    *string          `json:"exception"`
	ExcInfo      any              `json:"exc_info"`
	State        string           `json:"state"`
	TimeStarted  map[string]int64 `json:"time_started"`
	TimeFinished map[string]int64 `json:"time_finished"`
}

// JobProgress represents job progress information
type JobProgress struct {
	Percent     float64 `json:"percent"`
	Description string  `json:"description"`
	Extra       any     `json:"extra"`
}

// JobState represents possible job states
type JobState string

const (
	JobStateWaiting JobState = "WAITING"
	JobStateRunning JobState = "RUNNING"
	JobStateSuccess JobState = "SUCCESS"
	JobStateFailed  JobState = "FAILED"
	JobStateAborted JobState = "ABORTED"
)

// List returns all jobs
func (j *JobClient) List(ctx context.Context) ([]Job, error) {
	var result []Job
	err := j.client.Call(ctx, "core.get_jobs", []any{}, &result)
	return result, err
}

// Get returns a specific job by ID
func (j *JobClient) Get(ctx context.Context, id int) (*Job, error) {
	var result []Job
	err := j.client.Call(ctx, "core.get_jobs", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("job", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// IsCompleted checks if a job has completed (success or failed)
func (j *Job) IsCompleted() bool {
	state := JobState(j.State)
	return state == JobStateSuccess || state == JobStateFailed || state == JobStateAborted
}

// IsRunning checks if a job is currently running
func (j *Job) IsRunning() bool {
	return JobState(j.State) == JobStateRunning
}

// IsSuccessful checks if a job completed successfully
func (j *Job) IsSuccessful() bool {
	return JobState(j.State) == JobStateSuccess
}

// IsFailed checks if a job failed
func (j *Job) IsFailed() bool {
	state := JobState(j.State)
	return state == JobStateFailed || state == JobStateAborted
}

// Wait waits for a job to complete and returns the final job result
func (j *JobClient) Wait(ctx context.Context, jobID int) (*Job, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			job, err := j.Get(ctx, jobID)
			if err != nil {
				return nil, fmt.Errorf("get job %d: %w", jobID, err)
			}

			if job.IsCompleted() {
				if job.IsFailed() {
					if job.Error != nil {
						return job, fmt.Errorf("job %d failed: %s", jobID, *job.Error)
					}
					if job.Exception != nil {
						return job, fmt.Errorf("job %d failed with exception: %s", jobID, *job.Exception)
					}
					return job, fmt.Errorf("job %d failed", jobID)
				}
				return job, nil
			}
		}
	}
}
