package nm

import (
	"log"
	"net"

	"github.com/the-maldridge/locksmith/internal/models"
)

func (nm *NetworkManager) configurePeer(wnet *Network, peer models.Peer) {
	// Assign address
	for _, h := range wnet.IPAM {
		netinfo := nm.ipam[h].NetInfo()

		ip, err := nm.ipam[h].Assign(peer)
		if err != nil {
			log.Printf("Error assigning address on net '%s': '%s'", wnet.ID, err)
			continue
		}
		peerAddr := net.IPNet{
			IP:   ip,
			Mask: netinfo.ClientMask,
		}
		if wnet.AddressTable[peer.PubKey] == nil {
			wnet.AddressTable[peer.PubKey] = make(map[string]struct{})
		}

		wnet.AddressTable[peer.PubKey][peerAddr.String()] = struct{}{}
		log.Printf("Peer '%s' has been issued address '%s' on net '%s'", peer.PubKey, ip, wnet.ID)
	}
}

func (nm *NetworkManager) deconfigurePeer(wnet *Network, peer models.Peer) {
	// Remove address
	for _, h := range wnet.IPAM {
		netinfo := nm.ipam[h].NetInfo()

		addr, err := nm.ipam[h].Release(peer)
		if err != nil {
			log.Printf("Error releasing address on net '%s': '%s'", wnet.ID, err)
			continue
		}
		peerAddr := net.IPNet{
			IP:   addr,
			Mask: netinfo.ClientMask,
		}
		delete(wnet.AddressTable[peer.PubKey], peerAddr.String())
		log.Printf("Address '%s' on net '%s' has been released", addr, wnet.ID)
	}

	for addr := range wnet.AddressTable[peer.PubKey] {
		log.Printf("Address '%s'  on net '%s' leaked!", addr, wnet.ID)
	}
	if len(wnet.AddressTable[peer.PubKey]) == 0 {
		delete(wnet.AddressTable, peer.PubKey)
	}

	wnet.StagedPeers[peer.PubKey] = peer
	delete(wnet.ApprovedPeers, peer.PubKey)
	delete(wnet.ApprovalExpirations, peer.PubKey)
}

// GenerateConfigForPeer assembles the peer's half of the
// configuration data.  This is not stored persistently as the
// underlying information may change frequently.
func (nm *NetworkManager) GenerateConfigForPeer(netID string, pubkey string) (models.PeerConfig, error) {
	net, err := nm.GetNet(netID)
	if err != nil {
		return models.PeerConfig{}, err
	}

	c := models.PeerConfig{}

	c.DNS = net.DNS
	c.AllowedIPs = net.AllowedIPs

	for _, h := range net.IPAM {
		netinfo := nm.ipam[h].NetInfo()
		for _, resolver := range netinfo.DNS {
			c.DNS = append(c.DNS, resolver.String())
		}
		c.Search = append(c.Search, netinfo.Search...)
	}

	_, c.Staged = net.StagedPeers[pubkey]
	_, c.Approved = net.ApprovedPeers[pubkey]
	_, c.Active = net.ActivePeers[pubkey]

	// As a final step, fill in addresses if any are present.
	addrs, ok := net.AddressTable[pubkey]
	if !ok {
		return c, nil
	}

	for addr := range addrs {
		c.Addresses = append(c.Addresses, addr)
	}

	return c, nil
}
