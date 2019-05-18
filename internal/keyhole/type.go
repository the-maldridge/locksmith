package keyhole

import (
	"golang.zx2c4.com/wireguard/wgctrl"
)

// Keyhole is a convenience type to bind server components to.
type Keyhole struct {
	wg *wgctrl.Client
	devices []string
}

// InterfaceInfo fully describes an interface in the information that
// can be publicly shared.
type InterfaceInfo struct {
	Name        string
	PublicKey   string
	ActivePeers []string
}
