package nm

import (
	"github.com/the-maldridge/locksmith/internal/models"
)

// NetworkManager manages all the networks that are currently setup
// and live.
type NetworkManager struct {
	networks []Network

	preApproveHooks []PreApproveHook

	knownPeers map[string][]models.Client
}

// Network represents a network from the configuration.
type Network struct {
	Name      string
	ID        string
	Interface string

	PreApproveHooks []string
}

// PreApproveHook represents a hook that gets called during the
// attempt to register the client at all.  These are designed to abort
// early if for example the owner isn't known to the system.
type PreApproveHook func(string, models.Client) error
