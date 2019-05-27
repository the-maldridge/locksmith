package dummy

import (
	"net"

	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm/ipam"
)

// The Dummy implementation returns the same address for all peers.
type Dummy struct{}

func init() {
	ipam.Register("dummy", newAddresser)
}

func newAddresser() (ipam.IPAM, error) {
	return &Dummy{}, nil
}

// NetInfo provides static information about the network.
func (*Dummy) NetInfo() ipam.NetInfo {
	_, n, _ := net.ParseCIDR("192.168.0.0/24")
	return ipam.NetInfo{
		Network:    *n,
		DNS:        []net.IP{net.ParseIP("192.168.0.1")},
		Search:     []string{".local"},
		ClientMask: net.CIDRMask(32, 32),
	}
}

// Assign always assigns 192.168.0.2
func (*Dummy) Assign(models.NetState, models.Peer) (net.IP, error) {
	return net.ParseIP("192.168.0.2"), nil
}

// Release releases the specified address back into the pool.
func (*Dummy) Release(models.Peer) (net.IP, error) {
	return net.ParseIP("192.168.0.2"), nil
}
