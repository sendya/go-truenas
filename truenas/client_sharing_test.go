package truenas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test data for AFP shares
var (
	TestAFPShare = AFPShare{
		ID:               1,
		Path:             "/mnt/tank/afp-share",
		Home:             false,
		Name:             "test-afp-share",
		Comment:          "Test AFP Share",
		Allow:            []string{"user1", "user2"},
		Deny:             []string{"user3"},
		RO:               []string{"ro-group"},
		RW:               []string{"rw-group"},
		TimeMachine:      true,
		TimeMachineQuota: 1000,
		NoDev:            false,
		NoStat:           false,
		UPriv:            true,
		FPerm:            "644",
		DPerm:            "755",
		UMask:            "022",
		HostsAllow:       []string{"192.168.1.0/24"},
		HostsDeny:        []string{"192.168.2.0/24"},
		VUID:             Ptr("test-vuid"),
		AuxParams:        "test aux params",
		Enabled:          true,
	}

	TestAFPShareRequest = AFPShareRequest{
		Path:             "/mnt/tank/afp-share",
		Home:             false,
		Name:             "test-afp-share",
		Comment:          "Test AFP Share",
		Allow:            []string{"user1", "user2"},
		Deny:             []string{"user3"},
		RO:               []string{"ro-group"},
		RW:               []string{"rw-group"},
		TimeMachine:      true,
		TimeMachineQuota: 1000,
		NoDev:            false,
		NoStat:           false,
		UPriv:            true,
		FPerm:            "644",
		DPerm:            "755",
		UMask:            "022",
		HostsAllow:       []string{"192.168.1.0/24"},
		HostsDeny:        []string{"192.168.2.0/24"},
		VUID:             Ptr("test-vuid"),
		AuxParams:        "test aux params",
		Enabled:          true,
	}
)

// Test data for NFS shares
var (
	TestNFSShare = NFSShare{
		ID:           1,
		Path:         "/mnt/tank/nfs-share",
		Aliases:      []string{"alias1", "alias2"},
		Comment:      "Test NFS Share",
		Networks:     []string{"192.168.1.0/24"},
		Hosts:        []string{"host1", "host2"},
		RO:           false,
		MapRootUser:  Ptr("root"),
		MapRootGroup: Ptr("wheel"),
		MapAllUser:   Ptr("nobody"),
		MapAllGroup:  Ptr("nogroup"),
		Security:     []string{"sys", "krb5"},
		Enabled:      true,
		Locked:       false,
	}

	TestNFSShareRequest = NFSShareRequest{
		Path:         "/mnt/tank/nfs-share",
		Comment:      "Test NFS Share",
		Networks:     []string{"192.168.1.0/24"},
		Hosts:        []string{"host1", "host2"},
		RO:           false,
		MapRootUser:  Ptr("root"),
		MapRootGroup: Ptr("wheel"),
		MapAllUser:   Ptr("nobody"),
		MapAllGroup:  Ptr("nogroup"),
		Security:     []string{"sys", "krb5"},
		Enabled:      true,
	}
)

