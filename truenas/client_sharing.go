package truenas

import (
	"context"
	"fmt"
)

// SharingClient provides methods for managing file shares across all protocols
type SharingClient struct {
	client *Client
	AFP    *SharingAFPClient
	NFS    *SharingNFSClient
	SMB    *SharingSMBClient
	WebDAV *SharingWebDAVClient
}

// NewSharingClient creates a new sharing client
func NewSharingClient(client *Client) *SharingClient {
	return &SharingClient{
		client: client,
		AFP:    NewSharingAFPClient(client),
		NFS:    NewSharingNFSClient(client),
		SMB:    NewSharingSMBClient(client),
		WebDAV: NewSharingWebDAVClient(client),
	}
}

// AFP (Apple Filing Protocol) Client

// SharingAFPClient provides methods for AFP share management
type SharingAFPClient struct {
	client *Client
}

// NewSharingAFPClient creates a new AFP sharing client
func NewSharingAFPClient(client *Client) *SharingAFPClient {
	return &SharingAFPClient{client: client}
}

// AFPShare represents an AFP share configuration
type AFPShare struct {
	ID               int      `json:"id"`
	Path             string   `json:"path"`
	Home             bool     `json:"home"`
	Name             string   `json:"name"`
	Comment          string   `json:"comment"`
	Allow            []string `json:"allow"`
	Deny             []string `json:"deny"`
	RO               []string `json:"ro"`
	RW               []string `json:"rw"`
	TimeMachine      bool     `json:"timemachine"`
	TimeMachineQuota int      `json:"timemachine_quota"`
	NoDev            bool     `json:"nodev"`
	NoStat           bool     `json:"nostat"`
	UPriv            bool     `json:"upriv"`
	FPerm            string   `json:"fperm"`
	DPerm            string   `json:"dperm"`
	UMask            string   `json:"umask"`
	HostsAllow       []string `json:"hostsallow"`
	HostsDeny        []string `json:"hostsdeny"`
	VUID             *string  `json:"vuid"`
	AuxParams        string   `json:"auxparams"`
	Enabled          bool     `json:"enabled"`
}

// AFPShareRequest represents parameters for creating/updating AFP shares
type AFPShareRequest struct {
	Path             string   `json:"path"`
	Home             bool     `json:"home"`
	Name             string   `json:"name"`
	Comment          string   `json:"comment"`
	Allow            []string `json:"allow"`
	Deny             []string `json:"deny"`
	RO               []string `json:"ro"`
	RW               []string `json:"rw"`
	TimeMachine      bool     `json:"timemachine"`
	TimeMachineQuota int      `json:"timemachine_quota"`
	NoDev            bool     `json:"nodev"`
	NoStat           bool     `json:"nostat"`
	UPriv            bool     `json:"upriv"`
	FPerm            string   `json:"fperm"`
	DPerm            string   `json:"dperm"`
	UMask            string   `json:"umask"`
	HostsAllow       []string `json:"hostsallow"`
	HostsDeny        []string `json:"hostsdeny"`
	VUID             *string  `json:"vuid,omitempty"`
	AuxParams        string   `json:"auxparams"`
	Enabled          bool     `json:"enabled"`
}

// List returns all AFP shares
func (a *SharingAFPClient) List(ctx context.Context) ([]AFPShare, error) {
	var result []AFPShare
	err := a.client.Call(ctx, "sharing.afp.query", []any{}, &result)
	return result, err
}

