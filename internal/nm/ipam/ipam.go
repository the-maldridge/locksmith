package ipam

import (
	"net"

	"github.com/the-maldridge/locksmith/internal/models"
)

// An IPAM assigns addresses.  AssignAddress will return exactly one
// specified address and will not modify the supplied peer.  It is up
// to the network to actually assign the address to the peer.  It is
// also up to the peer to the network to release addresses from a peer
// prior to notifying the Addresser that they have been released.  The
// intent here is that each network may have potentially many
// addressers operating in v4 and v6 mode, and the addressers may not
// be process-local, instead consulting an external IPAM solution to
// derive the addresses for the given peer and network.
type IPAM interface {
	NetInfo() NetInfo
	Assign(models.NetState, models.Peer) (net.IP, error)
	Release(models.Peer) (net.IP, error)
}

// An Factory returns a ready to use addresser
type Factory func() (IPAM, error)

// NetInfo provides information about the particular network an
// Addresser is responsible for.
type NetInfo struct {
	Network    net.IPNet
	DNS        []net.IP
	Search     []string
	ClientMask net.IPMask
}

var (
	impls map[string]Factory
)

func init() {
	impls = make(map[string]Factory)
}

// Register registers the addresser factory into the list that can be
// recalled later.
func Register(name string, f Factory) {
	if _, ok := impls[name]; ok {
		// Already registered
		return
	}
	impls[name] = f
}

// Initialize safely attempts to initialize an addresser or returns an
// ErrUnknownAddresser trying.
func Initialize(name string) (IPAM, error) {
	a, ok := impls[name]
	if !ok {
		return nil, ErrUnknownAddresser
	}
	return a()
}
