package truenas

import (
	"context"
	"encoding/json"
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

// AppState represents the current state of an application according to TrueNAS API
type AppState string

const (
	AppStateCrashed   AppState = "CRASHED"
	AppStateDeploying AppState = "DEPLOYING"
	AppStateRunning   AppState = "RUNNING"
	AppStateStopped   AppState = "STOPPED"
)

// App represents a TrueNAS application according to official API documentation
type App struct {
	Name             string                 `json:"name"`
	ID               string                 `json:"id"`
	State            AppState               `json:"state"`
	UpgradeAvailable bool                   `json:"upgrade_available"`
	HumanVersion     string                 `json:"human_version"`
	Version          string                 `json:"version"`
	Metadata         map[string]interface{} `json:"metadata"`
	ActiveWorkloads  *AppActiveWorkloads    `json:"active_workloads"`
}

// AppActiveWorkloads represents active workloads according to TrueNAS API documentation
type AppActiveWorkloads struct {
	Containers       int                  `json:"containers"`
	UsedPorts        []AppUsedPort        `json:"used_ports"`
	ContainerDetails []AppContainerDetail `json:"container_details"`
	Volumes          []AppVolume          `json:"volumes"`
}

// AppUsedPort represents a port used by the app
type AppUsedPort struct {
	ContainerPort int32         `json:"container_port"`
	Protocol      string        `json:"protocol"`
	HostPorts     []AppHostPort `json:"host_ports"`
}

// AppHostPort represents a host port mapping
type AppHostPort struct {
	HostPort int32  `json:"host_port"`
	HostIP   string `json:"host_ip"`
}

// AppContainerDetail represents detailed container information
type AppContainerDetail struct {
	ID           string        `json:"id"`
	ServiceName  string        `json:"service_name"`
	Image        string        `json:"image"`
	PortConfig   []AppUsedPort `json:"port_config"`
	State        string        `json:"state"` // "running", "starting", "exited"
	VolumeMounts []AppVolume   `json:"volume_mounts"`
}

// AppVolume represents volume information for the app
type AppVolume struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Type        string `json:"type"`
}
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

// AppStatsOptions represents the optional parameters for app.stats
type AppStatsOptions struct {
	Interval int `json:"interval,omitempty"`
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

// AppStats represents statistics for a TrueNAS application
type AppStats struct {
	AppName  string            `json:"app_name"`
	CPUUsage float64           `json:"cpu_usage"`
	Memory   int64             `json:"memory"`
	Networks []AppStatsNetwork `json:"networks"`
	Blkio    *AppStatsBlkio    `json:"blkio"`
}

// AppStatsNetwork represents per-interface network statistics for an app
type AppStatsNetwork struct {
	InterfaceName string `json:"interface_name"`
	RXBytes       int64  `json:"rx_bytes"`
	TXBytes       int64  `json:"tx_bytes"`
}

// AppStatsBlkio represents block I/O statistics for an app
type AppStatsBlkio struct {
	Read  int64 `json:"read"`
	Write int64 `json:"write"`
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
func (a *AppClient) Get(ctx context.Context, name string, extra map[string]any) (*App, error) {
	var result []App

	params := make([]any, 0, 2)

	// filter by name
	params = append(params, []any{[]any{"name", "=", name}})

	if len(extra) > 0 {
		params = append(params, map[string]any{
			"extra": extra,
		})
	}
	err := a.client.Call(ctx, "app.query", params, &result)
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

// ListCrashed returns all applications in crashed state
func (a *AppClient) ListCrashed(ctx context.Context) ([]App, error) {
	return a.QueryByState(ctx, AppStateCrashed)
}

// Stats retrieves statistics for all applications
func (a *AppClient) SubscribeStats(ctx context.Context, fn func([]AppStats) error) error {
	return a.client.Subscribe.Subscribe(ctx, "app.stats", func(m Message) error {
		var result []AppStats
		_ = json.Unmarshal(m.Fields, &result)
		return fn(result)
	})
}

func (a *AppClient) UnsubscribeStats(ctx context.Context) error {
	return a.client.Subscribe.Unsubscribe(ctx, "app.stats")
}
