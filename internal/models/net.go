package models

import (
	"time"
)

// NetConfig represents the static configuration data that's needed to
// fully define a network's configured state.  Notably, there is no
// information on peers in here.
type NetConfig struct {
	Name           string
	ID             string
	Interface      string
	Driver         string
	ApproveMode    string
	ApproveExpiry  time.Duration
	ActivateMode   string
	ActivateExpiry time.Duration

	PreApproveHooks []string
	IPAM            []string

	DNS        []string
	AllowedIPs []string
}

// NetState holds all the runtime state needed to compute anything for
// a network.  This allows a seperation between static configuration
// data and dynamic state data.
type NetState struct {
	ApprovalExpirations   map[string]time.Time
	ActivationExpirations map[string]time.Time

	AddressTable map[string]map[string]struct{}

	StagedPeers   map[string]Peer
	ApprovedPeers map[string]Peer
	ActivePeers   map[string]Peer
}

// Initialize is used to initialize the maps of the structure.
func (n *NetState) Initialize() {
	if n.ApprovalExpirations == nil {
		n.ApprovalExpirations = make(map[string]time.Time)
	}
	if n.ActivationExpirations == nil {
		n.ActivationExpirations = make(map[string]time.Time)
	}

	if n.AddressTable == nil {
		n.AddressTable = make(map[string]map[string]struct{})
	}

	if n.StagedPeers == nil {
		n.StagedPeers = make(map[string]Peer)
	}
	if n.ApprovedPeers == nil {
		n.ApprovedPeers = make(map[string]Peer)
	}
	if n.ActivePeers == nil {
		n.ActivePeers = make(map[string]Peer)
	}
}
