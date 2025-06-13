package truenas

import (
	"context"
	"fmt"
	"time"
)

// AlertClient provides methods for alert management
type AlertClient struct {
	client *Client
}

// NewAlertClient creates a new alert client
func NewAlertClient(client *Client) *AlertClient {
	return &AlertClient{client: client}
}

// Alert represents a system alert
type Alert struct {
	UUID           string    `json:"uuid"`
	Source         string    `json:"source"`
	Klass          string    `json:"klass"`
	Args           any       `json:"args"`
	Node           string    `json:"node"`
	Key            string    `json:"key"`
	DateTime       time.Time `json:"datetime"`
	LastOccurrence time.Time `json:"last_occurrence"`
	Dismissed      bool      `json:"dismissed"`
	Mail           any       `json:"mail"`
	Text           string    `json:"text"`
	Level          string    `json:"level"`
	OneShot        bool      `json:"one_shot"`
	Formatted      string    `json:"formatted"`
}

// AlertCategory represents an alert category/class
type AlertCategory struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Level            string `json:"level"`
	Category         string `json:"category"`
	Description      string `json:"description"`
	ProactiveSupport bool   `json:"proactive_support"`
}

// AlertPolicy represents alert frequency policy
type AlertPolicy struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AlertService represents an alert service configuration
type AlertService struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes"`
	Level      string         `json:"level"`
	Enabled    bool           `json:"enabled"`
}

// AlertServiceCreateRequest represents parameters for alertservice.create
type AlertServiceCreateRequest struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes"`
	Level      string         `json:"level"`
	Enabled    bool           `json:"enabled"`
}

// AlertServiceUpdateRequest represents parameters for alertservice.update
type AlertServiceUpdateRequest struct {
	Name       string         `json:"name,omitempty"`
	Type       string         `json:"type,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
	Level      string         `json:"level,omitempty"`
	Enabled    bool           `json:"enabled,omitempty"`
}

// AlertClassesUpdateRequest represents parameters for alertclasses.update
type AlertClassesUpdateRequest struct {
	Classes map[string]any `json:"classes"`
}

// AlertLevel represents alert severity levels
type AlertLevel string

const (
	AlertLevelInfo      AlertLevel = "INFO"
	AlertLevelNotice    AlertLevel = "NOTICE"
	AlertLevelWarning   AlertLevel = "WARNING"
	AlertLevelError     AlertLevel = "ERROR"
	AlertLevelCritical  AlertLevel = "CRITICAL"
	AlertLevelAlert     AlertLevel = "ALERT"
	AlertLevelEmergency AlertLevel = "EMERGENCY"
)

// Alert Management

// List returns all alerts (active and dismissed)
func (a *AlertClient) List(ctx context.Context) ([]Alert, error) {
	var result []Alert
	err := a.client.Call(ctx, "alert.list", []any{}, &result)
	return result, err
}

// Dismiss dismisses an alert by UUID
func (a *AlertClient) Dismiss(ctx context.Context, uuid string) error {
	return a.client.Call(ctx, "alert.dismiss", []any{uuid}, nil)
}

// Restore restores a dismissed alert by UUID
func (a *AlertClient) Restore(ctx context.Context, uuid string) error {
	return a.client.Call(ctx, "alert.restore", []any{uuid}, nil)
}

// ListCategories returns all available alert categories/classes
func (a *AlertClient) ListCategories(ctx context.Context) ([]AlertCategory, error) {
	var result []AlertCategory
	err := a.client.Call(ctx, "alert.list_categories", []any{}, &result)
	return result, err
}

// ListPolicies returns all available alert policies
func (a *AlertClient) ListPolicies(ctx context.Context) ([]AlertPolicy, error) {
	var result []AlertPolicy
	err := a.client.Call(ctx, "alert.list_policies", []any{}, &result)
	return result, err
}

// Alert Classes Configuration

// GetAlertClassesConfig returns current alert classes configuration
func (a *AlertClient) GetAlertClassesConfig(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := a.client.Call(ctx, "alertclasses.config", []any{}, &result)
	return result, err
}

// UpdateAlertClasses updates alert classes configuration
func (a *AlertClient) UpdateAlertClasses(ctx context.Context, req *AlertClassesUpdateRequest) (map[string]any, error) {
	var result map[string]any
	err := a.client.Call(ctx, "alertclasses.update", []any{*req}, &result)
	return result, err
}

// Alert Services Management

// AlertServiceClient provides methods for alert service management
type AlertServiceClient struct {
	client *Client
}

// NewAlertServiceClient creates a new alert service client
func NewAlertServiceClient(client *Client) *AlertServiceClient {
	return &AlertServiceClient{client: client}
}

// List returns all alert services
func (s *AlertServiceClient) List(ctx context.Context) ([]AlertService, error) {
	var result []AlertService
	err := s.client.Call(ctx, "alertservice.query", []any{}, &result)
	return result, err
}

// Get returns a specific alert service by ID
func (s *AlertServiceClient) Get(ctx context.Context, id int) (*AlertService, error) {
	var result []AlertService
	err := s.client.Call(ctx, "alertservice.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("alert_service", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new alert service
func (s *AlertServiceClient) Create(ctx context.Context, req *AlertServiceCreateRequest) (*AlertService, error) {
	var result AlertService
	err := s.client.Call(ctx, "alertservice.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing alert service
func (s *AlertServiceClient) Update(ctx context.Context, id int, req *AlertServiceUpdateRequest) (*AlertService, error) {
	var result AlertService
	err := s.client.Call(ctx, "alertservice.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes an alert service
func (s *AlertServiceClient) Delete(ctx context.Context, id int) error {
	return s.client.Call(ctx, "alertservice.delete", []any{id}, nil)
}

// Test sends a test alert using the specified service configuration
func (s *AlertServiceClient) Test(ctx context.Context, req *AlertServiceCreateRequest) error {
	return s.client.Call(ctx, "alertservice.test", []any{*req}, nil)
}

// ListTypes returns all available alert service types
func (s *AlertServiceClient) ListTypes(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := s.client.Call(ctx, "alertservice.list_types", []any{}, &result)
	return result, err
}
