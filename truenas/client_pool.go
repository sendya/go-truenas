package truenas

import (
	"context"
	"fmt"
)

// PoolStatus represents pool status values
type PoolStatus string

const (
	PoolStatusOnline   PoolStatus = "ONLINE"
	PoolStatusDegraded PoolStatus = "DEGRADED"
	PoolStatusFaulted  PoolStatus = "FAULTED"
	PoolStatusOffline  PoolStatus = "OFFLINE"
	PoolStatusUnavail  PoolStatus = "UNAVAIL"
	PoolStatusRemoved  PoolStatus = "REMOVED"
)

// PoolScrubAction represents scrub action values
type PoolScrubAction string

const (
	PoolScrubActionStart PoolScrubAction = "START"
	PoolScrubActionStop  PoolScrubAction = "STOP"
	PoolScrubActionPause PoolScrubAction = "PAUSE"
)

// VDevType represents virtual device types
type VDevType string

const (
	VDevTypeDisk   VDevType = "DISK"
	VDevTypeStripe VDevType = "STRIPE"
	VDevTypeMirror VDevType = "MIRROR"
	VDevTypeRaidz  VDevType = "RAIDZ"
	VDevTypeRaidz1 VDevType = "RAIDZ1"
	VDevTypeRaidz2 VDevType = "RAIDZ2"
	VDevTypeRaidz3 VDevType = "RAIDZ3"
	VDevTypeSpare  VDevType = "SPARE"
	VDevTypeLog    VDevType = "LOG"
	VDevTypeCache  VDevType = "CACHE"
)

// PoolClient provides methods for pool management
type PoolClient struct {
	client *Client
}

// NewPoolClient creates a new pool client
func NewPoolClient(client *Client) *PoolClient {
	return &PoolClient{client: client}
}

// Pool represents a storage pool
type Pool struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	GUID         string        `json:"guid"`
	Status       PoolStatus    `json:"status"`
	Path         string        `json:"path"`
	Scan         *PoolScan     `json:"scan"`
	Healthy      bool          `json:"healthy"`
	Warning      bool          `json:"warning"`
	Topology     *PoolTopology `json:"topology"`
	Encrypt      int           `json:"encrypt"`
	IsUpgraded   bool          `json:"is_upgraded"`
	StatusDetail string        `json:"status_detail"`
	Autotrim     *PoolProperty `json:"autotrim"`
	IsDecrypted  bool          `json:"is_decrypted"`
}

// PoolScan represents pool scan information
type PoolScan struct {
	Function       string       `json:"function"`
	State          string       `json:"state"`
	StartTime      *TrueNASTime `json:"start_time"`
	EndTime        *TrueNASTime `json:"end_time"`
	Percentage     float64      `json:"percentage"`
	BytesToProcess *int64       `json:"bytes_to_process"`
	BytesProcessed *int64       `json:"bytes_processed"`
	BytesIssued    *int64       `json:"bytes_issued"`
	Errors         int          `json:"errors"`
}

// PoolTopology represents pool topology structure (for reading)
type PoolTopology struct {
	Data    []VDev `json:"data"`
	Cache   []VDev `json:"cache,omitempty"`
	Log     []VDev `json:"log,omitempty"`
	Spare   []VDev `json:"spare,omitempty"`
	Special []VDev `json:"special,omitempty"`
	Dedup   []VDev `json:"dedup,omitempty"`
}

// PoolTopologyCreate represents pool topology structure for creation
type PoolTopologyCreate struct {
	Data    []VDevCreate `json:"data"`
	Cache   []VDevCreate `json:"cache,omitempty"`
	Log     []VDevCreate `json:"log,omitempty"`
	Spare   []VDevCreate `json:"spare,omitempty"`
	Special []VDevCreate `json:"special,omitempty"`
	Dedup   []VDevCreate `json:"dedup,omitempty"`
}

// VDev represents a virtual device (for reading topology)
type VDev struct {
	Type     VDevType   `json:"type"`
	Children []VDev     `json:"children,omitempty"`
	Disk     *string    `json:"disk,omitempty"`
	GUID     *string    `json:"guid,omitempty"`
	Status   *string    `json:"status,omitempty"`
	Stats    *VDevStats `json:"stats,omitempty"`
}

// VDevCreate represents a virtual device for pool creation
type VDevCreate struct {
	Type  VDevType `json:"type"`
	Disks []string `json:"disks"`
}

// VDevStats represents virtual device statistics
type VDevStats struct {
	ReadErrors     int64 `json:"read_errors"`
	WriteErrors    int64 `json:"write_errors"`
	ChecksumErrors int64 `json:"checksum_errors"`
	Timestamp      int64 `json:"timestamp"`
	Configured     bool  `json:"configured"`
}

