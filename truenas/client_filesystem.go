package truenas

import (
	"context"
	"time"
)

// FilesystemClient provides methods for filesystem management
type FilesystemClient struct {
	client *Client
}

// NewFilesystemClient creates a new filesystem client
func NewFilesystemClient(client *Client) *FilesystemClient {
	return &FilesystemClient{client: client}
}

// FilesystemStat represents filesystem stat information
type FilesystemStat struct {
	Size       int64     `json:"size"`
	Mode       int       `json:"mode"`
	UID        int       `json:"uid"`
	GID        int       `json:"gid"`
	Atime      time.Time `json:"atime"`
	Mtime      time.Time `json:"mtime"`
	Ctime      time.Time `json:"ctime"`
	Dev        int64     `json:"dev"`
	Inode      int64     `json:"inode"`
	Nlink      int       `json:"nlink"`
	User       string    `json:"user"`
	Group      string    `json:"group"`
	Acl        bool      `json:"acl"`
	IsFile     bool      `json:"is_file"`
	IsDir      bool      `json:"is_dir"`
	IsSymlink  bool      `json:"is_symlink"`
	IsCharDev  bool      `json:"is_char_device"`
	IsBlockDev bool      `json:"is_block_device"`
	IsFIFO     bool      `json:"is_fifo"`
	IsSocket   bool      `json:"is_socket"`
	RealPath   string    `json:"realpath"`
}

// FilesystemStatfs represents filesystem statistics
type FilesystemStatfs struct {
	FreeBytes  int64    `json:"free_bytes"`
	AvailBytes int64    `json:"avail_bytes"`
	TotalBytes int64    `json:"total_bytes"`
	TotalFiles int64    `json:"total_files"`
	FreeFiles  int64    `json:"free_files"`
	NameMax    int      `json:"name_max"`
	Fstype     string   `json:"fstype"`
	Flags      []string `json:"flags"`
}

// DirEntry represents a directory entry
type DirEntry struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	RealPath string    `json:"realpath"`
	Type     string    `json:"type"`
	Size     int64     `json:"size"`
	Mode     int       `json:"mode"`
	UID      int       `json:"uid"`
	GID      int       `json:"gid"`
	Mtime    time.Time `json:"mtime"`
	HasACL   bool      `json:"acl"`
}

// ACL represents an Access Control List
type ACL struct {
	ACLType string     `json:"acltype"`
	UID     int        `json:"uid"`
	GID     int        `json:"gid"`
	ACL     []ACLEntry `json:"acl"`
	Trivial bool       `json:"trivial"`
	Path    string     `json:"path"`
}

// ACLEntry represents a single ACL entry
type ACLEntry struct {
	Tag   string `json:"tag"`
	ID    *int   `json:"id,omitempty"`
	Type  string `json:"type"`
	Perms any    `json:"perms"` // Can be simplified string or detailed object
	Flags any    `json:"flags,omitempty"`
	Who   string `json:"who,omitempty"`
}

// NFS41Flags represents NFSv4.1 ACL flags
type NFS41Flags struct {
	Autoinherit bool `json:"autoinherit"`
	Protected   bool `json:"protected"`
}

// ChownRequest represents parameters for filesystem.chown
type ChownRequest struct {
	Path    string       `json:"path"`
	UID     *int         `json:"uid,omitempty"`
	GID     *int         `json:"gid,omitempty"`
	Options ChownOptions `json:"options"`
}

// ChownOptions represents options for chown operation
type ChownOptions struct {
	Recursive bool `json:"recursive"`
	Traverse  bool `json:"traverse"`
}

// SetACLRequest represents parameters for filesystem.setacl
type SetACLRequest struct {
	Path       string        `json:"path"`
	UID        *int          `json:"uid,omitempty"`
	GID        *int          `json:"gid,omitempty"`
	DACL       []ACLEntry    `json:"dacl"`
	NFS41Flags *NFS41Flags   `json:"nfs41_flags,omitempty"`
	ACLType    ACLType       `json:"acltype"`
	Options    SetACLOptions `json:"options"`
}

// SetACLOptions represents options for setacl operation
type SetACLOptions struct {
	StripACL     bool `json:"stripacl"`
	Recursive    bool `json:"recursive"`
	Traverse     bool `json:"traverse"`
	Canonicalize bool `json:"canonicalize"`
}

// SetPermRequest represents parameters for filesystem.setperm
type SetPermRequest struct {
	Path    string         `json:"path"`
	Mode    *string        `json:"mode,omitempty"`
	UID     *int           `json:"uid,omitempty"`
	GID     *int           `json:"gid,omitempty"`
	Options SetPermOptions `json:"options"`
}

// SetPermOptions represents options for setperm operation
type SetPermOptions struct {
	StripACL  bool `json:"stripacl"`
	Recursive bool `json:"recursive"`
	Traverse  bool `json:"traverse"`
}

// PutFileOptions represents options for filesystem.put
type PutFileOptions struct {
	Append bool `json:"append"`
	Mode   *int `json:"mode,omitempty"`
}

// DefaultACLType represents ACL template types
type DefaultACLType string

const (
	DefaultACLTypeOpen       DefaultACLType = "OPEN"
	DefaultACLTypeRestricted DefaultACLType = "RESTRICTED"
	DefaultACLTypeHome       DefaultACLType = "HOME"
	DefaultACLTypeDomainHome DefaultACLType = "DOMAIN_HOME"
)