// Test data for SMB shares
var (
	TestSMBShare = SMBShare{
		ID:               1,
		Purpose:          SMBPurposeDefaultShare,
		Path:             "/mnt/tank/smb-share",
		PathSuffix:       "",
		Home:             false,
		Name:             "test-smb-share",
		Comment:          "Test SMB Share",
		RO:               false,
		Browsable:        true,
		TimeMachine:      false,
		RecycleBin:       false,
		GuestOK:          false,
		ABE:              false,
		HostsAllow:       []string{"192.168.1.0/24"},
		HostsDeny:        []string{"192.168.2.0/24"},
		AAPLNameMangling: false,
		ACL:              true,
		DurableHandle:    false,
		ShadowCopy:       true,
		Streams:          true,
		FSRVP:            false,
		AuxSMBConf:       "test smb config",
		Enabled:          true,
	}

	TestSMBShareRequest = SMBShareRequest{
		Purpose:          SMBPurposeDefaultShare,
		Path:             "/mnt/tank/smb-share",
		PathSuffix:       "",
		Home:             false,
		Name:             "test-smb-share",
		Comment:          "Test SMB Share",
		RO:               false,
		Browsable:        true,
		TimeMachine:      false,
		RecycleBin:       false,
		GuestOK:          false,
		ABE:              false,
		HostsAllow:       []string{"192.168.1.0/24"},
		HostsDeny:        []string{"192.168.2.0/24"},
		AAPLNameMangling: false,
		ACL:              true,
		DurableHandle:    false,
		ShadowCopy:       true,
		Streams:          true,
		FSRVP:            false,
		AuxSMBConf:       "test smb config",
		Enabled:          true,
	}

	TestSMBPresets = []SMBPreset{
		{
			Name:        "Default Share",
			Description: "Default SMB share configuration",
			Config:      map[string]any{"browsable": true, "ro": false},
		},
		{
			Name:        "Time Machine",
			Description: "Optimized for Apple Time Machine backups",
			Config:      map[string]any{"timemachine": true, "browsable": false},
		},
	}
)

// Test data for WebDAV shares
var (
	TestWebDAVShare = WebDAVShare{
		ID:      1,
		Perm:    true,
		RO:      false,
		Comment: "Test WebDAV Share",
		Name:    "test-webdav-share",
		Path:    "/mnt/tank/webdav-share",
		Enabled: true,
	}

	TestWebDAVShareRequest = WebDAVShareRequest{
		Perm:    true,
		RO:      false,
		Comment: "Test WebDAV Share",
		Name:    "test-webdav-share",
		Path:    "/mnt/tank/webdav-share",
		Enabled: true,
	}
)

// AFP Sharing Client Tests

func TestSharingAFPClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockShares := []AFPShare{TestAFPShare, {
		ID:      2,
		Path:    "/mnt/tank/afp-share2",
		Name:    "test-afp-share2",
		Comment: "Test AFP Share 2",
		Enabled: false,
	}}
	server.SetResponse("sharing.afp.query", mockShares)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.AFP.List(ctx)
	require.NoError(t, err)
	assert.Len(t, shares, 2)
	assert.Equal(t, "test-afp-share", shares[0].Name)
	assert.Equal(t, "test-afp-share2", shares[1].Name)
}

func TestSharingAFPClient_List_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.afp.query", 500, "AFP service unavailable")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.AFP.List(ctx)
	assert.Error(t, err)
	assert.Nil(t, shares)
	assert.Contains(t, err.Error(), "AFP service unavailable")
}

func TestSharingAFPClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.afp.query", []AFPShare{TestAFPShare})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.AFP.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "test-afp-share", share.Name)
	assert.Equal(t, "/mnt/tank/afp-share", share.Path)
	assert.True(t, share.TimeMachine)
}

func TestSharingAFPClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.afp.query", []AFPShare{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.AFP.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, share)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestSharingAFPClient_Get_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.afp.query", 500, "Database error")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.AFP.Get(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, share)
	assert.Contains(t, err.Error(), "Database error")
}

func TestSharingAFPClient_Create(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.afp.create", TestAFPShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.AFP.Create(ctx, &TestAFPShareRequest)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "test-afp-share", share.Name)
	assert.Equal(t, "/mnt/tank/afp-share", share.Path)
	assert.True(t, share.TimeMachine)
}

func TestSharingAFPClient_Create_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.afp.create", 400, "Invalid path")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.AFP.Create(ctx, &TestAFPShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Invalid path")
}

func TestSharingAFPClient_Update(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	updatedShare := TestAFPShare
	updatedShare.Comment = "Updated AFP Share"
	server.SetResponse("sharing.afp.update", updatedShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	updateReq := TestAFPShareRequest
	updateReq.Comment = "Updated AFP Share"

	ctx := NewTestContext(t)
	share, err := client.Sharing.AFP.Update(ctx, 1, &updateReq)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "Updated AFP Share", share.Comment)
}