// EncryptionOptions represents encryption configuration
type EncryptionOptions struct {
	Generate    bool    `json:"generate,omitempty"`
	Pbkdf2iters *int    `json:"pbkdf2iters,omitempty"`
	Algorithm   *string `json:"algorithm,omitempty"`
	Passphrase  *string `json:"passphrase,omitempty"`
	Key         *string `json:"key,omitempty"`
	KeyFormat   *string `json:"key_format,omitempty"`
	KeyLocation *string `json:"key_location,omitempty"`
}

// PoolCreateRequest represents parameters for pool.create
type PoolCreateRequest struct {
	Name              string             `json:"name"`
	Encryption        bool               `json:"encryption,omitempty"`
	Topology          PoolTopologyCreate `json:"topology"`
	EncryptionOptions *EncryptionOptions `json:"encryption_options,omitempty"`
	Deduplication     *string            `json:"deduplication,omitempty"`
	Checksum          *string            `json:"checksum,omitempty"`
}

// PoolUpdateRequest represents parameters for pool.update
type PoolUpdateRequest struct {
	Topology          *PoolTopology      `json:"topology,omitempty"`
	EncryptionOptions *EncryptionOptions `json:"encryption_options,omitempty"`
	Autotrim          *string            `json:"autotrim,omitempty"`
	Comments          *string            `json:"comments,omitempty"`
	Deduplication     *string            `json:"deduplication,omitempty"`
}

// PoolScrubRequest represents parameters for pool.scrub
type PoolScrubRequest struct {
	Action    PoolScrubAction `json:"action"`
	Threshold int             `json:"threshold,omitempty"`
}

// PoolExportRequest represents parameters for pool.export
type PoolExportRequest struct {
	Cascade        bool `json:"cascade,omitempty"`
	RestartService bool `json:"restart_services,omitempty"`
	Destroy        bool `json:"destroy,omitempty"`
}

// PoolImportRequest represents parameters for pool.import
type PoolImportRequest struct {
	GUID              string `json:"guid,omitempty"`
	Name              string `json:"name,omitempty"`
	Passphrase        string `json:"passphrase,omitempty"`
	EnableAttachments bool   `json:"enable_attachments"`
}

