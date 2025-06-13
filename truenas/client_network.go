package truenas

import (
	"context"
	"fmt"
)

// InterfaceType represents network interface types
type InterfaceType string

const (
	InterfaceTypePhysical InterfaceType = "PHYSICAL"
	InterfaceTypeBridge   InterfaceType = "BRIDGE"
	InterfaceTypeLinkAgg  InterfaceType = "LINK_AGGREGATION"
	InterfaceTypeVLAN     InterfaceType = "VLAN"
)

// LAGProtocol represents link aggregation protocols
type LAGProtocol string

const (
	LAGProtocolLACP        LAGProtocol = "LACP"
	LAGProtocolFailover    LAGProtocol = "FAILOVER"
	LAGProtocolLoadBalance LAGProtocol = "LOADBALANCE"
	LAGProtocolRoundRobin  LAGProtocol = "ROUNDROBIN"
	LAGProtocolNone        LAGProtocol = "NONE"
)

// NetworkClient provides methods for network management
type NetworkClient struct {
	client *Client
}

// NewNetworkClient creates a new network client
func NewNetworkClient(client *Client) *NetworkClient {
	return &NetworkClient{client: client}
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	ID          int                     `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Type        InterfaceType           `json:"type"`
	IPV4DHCP    bool                    `json:"ipv4_dhcp"`
	IPV6Auto    bool                    `json:"ipv6_auto"`
	Options     string                  `json:"options"`
	MTU         int                     `json:"mtu"`
	Disable     bool                    `json:"disable_offload_capabilities"`
	Aliases     []NetworkInterfaceAlias `json:"aliases"`
	HasDefault  bool                    `json:"has_default_route"`
	State       NetworkInterfaceState   `json:"state"`

	// Bridge-specific fields
	BridgeMembers []string `json:"bridge_members,omitempty"`

	// LAG-specific fields
	LagProtocol LAGProtocol `json:"lag_protocol,omitempty"`
	LagPorts    []string    `json:"lag_ports,omitempty"`

	// VLAN-specific fields
	VlanParent string `json:"vlan_parent_interface,omitempty"`
	VlanTag    int    `json:"vlan_tag,omitempty"`
	VlanPcp    int    `json:"vlan_pcp,omitempty"`

	// Failover fields
	FailoverCritical       bool     `json:"failover_critical,omitempty"`
	FailoverGroup          int      `json:"failover_group,omitempty"`
	FailoverVhid           int      `json:"failover_vhid,omitempty"`
	FailoverAliases        []string `json:"failover_aliases,omitempty"`
	FailoverVirtualAliases []string `json:"failover_virtual_aliases,omitempty"`
}

// NetworkInterfaceAlias represents an interface alias/additional IP
type NetworkInterfaceAlias struct {
	Type      string `json:"type"`
	Address   string `json:"address"`
	Netmask   int    `json:"netmask"`
	Broadcast string `json:"broadcast,omitempty"`
}

// NetworkInterfaceState represents the current state of an interface
type NetworkInterfaceState struct {
	Name               string           `json:"name"`
	OrigName           string           `json:"orig_name"`
	Description        string           `json:"description"`
	MTU                int              `json:"mtu"`
	Cloned             bool             `json:"cloned"`
	Flags              []string         `json:"flags"`
	Nd6Flags           []string         `json:"nd6_flags"`
	Capabilities       []string         `json:"capabilities"`
	LinkState          string           `json:"link_state"`
	MediaType          string           `json:"media_type"`
	MediaSubtype       string           `json:"media_subtype"`
	ActiveMediaType    string           `json:"active_media_type"`
	ActiveMediaSubtype string           `json:"active_media_subtype"`
	Supported          []string         `json:"supported_media"`
	MediaOptions       []string         `json:"media_options"`
	LinkAddress        string           `json:"link_address"`
	Addresses          []NetworkAddress `json:"addresses"`
}

// NetworkAddress represents an IP address on an interface
type NetworkAddress struct {
	Type      string `json:"type"`
	Address   string `json:"address"`
	Netmask   int    `json:"netmask"`
	Broadcast string `json:"broadcast,omitempty"`
}

// NetworkConfiguration represents global network configuration
type NetworkConfiguration struct {
	ID                  int                  `json:"id"`
	Hostname            string               `json:"hostname"`
	HostnameB           string               `json:"hostname_b"`
	HostnameVirtual     string               `json:"hostname_virtual"`
	Domain              string               `json:"domain"`
	Domains             []string             `json:"domains"`
	IPv4Gateway         string               `json:"ipv4gateway"`
	IPv6Gateway         string               `json:"ipv6gateway"`
	Nameserver1         string               `json:"nameserver1"`
	Nameserver2         string               `json:"nameserver2"`
	Nameserver3         string               `json:"nameserver3"`
	HTTPProxy           string               `json:"httpproxy"`
	NetwaitSummary      any                  `json:"netwait_summary"`
	ActivityType        string               `json:"activity_type"`
	ActivityInterval    int                  `json:"activity_interval"`
	HostnameDB          []string             `json:"hostname_db"`
	HostnameLocal       []string             `json:"hostname_local"`
	ServiceAnnouncement *ServiceAnnouncement `json:"service_announcement,omitempty"`
	NetwaitEnabled      bool                 `json:"netwait_enabled"`
	NetwaitIP           []string             `json:"netwait_ip"`
	Hosts               string               `json:"hosts"`
}

// ServiceAnnouncement represents service announcement configuration
type ServiceAnnouncement struct {
	Netbios bool `json:"netbios"`
	MDNS    bool `json:"mdns"`
	WSD     bool `json:"wsd"`
}

// StaticRoute represents a static network route
type StaticRoute struct {
	ID          int    `json:"id"`
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Description string `json:"description"`
}

// NetworkInterfaceCreateRequest represents parameters for interface.create
type NetworkInterfaceCreateRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Type        InterfaceType           `json:"type"`
	IPV4DHCP    bool                    `json:"ipv4_dhcp,omitempty"`
	IPV6Auto    bool                    `json:"ipv6_auto,omitempty"`
	Options     string                  `json:"options,omitempty"`
	MTU         int                     `json:"mtu,omitempty"`
	Disable     bool                    `json:"disable_offload_capabilities,omitempty"`
	Aliases     []NetworkInterfaceAlias `json:"aliases,omitempty"`

	// Bridge-specific fields
	BridgeMembers []string `json:"bridge_members,omitempty"`

	// LAG-specific fields
	LagProtocol *LAGProtocol `json:"lag_protocol,omitempty"`
	LagPorts    []string     `json:"lag_ports,omitempty"`

	// VLAN-specific fields
	VlanParent *string `json:"vlan_parent_interface,omitempty"`
	VlanTag    *int    `json:"vlan_tag,omitempty"`
	VlanPcp    *int    `json:"vlan_pcp,omitempty"`

	// Failover fields
	FailoverCritical       *bool    `json:"failover_critical,omitempty"`
	FailoverGroup          *int     `json:"failover_group,omitempty"`
	FailoverVhid           *int     `json:"failover_vhid,omitempty"`
	FailoverAliases        []string `json:"failover_aliases,omitempty"`
	FailoverVirtualAliases []string `json:"failover_virtual_aliases,omitempty"`
}

// NetworkInterfaceUpdateRequest represents parameters for interface.update
type NetworkInterfaceUpdateRequest struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	IPV4DHCP    *bool                   `json:"ipv4_dhcp,omitempty"`
	IPV6Auto    *bool                   `json:"ipv6_auto,omitempty"`
	Options     *string                 `json:"options,omitempty"`
	MTU         *int                    `json:"mtu,omitempty"`
	Disable     *bool                   `json:"disable_offload_capabilities,omitempty"`
	Aliases     []NetworkInterfaceAlias `json:"aliases,omitempty"`

	// Bridge-specific fields
	BridgeMembers []string `json:"bridge_members,omitempty"`

	// LAG-specific fields
	LagProtocol *LAGProtocol `json:"lag_protocol,omitempty"`
	LagPorts    []string     `json:"lag_ports,omitempty"`

	// VLAN-specific fields
	VlanParent *string `json:"vlan_parent_interface,omitempty"`
	VlanTag    *int    `json:"vlan_tag,omitempty"`
	VlanPcp    *int    `json:"vlan_pcp,omitempty"`

	// Failover fields
	FailoverCritical       *bool    `json:"failover_critical,omitempty"`
	FailoverGroup          *int     `json:"failover_group,omitempty"`
	FailoverVhid           *int     `json:"failover_vhid,omitempty"`
	FailoverAliases        []string `json:"failover_aliases,omitempty"`
	FailoverVirtualAliases []string `json:"failover_virtual_aliases,omitempty"`
}

// StaticRouteCreateRequest represents parameters for staticroute.create
type StaticRouteCreateRequest struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Description string `json:"description,omitempty"`
}

// Interface Management

// ListInterfaces returns all network interfaces
func (n *NetworkClient) ListInterfaces(ctx context.Context) ([]NetworkInterface, error) {
	var result []NetworkInterface
	err := n.client.Call(ctx, "interface.query", []any{}, &result)
	return result, err
}

// GetInterface returns a specific interface by ID
func (n *NetworkClient) GetInterface(ctx context.Context, id int) (*NetworkInterface, error) {
	var result []NetworkInterface
	err := n.client.Call(ctx, "interface.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("network_interface", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// GetInterfaceByName returns a specific interface by name
func (n *NetworkClient) GetInterfaceByName(ctx context.Context, name string) (*NetworkInterface, error) {
	var result []NetworkInterface
	err := n.client.Call(ctx, "interface.query", []any{[]any{[]any{"name", "=", name}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("network_interface", fmt.Sprintf("name %s", name))
	}
	return &result[0], nil
}

// CreateInterface creates a new network interface
func (n *NetworkClient) CreateInterface(ctx context.Context, req *NetworkInterfaceCreateRequest) (*NetworkInterface, error) {
	var result NetworkInterface
	err := n.client.Call(ctx, "interface.create", []any{*req}, &result)
	return &result, err
}

// UpdateInterface updates an existing interface
func (n *NetworkClient) UpdateInterface(ctx context.Context, id int, req *NetworkInterfaceUpdateRequest) (*NetworkInterface, error) {
	var result NetworkInterface
	err := n.client.Call(ctx, "interface.update", []any{id, *req}, &result)
	return &result, err
}

// DeleteInterface deletes an interface
func (n *NetworkClient) DeleteInterface(ctx context.Context, id int) error {
	return n.client.Call(ctx, "interface.delete", []any{id}, nil)
}

// Global Network Configuration

// GetConfiguration returns global network configuration
func (n *NetworkClient) GetConfiguration(ctx context.Context) (*NetworkConfiguration, error) {
	var result NetworkConfiguration
	err := n.client.Call(ctx, "network.configuration.config", []any{}, &result)
	return &result, err
}

// UpdateConfiguration updates global network configuration
func (n *NetworkClient) UpdateConfiguration(ctx context.Context, config *NetworkConfiguration) (*NetworkConfiguration, error) {
	var result NetworkConfiguration
	err := n.client.Call(ctx, "network.configuration.update", []any{*config}, &result)
	return &result, err
}

// Static Routes

// ListStaticRoutes returns all static routes
func (n *NetworkClient) ListStaticRoutes(ctx context.Context) ([]StaticRoute, error) {
	var result []StaticRoute
	err := n.client.Call(ctx, "staticroute.query", []any{}, &result)
	return result, err
}

// GetStaticRoute returns a specific static route by ID
func (n *NetworkClient) GetStaticRoute(ctx context.Context, id int) (*StaticRoute, error) {
	var result []StaticRoute
	err := n.client.Call(ctx, "staticroute.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("static_route", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// CreateStaticRoute creates a new static route
func (n *NetworkClient) CreateStaticRoute(ctx context.Context, req StaticRouteCreateRequest) (*StaticRoute, error) {
	var result StaticRoute
	err := n.client.Call(ctx, "staticroute.create", []any{req}, &result)
	return &result, err
}

// UpdateStaticRoute updates an existing static route
func (n *NetworkClient) UpdateStaticRoute(ctx context.Context, id int, req StaticRouteCreateRequest) (*StaticRoute, error) {
	var result StaticRoute
	err := n.client.Call(ctx, "staticroute.update", []any{id, req}, &result)
	return &result, err
}

// DeleteStaticRoute deletes a static route
func (n *NetworkClient) DeleteStaticRoute(ctx context.Context, id int) error {
	return n.client.Call(ctx, "staticroute.delete", []any{id}, nil)
}

// Utility Methods

// GetInterfaceChoices returns available interface choices
func (n *NetworkClient) GetInterfaceChoices(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	err := n.client.Call(ctx, "interface.choices", []any{}, &result)
	return result, err
}

// HasPendingChanges checks if there are pending network changes
func (n *NetworkClient) HasPendingChanges(ctx context.Context) (bool, error) {
	var result bool
	err := n.client.Call(ctx, "interface.has_pending_changes", []any{}, &result)
	return result, err
}

// CommitPendingChanges commits pending network changes
func (n *NetworkClient) CommitPendingChanges(ctx context.Context, rollback bool) error {
	params := []any{}
	if rollback {
		params = append(params, map[string]any{"rollback": true})
	}
	return n.client.Call(ctx, "interface.commit", params, nil)
}

// CheckinWaiting checks if network changes are waiting for checkin
func (n *NetworkClient) CheckinWaiting(ctx context.Context) (bool, error) {
	var result bool
	err := n.client.Call(ctx, "interface.checkin_waiting", []any{}, &result)
	return result, err
}

// Checkin confirms network changes
func (n *NetworkClient) Checkin(ctx context.Context) error {
	return n.client.Call(ctx, "interface.checkin", []any{}, nil)
}

// RollbackPendingChanges rolls back pending network changes
func (n *NetworkClient) RollbackPendingChanges(ctx context.Context) error {
	return n.client.Call(ctx, "interface.rollback", []any{}, nil)
}