func TestSharingAFPClient_Update_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.afp.update", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.AFP.Update(ctx, 999, &TestAFPShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Share not found")
}

func TestSharingAFPClient_Delete(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.afp.delete", nil)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.AFP.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestSharingAFPClient_Delete_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.afp.delete", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.AFP.Delete(ctx, 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Share not found")
}

// NFS Sharing Client Tests

func TestSharingNFSClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockShares := []NFSShare{TestNFSShare, {
		ID:      2,
		Path:    "/mnt/tank/nfs-share2",
		Comment: "Test NFS Share 2",
		Enabled: false,
	}}
	server.SetResponse("sharing.nfs.query", mockShares)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.NFS.List(ctx)
	require.NoError(t, err)
	assert.Len(t, shares, 2)
	assert.Equal(t, "/mnt/tank/nfs-share", shares[0].Path)
	assert.Equal(t, "/mnt/tank/nfs-share2", shares[1].Path)
}

func TestSharingNFSClient_List_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.nfs.query", 500, "NFS service unavailable")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.NFS.List(ctx)
	assert.Error(t, err)
	assert.Nil(t, shares)
	assert.Contains(t, err.Error(), "NFS service unavailable")
}

func TestSharingNFSClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.nfs.query", []NFSShare{TestNFSShare})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.NFS.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "/mnt/tank/nfs-share", share.Path)
	assert.Equal(t, "Test NFS Share", share.Comment)
	assert.Len(t, share.Networks, 1)
}

func TestSharingNFSClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.nfs.query", []NFSShare{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.NFS.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, share)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestSharingNFSClient_Get_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.nfs.query", 500, "Database error")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.NFS.Get(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, share)
	assert.Contains(t, err.Error(), "Database error")
}

func TestSharingNFSClient_Create(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.nfs.create", TestNFSShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.NFS.Create(ctx, &TestNFSShareRequest)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "/mnt/tank/nfs-share", share.Path)
	assert.Equal(t, "Test NFS Share", share.Comment)
}

func TestSharingNFSClient_Create_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.nfs.create", 400, "Invalid network address")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.NFS.Create(ctx, &TestNFSShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Invalid network address")
}

func TestSharingNFSClient_Update(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	updatedShare := TestNFSShare
	updatedShare.Comment = "Updated NFS Share"
	server.SetResponse("sharing.nfs.update", updatedShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	updateReq := TestNFSShareRequest
	updateReq.Comment = "Updated NFS Share"

	ctx := NewTestContext(t)
	share, err := client.Sharing.NFS.Update(ctx, 1, &updateReq)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "Updated NFS Share", share.Comment)
}

func TestSharingNFSClient_Update_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.nfs.update", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.NFS.Update(ctx, 999, &TestNFSShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Share not found")
}

func TestSharingNFSClient_Delete(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.nfs.delete", nil)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.NFS.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestSharingNFSClient_Delete_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.nfs.delete", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.NFS.Delete(ctx, 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Share not found")
}

func TestSharingNFSClient_GetHumanIdentifier(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.nfs.human_identifier", "/mnt/tank/nfs-share")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	identifier, err := client.Sharing.NFS.GetHumanIdentifier(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "/mnt/tank/nfs-share", identifier)
}

func TestSharingNFSClient_GetHumanIdentifier_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.nfs.human_identifier", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	identifier, err := client.Sharing.NFS.GetHumanIdentifier(ctx, 999)
	assert.Error(t, err)
	assert.Empty(t, identifier)
	assert.Contains(t, err.Error(), "Share not found")
}

// SMB Sharing Client Tests

func TestSharingSMBClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockShares := []SMBShare{TestSMBShare, {
		ID:      2,
		Path:    "/mnt/tank/smb-share2",
		Name:    "test-smb-share2",
		Comment: "Test SMB Share 2",
		Enabled: false,
	}}
	server.SetResponse("sharing.smb.query", mockShares)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.SMB.List(ctx)
	require.NoError(t, err)
	assert.Len(t, shares, 2)
	assert.Equal(t, "test-smb-share", shares[0].Name)
	assert.Equal(t, "test-smb-share2", shares[1].Name)
}

func TestSharingSMBClient_List_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.smb.query", 500, "SMB service unavailable")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.SMB.List(ctx)
	assert.Error(t, err)
	assert.Nil(t, shares)
	assert.Contains(t, err.Error(), "SMB service unavailable")
}

func TestSharingSMBClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.smb.query", []SMBShare{TestSMBShare})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.SMB.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "test-smb-share", share.Name)
	assert.Equal(t, "/mnt/tank/smb-share", share.Path)
	assert.Equal(t, SMBPurposeDefaultShare, share.Purpose)
}

func TestSharingSMBClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.smb.query", []SMBShare{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.SMB.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, share)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestSharingSMBClient_Get_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.smb.query", 500, "Database error")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.SMB.Get(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, share)
	assert.Contains(t, err.Error(), "Database error")
}

func TestSharingSMBClient_Create(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.smb.create", TestSMBShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.SMB.Create(ctx, &TestSMBShareRequest)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "test-smb-share", share.Name)
	assert.Equal(t, "/mnt/tank/smb-share", share.Path)
}

func TestSharingSMBClient_Create_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.smb.create", 400, "Invalid share name")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.SMB.Create(ctx, &TestSMBShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Invalid share name")
}

func TestSharingSMBClient_Update(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	updatedShare := TestSMBShare
	updatedShare.Comment = "Updated SMB Share"
	server.SetResponse("sharing.smb.update", updatedShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	updateReq := TestSMBShareRequest
	updateReq.Comment = "Updated SMB Share"

	ctx := NewTestContext(t)
	share, err := client.Sharing.SMB.Update(ctx, 1, &updateReq)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "Updated SMB Share", share.Comment)
}

func TestSharingSMBClient_Update_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.smb.update", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.SMB.Update(ctx, 999, &TestSMBShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Share not found")
}

func TestSharingSMBClient_Delete(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.smb.delete", nil)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.SMB.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestSharingSMBClient_Delete_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.smb.delete", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.SMB.Delete(ctx, 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Share not found")
}

func TestSharingSMBClient_GetPresets(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.smb.presets", TestSMBPresets)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	presets, err := client.Sharing.SMB.GetPresets(ctx)
	require.NoError(t, err)
	assert.Len(t, presets, 2)
	assert.Equal(t, "Default Share", presets[0].Name)
	assert.Equal(t, "Time Machine", presets[1].Name)
}

func TestSharingSMBClient_GetPresets_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.smb.presets", 500, "Service unavailable")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	presets, err := client.Sharing.SMB.GetPresets(ctx)
	assert.Error(t, err)
	assert.Nil(t, presets)
	assert.Contains(t, err.Error(), "Service unavailable")
}

// WebDAV Sharing Client Tests

func TestSharingWebDAVClient_List(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	mockShares := []WebDAVShare{TestWebDAVShare, {
		ID:      2,
		Name:    "test-webdav-share2",
		Path:    "/mnt/tank/webdav-share2",
		Comment: "Test WebDAV Share 2",
		Enabled: false,
	}}
	server.SetResponse("sharing.webdav.query", mockShares)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.WebDAV.List(ctx)
	require.NoError(t, err)
	assert.Len(t, shares, 2)
	assert.Equal(t, "test-webdav-share", shares[0].Name)
	assert.Equal(t, "test-webdav-share2", shares[1].Name)
}

func TestSharingWebDAVClient_List_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.webdav.query", 500, "WebDAV service unavailable")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	shares, err := client.Sharing.WebDAV.List(ctx)
	assert.Error(t, err)
	assert.Nil(t, shares)
	assert.Contains(t, err.Error(), "WebDAV service unavailable")
}