// PoolImportFindResult represents a pool available for import
type PoolImportFindResult struct {
	GUID   string `json:"guid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// PoolProcess represents a process using a pool
type PoolProcess struct {
	PID     int    `json:"pid"`
	Name    string `json:"name"`
	CmdLine string `json:"cmdline"`
	Service string `json:"service"`
}

// PoolProperty represents a pool property with its value and metadata
type PoolProperty struct {
	Parsed   any    `json:"parsed"`
	RawValue string `json:"rawvalue"`
	Value    string `json:"value"`
	Source   string `json:"source"`
}

// PoolScrubTask represents a scheduled scrub task
type PoolScrubTask struct {
	ID          int          `json:"id"`
	Pool        int          `json:"pool"`
	Threshold   int          `json:"threshold"`
	Description string       `json:"description"`
	Schedule    CronSchedule `json:"schedule"`
	Enabled     bool         `json:"enabled"`
}

// CronSchedule represents a cron schedule
type CronSchedule struct {
	Minute string `json:"minute"`
	Hour   string `json:"hour"`
	Dom    string `json:"dom"`
	Month  string `json:"month"`
	Dow    string `json:"dow"`
}

// PoolScrubTaskRequest represents parameters for creating/updating scrub tasks
type PoolScrubTaskRequest struct {
	Pool        int          `json:"pool"`
	Threshold   int          `json:"threshold"`
	Description string       `json:"description"`
	Schedule    CronSchedule `json:"schedule"`
	Enabled     bool         `json:"enabled"`
}

// List returns all storage pools
func (p *PoolClient) List(ctx context.Context) ([]Pool, error) {
	var result []Pool
	err := p.client.Call(ctx, "pool.query", []any{}, &result)
	return result, err
}

// Get returns a specific pool by ID
func (p *PoolClient) Get(ctx context.Context, id int) (*Pool, error) {
	var result []Pool
	err := p.client.Call(ctx, "pool.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("pool", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// GetByName returns a specific pool by name
func (p *PoolClient) GetByName(ctx context.Context, name string) (*Pool, error) {
	var result []Pool
	err := p.client.Call(ctx, "pool.query", []any{[]any{[]any{"name", "=", name}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("pool", fmt.Sprintf("name %s", name))
	}
	return &result[0], nil
}

// Create creates a new storage pool
func (p *PoolClient) Create(ctx context.Context, req PoolCreateRequest) (*Pool, error) {
	var result Pool
	err := p.client.CallJob(ctx, "pool.create", []any{req}, &result)
	return &result, err
}

// Update updates an existing pool
func (p *PoolClient) Update(ctx context.Context, id int, req PoolUpdateRequest) (*Pool, error) {
	var result Pool
	err := p.client.CallJob(ctx, "pool.update", []any{id, req}, &result)
	return &result, err
}

// Delete permanently destroys a pool and all its data
func (p *PoolClient) Delete(ctx context.Context, id int, cascade bool) error {
	options := map[string]any{
		"destroy": true,
	}
	if cascade {
		options["cascade"] = true
	}
	return p.client.CallJob(ctx, "pool.export", []any{id, options}, nil)
}

// Export exports a pool
func (p *PoolClient) Export(ctx context.Context, id int, req PoolExportRequest) error {
	return p.client.CallJob(ctx, "pool.export", []any{id, req}, nil)
}

// Import imports a pool
func (p *PoolClient) Import(ctx context.Context, req PoolImportRequest) (*Pool, error) {
	var result Pool
	err := p.client.CallJob(ctx, "pool.import_pool", []any{req}, &result)
	return &result, err
}

// FindImportablePools returns pools available for import
func (p *PoolClient) FindImportablePools(ctx context.Context) ([]PoolImportFindResult, error) {
	var result []PoolImportFindResult
	err := p.client.CallJob(ctx, "pool.import_find", []any{}, &result)
	return result, err
}

// Scrub starts, stops, or pauses a pool scrub operation
func (p *PoolClient) Scrub(ctx context.Context, id int, action PoolScrubAction) error {
	options := map[string]any{
		"action": action,
	}
	return p.client.CallJob(ctx, "pool.scrub", []any{id, options}, nil)
}

// GetProcesses returns processes using the pool
func (p *PoolClient) GetProcesses(ctx context.Context, id int) ([]PoolProcess, error) {
	var result []PoolProcess
	err := p.client.Call(ctx, "pool.processes", []any{id}, &result)
	return result, err
}

// Scrub Management Methods

// ListScrubTasks returns all scheduled scrub tasks
func (p *PoolClient) ListScrubTasks(ctx context.Context) ([]PoolScrubTask, error) {
	var result []PoolScrubTask
	err := p.client.Call(ctx, "pool.scrub.query", []any{}, &result)
	return result, err
}

// GetScrubTask returns a specific scrub task by ID
func (p *PoolClient) GetScrubTask(ctx context.Context, id int) (*PoolScrubTask, error) {
	var result []PoolScrubTask
	err := p.client.Call(ctx, "pool.scrub.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("scrub_task", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// GetScrubTasksByPool returns all scrub tasks for a specific pool
func (p *PoolClient) GetScrubTasksByPool(ctx context.Context, poolID int) ([]PoolScrubTask, error) {
	var result []PoolScrubTask
	err := p.client.Call(ctx, "pool.scrub.query", []any{[]any{[]any{"pool", "=", poolID}}}, &result)
	return result, err
}

// CreateScrubTask creates a new scheduled scrub task
func (p *PoolClient) CreateScrubTask(ctx context.Context, req PoolScrubTaskRequest) (*PoolScrubTask, error) {
	var result PoolScrubTask
	err := p.client.Call(ctx, "pool.scrub.create", []any{req}, &result)
	return &result, err
}

// UpdateScrubTask updates an existing scrub task
func (p *PoolClient) UpdateScrubTask(ctx context.Context, id int, req PoolScrubTaskRequest) (*PoolScrubTask, error) {
	var result PoolScrubTask
	err := p.client.Call(ctx, "pool.scrub.update", []any{id, req}, &result)
	return &result, err
}

// DeleteScrubTask deletes a scheduled scrub task
func (p *PoolClient) DeleteScrubTask(ctx context.Context, id int) error {
	return p.client.Call(ctx, "pool.scrub.delete", []any{id}, nil)
}

// RunScrub runs a scrub operation on a pool and waits for completion
func (p *PoolClient) RunScrub(ctx context.Context, poolName, action string) error {
	return p.client.CallJob(ctx, "pool.scrub.scrub", []any{poolName, action}, nil)
}

// RunScrubAsync runs a scrub operation on a pool and returns the job ID for monitoring
func (p *PoolClient) RunScrubAsync(ctx context.Context, poolName, action string) (int, error) {
	var result int
	err := p.client.Call(ctx, "pool.scrub.scrub", []any{poolName, action}, &result)
	return result, err
}
