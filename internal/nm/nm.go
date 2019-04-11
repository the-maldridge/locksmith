package nm

import (
	"log"
	"strings"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/models"
)

func init() {
	viper.SetDefault("nm.store.impl", "json")
}

// New returns an initialized instance ready to go.
func New() (NetworkManager, error) {
	nm := NetworkManager{}

	s, err := InitializeStore(viper.GetString("nm.store.impl"))
	if err != nil {
		return NetworkManager{}, err
	}
	nm.s = s

	nm.networks = parseNetworkConfig()

	for i := range nm.networks {
		nm.networks[i].StagedPeers = make(map[string]models.Peer)
		nm.networks[i].ApprovedPeers = make(map[string]models.Peer)
		nm.networks[i].ActivePeers = make(map[string]models.Peer)
	}
	nm.networks = nm.loadPeers(nm.networks)

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
func (nm *NetworkManager) AttemptNetworkRegistration(netID string, client models.Peer) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	for i := range net.PreApproveHooks {
		if err := nm.RunPreApproveHook(net.PreApproveHooks[i], netID, client); err != nil {
			return err
		}
	}

	if err := nm.stagePeer(netID, client); err != nil {
		return err
	}

	return nm.s.PutNetwork(net)
}

// stagePeer takes a pre-approved peer and stages them.  If the
// ApproveMode is set to AUTO then the peer is not staged and is
// instead added directly to ApprovedPeers.
func (nm *NetworkManager) stagePeer(netID string, client models.Peer) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	net.StagedPeers[client.PubKey] = client
	log.Printf("Network '%s' has staged peer '%s'",
		net.Name,
		client.PubKey)

	if err := nm.s.PutNetwork(net); err != nil {
		return err
	}

	if strings.ToUpper(net.ApproveMode) == "AUTO" {
		log.Printf("Network '%s' is automatically approving peer '%s'",
			net.Name,
			client.PubKey)
		return nm.ApprovePeer(netID, client.PubKey)
	}
	return nil
}

// ApprovePeer looks for a pubkey in StagedPeers and puts it into
// ApprovedPeers.
func (nm *NetworkManager) ApprovePeer(netID, pubkey string) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	peer, ok := net.StagedPeers[pubkey]
	if !ok {
		return ErrUnknownPeer
	}

	net.ApprovedPeers[pubkey] = peer
	delete(net.StagedPeers, pubkey)

	if err := nm.s.PutNetwork(net); err != nil {
		return err
	}

	if strings.ToUpper(net.ActivateMode) == "AUTO" {
		log.Printf("Network '%s' is automatically activating peer '%s'",
			net.Name,
			pubkey)
		return nm.ActivatePeer(netID, pubkey)
	}
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

	peer, ok := net.ApprovedPeers[pubkey]
	if !ok {
		return ErrUnknownPeer
	}

	net.ActivePeers[peer.PubKey] = peer
	go net.Sync()
	return nm.s.PutNetwork(net)
}

func parseNetworkConfig() []Network {
	var out []Network
	if err := viper.UnmarshalKey("Network", &out); err != nil {
		log.Printf("Error loading networks: %s", err)
		return nil
	}
	return out
}

func (nm *NetworkManager) loadPeers(n []Network) []Network {
	for i := range nm.networks {
		t, err := nm.s.GetNetwork(n[i].ID)
		if err != nil {
			log.Println("Error reloading network:", n[i].ID)
			continue
		}
		n[i].StagedPeers = t.StagedPeers
		n[i].ApprovedPeers = t.ApprovedPeers
		n[i].ActivePeers = t.ActivePeers

		log.Printf("Network '%s' loaded with %d staged, %d approved, %d active",
			n[i].ID,
			len(t.StagedPeers),
			len(t.ApprovedPeers),
			len(t.ActivePeers),
		)
	}
	return n
}
