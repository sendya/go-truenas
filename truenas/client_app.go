package truenas

import (
	"context"
	"fmt"
)

// AppClient provides methods for application management
type AppClient struct {
	client *Client
}

// NewAppClient creates a new app client
func NewAppClient(client *Client) *AppClient {
	return &AppClient{client: client}
}

// AppState represents the current state of an application
type AppState string

const (
	AppStateRunning   AppState = "RUNNING"
	AppStateStopped   AppState = "STOPPED"
	AppStateDeploying AppState = "DEPLOYING"
	AppStateError     AppState = "ERROR"
	AppStateUpgrading AppState = "UPGRADING"
	AppStatePending   AppState = "PENDING"
)

// App represents a TrueNAS application
type App struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	State        AppState               `json:"state"`
	Version      string                 `json:"version"`
	Upgrade      *bool                  `json:"upgrade,omitempty"`
	ChartName    string                 `json:"chart_name"`
	Namespace    string                 `json:"namespace"`
	Catalog      string                 `json:"catalog"`
	CatalogTrain string                 `json:"catalog_train"`
	Config       map[string]interface{} `json:"config,omitempty"`
	History      []AppHistory           `json:"history,omitempty"`
	Resources    *AppResources          `json:"resources,omitempty"`
	Metadata     *AppMetadata           `json:"metadata,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Icon         string                 `json:"icon,omitempty"`
	Maintainers  []AppMaintainer        `json:"maintainers,omitempty"`
	Sources      []string               `json:"sources,omitempty"`
	Home         string                 `json:"home,omitempty"`
	Keywords     []string               `json:"keywords,omitempty"`
}

// AppHistory represents an app's deployment history entry
type AppHistory struct {
	Revision    int                    `json:"revision"`
	UpdatedAt   *TrueNASTime           `json:"updated"`
	Status      string                 `json:"status"`
	Chart       string                 `json:"chart"`
	AppVersion  string                 `json:"app_version"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// AppResources represents resource usage information for an app
type AppResources struct {
	Containers []AppContainer `json:"containers,omitempty"`
	Storage    []AppStorage   `json:"storage,omitempty"`
}

// AppContainer represents a container within an app
type AppContainer struct {
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	CPU     *AppResourceUsage `json:"cpu,omitempty"`
	Memory  *AppResourceUsage `json:"memory,omitempty"`
	Restart int               `json:"restart_count"`
}

// AppStorage represents storage information for an app
type AppStorage struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Path      string `json:"path,omitempty"`
	Size      string `json:"size,omitempty"`
	Used      string `json:"used,omitempty"`
	Available string `json:"available,omitempty"`
}

// AppResourceUsage represents CPU/Memory usage
type AppResourceUsage struct {
	Used  string `json:"used"`
	Limit string `json:"limit,omitempty"`
}

// AppMetadata represents app metadata information
type AppMetadata struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	AppVersion  string            `json:"app_version"`
	Description string            `json:"description"`
	Home        string            `json:"home"`
	Sources     []string          `json:"sources,omitempty"`
	Maintainers []AppMaintainer   `json:"maintainers,omitempty"`
	Icon        string            `json:"icon,omitempty"`
	Keywords    []string          `json:"keywords,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// AppMaintainer represents an app maintainer
type AppMaintainer struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// AppQueryOptions represents query options for filtering apps
type AppQueryOptions struct {
	ExtraOptions map[string]interface{} `json:"extra,omitempty"`
}

// AppCreateRequest represents parameters for app.create
type AppCreateRequest struct {
	ReleaseName  string                 `json:"release_name"`
	ChartRelease string                 `json:"chart_release"`
	Values       map[string]interface{} `json:"values,omitempty"`
}

// AppUpdateRequest represents parameters for app.update
type AppUpdateRequest struct {
	Values map[string]interface{} `json:"values,omitempty"`
}

// List returns all applications
func (a *AppClient) List(ctx context.Context) ([]App, error) {
	var result []App
	err := a.client.Call(ctx, "app.query", []any{}, &result)
	return result, err
}

// ListWithOptions returns applications with custom query options
func (a *AppClient) ListWithOptions(ctx context.Context, options *AppQueryOptions) ([]App, error) {
	var result []App
	params := []any{}
	if options != nil {
		params = append(params, []any{}, options)
	}
	err := a.client.Call(ctx, "app.query", params, &result)
	return result, err
}

// Get returns a specific application by name
func (a *AppClient) Get(ctx context.Context, name string) (*App, error) {
	var result []App
	err := a.client.Call(ctx, "app.query", []any{[]any{[]any{"name", "=", name}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("app", fmt.Sprintf("name %s", name))
	}
	return &result[0], nil
}

// GetByID returns a specific application by ID
func (a *AppClient) GetByID(ctx context.Context, id string) (*App, error) {
	var result []App
	err := a.client.Call(ctx, "app.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("app", fmt.Sprintf("ID %s", id))
	}
	return &result[0], nil
}

// QueryByState returns applications filtered by state
func (a *AppClient) QueryByState(ctx context.Context, state AppState) ([]App, error) {
	var result []App
	err := a.client.Call(ctx, "app.query", []any{[]any{[]any{"state", "=", string(state)}}}, &result)
	return result, err
}

// QueryByCatalog returns applications filtered by catalog
func (a *AppClient) QueryByCatalog(ctx context.Context, catalog string) ([]App, error) {
	var result []App
	err := a.client.Call(ctx, "app.query", []any{[]any{[]any{"catalog", "=", catalog}}}, &result)
	return result, err
}

// QueryWithFilters returns applications with custom filters
// filters should be in the format: [["field", "operator", "value"], ...]
// Example: [["state", "=", "RUNNING"], ["catalog", "=", "TRUENAS"]]
func (a *AppClient) QueryWithFilters(ctx context.Context, filters [][]any, options *AppQueryOptions) ([]App, error) {
	var result []App
	params := []any{filters}
	if options != nil {
		params = append(params, options)
	}
	err := a.client.Call(ctx, "app.query", params, &result)
	return result, err
}

// ListRunning returns all running applications
func (a *AppClient) ListRunning(ctx context.Context) ([]App, error) {
	return a.QueryByState(ctx, AppStateRunning)
}

// ListStopped returns all stopped applications
func (a *AppClient) ListStopped(ctx context.Context) ([]App, error) {
	return a.QueryByState(ctx, AppStateStopped)
}

// ListDeploying returns all deploying applications
func (a *AppClient) ListDeploying(ctx context.Context) ([]App, error) {
	return a.QueryByState(ctx, AppStateDeploying)
}

// ListWithErrors returns all applications in error state
func (a *AppClient) ListWithErrors(ctx context.Context) ([]App, error) {
	return a.QueryByState(ctx, AppStateError)
}
