package linearv4

import (
	"log"
	"net"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm/ipam"
)

// LinearV4 allocates IPv4 addresses linearly.  It does so based on a
// simple map of addresses that are in use.
type LinearV4 struct {
	info     ipam.NetInfo
	reserved []string
}

func init() {
	ipam.Register("linearv4", newAddresser)
}

func newAddresser(pool string) (ipam.IPAM, error) {
	cfg := viper.Sub("nm.ipam.linearv4." + pool)

	_, n, _ := net.ParseCIDR(cfg.GetString("cidr"))

	dns := make([]net.IP, len(cfg.GetStringSlice("dns")))
	for i, s := range cfg.GetStringSlice("dns") {
		dns[i] = net.ParseIP(s)
	}

	x := &LinearV4{}
	x.info = ipam.NetInfo{
		Network:    *n,
		Search:     cfg.GetStringSlice("Search"),
		DNS:        dns,

		// On IPv4 the network is always going to issue a /32
		// to a client.  This doesn't permit you to hide a
		// network behind it (without NAT) but it does make it
		// really easy to have roaming nodes.
		ClientMask: net.CIDRMask(32, 32),
	}
	x.reserved = cfg.GetStringSlice("reserved")

	log.Println(x)
	return x, nil
}

// NetInfo provides static information about the network.
func (l *LinearV4) NetInfo() ipam.NetInfo {
	return l.info
}

// Assign always assigns 192.168.0.2
func (*LinearV4) Assign(models.Peer) (net.IP, error) {
	return net.ParseIP("192.168.0.2"), nil
}

// Release releases the specified address back into the pool.
func (*LinearV4) Release(models.Peer) (net.IP, error) {
	return net.ParseIP("192.168.0.2"), nil
}
