package truenas

import (
	"context"
	"fmt"
)

// VMClient provides methods for virtual machine management
type VMClient struct {
	client *Client
}

// NewVMClient creates a new VM client
func NewVMClient(client *Client) *VMClient {
	return &VMClient{client: client}
}

// VM represents a virtual machine
type VM struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	VCPUs           int          `json:"vcpus"`
	Cores           int          `json:"cores"`
	Threads         int          `json:"threads"`
	Memory          int          `json:"memory"`
	Bootloader      VMBootloader `json:"bootloader"`
	GrubConfig      string       `json:"grubconfig"`
	Devices         []VMDevice   `json:"devices"`
	Autostart       bool         `json:"autostart"`
	Time            VMTime       `json:"time"`
	ShutdownTimeout int          `json:"shutdown_timeout"`
	Status          VMStatus     `json:"status"`
}

// VMStatus represents VM runtime status
type VMStatus struct {
	State VMState `json:"state"`
	PID   int     `json:"pid,omitempty"`
}

// VMDevice represents a VM device
type VMDevice struct {
	ID         int            `json:"id"`
	DType      VMDeviceType   `json:"dtype"`
	VM         int            `json:"vm"`
	Attributes map[string]any `json:"attributes"`
	Order      int            `json:"order"`
}

// VMCreateRequest represents parameters for vm.create
type VMCreateRequest struct {
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	VCPUs           int          `json:"vcpus"`
	Cores           int          `json:"cores,omitempty"`
	Threads         int          `json:"threads,omitempty"`
	Memory          int          `json:"memory"`
	Bootloader      VMBootloader `json:"bootloader,omitempty"`
	GrubConfig      string       `json:"grubconfig,omitempty"`
	Devices         []VMDevice   `json:"devices,omitempty"`
	Autostart       *bool        `json:"autostart,omitempty"`
	Time            VMTime       `json:"time,omitempty"`
	ShutdownTimeout int          `json:"shutdown_timeout,omitempty"`
}

// VMUpdateRequest represents parameters for vm.update
type VMUpdateRequest struct {
	Name            string       `json:"name,omitempty"`
	Description     string       `json:"description,omitempty"`
	VCPUs           int          `json:"vcpus,omitempty"`
	Cores           int          `json:"cores,omitempty"`
	Threads         int          `json:"threads,omitempty"`
	Memory          int          `json:"memory,omitempty"`
	Bootloader      VMBootloader `json:"bootloader,omitempty"`
	GrubConfig      string       `json:"grubconfig,omitempty"`
	Devices         []VMDevice   `json:"devices,omitempty"`
	Autostart       *bool        `json:"autostart,omitempty"`
	Time            VMTime       `json:"time,omitempty"`
	ShutdownTimeout int          `json:"shutdown_timeout,omitempty"`
}

// VMDeleteRequest represents parameters for vm.delete
type VMDeleteRequest struct {
	Zvols *bool `json:"zvols,omitempty"`
	Force *bool `json:"force,omitempty"`
}

// VMStartRequest represents parameters for vm.start
type VMStartRequest struct {
	Overcommit *bool `json:"overcommit,omitempty"`
}

// VMStopRequest represents parameters for vm.stop
type VMStopRequest struct {
	Force             *bool `json:"force,omitempty"`
	ForceAfterTimeout *bool `json:"force_after_timeout,omitempty"`
}

// VMDeviceCreateRequest represents parameters for vm.device.create
type VMDeviceCreateRequest struct {
	DType      VMDeviceType   `json:"dtype"`
	VM         int            `json:"vm"`
	Attributes map[string]any `json:"attributes"`
	Order      int            `json:"order,omitempty"`
}

// VMDeviceDeleteRequest represents parameters for vm.device.delete
type VMDeviceDeleteRequest struct {
	Zvol    *bool `json:"zvol,omitempty"`
	RawFile *bool `json:"raw_file,omitempty"`
}

// VMMemoryInfo represents VM memory usage information
type VMMemoryInfo struct {
	RNP  int `json:"RNP"`  // Running but not provisioned
	PRD  int `json:"PRD"`  // Provisioned but not running
	RPRD int `json:"RPRD"` // Running and provisioned
}

// VMBootloader represents available bootloader types
type VMBootloader string

const (
	VMBootloaderUEFI    VMBootloader = "UEFI"
	VMBootloaderUEFICSM VMBootloader = "UEFI_CSM"
	VMBootloaderGRUB    VMBootloader = "GRUB"
)

