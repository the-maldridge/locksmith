package nm

import (
	"log"
	"strings"
	"time"

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

	useExpiry := false
	for i := range nm.networks {
		nm.networks[i].StagedPeers = make(map[string]models.Peer)
		nm.networks[i].ApprovedPeers = make(map[string]models.Peer)
		nm.networks[i].ActivePeers = make(map[string]models.Peer)
		nm.networks[i].ApprovalExpirations = make(map[string]time.Time)
		nm.networks[i].ActivationExpirations = make(map[string]time.Time)

		if nm.networks[i].ApproveExpiry != 0 || nm.networks[i].ActivateExpiry != 0 {
			useExpiry = true
		}

		// Load Addressers for this network
		for j := range nm.networks[i].AddrHandlers {
			a, err := InitializeAddresser(nm.networks[i].AddrHandlers[j])
			if err != nil {
				log.Printf("Network '%s' error during AddrHandler initialization: '%s'", err)
				continue
			}
			log.Printf("Network '%s' using '%s' addresser", nm.networks[i].ID, nm.networks[i].AddrHandlers[j])
			nm.networks[i].Addressers = append(nm.networks[i].Addressers, a)
		}
	}
	nm.networks = nm.loadPeers(nm.networks)

	// Launch the expiration timer
	if useExpiry {
		nm.ProcessExpirations()
		go nm.expirationTimer()
	}

	return nm, nil
}

// GetNet returns the network if known or an error if not known.
func (nm *NetworkManager) GetNet(id string) (*Network, error) {
	for i := range nm.networks {
		if nm.networks[i].ID == id {
			return &nm.networks[i], nil
		}
	}
	return &Network{}, ErrUnknownNetwork
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

	return nm.s.PutNetwork(*net)
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

	if err := nm.s.PutNetwork(*net); err != nil {
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

	if net.ApproveExpiry != 0 {
		// Approvals expire for this network
		net.ApprovalExpirations[pubkey] = time.Now().Add(net.ApproveExpiry)
	}

	if err := nm.s.PutNetwork(*net); err != nil {
		return err
	}

	log.Printf("Network '%s' has approved peer '%s'",
		net.Name,
		pubkey)

	if strings.ToUpper(net.ActivateMode) == "AUTO" {
		log.Printf("Network '%s' is automatically activating peer '%s'",
			net.Name,
			pubkey)
		return nm.ActivatePeer(netID, pubkey)
	}
	return nil
}

// DisapprovePeer removes a peer from the approval set and deactivates
// them as well, just in case.
func (nm *NetworkManager) DisapprovePeer(netID, pubkey string) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	peer, ok := net.ApprovedPeers[pubkey]
	if !ok {
		return ErrUnknownPeer
	}

	net.StagedPeers[pubkey] = peer
	delete(net.ApprovedPeers, pubkey)
	delete(net.ApprovalExpirations, pubkey)

	if err := nm.s.PutNetwork(*net); err != nil {
		return err
	}

	log.Printf("Network '%s' has disapproved peer '%s'",
		net.Name,
		pubkey)
	return nm.DeactivatePeer(netID, pubkey)
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

	if net.ActivateExpiry != 0 {
		// Activation expiry is active for this network.
		net.ActivationExpirations[pubkey] = time.Now().Add(net.ActivateExpiry)
	}

	net.ActivePeers[peer.PubKey] = peer
	go net.Sync()
	log.Printf("Network '%s' has activated peer '%s'",
		net.Name,
		pubkey)
	return nm.s.PutNetwork(*net)
}

// DeactivatePeer is used to immediately remove a peer from the active
// set.
func (nm *NetworkManager) DeactivatePeer(netID, pubkey string) error {
	net, err := nm.GetNet(netID)
	if err != nil {
		return err
	}

	pri := len(net.ActivePeers)
	delete(net.ActivePeers, pubkey)
	delete(net.ActivationExpirations, pubkey)
	if pri != len(net.ActivePeers) {
		go net.Sync()
	}
	log.Printf("Network '%s' has deactivated peer '%s'",
		net.Name,
		pubkey)
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

		n[i].ApprovalExpirations = t.ApprovalExpirations
		n[i].ActivationExpirations = t.ActivationExpirations

		log.Printf("Network '%s' loaded with %d staged, %d approved, %d active",
			n[i].ID,
			len(t.StagedPeers),
			len(t.ApprovedPeers),
			len(t.ActivePeers),
		)
	}
	return n
}
