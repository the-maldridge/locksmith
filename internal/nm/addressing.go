package nm

import (
	"net"

	"github.com/the-maldridge/locksmith/internal/models"
)

// An Addresser assigns addresses.  AssignAddress will return exactly
// one specified address and will not modify the supplied peer.  It is
// up to the network to actually assign the address to the peer.  It
// is also up to the peer to the network to release addresses from a
// peer prior to notifying the Addresser that they have been released.
// The intent here is that each network may have potentially many
// addressers operating in v4 and v6 mode, and the addressers may not
// be process-local, instead consulting an external IPAM solution to
// derive the addresses for the given peer and network.
type Addresser interface {
	AssignAddress(Network, models.Peer) (net.IP, error)
	ReleaseAddress(models.Peer) error
}

// An AddresserFactory returns a ready to use addresser
type AddresserFactory func() (Addresser, error)

var (
	addressers map[string]AddresserFactory
)

func init() {
	addressers = make(map[string]AddresserFactory)
}

// RegisterAddresser registers the addresser factory into the list
// that can be recalled later.
func RegisterAddresser(name string, f AddresserFactory) {
	if _, ok := addressers[name]; ok {
		// Already registered
		return
	}
	addressers[name] = f
}

// InitializeAddresser safely attempts to initialize an addresser or
// returns an ErrUnknownAddresser trying.
func InitializeAddresser(name string) (Addresser, error) {
	a, ok := addressers[name]
	if !ok {
		return nil, ErrUnknownAddresser
	}
	return a()
}
