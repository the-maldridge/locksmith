package nm

import (
	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm/ipam"
	"github.com/the-maldridge/locksmith/internal/nm/state"
)

// NetworkManager manages all the networks that are currently setup
// and live.
type NetworkManager struct {
	state.Store

	networks        []models.NetConfig
	preApproveHooks []PreApproveHook
	ipam            map[string]ipam.IPAM
}

// Network represents a network from the configuration.
type Network struct {
	models.NetConfig
	models.NetState
}

// PreApproveHook represents a hook that gets called during the
// attempt to register the client at all.  These are designed to abort
// early if for example the owner isn't known to the system.
type PreApproveHook func(string, models.Peer) error
