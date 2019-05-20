package keyhole

import (
	"golang.zx2c4.com/wireguard/wgctrl"
)

// Keyhole is a convenience type to bind server components to.
type Keyhole struct {
	wg      *wgctrl.Client
	devices []string
}

// Client is a convenience wrapper on top of RPC to make access a bit
// more straightforward inside the network manager.
type Client struct {
	server string
}

// InterfaceInfo fully describes an interface in the information that
// can be publicly shared.
type InterfaceInfo struct {
	Name        string
	PublicKey   string
	ActivePeers []string
}

// InterfaceConfig contains the information needed by keyhole to
// configure a single interface.
type InterfaceConfig struct {
	Name  string
	Peers []Peer
}

// A Peer represents the minimal information needed by keyhole to
// compute configure a peer for an interface.
type Peer struct {
	Pubkey     string
	AllowedIPs []string
}
