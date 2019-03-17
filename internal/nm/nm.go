package nm

import (
	"log"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/models"
)

// New returns an initialized instance ready to go.
func New() (NetworkManager, error) {
	nm := NetworkManager{}

	nm.networks = parseNetworkConfig()
	nm.knownPeers = make(map[string][]models.Client)

	return nm, nil
}

// GetNet returns the network if known or an error if not known.
func (nm *NetworkManager) GetNet(id string) (Network, error) {
	for i := range nm.networks {
		if nm.networks[i].ID == id {
			return nm.networks[i], nil
		}
	}
	return Network{}, ErrUnknownNetwork
}

// AttemptNetworkRegistration as the name implies attempts to register
// the client into the named network.  It checks the pre-approve hooks
// to make sure we don't need to reject the client right away, and
// after we're confident they're OK, then we add them to a staged
// clients list.  Staged clients are de-staged asynchronously from
// this request.
func (nm *NetworkManager) AttemptNetworkRegistration(netID string, client models.Client) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	for i := range net.PreApproveHooks {
		if err := nm.RunPreApproveHook(net.PreApproveHooks[i], netID, client); err != nil {
			return err
		}
	}
	return nil
}

func parseNetworkConfig() []Network {
	var out []Network
	if err := viper.UnmarshalKey("Network", &out); err != nil {
		log.Printf("Error loading networks: %s", err)
		return nil
	}
	return out
}
