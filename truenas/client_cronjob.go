package truenas

import (
	"context"
	"fmt"
)

// CronjobClient provides methods for cronjob management
type CronjobClient struct {
	client *Client
}

// NewCronjobClient creates a new cronjob client
func NewCronjobClient(client *Client) *CronjobClient {
	return &CronjobClient{client: client}
}

// Cronjob represents a scheduled cron job
type Cronjob struct {
	ID          int      `json:"id"`
	Enabled     bool     `json:"enabled"`
	Stderr      bool     `json:"stderr"`
	Stdout      bool     `json:"stdout"`
	Schedule    Schedule `json:"schedule"`
	Command     string   `json:"command"`
	Description string   `json:"description"`
	User        string   `json:"user"`
}

// Schedule represents the cron schedule configuration
type Schedule struct {
	Minute string `json:"minute"`
	Hour   string `json:"hour"`
	DOM    string `json:"dom"` // Day of month
	Month  string `json:"month"`
	DOW    string `json:"dow"` // Day of week
}

// CronjobCreateRequest represents parameters for cronjob.create
type CronjobCreateRequest struct {
	Enabled     bool     `json:"enabled"`
	Stderr      bool     `json:"stderr"`
	Stdout      bool     `json:"stdout"`
	Schedule    Schedule `json:"schedule"`
	Command     string   `json:"command"`
	Description string   `json:"description"`
	User        string   `json:"user"`
}

// CronjobUpdateRequest represents parameters for cronjob.update
type CronjobUpdateRequest struct {
	Enabled     *bool     `json:"enabled,omitempty"`
	Stderr      *bool     `json:"stderr,omitempty"`
	Stdout      *bool     `json:"stdout,omitempty"`
	Schedule    *Schedule `json:"schedule,omitempty"`
	Command     *string   `json:"command,omitempty"`
	Description *string   `json:"description,omitempty"`
	User        *string   `json:"user,omitempty"`
}

// List returns all cronjobs
func (c *CronjobClient) List(ctx context.Context) ([]Cronjob, error) {
	var result []Cronjob
	err := c.client.Call(ctx, "cronjob.query", []any{}, &result)
	return result, err
}

// Get returns a specific cronjob by ID
func (c *CronjobClient) Get(ctx context.Context, id int) (*Cronjob, error) {
	var result []Cronjob
	err := c.client.Call(ctx, "cronjob.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("cronjob", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new cronjob
func (c *CronjobClient) Create(ctx context.Context, req *CronjobCreateRequest) (*Cronjob, error) {
	var result Cronjob
	err := c.client.Call(ctx, "cronjob.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing cronjob
func (c *CronjobClient) Update(ctx context.Context, id int, req *CronjobUpdateRequest) (*Cronjob, error) {
	var result Cronjob
	err := c.client.Call(ctx, "cronjob.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes a cronjob
func (c *CronjobClient) Delete(ctx context.Context, id int) error {
	return c.client.Call(ctx, "cronjob.delete", []any{id}, nil)
}

// Run executes a cronjob immediately (asynchronous job)
func (c *CronjobClient) Run(ctx context.Context, id int, skipDisabled bool) error {
	return c.client.CallJob(ctx, "cronjob.run", []any{id, skipDisabled}, nil)
}

// Helper methods for creating common schedules

// NewDailySchedule creates a schedule that runs daily at the specified hour and minute
func NewDailySchedule(hour, minute string) Schedule {
	return Schedule{
		Minute: minute,
		Hour:   hour,
		DOM:    "*",
		Month:  "*",
		DOW:    "*",
	}
}

// NewWeeklySchedule creates a schedule that runs weekly on the specified day, hour, and minute
func NewWeeklySchedule(dayOfWeek, hour, minute string) Schedule {
	return Schedule{
		Minute: minute,
		Hour:   hour,
		DOM:    "*",
		Month:  "*",
		DOW:    dayOfWeek,
	}
}

// NewMonthlySchedule creates a schedule that runs monthly on the specified day, hour, and minute
func NewMonthlySchedule(dayOfMonth, hour, minute string) Schedule {
	return Schedule{
		Minute: minute,
		Hour:   hour,
		DOM:    dayOfMonth,
		Month:  "*",
		DOW:    "*",
	}
}

// NewHourlySchedule creates a schedule that runs hourly at the specified minute
func NewHourlySchedule(minute string) Schedule {
	return Schedule{
		Minute: minute,
		Hour:   "*",
		DOM:    "*",
		Month:  "*",
		DOW:    "*",
	}
}

// NewCustomSchedule creates a schedule with custom cron expressions
func NewCustomSchedule(minute, hour, dom, month, dow string) Schedule {
	return Schedule{
		Minute: minute,
		Hour:   hour,
		DOM:    dom,
		Month:  month,
		DOW:    dow,
	}
}
