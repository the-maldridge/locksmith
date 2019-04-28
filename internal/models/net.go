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
	ApproveMode    string
	ApproveExpiry  time.Duration
	ActivateMode   string
	ActivateExpiry time.Duration

	PreApproveHooks []string
	AddrHandlers    []string
}

// NetState holds all the runtime state needed to compute anything for
// a network.  This allows a seperation between static configuration
// data and dynamic state data.
type NetState struct {
	ApprovalExpirations   map[string]time.Time
	ActivationExpirations map[string]time.Time

	AddressTable map[string]Peer

	StagedPeers   map[string]Peer
	ApprovedPeers map[string]Peer
	ActivePeers   map[string]Peer
}