// Get returns a specific AFP share by ID
func (a *SharingAFPClient) Get(ctx context.Context, id int) (*AFPShare, error) {
	var result []AFPShare
	err := a.client.Call(ctx, "sharing.afp.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("afp_share", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new AFP share
func (a *SharingAFPClient) Create(ctx context.Context, req *AFPShareRequest) (*AFPShare, error) {
	var result AFPShare
	err := a.client.Call(ctx, "sharing.afp.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing AFP share
func (a *SharingAFPClient) Update(ctx context.Context, id int, req *AFPShareRequest) (*AFPShare, error) {
	var result AFPShare
	err := a.client.Call(ctx, "sharing.afp.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes an AFP share
func (a *SharingAFPClient) Delete(ctx context.Context, id int) error {
	return a.client.Call(ctx, "sharing.afp.delete", []any{id}, nil)
}

// NFS (Network File System) Client

// SharingNFSClient provides methods for NFS share management
type SharingNFSClient struct {
	client *Client
}

// NewSharingNFSClient creates a new NFS sharing client
func NewSharingNFSClient(client *Client) *SharingNFSClient {
	return &SharingNFSClient{client: client}
}

// NFSShare represents an NFS share configuration
type NFSShare struct {
	ID           int      `json:"id"`
	Path         string   `json:"path"`
	Aliases      []string `json:"aliases"`
	Comment      string   `json:"comment"`
	Networks     []string `json:"networks"`
	Hosts        []string `json:"hosts"`
	RO           bool     `json:"ro"`
	MapRootUser  *string  `json:"maproot_user"`
	MapRootGroup *string  `json:"maproot_group"`
	MapAllUser   *string  `json:"mapall_user"`
	MapAllGroup  *string  `json:"mapall_group"`
	Security     []string `json:"security"`
	Enabled      bool     `json:"enabled"`
	Locked       bool     `json:"locked"`
}

// NFSShareRequest represents parameters for creating/updating NFS shares
type NFSShareRequest struct {
	Path         string   `json:"path"`
	Comment      string   `json:"comment,omitempty"`
	Networks     []string `json:"networks,omitempty"`
	Hosts        []string `json:"hosts,omitempty"`
	RO           bool     `json:"ro,omitempty"`
	MapRootUser  *string  `json:"maproot_user,omitempty"`
	MapRootGroup *string  `json:"maproot_group,omitempty"`
	MapAllUser   *string  `json:"mapall_user,omitempty"`
	MapAllGroup  *string  `json:"mapall_group,omitempty"`
	Security     []string `json:"security"`
	Enabled      bool     `json:"enabled"`
}

// List returns all NFS shares
func (n *SharingNFSClient) List(ctx context.Context) ([]NFSShare, error) {
	var result []NFSShare
	err := n.client.Call(ctx, "sharing.nfs.query", []any{}, &result)
	return result, err
}

// Get returns a specific NFS share by ID
func (n *SharingNFSClient) Get(ctx context.Context, id int) (*NFSShare, error) {
	var result []NFSShare
	err := n.client.Call(ctx, "sharing.nfs.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("nfs_share", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new NFS share
func (n *SharingNFSClient) Create(ctx context.Context, req *NFSShareRequest) (*NFSShare, error) {
	var result NFSShare
	err := n.client.Call(ctx, "sharing.nfs.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing NFS share
func (n *SharingNFSClient) Update(ctx context.Context, id int, req *NFSShareRequest) (*NFSShare, error) {
	var result NFSShare
	err := n.client.Call(ctx, "sharing.nfs.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes an NFS share
func (n *SharingNFSClient) Delete(ctx context.Context, id int) error {
	return n.client.Call(ctx, "sharing.nfs.delete", []any{id}, nil)
}

// GetHumanIdentifier returns a human-readable identifier for an NFS share
func (n *SharingNFSClient) GetHumanIdentifier(ctx context.Context, id int) (string, error) {
	var result string
	err := n.client.Call(ctx, "sharing.nfs.human_identifier", []any{id}, &result)
	return result, err
}

// SMB (Server Message Block) Client

// SharingSMBClient provides methods for SMB share management
type SharingSMBClient struct {
	client *Client
}

// NewSharingSMBClient creates a new SMB sharing client
func NewSharingSMBClient(client *Client) *SharingSMBClient {
	return &SharingSMBClient{client: client}
}

// SMBPurpose represents SMB share purpose presets
type SMBPurpose string

const (
	SMBPurposeNoPreset            SMBPurpose = "NO_PRESET"
	SMBPurposeDefaultShare        SMBPurpose = "DEFAULT_SHARE"
	SMBPurposeEnhancedTimeMachine SMBPurpose = "ENHANCED_TIMEMACHINE"
	SMBPurposeMultiProtocolAFP    SMBPurpose = "MULTI_PROTOCOL_AFP"
	SMBPurposeMultiProtocolNFS    SMBPurpose = "MULTI_PROTOCOL_NFS"
	SMBPurposePrivateDatasets     SMBPurpose = "PRIVATE_DATASETS"
	SMBPurposeWormDropbox         SMBPurpose = "WORM_DROPBOX"
)

// SMBShare represents an SMB share configuration
type SMBShare struct {
	ID               int        `json:"id"`
	Purpose          SMBPurpose `json:"purpose"`
	Path             string     `json:"path"`
	PathSuffix       string     `json:"path_suffix"`
	Home             bool       `json:"home"`
	Name             string     `json:"name"`
	Comment          string     `json:"comment"`
	RO               bool       `json:"ro"`
	Browsable        bool       `json:"browsable"`
	TimeMachine      bool       `json:"timemachine"`
	RecycleBin       bool       `json:"recyclebin"`
	GuestOK          bool       `json:"guestok"`
	ABE              bool       `json:"abe"`
	HostsAllow       []string   `json:"hostsallow"`
	HostsDeny        []string   `json:"hostsdeny"`
	AAPLNameMangling bool       `json:"aapl_name_mangling"`
	ACL              bool       `json:"acl"`
	DurableHandle    bool       `json:"durablehandle"`
	ShadowCopy       bool       `json:"shadowcopy"`
	Streams          bool       `json:"streams"`
	FSRVP            bool       `json:"fsrvp"`
	AuxSMBConf       string     `json:"auxsmbconf"`
	Enabled          bool       `json:"enabled"`
}

// SMBShareRequest represents parameters for creating/updating SMB shares
type SMBShareRequest struct {
	Purpose          SMBPurpose `json:"purpose"`
	Path             string     `json:"path"`
	PathSuffix       string     `json:"path_suffix"`
	Home             bool       `json:"home"`
	Name             string     `json:"name"`
	Comment          string     `json:"comment"`
	RO               bool       `json:"ro"`
	Browsable        bool       `json:"browsable"`
	TimeMachine      bool       `json:"timemachine"`
	RecycleBin       bool       `json:"recyclebin"`
	GuestOK          bool       `json:"guestok"`
	ABE              bool       `json:"abe"`
	HostsAllow       []string   `json:"hostsallow"`
	HostsDeny        []string   `json:"hostsdeny"`
	AAPLNameMangling bool       `json:"aapl_name_mangling"`
	ACL              bool       `json:"acl"`
	DurableHandle    bool       `json:"durablehandle"`
	ShadowCopy       bool       `json:"shadowcopy"`
	Streams          bool       `json:"streams"`
	FSRVP            bool       `json:"fsrvp"`
	AuxSMBConf       string     `json:"auxsmbconf"`
	Enabled          bool       `json:"enabled"`
}

// SMBPreset represents an SMB configuration preset
type SMBPreset struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Config      any    `json:"config"`
}

// List returns all SMB shares
func (s *SharingSMBClient) List(ctx context.Context) ([]SMBShare, error) {
	var result []SMBShare
	err := s.client.Call(ctx, "sharing.smb.query", []any{}, &result)
	return result, err
}

// Get returns a specific SMB share by ID
func (s *SharingSMBClient) Get(ctx context.Context, id int) (*SMBShare, error) {
	var result []SMBShare
	err := s.client.Call(ctx, "sharing.smb.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("smb_share", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new SMB share
func (s *SharingSMBClient) Create(ctx context.Context, req *SMBShareRequest) (*SMBShare, error) {
	var result SMBShare
	err := s.client.Call(ctx, "sharing.smb.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing SMB share
func (s *SharingSMBClient) Update(ctx context.Context, id int, req *SMBShareRequest) (*SMBShare, error) {
	var result SMBShare
	err := s.client.Call(ctx, "sharing.smb.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes an SMB share (forcibly disconnects clients)
func (s *SharingSMBClient) Delete(ctx context.Context, id int) error {
	return s.client.Call(ctx, "sharing.smb.delete", []any{id}, nil)
}

// GetPresets returns available SMB configuration presets
func (s *SharingSMBClient) GetPresets(ctx context.Context) ([]SMBPreset, error) {
	var result []SMBPreset
	err := s.client.Call(ctx, "sharing.smb.presets", []any{}, &result)
	return result, err
}

// WebDAV Client

// SharingWebDAVClient provides methods for WebDAV share management
type SharingWebDAVClient struct {
	client *Client
}

// NewSharingWebDAVClient creates a new WebDAV sharing client
func NewSharingWebDAVClient(client *Client) *SharingWebDAVClient {
	return &SharingWebDAVClient{client: client}
}

// WebDAVShare represents a WebDAV share configuration
type WebDAVShare struct {
	ID      int    `json:"id"`
	Perm    bool   `json:"perm"`
	RO      bool   `json:"ro"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Enabled bool   `json:"enabled"`
}

// WebDAVShareRequest represents parameters for creating/updating WebDAV shares
type WebDAVShareRequest struct {
	Perm    bool   `json:"perm"`
	RO      bool   `json:"ro"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Enabled bool   `json:"enabled"`
}

// List returns all WebDAV shares
func (w *SharingWebDAVClient) List(ctx context.Context) ([]WebDAVShare, error) {
	var result []WebDAVShare
	err := w.client.Call(ctx, "sharing.webdav.query", []any{}, &result)
	return result, err
}

// Get returns a specific WebDAV share by ID
func (w *SharingWebDAVClient) Get(ctx context.Context, id int) (*WebDAVShare, error) {
	var result []WebDAVShare
	err := w.client.Call(ctx, "sharing.webdav.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("webdav_share", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new WebDAV share
func (w *SharingWebDAVClient) Create(ctx context.Context, req *WebDAVShareRequest) (*WebDAVShare, error) {
	var result WebDAVShare
	err := w.client.Call(ctx, "sharing.webdav.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing WebDAV share
func (w *SharingWebDAVClient) Update(ctx context.Context, id int, req *WebDAVShareRequest) (*WebDAVShare, error) {
	var result WebDAVShare
	err := w.client.Call(ctx, "sharing.webdav.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes a WebDAV share
func (w *SharingWebDAVClient) Delete(ctx context.Context, id int) error {
	return w.client.Call(ctx, "sharing.webdav.delete", []any{id}, nil)
}