func TestSharingWebDAVClient_Get(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.webdav.query", []WebDAVShare{TestWebDAVShare})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.WebDAV.Get(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "test-webdav-share", share.Name)
	assert.Equal(t, "/mnt/tank/webdav-share", share.Path)
	assert.True(t, share.Perm)
	assert.False(t, share.RO)
}

func TestSharingWebDAVClient_Get_NotFound(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.webdav.query", []WebDAVShare{})

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.WebDAV.Get(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, share)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

func TestSharingWebDAVClient_Get_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.webdav.query", 500, "Database error")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.WebDAV.Get(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, share)
	assert.Contains(t, err.Error(), "Database error")
}

func TestSharingWebDAVClient_Create(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.webdav.create", TestWebDAVShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.WebDAV.Create(ctx, &TestWebDAVShareRequest)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "test-webdav-share", share.Name)
	assert.Equal(t, "/mnt/tank/webdav-share", share.Path)
}

func TestSharingWebDAVClient_Create_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.webdav.create", 400, "Invalid configuration")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.WebDAV.Create(ctx, &TestWebDAVShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Invalid configuration")
}

func TestSharingWebDAVClient_Update(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	updatedShare := TestWebDAVShare
	updatedShare.Comment = "Updated WebDAV Share"
	server.SetResponse("sharing.webdav.update", updatedShare)

	client := server.CreateTestClient(t)
	defer client.Close()

	updateReq := TestWebDAVShareRequest
	updateReq.Comment = "Updated WebDAV Share"

	ctx := NewTestContext(t)
	share, err := client.Sharing.WebDAV.Update(ctx, 1, &updateReq)
	require.NoError(t, err)
	require.NotNil(t, share)
	assert.Equal(t, "Updated WebDAV Share", share.Comment)
}

func TestSharingWebDAVClient_Update_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.webdav.update", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	share, err := client.Sharing.WebDAV.Update(ctx, 999, &TestWebDAVShareRequest)
	assert.Error(t, err)
	assert.NotNil(t, share) // API returns empty struct even on error
	assert.Contains(t, err.Error(), "Share not found")
}

func TestSharingWebDAVClient_Delete(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetResponse("sharing.webdav.delete", nil)

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.WebDAV.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestSharingWebDAVClient_Delete_Error(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	server.SetError("sharing.webdav.delete", 404, "Share not found")

	client := server.CreateTestClient(t)
	defer client.Close()

	ctx := NewTestContext(t)
	err := client.Sharing.WebDAV.Delete(ctx, 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Share not found")
}

// Test NewSharingClient function
func TestNewSharingClient(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	sharingClient := NewSharingClient(client)
	require.NotNil(t, sharingClient)
	require.NotNil(t, sharingClient.AFP)
	require.NotNil(t, sharingClient.NFS)
	require.NotNil(t, sharingClient.SMB)
	require.NotNil(t, sharingClient.WebDAV)
}

// Test individual client constructors
func TestNewSharingAFPClient(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	afpClient := NewSharingAFPClient(client)
	require.NotNil(t, afpClient)
	assert.Equal(t, client, afpClient.client)
}

func TestNewSharingNFSClient(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	nfsClient := NewSharingNFSClient(client)
	require.NotNil(t, nfsClient)
	assert.Equal(t, client, nfsClient.client)
}

func TestNewSharingSMBClient(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	smbClient := NewSharingSMBClient(client)
	require.NotNil(t, smbClient)
	assert.Equal(t, client, smbClient.client)
}

func TestNewSharingWebDAVClient(t *testing.T) {
	t.Parallel()
	server := NewTestServer(t)
	defer server.Close()

	client := server.CreateTestClient(t)
	defer client.Close()

	webdavClient := NewSharingWebDAVClient(client)
	require.NotNil(t, webdavClient)
	assert.Equal(t, client, webdavClient.client)
}
