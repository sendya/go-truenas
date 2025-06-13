package truenas

import (
	"context"
)

// BootClient provides methods for boot pool management
type BootClient struct {
	client *Client
}

// NewBootClient creates a new boot client
func NewBootClient(client *Client) *BootClient {
	return &BootClient{client: client}
}

// BootDisk represents a disk in the boot pool
type BootDisk struct {
	Name      string `json:"name"`
	Label     string `json:"label"`
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	Status    string `json:"status"`
	Serial    string `json:"serial"`
	Model     string `json:"model"`
	Type      string `json:"type"`
	Available bool   `json:"available"`
}

// BootState represents the current state of the boot pool
type BootState struct {
	Name       string         `json:"name"`
	ID         string         `json:"id"`
	GUID       string         `json:"guid"`
	Hostname   string         `json:"hostname"`
	Status     string         `json:"status"`
	Scan       any            `json:"scan"`
	Properties map[string]any `json:"properties"`
	Groups     []BootVdev     `json:"groups"`
	Topology   BootTopology   `json:"topology"`
	Healthy    bool           `json:"healthy"`
	Warning    bool           `json:"warning"`
	Unknown    bool           `json:"unknown"`
}

// BootVdev represents a virtual device in the boot pool
type BootVdev struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Status   string     `json:"status"`
	Stats    any        `json:"stats"`
	Children []BootVdev `json:"children"`
	Device   string     `json:"device,omitempty"`
	Disk     string     `json:"disk,omitempty"`
	Path     string     `json:"path,omitempty"`
}

// BootTopology represents the topology of the boot pool
type BootTopology struct {
	Data    []BootVdev `json:"data"`
	Log     []BootVdev `json:"log"`
	Cache   []BootVdev `json:"cache"`
	Spare   []BootVdev `json:"spare"`
	Special []BootVdev `json:"special"`
	Dedup   []BootVdev `json:"dedup"`
}

// BootAttachRequest represents parameters for boot.attach
type BootAttachRequest struct {
	Device string `json:"dev"`
	Expand bool   `json:"expand,omitempty"`
}

// GetDisks returns all disks in the boot pool
func (b *BootClient) GetDisks(ctx context.Context) ([]BootDisk, error) {
	var result []BootDisk
	err := b.client.Call(ctx, "boot.get_disks", []any{}, &result)
	return result, err
}

// GetState returns the current state of the boot pool
func (b *BootClient) GetState(ctx context.Context) (*BootState, error) {
	var result BootState
	err := b.client.Call(ctx, "boot.get_state", []any{}, &result)
	return &result, err
}

// Attach attaches a disk to the boot pool (converts stripe to mirror)
func (b *BootClient) Attach(ctx context.Context, device string, expand bool) error {
	options := map[string]any{}
	if expand {
		options["expand"] = true
	}
	return b.client.CallJob(ctx, "boot.attach", []any{device, options}, nil)
}

// Detach detaches a device from the boot pool
func (b *BootClient) Detach(ctx context.Context, device string) error {
	return b.client.Call(ctx, "boot.detach", []any{device}, nil)
}

// Replace replaces a device in the boot pool
func (b *BootClient) Replace(ctx context.Context, label, device string) error {
	return b.client.Call(ctx, "boot.replace", []any{label, device}, nil)
}

// Scrub starts a scrub operation on the boot pool
func (b *BootClient) Scrub(ctx context.Context) error {
	return b.client.CallJob(ctx, "boot.scrub", []any{}, nil)
}

// GetScrubInterval returns the automatic scrub interval in days
func (b *BootClient) GetScrubInterval(ctx context.Context) (int, error) {
	var result int
	err := b.client.Call(ctx, "boot.get_scrub_interval", []any{}, &result)
	return result, err
}

// SetScrubInterval sets the automatic scrub interval in days
func (b *BootClient) SetScrubInterval(ctx context.Context, interval int) error {
	return b.client.Call(ctx, "boot.set_scrub_interval", []any{interval}, nil)
}
