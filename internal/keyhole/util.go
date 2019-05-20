package keyhole

import (
	"log"
	"net"

	"github.com/deckarep/golang-set"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// generateConfig does as its name implies.  What it does is get the
// current configuration for an interface, then figure out what steps
// are needed to turn that configuration into the desired one.  As
// built this has a data race, wherein the state that's used to
// compute the desired state can become stale.  Given that this race
// lasts only a few milliseconds and it would requite a mutex or
// atomics to remove, it will stay for the time being.  This function
// is also unnecessarily complex as it generates the configuration.
// Its a good candidate to be broken up.
func (k *Keyhole) generateConfig(ic InterfaceConfig) (*wgtypes.Config, error) {
	baseCfg, err := k.wg.Device(ic.Name)
	if err != nil {
		return nil, err
	}

	// Get sets of the keys
	activePeers := mapset.NewSet()
	for _, p := range baseCfg.Peers {
		activePeers.Add(p.PublicKey.String())
	}
	wantPeers := mapset.NewSet()
	for _, p := range ic.Peers {
		wantPeers.Add(p.Pubkey)
	}

	delPeers := activePeers.Difference(wantPeers)
	newPeers := wantPeers.Difference(activePeers)
	chkPeers := activePeers.Intersect(wantPeers)

	update := []wgtypes.PeerConfig{}
	for p := range delPeers.Iterator().C {
		s, ok := p.(string)
		if !ok {
			log.Println("Not a string from a string set!", p)
			continue
		}

		k, err := wgtypes.ParseKey(s)
		if err != nil {
			log.Println("Unparsable key that was on adapter:", p)
			continue
		}
		update = append(update, wgtypes.PeerConfig{
			PublicKey: k,
			Remove:    true,
		})
	}

	for p := range newPeers.Iterator().C {
		s, ok := p.(string)
		if !ok {
			log.Println("Not a key!", p)
			continue
		}

		for _, peer := range ic.Peers {
			if peer.Pubkey == s {
				wgpeer := wgtypes.PeerConfig{}

				// Stick the key in
				k, err := wgtypes.ParseKey(peer.Pubkey)
				if err != nil {
					log.Println("Unparsable key in update:", peer.Pubkey)
				}
				wgpeer.PublicKey = k

				for _, cidrStr := range peer.AllowedIPs {
					_, cidr, err := net.ParseCIDR(cidrStr)
					if err != nil {
						log.Println("Unparsable cidr in peer:", cidrStr)
						continue
					}
					wgpeer.AllowedIPs = append(wgpeer.AllowedIPs, *cidr)
				}

				// Add the new peer
				update = append(update, wgpeer)
				break
			}
		}
	}

	log.Printf("Interface '%s': Delete %d Add %d Check %d",
		ic.Name,
		delPeers.Cardinality(),
		newPeers.Cardinality(),
		chkPeers.Cardinality())

	return &wgtypes.Config{Peers: update}, nil
}
