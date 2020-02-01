package linearv4

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm/ipam"
)

// LinearV4 allocates IPv4 addresses linearly.  It does so based on a
// simple map of addresses that are in use.
type LinearV4 struct {
	info     ipam.NetInfo
	reserved []string

	fwd map[string]string
	rev map[string]string

	ip net.IP

	savepath string
}

func init() {
	ipam.Register("linearv4", newAddresser)
}

func newAddresser(pool string) (ipam.IPAM, error) {
	cfg := viper.Sub("nm.ipam.linearv4." + pool)

	ip, n, _ := net.ParseCIDR(cfg.GetString("cidr"))

	dns := make([]net.IP, len(cfg.GetStringSlice("dns")))
	for i, s := range cfg.GetStringSlice("dns") {
		dns[i] = net.ParseIP(s)
	}

	x := &LinearV4{}
	x.info = ipam.NetInfo{
		Network: *n,
		Search:  cfg.GetStringSlice("Search"),
		DNS:     dns,

		// On IPv4 the network is always going to issue a /32
		// to a client.  This doesn't permit you to hide a
		// network behind it (without NAT) but it does make it
		// really easy to have roaming nodes.
		ClientMask: net.CIDRMask(32, 32),
	}
	x.reserved = cfg.GetStringSlice("reserved")

	x.fwd = make(map[string]string)
	x.rev = make(map[string]string)

	// Reserve some addresses
	for _, a := range x.reserved {
		x.fwd[a] = "RESERVED"
	}
	x.fwd[ip.String()] = "RESERVED"

	x.ip = ip

	x.savepath = filepath.Join(viper.GetString("core.home"), "linearv4", pool+".json")
	x.loadState()

	return x, nil
}

// NetInfo provides static information about the network.
func (l *LinearV4) NetInfo() ipam.NetInfo {
	return l.info
}

// Assign always assigns 192.168.0.2
func (l *LinearV4) Assign(p models.Peer) (net.IP, error) {
	ip := l.ip
	for ip := ip.Mask(l.info.Network.Mask); l.info.Network.Contains(ip); inc(ip) {
		if _, used := l.fwd[ip.String()]; used {
			continue
		}
		// address isn't assigned yet, assign and return.
		l.fwd[ip.String()] = p.PubKey
		l.rev[p.PubKey] = ip.String()
		l.putState()
		return ip, nil
	}
	return nil, errors.New("no addresses are available")
}

// Release releases the specified address back into the pool.
func (l *LinearV4) Release(p models.Peer) (net.IP, error) {
	addr := l.rev[p.PubKey]
	delete(l.rev, p.PubKey)
	delete(l.fwd, addr)
	l.putState()
	return net.ParseIP(addr), nil
}

func (l *LinearV4) putState() {
	state := struct {
		Fwd map[string]string
		Rev map[string]string
	}{
		Fwd: l.fwd,
		Rev: l.rev,
	}

	blob, _ := json.Marshal(state)
	ioutil.WriteFile(l.savepath, blob, 0640)
}

func (l *LinearV4) loadState() {
	state := struct {
		Fwd map[string]string
		Rev map[string]string
	}{
		Fwd: make(map[string]string),
		Rev: make(map[string]string),
	}

	in, err := ioutil.ReadFile(l.savepath)
	if err != nil {
		return
	}
	json.Unmarshal(in, &state)
	l.fwd = state.Fwd
	l.rev = state.Rev
}

// borrowed from https://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
