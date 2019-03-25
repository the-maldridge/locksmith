package nm

import (
	"log"
	"strings"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/models"
)

// New returns an initialized instance ready to go.
func New() (NetworkManager, error) {
	nm := NetworkManager{}

	nm.networks = parseNetworkConfig()

	for i := range nm.networks {
		nm.networks[i].StagedPeers = make(map[string]models.Client)
		nm.networks[i].ApprovedPeers = make(map[string]models.Client)
		nm.networks[i].ActivePeers = make(map[string]models.Client)
	}

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

	return nm.stagePeer(netID, client)
}

// stagePeer takes a pre-approved peer and stages them.  If the
// ApproveMode is set to AUTO then the peer is not staged and is
// instead added directly to ApprovedPeers.
func (nm *NetworkManager) stagePeer(netID string, client models.Client) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	if strings.ToUpper(net.ApproveMode) == "AUTO" {
		// Automatic approval, add directly to approved peers.
		net.ApprovedPeers[client.PubKey] = client
		log.Println(net)
		log.Println("Automatic approval, peer auto-approved")
		return nil
	}
	net.StagedPeers[client.PubKey] = client
	log.Println(net)
	log.Println("The peer has been staged")
	return nil
}

// ActivatePeer recalls the peer from the ApprovedPeers map and
// attempts to activate it.  If the peer is not approved an error is
// returned.
func (nm *NetworkManager) ActivatePeer(netID, pubkey string) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	log.Println(net.ApprovedPeers)
	log.Println(pubkey)

	peer, ok := net.ApprovedPeers[pubkey]
	if !ok {
		return ErrUnknownPeer
	}

	net.ActivePeers[peer.PubKey] = peer
	return net.Sync()
}

func parseNetworkConfig() []Network {
	var out []Network
	if err := viper.UnmarshalKey("Network", &out); err != nil {
		log.Printf("Error loading networks: %s", err)
		return nil
	}
	return out
}