// ShareType represents share types for ACL templates
type ShareType string

const (
	ShareTypeNone ShareType = "NONE"
	ShareTypeAFP  ShareType = "AFP"
	ShareTypeSMB  ShareType = "SMB"
	ShareTypeNFS  ShareType = "NFS"
)

// ACLType represents ACL types
type ACLType string

const (
	ACLTypeNFS4    ACLType = "NFS4"
	ACLTypePOSIX1E ACLType = "POSIX1E"
	ACLTypeRICH    ACLType = "RICH"
)

// Basic filesystem operations

// Stat returns filesystem stat information for a path
func (f *FilesystemClient) Stat(ctx context.Context, path string) (*FilesystemStat, error) {
	var result FilesystemStat
	err := f.client.Call(ctx, "filesystem.stat", []any{path}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Statfs returns filesystem statistics for a path
func (f *FilesystemClient) Statfs(ctx context.Context, path string) (*FilesystemStatfs, error) {
	var result FilesystemStatfs
	err := f.client.Call(ctx, "filesystem.statfs", []any{path}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDir returns directory contents
func (f *FilesystemClient) ListDir(ctx context.Context, path string) ([]DirEntry, error) {
	var result []DirEntry
	err := f.client.Call(ctx, "filesystem.listdir", []any{path}, &result)
	return result, err
}

// ACL operations

// GetACL returns the ACL for a path
func (f *FilesystemClient) GetACL(ctx context.Context, path string, simplified bool) (*ACL, error) {
	var result ACL
	err := f.client.Call(ctx, "filesystem.getacl", []any{path, simplified}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SetACL sets the ACL for a path (asynchronous job)
func (f *FilesystemClient) SetACL(ctx context.Context, req *SetACLRequest) error {
	return f.client.CallJob(ctx, "filesystem.setacl", []any{
		req.Path, req.UID, req.GID, req.DACL, req.NFS41Flags, req.ACLType, req.Options,
	}, nil)
}

// IsACLTrivial checks if the ACL can be expressed as a simple file mode
func (f *FilesystemClient) IsACLTrivial(ctx context.Context, path string) (bool, error) {
	var result bool
	err := f.client.Call(ctx, "filesystem.acl_is_trivial", []any{path}, &result)
	return result, err
}

// GetDefaultACL returns a default ACL template
func (f *FilesystemClient) GetDefaultACL(ctx context.Context, aclType DefaultACLType, shareType ShareType) (*ACL, error) {
	var result ACL
	err := f.client.Call(ctx, "filesystem.get_default_acl", []any{string(aclType), string(shareType)}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDefaultACLChoices returns available default ACL types
func (f *FilesystemClient) GetDefaultACLChoices(ctx context.Context) ([]string, error) {
	var result []string
	err := f.client.Call(ctx, "filesystem.default_acl_choices", []any{}, &result)
	return result, err
}

// Permission operations

// SetPermissions sets permissions for a path (asynchronous job)
func (f *FilesystemClient) SetPermissions(ctx context.Context, req *SetPermRequest) error {
	return f.client.CallJob(ctx, "filesystem.setperm", []any{
		req.Path, req.Mode, req.UID, req.GID, req.Options,
	}, nil)
}

// ChangeOwner changes owner/group of a path (asynchronous job)
func (f *FilesystemClient) ChangeOwner(ctx context.Context, req *ChownRequest) error {
	return f.client.CallJob(ctx, "filesystem.chown", []any{
		req.Path, req.UID, req.GID, req.Options,
	}, nil)
}

// File operations

// GetFile downloads a file (asynchronous job with download support)
func (f *FilesystemClient) GetFile(ctx context.Context, path string) error {
	return f.client.CallJob(ctx, "filesystem.get", []any{path}, nil)
}

// PutFile uploads a file (asynchronous job with upload support)
func (f *FilesystemClient) PutFile(ctx context.Context, path string, options *PutFileOptions) error {
	if options == nil {
		options = &PutFileOptions{}
	}
	return f.client.CallJob(ctx, "filesystem.put", []any{path, *options}, nil)
}

// Helper methods for common operations

// CreateDefaultACL creates a default ACL for a given purpose
func (f *FilesystemClient) CreateDefaultACL(ctx context.Context, aclType DefaultACLType) (*ACL, error) {
	return f.GetDefaultACL(ctx, aclType, ShareTypeNone)
}

// CreateShareACL creates a default ACL for a specific share type
func (f *FilesystemClient) CreateShareACL(ctx context.Context, shareType ShareType) (*ACL, error) {
	return f.GetDefaultACL(ctx, DefaultACLTypeOpen, shareType)
}

// SetSimplePermissions sets simple octal permissions
func (f *FilesystemClient) SetSimplePermissions(ctx context.Context, path, mode string, recursive bool) error {
	req := &SetPermRequest{
		Path: path,
		Mode: &mode,
		Options: SetPermOptions{
			Recursive: recursive,
			StripACL:  true,
		},
	}
	return f.SetPermissions(ctx, req)
}

// SetOwnership sets file ownership
func (f *FilesystemClient) SetOwnership(ctx context.Context, path string, uid, gid *int, recursive bool) error {
	req := &ChownRequest{
		Path: path,
		UID:  uid,
		GID:  gid,
		Options: ChownOptions{
			Recursive: recursive,
		},
	}
	return f.ChangeOwner(ctx, req)
}
