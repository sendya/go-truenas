package truenas

import (
	"context"
	"fmt"
)

// DiskClient provides methods for disk management
type DiskClient struct {
	client *Client
}

// NewDiskClient creates a new disk client
func NewDiskClient(client *Client) *DiskClient {
	return &DiskClient{client: client}
}

// DiskType represents disk types
type DiskType string

const (
	DiskTypeSSD DiskType = "SSD"
	DiskTypeHDD DiskType = "HDD"
)

// Disk represents a physical disk device
type Disk struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Devname         string         `json:"devname"`
	Identifier      string         `json:"identifier"`
	Serial          string         `json:"serial"`
	Size            int64          `json:"size"`
	Model           string         `json:"model"`
	Type            DiskType       `json:"type"`
	Rotationrate    *int           `json:"rotationrate"`
	HDDStandby      HDDStandby     `json:"hddstandby"`
	HDDStandbyForce bool           `json:"hddstandby_force"`
	AdvPowerMgmt    AdvPowerMgmt   `json:"advpowermgmt"`
	AcousticLevel   AcousticLevel  `json:"acousticlevel"`
	SmartEnabled    bool           `json:"togglesmart"`
	SmartOptions    string         `json:"smartoptions"`
	Critical        *int           `json:"critical"`
	Informational   *int           `json:"informational"`
	Difference      *int           `json:"difference"`
	Description     string         `json:"description"`
	Passwd          string         `json:"passwd"`
	Enclosure       *DiskEnclosure `json:"enclosure"`
	Pool            *string        `json:"pool,omitempty"`
	Expired         bool           `json:"expired,omitempty"`
}

// DiskEnclosure represents disk enclosure information
type DiskEnclosure struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slot   int    `json:"slot"`
	Number int    `json:"number"`
}

// DiskUpdateRequest represents parameters for disk.update
type DiskUpdateRequest struct {
	ToggleSmart     *bool          `json:"togglesmart,omitempty"`
	AcousticLevel   *AcousticLevel `json:"acousticlevel,omitempty"`
	AdvPowerMgmt    *AdvPowerMgmt  `json:"advpowermgmt,omitempty"`
	Description     *string        `json:"description,omitempty"`
	HDDStandby      *HDDStandby    `json:"hddstandby,omitempty"`
	HDDStandbyForce *bool          `json:"hddstandby_force,omitempty"`
	Passwd          *string        `json:"passwd,omitempty"`
	SmartOptions    *string        `json:"smartoptions,omitempty"`
	Critical        *int           `json:"critical,omitempty"`
	Informational   *int           `json:"informational,omitempty"`
	Difference      *int           `json:"difference,omitempty"`
	Enclosure       *DiskEnclosure `json:"enclosure,omitempty"`
}

// DiskQueryOptions represents additional options for disk.query
type DiskQueryOptions struct {
	IncludeExpired bool `json:"include_expired,omitempty"`
	Passwords      bool `json:"passwords,omitempty"`
	Pools          bool `json:"pools,omitempty"`
}

// EncryptedDevice represents an encrypted disk device
type EncryptedDevice struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Encrypted bool   `json:"encrypted"`
	Provider  string `json:"provider"`
	Decrypted bool   `json:"decrypted"`
}

// UnusedDisk represents an unused disk
type UnusedDisk struct {
	Name       string      `json:"name"`
	Devname    string      `json:"devname"`
	Size       int64       `json:"size"`
	Serial     string      `json:"serial"`
	Model      string      `json:"model"`
	Type       DiskType    `json:"type"`
	Partitions []Partition `json:"partitions,omitempty"`
	Driver     string      `json:"driver,omitempty"`
}