// VMTime represents VM time synchronization modes
type VMTime string

const (
	VMTimeLocal VMTime = "LOCAL"
	VMTimeUTC   VMTime = "UTC"
)

// VMState represents VM running states
type VMState string

const (
	VMStateRunning VMState = "RUNNING"
	VMStateStopped VMState = "STOPPED"
)

// VMDeviceType represents VM device types
type VMDeviceType string

const (
	VMDeviceTypeNIC   VMDeviceType = "NIC"
	VMDeviceTypeDisk  VMDeviceType = "DISK"
	VMDeviceTypeCDROM VMDeviceType = "CDROM"
	VMDeviceTypePCI   VMDeviceType = "PCI"
	VMDeviceTypeVNC   VMDeviceType = "VNC"
	VMDeviceTypeRAW   VMDeviceType = "RAW"
)

// List returns all VMs
func (v *VMClient) List(ctx context.Context) ([]VM, error) {
	var result []VM
	err := v.client.Call(ctx, "vm.query", []any{}, &result)
	return result, err
}

// Get returns a specific VM by ID
func (v *VMClient) Get(ctx context.Context, id int) (*VM, error) {
	var result []VM
	err := v.client.Call(ctx, "vm.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("vm", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new VM
func (v *VMClient) Create(ctx context.Context, req *VMCreateRequest) (*VM, error) {
	var result VM
	err := v.client.Call(ctx, "vm.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing VM
func (v *VMClient) Update(ctx context.Context, id int, req *VMUpdateRequest) (*VM, error) {
	var result VM
	err := v.client.Call(ctx, "vm.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes a VM
func (v *VMClient) Delete(ctx context.Context, id int, req *VMDeleteRequest) error {
	params := []any{id}
	if req != nil {
		params = append(params, *req)
	}
	return v.client.Call(ctx, "vm.delete", params, nil)
}

// Clone clones a VM
func (v *VMClient) Clone(ctx context.Context, id int, name string) (*VM, error) {
	var result VM
	params := []any{id}
	if name != "" {
		params = append(params, name)
	}
	err := v.client.Call(ctx, "vm.clone", params, &result)
	return &result, err
}

// VM Control Operations

// Start starts a VM
func (v *VMClient) Start(ctx context.Context, id int, req *VMStartRequest) error {
	params := []any{id}
	if req != nil {
		params = append(params, *req)
	}
	return v.client.Call(ctx, "vm.start", params, nil)
}

// Stop stops a VM gracefully
func (v *VMClient) Stop(ctx context.Context, id int, req *VMStopRequest) error {
	params := []any{id}
	if req != nil {
		params = append(params, *req)
	}
	return v.client.CallJob(ctx, "vm.stop", params, nil)
}

// PowerOff forcefully powers off a VM
func (v *VMClient) PowerOff(ctx context.Context, id int) error {
	return v.client.Call(ctx, "vm.poweroff", []any{id}, nil)
}

// Restart restarts a VM
func (v *VMClient) Restart(ctx context.Context, id int) error {
	return v.client.CallJob(ctx, "vm.restart", []any{id}, nil)
}

// GetStatus returns the current status of a VM
func (v *VMClient) GetStatus(ctx context.Context, id int) (*VMStatus, error) {
	var result VMStatus
	err := v.client.Call(ctx, "vm.status", []any{id}, &result)
	return &result, err
}

// VM Information Methods

// GetFlags returns CPU flags for bhyve
func (v *VMClient) GetFlags(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := v.client.Call(ctx, "vm.flags", []any{}, &result)
	return result, err
}

// GetAvailableMemory returns available memory for VMs
func (v *VMClient) GetAvailableMemory(ctx context.Context, overcommit bool) (int, error) {
	var result int
	err := v.client.Call(ctx, "vm.get_available_memory", []any{overcommit}, &result)
	return result, err
}

// GetMemoryInUse returns memory usage information
func (v *VMClient) GetMemoryInUse(ctx context.Context) (*VMMemoryInfo, error) {
	var result VMMemoryInfo
	err := v.client.Call(ctx, "vm.get_vmemory_in_use", []any{}, &result)
	return &result, err
}

// GetAttachedInterfaces returns attached physical interfaces for a VM
func (v *VMClient) GetAttachedInterfaces(ctx context.Context, id int) ([]string, error) {
	var result []string
	err := v.client.Call(ctx, "vm.get_attached_iface", []any{id}, &result)
	return result, err
}

// GetConsole returns console device path for a VM
func (v *VMClient) GetConsole(ctx context.Context, id int) (string, error) {
	var result string
	err := v.client.Call(ctx, "vm.get_console", []any{id}, &result)
	return result, err
}

// VNC Methods

// GetVNC returns VNC devices for a VM
func (v *VMClient) GetVNC(ctx context.Context, id int) ([]map[string]any, error) {
	var result []map[string]any
	err := v.client.Call(ctx, "vm.get_vnc", []any{id}, &result)
	return result, err
}

// GetVNCWeb returns VNC web URLs for a VM
func (v *VMClient) GetVNCWeb(ctx context.Context, id int, host string) ([]string, error) {
	var result []string
	params := []any{id}
	if host != "" {
		params = append(params, host)
	}
	err := v.client.Call(ctx, "vm.get_vnc_web", params, &result)
	return result, err
}

// GetVNCIPv4 returns available IPv4 addresses for VNC
func (v *VMClient) GetVNCIPv4(ctx context.Context) ([]string, error) {
	var result []string
	err := v.client.Call(ctx, "vm.get_vnc_ipv4", []any{}, &result)
	return result, err
}

// GetVNCPortWizard returns VNC port configuration
func (v *VMClient) GetVNCPortWizard(ctx context.Context) (any, error) {
	var result any
	err := v.client.Call(ctx, "vm.vnc_port_wizard", []any{}, &result)
	return result, err
}

// Utility Methods

// GenerateRandomMAC generates a random MAC address
func (v *VMClient) GenerateRandomMAC(ctx context.Context) (string, error) {
	var result string
	err := v.client.Call(ctx, "vm.random_mac", []any{}, &result)
	return result, err
}

// IdentifyHypervisor checks hypervisor compatibility
func (v *VMClient) IdentifyHypervisor(ctx context.Context) (bool, error) {
	var result bool
	err := v.client.Call(ctx, "vm.identify_hypervisor", []any{}, &result)
	return result, err
}

// VM Device Management

// VMDeviceClient provides methods for VM device management
type VMDeviceClient struct {
	client *Client
}

// NewVMDeviceClient creates a new VM device client
func NewVMDeviceClient(client *Client) *VMDeviceClient {
	return &VMDeviceClient{client: client}
}

// ListDevices returns all VM devices
func (d *VMDeviceClient) List(ctx context.Context) ([]VMDevice, error) {
	var result []VMDevice
	err := d.client.Call(ctx, "vm.device.query", []any{}, &result)
	return result, err
}

// GetDevice returns a specific VM device by ID
func (d *VMDeviceClient) Get(ctx context.Context, id int) (*VMDevice, error) {
	var result []VMDevice
	err := d.client.Call(ctx, "vm.device.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("vm_device", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// CreateDevice creates a new VM device
func (d *VMDeviceClient) Create(ctx context.Context, req *VMDeviceCreateRequest) (*VMDevice, error) {
	var result VMDevice
	err := d.client.Call(ctx, "vm.device.create", []any{*req}, &result)
	return &result, err
}

// UpdateDevice updates an existing VM device
func (d *VMDeviceClient) Update(ctx context.Context, id int, req *VMDeviceCreateRequest) (*VMDevice, error) {
	var result VMDevice
	err := d.client.Call(ctx, "vm.device.update", []any{id, *req}, &result)
	return &result, err
}

// DeleteDevice deletes a VM device
func (d *VMDeviceClient) Delete(ctx context.Context, id int, req *VMDeviceDeleteRequest) error {
	params := []any{id}
	if req != nil {
		params = append(params, *req)
	}
	return d.client.Call(ctx, "vm.device.delete", params, nil)
}

// GetNICAttachChoices returns available NIC attach choices
func (d *VMDeviceClient) GetNICAttachChoices(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := d.client.Call(ctx, "vm.device.nic_attach_choices", []any{}, &result)
	return result, err
}

// GetPPTDevChoices returns available PCI passthrough device choices
func (d *VMDeviceClient) GetPPTDevChoices(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := d.client.Call(ctx, "vm.device.pptdev_choices", []any{}, &result)
	return result, err
}

// GetVNCBindChoices returns available VNC bind choices
func (d *VMDeviceClient) GetVNCBindChoices(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := d.client.Call(ctx, "vm.device.vnc_bind_choices", []any{}, &result)
	return result, err
}
