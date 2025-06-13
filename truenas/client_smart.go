package truenas

import (
	"context"
	"fmt"
)

// SmartClient provides methods for SMART monitoring and testing
type SmartClient struct {
	client *Client
}

// NewSmartClient creates a new SMART client
func NewSmartClient(client *Client) *SmartClient {
	return &SmartClient{client: client}
}

// SmartConfig represents SMART service configuration
type SmartConfig struct {
	Interval      int    `json:"interval"`
	Powermode     string `json:"powermode"`
	Difference    int    `json:"difference"`
	Informational int    `json:"informational"`
	Critical      int    `json:"critical"`
}

// SmartTest represents a SMART test task
type SmartTest struct {
	ID       int               `json:"id"`
	Schedule SmartTestSchedule `json:"schedule"`
	Desc     string            `json:"desc"`
	AllDisks bool              `json:"all_disks"`
	Disks    []string          `json:"disks"`
	Type     string            `json:"type"`
}

// SmartTestSchedule represents the cron schedule for a SMART test
type SmartTestSchedule struct {
	Minute string `json:"minute,omitempty"`
	Hour   string `json:"hour"`
	DOM    string `json:"dom"`
	Month  string `json:"month"`
	DOW    string `json:"dow"`
}

// SmartTestCreateRequest represents parameters for smart.test.create
type SmartTestCreateRequest struct {
	Schedule SmartTestSchedule `json:"schedule"`
	Desc     string            `json:"desc,omitempty"`
	AllDisks bool              `json:"all_disks"`
	Disks    []string          `json:"disks,omitempty"`
	Type     string            `json:"type"`
}

// SmartTestResult represents SMART test results for a disk
type SmartTestResult struct {
	Disk  string      `json:"disk"`
	Tests []SmartTest `json:"tests"`
}

// SmartTestDetail represents individual test details
type SmartTestDetail struct {
	Num             int     `json:"num"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
	StatusVerbose   string  `json:"status_verbose"`
	Remaining       float64 `json:"remaining"`
	Lifetime        int     `json:"lifetime"`
	LBAOfFirstError any     `json:"lba_of_first_error"`
	SegmentNumber   any     `json:"segment_number,omitempty"`
}

// SmartManualTestRequest represents parameters for running manual SMART tests
type SmartManualTestRequest struct {
	Disk string `json:"disk"`
	Type string `json:"type"`
}

// SmartAttributes represents SMART attributes for a disk
type SmartAttributes struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Threshold  int    `json:"threshold"`
	Type       string `json:"type"`
	Updated    string `json:"updated"`
	WhenFailed string `json:"when_failed"`
	RawValue   any    `json:"raw_value"`
}

// SmartTestType represents available SMART test types
type SmartTestType string

const (
	SmartTestTypeLong       SmartTestType = "LONG"
	SmartTestTypeShort      SmartTestType = "SHORT"
	SmartTestTypeConveyance SmartTestType = "CONVEYANCE"
	SmartTestTypeOffline    SmartTestType = "OFFLINE"
)

// SmartPowerMode represents SMART power modes
type SmartPowerMode string

const (
	SmartPowerModeNever   SmartPowerMode = "NEVER"
	SmartPowerModeSleep   SmartPowerMode = "SLEEP"
	SmartPowerModeStandby SmartPowerMode = "STANDBY"
	SmartPowerModeIdle    SmartPowerMode = "IDLE"
)

// GetConfig returns SMART service configuration
func (s *SmartClient) GetConfig(ctx context.Context) (*SmartConfig, error) {
	var result SmartConfig
	err := s.client.Call(ctx, "smart.config", []any{}, &result)
	return &result, err
}

// UpdateConfig updates SMART service configuration
func (s *SmartClient) UpdateConfig(ctx context.Context, config *SmartConfig) (*SmartConfig, error) {
	var result SmartConfig
	err := s.client.Call(ctx, "smart.update", []any{*config}, &result)
	return &result, err
}

// SMART Test Management

// ListTests returns all SMART test tasks
func (s *SmartClient) ListTests(ctx context.Context) ([]SmartTest, error) {
	var result []SmartTest
	err := s.client.Call(ctx, "smart.test.query", []any{}, &result)
	return result, err
}

// GetTest returns a specific SMART test by ID
func (s *SmartClient) GetTest(ctx context.Context, id int) (*SmartTest, error) {
	var result []SmartTest
	err := s.client.Call(ctx, "smart.test.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("smart_test", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// CreateTest creates a new SMART test task
func (s *SmartClient) CreateTest(ctx context.Context, req *SmartTestCreateRequest) (*SmartTest, error) {
	var result SmartTest
	err := s.client.Call(ctx, "smart.test.create", []any{*req}, &result)
	return &result, err
}

// UpdateTest updates an existing SMART test task
func (s *SmartClient) UpdateTest(ctx context.Context, id int, req *SmartTestCreateRequest) (*SmartTest, error) {
	var result SmartTest
	err := s.client.Call(ctx, "smart.test.update", []any{id, *req}, &result)
	return &result, err
}

// DeleteTest deletes a SMART test task
func (s *SmartClient) DeleteTest(ctx context.Context, id int) error {
	return s.client.Call(ctx, "smart.test.delete", []any{id}, nil)
}

// GetDiskChoices returns available disk choices for SMART tests
func (s *SmartClient) GetDiskChoices(ctx context.Context, fullDisk bool) (any, error) {
	var result any
	err := s.client.Call(ctx, "smart.test.disk_choices", []any{fullDisk}, &result)
	return result, err
}

// RunManualTest runs manual SMART tests for specified disks
func (s *SmartClient) RunManualTest(ctx context.Context, tests []SmartManualTestRequest) error {
	return s.client.Call(ctx, "smart.test.manual_test", []any{tests}, nil)
}

// SMART Test Results

// GetAllTestResults returns SMART test results for all disks
func (s *SmartClient) GetAllTestResults(ctx context.Context) ([]SmartTestResult, error) {
	var result []SmartTestResult
	err := s.client.Call(ctx, "smart.test.results", []any{}, &result)
	return result, err
}

// GetDiskTestResults returns SMART test results for a specific disk
func (s *SmartClient) GetDiskTestResults(ctx context.Context, diskName string) (*SmartTestResult, error) {
	var result SmartTestResult
	err := s.client.Call(ctx, "smart.test.results", []any{[]any{[]any{"disk", "=", diskName}}, map[string]any{"get": true}}, &result)
	return &result, err
}

// SMART Attributes

// GetDiskAttributes returns SMART attributes for a specific disk
func (s *SmartClient) GetDiskAttributes(ctx context.Context, diskName string) ([]SmartAttributes, error) {
	var result []SmartAttributes
	err := s.client.Call(ctx, "disk.smart_attributes", []any{diskName}, &result)
	return result, err
}
