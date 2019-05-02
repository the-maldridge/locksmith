package dummy

import (
	"net"

	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm/addresser"
)

// The Addresser returns the same address for all peers.
type Addresser struct{}

func init() {
	addresser.Register("dummy", newAddresser)
}

func newAddresser() (addresser.Addresser, error) {
	return &Addresser{}, nil
}

// Assign always assigns 1.1.1.1
func (*Addresser) Assign(models.NetState, models.Peer) (net.IP, error) {
	return net.ParseIP("1.1.1.1"), nil
}

// Release releases the specified address back into the pool.
func (*Addresser) Release(models.Peer) error {
	return nil
}
