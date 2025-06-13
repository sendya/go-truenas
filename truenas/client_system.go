package truenas

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// BootEnvActive represents boot environment active state
type BootEnvActive string

const (
	BootEnvActiveNone   BootEnvActive = "-"  // Not active
	BootEnvActiveNow    BootEnvActive = "N"  // Active now
	BootEnvActiveReboot BootEnvActive = "R"  // Active on reboot
	BootEnvActiveBoth   BootEnvActive = "NR" // Active now and on reboot
)

// UpdateStatus represents update status values
type UpdateStatus string

const (
	UpdateStatusAvailable   UpdateStatus = "AVAILABLE"
	UpdateStatusUnavailable UpdateStatus = "UNAVAILABLE"
	UpdateStatusDownloaded  UpdateStatus = "DOWNLOADED"
)

// TrueNASTime handles MongoDB-style date objects from TrueNAS API
type TrueNASTime struct {
	time.Time
}

// UnmarshalJSON handles both MongoDB-style dates {"$date": timestamp} and standard JSON dates
func (t *TrueNASTime) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as MongoDB-style date object first
	var mongoDate struct {
		Date int64 `json:"$date"`
	}

	if err := json.Unmarshal(data, &mongoDate); err == nil && mongoDate.Date != 0 {
		t.Time = time.Unix(mongoDate.Date/1000, (mongoDate.Date%1000)*1000000)
		return nil
	}

	// Fall back to standard JSON time parsing
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return fmt.Errorf("unable to unmarshal time: %w", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return fmt.Errorf("unable to parse time string %q: %w", timeStr, err)
	}

	t.Time = parsedTime
	return nil
}

// MarshalJSON marshals TrueNASTime as standard JSON time
func (t TrueNASTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.RFC3339))
}

// SystemClient provides methods for system management
type SystemClient struct {
	client *Client
}

// NewSystemClient creates a new system client
func NewSystemClient(client *Client) *SystemClient {
	return &SystemClient{client: client}
}

// SystemInfo represents detailed system information
type SystemInfo struct {
	Version              string      `json:"version"`
	BuildDate            TrueNASTime `json:"buildtime"`
	Hostname             string      `json:"hostname"`
	PhysicalMemory       uint64      `json:"physmem"`
	Model                string      `json:"model"`
	Cores                int         `json:"cores"`
	PhysicalCores        int         `json:"physical_cores"`
	LoadAvg              []float64   `json:"loadavg"`
	Uptime               string      `json:"uptime"`
	UptimeSeconds        float64     `json:"uptime_seconds"`
	SystemSerial         string      `json:"system_serial"`
	SystemProduct        string      `json:"system_product"`
	SystemProductVersion string      `json:"system_product_version"`
	License              any         `json:"license"`
	BootTime             TrueNASTime `json:"boottime"`
	DateTime             TrueNASTime `json:"datetime"`
	Timezone             string      `json:"timezone"`
	SystemManufacturer   string      `json:"system_manufacturer"`
	ECC                  bool        `json:"ecc_memory"`
}

// SystemGeneralConfig represents general system configuration
type SystemGeneralConfig struct {
	ID                  int      `json:"id"`
	UIAddress           []string `json:"ui_address"`
	UIV6Address         []string `json:"ui_v6address"`
	UIPort              int      `json:"ui_port"`
	UIHTTPSPort         int      `json:"ui_httpsport"`
	UIHTTPSProtocols    []string `json:"ui_httpsprotocols"`
	UIHTTPSRedirect     bool     `json:"ui_httpsredirect"`
	UIXFrameOptions     string   `json:"ui_x_frame_options"`
	UIAllowlist         []string `json:"ui_allowlist"`
	UIConsoleMsgEnabled bool     `json:"ui_consolemsg"`
	KBDMap              string   `json:"kbdmap"`
	Language            string   `json:"language"`
	Timezone            string   `json:"timezone"`
	CrashReporting      bool     `json:"crash_reporting"`
	UsageCollection     bool     `json:"usage_collection"`
	Birthday            any      `json:"birthday"`
	WizardShown         bool     `json:"wizardshown"`
	DSAuth              bool     `json:"ds_auth"`
}

// BootEnv represents boot environment information
type BootEnv struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Active     BootEnvActive `json:"active"`
	Mountpoint string        `json:"mountpoint"`
	Space      string        `json:"space"`
	Created    TrueNASTime   `json:"created"`
	Keep       bool          `json:"keep"`
	RawSpace   int64         `json:"rawspace"`
}

// UpdateConfig represents system update configuration
type UpdateConfig struct {
	AutoCheck bool   `json:"autocheck"`
	Train     string `json:"train"`
}

// UpdateInfo represents update information
type UpdateInfo struct {
	Status     UpdateStatus `json:"status"`
	Available  bool         `json:"available"`
	ChangeLog  string       `json:"changelog"`
	FileName   string       `json:"filename"`
	Version    string       `json:"version"`
	Notes      string       `json:"notes"`
	Notice     string       `json:"notice"`
	Downloaded bool         `json:"downloaded"`
	Error      *string      `json:"error,omitempty"`
}

// GetInfo returns system information
func (s *SystemClient) GetInfo(ctx context.Context) (*SystemInfo, error) {
	var result SystemInfo
	err := s.client.Call(ctx, "system.info", []any{}, &result)
	return &result, err
}