// Partition represents a disk partition
type Partition struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Type  string `json:"type"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

// SmartAttribute represents a S.M.A.R.T. attribute
type SmartAttribute struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Threshold  int    `json:"threshold"`
	Type       string `json:"type"`
	Updated    string `json:"updated"`
	WhenFailed string `json:"when_failed"`
	RawValue   int64  `json:"raw_value"`
}

// DiskTemperature represents disk temperature information
type DiskTemperature struct {
	Name        string  `json:"name"`
	Temperature *int    `json:"temperature"`
	Unit        string  `json:"unit"`
	Error       *string `json:"error,omitempty"`
}

// DecryptRequest represents parameters for disk.decrypt
type DecryptRequest struct {
	Devices    []string `json:"devices"`
	Passphrase *string  `json:"passphrase,omitempty"`
}

// WipeRequest represents parameters for disk.wipe
type WipeRequest struct {
	Device             string              `json:"dev"`
	Mode               WipeMode            `json:"mode"`
	SyncCache          bool                `json:"synccache"`
	SwapRemovalOptions *SwapRemovalOptions `json:"swap_removal_options,omitempty"`
}

// SwapRemovalOptions represents swap removal configuration
type SwapRemovalOptions struct {
	ConfigureSwap bool `json:"configure_swap"`
}

// Constants

// PowerMode represents S.M.A.R.T. power modes for temperature monitoring
type PowerMode string

const (
	PowerModeNever   PowerMode = "NEVER"
	PowerModeSleep   PowerMode = "SLEEP"
	PowerModeStandby PowerMode = "STANDBY"
	PowerModeIdle    PowerMode = "IDLE"
)

// AcousticLevel represents disk acoustic level settings
type AcousticLevel string

const (
	AcousticLevelDisabled AcousticLevel = "DISABLED"
	AcousticLevelMinimum  AcousticLevel = "MINIMUM"
	AcousticLevelMedium   AcousticLevel = "MEDIUM"
	AcousticLevelMaximum  AcousticLevel = "MAXIMUM"
)

// AdvPowerMgmt represents advanced power management settings
type AdvPowerMgmt string

const (
	AdvPowerMgmtDisabled AdvPowerMgmt = "DISABLED"
	AdvPowerMgmt1        AdvPowerMgmt = "1"
	AdvPowerMgmt64       AdvPowerMgmt = "64"
	AdvPowerMgmt127      AdvPowerMgmt = "127"
	AdvPowerMgmt128      AdvPowerMgmt = "128"
	AdvPowerMgmt192      AdvPowerMgmt = "192"
	AdvPowerMgmt254      AdvPowerMgmt = "254"
)

// HDDStandby represents HDD standby timeout settings
type HDDStandby string

const (
	HDDStandbyAlwaysOn HDDStandby = "ALWAYS ON"
	HDDStandby5        HDDStandby = "5"
	HDDStandby10       HDDStandby = "10"
	HDDStandby20       HDDStandby = "20"
	HDDStandby30       HDDStandby = "30"
	HDDStandby60       HDDStandby = "60"
	HDDStandby120      HDDStandby = "120"
	HDDStandby180      HDDStandby = "180"
	HDDStandby240      HDDStandby = "240"
	HDDStandby300      HDDStandby = "300"
	HDDStandby330      HDDStandby = "330"
)

// WipeMode represents disk wipe modes
type WipeMode string

const (
	WipeModeQuick      WipeMode = "QUICK"
	WipeModeFull       WipeMode = "FULL"
	WipeModeFullRandom WipeMode = "FULL_RANDOM"
)

// Basic disk operations

// List returns all disks
func (d *DiskClient) List(ctx context.Context) ([]Disk, error) {
	var result []Disk
	err := d.client.Call(ctx, "disk.query", []any{}, &result)
	return result, err
}

// ListWithOptions returns disks with additional options
func (d *DiskClient) ListWithOptions(ctx context.Context, opts *DiskQueryOptions) ([]Disk, error) {
	var result []Disk
	params := []any{[]any{}, map[string]any{"extra": opts}}
	err := d.client.Call(ctx, "disk.query", params, &result)
	return result, err
}

// Get returns a specific disk by ID
func (d *DiskClient) Get(ctx context.Context, id string) (*Disk, error) {
	var result []Disk
	err := d.client.Call(ctx, "disk.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("disk", fmt.Sprintf("ID %s", id))
	}
	return &result[0], nil
}

// Update updates disk configuration
func (d *DiskClient) Update(ctx context.Context, id string, req *DiskUpdateRequest) (*Disk, error) {
	var result Disk
	err := d.client.Call(ctx, "disk.update", []any{id, *req}, &result)
	return &result, err
}

// Encryption operations

// GetEncrypted returns all encrypted devices
func (d *DiskClient) GetEncrypted(ctx context.Context, includeUnused bool) ([]EncryptedDevice, error) {
	var result []EncryptedDevice
	options := map[string]any{"unused": includeUnused}
	err := d.client.Call(ctx, "disk.get_encrypted", []any{options}, &result)
	return result, err
}

// Decrypt decrypts encrypted devices
func (d *DiskClient) Decrypt(ctx context.Context, req *DecryptRequest) error {
	return d.client.CallJob(ctx, "disk.decrypt", []any{req.Devices, req.Passphrase}, nil)
}

// Utility operations

// GetUnused returns all unused disks
func (d *DiskClient) GetUnused(ctx context.Context, joinPartitions bool) ([]UnusedDisk, error) {
	var result []UnusedDisk
	err := d.client.Call(ctx, "disk.get_unused", []any{joinPartitions}, &result)
	return result, err
}

// LabelToDev converts disk label to device name
func (d *DiskClient) LabelToDev(ctx context.Context, label string) (string, error) {
	var result string
	err := d.client.Call(ctx, "disk.label_to_dev", []any{label}, &result)
	return result, err
}

// S.M.A.R.T. operations

// GetSmartAttributes returns S.M.A.R.T. attributes for a disk
func (d *DiskClient) GetSmartAttributes(ctx context.Context, deviceName string) ([]SmartAttribute, error) {
	var result []SmartAttribute
	err := d.client.Call(ctx, "disk.smart_attributes", []any{deviceName}, &result)
	return result, err
}

// GetTemperature returns temperature for a single disk
func (d *DiskClient) GetTemperature(ctx context.Context, deviceName string, powerMode PowerMode) (*DiskTemperature, error) {
	var result DiskTemperature
	err := d.client.Call(ctx, "disk.temperature", []any{deviceName, string(powerMode)}, &result)
	return &result, err
}

// GetTemperatures returns temperatures for multiple disks
func (d *DiskClient) GetTemperatures(ctx context.Context, deviceNames []string, powerMode PowerMode) ([]DiskTemperature, error) {
	var result []DiskTemperature
	err := d.client.Call(ctx, "disk.temperatures", []any{deviceNames, string(powerMode)}, &result)
	return result, err
}

// Power management operations

// Spindown spins down a disk
func (d *DiskClient) Spindown(ctx context.Context, deviceName string) error {
	return d.client.Call(ctx, "disk.spindown", []any{deviceName}, nil)
}

// Provisioning operations

// Overprovision configures disk overprovisioning
func (d *DiskClient) Overprovision(ctx context.Context, deviceName string, size int64) error {
	return d.client.Call(ctx, "disk.overprovision", []any{deviceName, size}, nil)
}

// Unoverprovision removes overprovisioning from disk
func (d *DiskClient) Unoverprovision(ctx context.Context, deviceName string) error {
	return d.client.Call(ctx, "disk.unoverprovision", []any{deviceName}, nil)
}

// Wipe operations

// Wipe performs a disk wipe operation (asynchronous job)
func (d *DiskClient) Wipe(ctx context.Context, req *WipeRequest) error {
	return d.client.CallJob(ctx, "disk.wipe", []any{req.Device, req.Mode, req.SyncCache, req.SwapRemovalOptions}, nil)
}

// SED operations

// GetSedDevName gets SED (Self-Encrypting Drive) device name
func (d *DiskClient) GetSedDevName(ctx context.Context, deviceName string) (string, error) {
	var result string
	err := d.client.Call(ctx, "disk.sed_dev_name", []any{deviceName}, &result)
	return result, err
}
