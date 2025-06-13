package truenas

import (
	"context"
	"fmt"
)

// ServiceClient provides methods for service management
type ServiceClient struct {
	client *Client
}

// NewServiceClient creates a new service client
func NewServiceClient(client *Client) *ServiceClient {
	return &ServiceClient{client: client}
}

// Service represents a system service
type Service struct {
	ID      int    `json:"id"`
	Service string `json:"service"`
	Enable  bool   `json:"enable"`
	State   string `json:"state"`
	PIDs    []int  `json:"pids"`
}

// ServiceUpdateRequest represents parameters for service.update
type ServiceUpdateRequest struct {
	Enable bool `json:"enable"`
}

// List returns all services
func (s *ServiceClient) List(ctx context.Context) ([]Service, error) {
	var result []Service
	err := s.client.Call(ctx, "service.query", []any{}, &result)
	return result, err
}

// Get returns a specific service by ID
func (s *ServiceClient) Get(ctx context.Context, id int) (*Service, error) {
	var result []Service
	err := s.client.Call(ctx, "service.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("service", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// GetByName returns a specific service by name
func (s *ServiceClient) GetByName(ctx context.Context, name string) (*Service, error) {
	var result []Service
	err := s.client.Call(ctx, "service.query", []any{[]any{[]any{"service", "=", name}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("service", fmt.Sprintf("name %s", name))
	}
	return &result[0], nil
}

// Update updates service configuration
func (s *ServiceClient) Update(ctx context.Context, id int, req ServiceUpdateRequest) (*Service, error) {
	var result any
	err := s.client.Call(ctx, "service.update", []any{id, req}, &result)
	if err != nil {
		return nil, err
	}

	// service.update returns just the service ID, so we need to fetch the full service
	return s.Get(ctx, id)
}

// Start starts a service
func (s *ServiceClient) Start(ctx context.Context, serviceName string, options ...map[string]any) error {
	params := []any{serviceName}
	if len(options) > 0 {
		params = append(params, options[0])
	}
	return s.client.Call(ctx, "service.start", params, nil)
}

// Stop stops a service
func (s *ServiceClient) Stop(ctx context.Context, serviceName string, options ...map[string]any) error {
	params := []any{serviceName}
	if len(options) > 0 {
		params = append(params, options[0])
	}
	return s.client.Call(ctx, "service.stop", params, nil)
}

// Restart restarts a service
func (s *ServiceClient) Restart(ctx context.Context, serviceName string, options ...map[string]any) error {
	params := []any{serviceName}
	if len(options) > 0 {
		params = append(params, options[0])
	}
	return s.client.Call(ctx, "service.restart", params, nil)
}

// Reload reloads a service configuration
func (s *ServiceClient) Reload(ctx context.Context, serviceName string, options ...map[string]any) error {
	params := []any{serviceName}
	if len(options) > 0 {
		params = append(params, options[0])
	}
	return s.client.Call(ctx, "service.reload", params, nil)
}

// Started checks if a service has been started
func (s *ServiceClient) Started(ctx context.Context, serviceName string) (bool, error) {
	var result bool
	err := s.client.Call(ctx, "service.started", []any{serviceName}, &result)
	return result, err
}

// SMB Service Methods

// SMBClient provides methods for SMB service management
type SMBClient struct {
	client *Client
}

// NewSMBClient creates a new SMB client
func NewSMBClient(client *Client) *SMBClient {
	return &SMBClient{client: client}
}

// SMBConfig represents SMB service configuration
type SMBConfig struct {
	NetBIOSName        string   `json:"netbiosname"`
	Workgroup          string   `json:"workgroup"`
	Description        string   `json:"description"`
	UnixCharset        string   `json:"unixcharset"`
	LogLevel           string   `json:"loglevel"`
	SyslogLevel        string   `json:"sysloglevel"`
	LocalMaster        bool     `json:"localmaster"`
	DomainLogons       bool     `json:"domainlogons"`
	TimeServer         bool     `json:"timeserver"`
	GuestAccount       string   `json:"guestaccount"`
	FileMap            string   `json:"filemask"`
	DirMap             string   `json:"dirmask"`
	NTLMv1Auth         bool     `json:"ntlmv1auth"`
	NullPasswords      bool     `json:"nullpw"`
	DebugPID           bool     `json:"debug_pid"`
	MaxLogSize         int      `json:"max_log_size"`
	UseSendfile        bool     `json:"use_sendfile"`
	AAAPLExtensions    bool     `json:"aapl_extensions"`
	EASupport          string   `json:"ea_support"`
	StoreDOSAttributes bool     `json:"store_dos_attributes"`
	HostsAllow         []string `json:"hostsallow"`
	HostsDeny          []string `json:"hostsdeny"`
	Admin              string   `json:"admin_group"`
	BindIP             []string `json:"bindip"`
	SMBEnable          bool     `json:"enable_smb1"`
	AuditEnable        bool     `json:"audit_enable"`
	AuditWatchList     []string `json:"audit_watch_list"`
	AuditIgnoreList    []string `json:"audit_ignore_list"`
}

// GetSMBConfig returns SMB service configuration
func (s *SMBClient) GetConfig(ctx context.Context) (*SMBConfig, error) {
	var result SMBConfig
	err := s.client.Call(ctx, "smb.config", []any{}, &result)
	return &result, err
}

// UpdateSMBConfig updates SMB service configuration
func (s *SMBClient) UpdateConfig(ctx context.Context, config *SMBConfig) (*SMBConfig, error) {
	var result SMBConfig
	err := s.client.Call(ctx, "smb.update", []any{*config}, &result)
	return &result, err
}

// NFS Service Methods

// NFSClient provides methods for NFS service management
type NFSClient struct {
	client *Client
}

// NewNFSClient creates a new NFS client
func NewNFSClient(client *Client) *NFSClient {
	return &NFSClient{client: client}
}

// NFSConfig represents NFS service configuration
type NFSConfig struct {
	V4             bool     `json:"v4"`
	V4V3Owner      bool     `json:"v4_v3owner"`
	V4KrbEnabled   bool     `json:"v4_krb"`
	V4Domain       string   `json:"v4_domain"`
	BindIP         []string `json:"bindip"`
	MountdPort     int      `json:"mountd_port"`
	RpcstatdPort   int      `json:"rpcstatd_port"`
	RpclockdPort   int      `json:"rpclockd_port"`
	Servers        int      `json:"servers"`
	UDPEnabled     bool     `json:"udp"`
	RPCGSSEnabled  bool     `json:"rpcgssd_enable"`
	UserdMaxGroups int      `json:"userd_manage_groups"`
}

// GetNFSConfig returns NFS service configuration
func (n *NFSClient) GetConfig(ctx context.Context) (*NFSConfig, error) {
	var result NFSConfig
	err := n.client.Call(ctx, "nfs.config", []any{}, &result)
	return &result, err
}

// UpdateNFSConfig updates NFS service configuration
func (n *NFSClient) UpdateConfig(ctx context.Context, config *NFSConfig) (*NFSConfig, error) {
	var result NFSConfig
	err := n.client.Call(ctx, "nfs.update", []any{*config}, &result)
	return &result, err
}

// SSH Service Methods

// SSHClient provides methods for SSH service management
type SSHClient struct {
	client *Client
}

// NewSSHClient creates a new SSH client
func NewSSHClient(client *Client) *SSHClient {
	return &SSHClient{client: client}
}

// SSHConfig represents SSH service configuration
type SSHConfig struct {
	TCPPort         []int    `json:"tcpport"`
	RootLogin       bool     `json:"rootlogin"`
	PasswordAuth    bool     `json:"passwordauth"`
	KerberosAuth    bool     `json:"kerberosauth"`
	TCPForwarding   bool     `json:"tcpfwd"`
	Compression     bool     `json:"compression"`
	SFTPLogLevel    string   `json:"sftp_log_level"`
	SFTPLogFacility string   `json:"sftp_log_facility"`
	WeakCiphers     []string `json:"weak_ciphers"`
	AuxParam        string   `json:"auxparam"`
	BindIface       []string `json:"bindiface"`
}

// GetSSHConfig returns SSH service configuration
func (s *SSHClient) GetConfig(ctx context.Context) (*SSHConfig, error) {
	var result SSHConfig
	err := s.client.Call(ctx, "ssh.config", []any{}, &result)
	return &result, err
}

// UpdateSSHConfig updates SSH service configuration
func (s *SSHClient) UpdateConfig(ctx context.Context, config *SSHConfig) (*SSHConfig, error) {
	var result SSHConfig
	err := s.client.Call(ctx, "ssh.update", []any{*config}, &result)
	return &result, err
}