// GetGeneralConfig returns general system configuration
func (s *SystemClient) GetGeneralConfig(ctx context.Context) (*SystemGeneralConfig, error) {
	var result SystemGeneralConfig
	err := s.client.Call(ctx, "system.general.config", []any{}, &result)
	return &result, err
}

// UpdateGeneralConfig updates general system configuration
func (s *SystemClient) UpdateGeneralConfig(ctx context.Context, config *SystemGeneralConfig) (*SystemGeneralConfig, error) {
	var result SystemGeneralConfig
	err := s.client.Call(ctx, "system.general.update", []any{*config}, &result)
	return &result, err
}

// Reboot reboots the system
func (s *SystemClient) Reboot(ctx context.Context, delay int) error {
	params := []any{}
	if delay > 0 {
		params = append(params, map[string]any{"delay": delay})
	}
	return s.client.CallJob(ctx, "system.reboot", params, nil)
}

// Shutdown shuts down the system
func (s *SystemClient) Shutdown(ctx context.Context, delay int) error {
	params := []any{}
	if delay > 0 {
		params = append(params, map[string]any{"delay": delay})
	}
	return s.client.CallJob(ctx, "system.shutdown", params, nil)
}

// Ready checks if the system is ready
func (s *SystemClient) Ready(ctx context.Context) (bool, error) {
	var result bool
	err := s.client.Call(ctx, "system.ready", []any{}, &result)
	return result, err
}

// GetVersion returns system version
func (s *SystemClient) GetVersion(ctx context.Context) (string, error) {
	var result string
	err := s.client.Call(ctx, "system.version", []any{}, &result)
	return result, err
}

// GetHostname returns system hostname
func (s *SystemClient) GetHostname(ctx context.Context) (string, error) {
	var result string
	err := s.client.Call(ctx, "system.hostname", []any{}, &result)
	return result, err
}

// SetHostname sets system hostname
func (s *SystemClient) SetHostname(ctx context.Context, hostname string) error {
	return s.client.Call(ctx, "system.hostname", []any{hostname}, nil)
}

// Boot Environment Methods

// ListBootEnvs returns all boot environments
func (s *SystemClient) ListBootEnvs(ctx context.Context) ([]BootEnv, error) {
	var result []BootEnv
	err := s.client.Call(ctx, "bootenv.query", []any{}, &result)
	return result, err
}

// CreateBootEnv creates a new boot environment
func (s *SystemClient) CreateBootEnv(ctx context.Context, name, source string) (*BootEnv, error) {
	var result BootEnv
	params := map[string]any{"name": name}
	if source != "" {
		params["source"] = source
	}
	err := s.client.Call(ctx, "bootenv.create", []any{params}, &result)
	return &result, err
}

// DeleteBootEnv deletes a boot environment
func (s *SystemClient) DeleteBootEnv(ctx context.Context, id string) error {
	return s.client.CallJob(ctx, "bootenv.delete", []any{id}, nil)
}

// ActivateBootEnv activates a boot environment
func (s *SystemClient) ActivateBootEnv(ctx context.Context, id string) error {
	return s.client.Call(ctx, "bootenv.activate", []any{id}, nil)
}

// SetBootEnvAttr sets boot environment attributes
func (s *SystemClient) SetBootEnvAttr(ctx context.Context, id string, attrs map[string]any) error {
	return s.client.Call(ctx, "bootenv.set_attribute", []any{id, attrs}, nil)
}

// Update Methods

// GetUpdateConfig returns update configuration
func (s *SystemClient) GetUpdateConfig(ctx context.Context) (*UpdateConfig, error) {
	var result UpdateConfig
	err := s.client.Call(ctx, "update.config", []any{}, &result)
	return &result, err
}

// CheckForUpdate checks for available updates
func (s *SystemClient) CheckForUpdate(ctx context.Context) (*UpdateInfo, error) {
	var result UpdateInfo
	err := s.client.Call(ctx, "update.check_available", []any{}, &result)
	return &result, err
}

// GetPendingUpdate returns pending update information
func (s *SystemClient) GetPendingUpdate(ctx context.Context) (*UpdateInfo, error) {
	var result UpdateInfo
	err := s.client.Call(ctx, "update.get_pending", []any{}, &result)
	return &result, err
}

// DownloadUpdate downloads available updates
func (s *SystemClient) DownloadUpdate(ctx context.Context) error {
	return s.client.CallJob(ctx, "update.download", []any{}, nil)
}

// ManualUpdate performs manual update from uploaded file
func (s *SystemClient) ManualUpdate(ctx context.Context, path string, rebootAfter bool) error {
	params := []any{path}
	if rebootAfter {
		params = append(params, map[string]any{"reboot_after": true})
	}
	return s.client.CallJob(ctx, "update.manual", params, nil)
}

// GetTrains returns available update trains
func (s *SystemClient) GetTrains(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := s.client.Call(ctx, "update.get_trains", []any{}, &result)
	return result, err
}

// SetTrain sets the update train
func (s *SystemClient) SetTrain(ctx context.Context, train string) error {
	return s.client.Call(ctx, "update.set_train", []any{train}, nil)
}
