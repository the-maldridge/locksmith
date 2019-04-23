package dummy

import (
	"net"

	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm"
)

// The Addresser returns the same address for all peers.
type Addresser struct{}

func init() {
	nm.RegisterAddresser("dummy", newAddresser)
}

func newAddresser() (nm.Addresser, error) {
	return &Addresser{}, nil
}

// AssignAddress always assigns
func (*Addresser) AssignAddress(nm.Network, models.Peer) (net.IP, error) {
	return net.ParseIP("1.1.1.1"), nil
}

// ReleaseAddress releases the specified address back into the pool.
func (*Addresser) ReleaseAddress(models.Peer) error {
	return nil
}
