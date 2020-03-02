package models

// Peer is an entity of some description that is requesting a key be
// added to the server.
type Peer struct {
	// Owner is the owner associated with this key.
	Owner string

	// OwnerLabel is a reference that is human readable for the
	// peer.  This label is not required to have a relation to the
	// owner, though it should generally be the name of the owner,
	// whereas the owner itself would be the username of the
	// owning user.
	OwnerLabel string

	// PubKey represents the key that is being requested to be
	// installed.
	PubKey string
}

// A PeerConfig is the configuration data that is fed back to a peer
// from the server.  This contains enough to initialize the peer with
// wg-quick.
type PeerConfig struct {
	PeerPubKey string

	Staged   bool
	Approved bool
	Active   bool

	DNS        []string
	AllowedIPs []string
	Addresses  []string
	Search     []string
}
