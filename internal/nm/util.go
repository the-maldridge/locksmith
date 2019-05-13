package nm

import (
	"log"
	"net"

	"github.com/the-maldridge/locksmith/internal/models"
)

func delAddr(ips []net.IP, tgt net.IP) []net.IP {
	var out []net.IP
	for _, i := range ips {
		if i.Equal(tgt) {
			continue
		}
		out = append(out, i)
	}
	return out
}

func (nm *NetworkManager) configurePeer(wnet *Network, peer *models.Peer) {
	// Assign address
	for _, h := range wnet.IPAM {
		netinfo := nm.ipam[h].NetInfo()

		ip, err := nm.ipam[h].Assign(wnet.NetState, *peer)
		if err != nil {
			log.Printf("Error assigning address on net '%s': '%s'", wnet.ID, err)
			continue
		}
		peer.Addresses = append(peer.Addresses, ip)

		peerAddr := net.IPNet{
			IP: ip,
			Mask: netinfo.ClientMask,
		}
		wnet.AddressTable[peerAddr.String()] = peer.PubKey
		log.Printf("Peer '%s' has been issued address '%s' on net '%s'", peer.PubKey, ip, wnet.ID)

		
		// Copy DNS settings for this network from IPAM
		for _, resolver := range netinfo.DNS {
			peer.DNS = append(peer.DNS, resolver)
		}

		// Copy DNS settings from config
		for _, resolver := range wnet.DNS {
			ip := net.ParseIP(resolver)
			if ip == nil {
				continue
			}
			peer.DNS = append(peer.DNS, ip)
		}

		// Assign AllowedIPs, these are the ones that are
		// defined by hand.
		for _, allow := range wnet.AllowedIPs {
			_, cidr, err := net.ParseCIDR(allow)
			if err != nil {
				continue
			}
			peer.AllowedIPs = append(peer.AllowedIPs, *cidr)
		}
	}
}

func (nm *NetworkManager) deconfigurePeer(wnet *Network, peer *models.Peer) {
	// Remove address
	for _, h := range wnet.IPAM {
		addr, err := nm.ipam[h].Release(*peer)
		if err != nil {
			log.Printf("Error releasing address on net '%s': '%s'", wnet.ID, err)
			continue
		}
		peer.Addresses = delAddr(peer.Addresses, addr)
		log.Printf("Address '%s' on net '%s' has been released", addr, wnet.ID)
	}
	for _, ip := range peer.Addresses {
		log.Printf("Peer '%s' on net '%s' left dangling address '%s'!", peer.PubKey, wnet.ID, ip)
	}

	// Dump from the address table
	for k, v := range wnet.AddressTable {
		if v == peer.PubKey {
			delete(wnet.AddressTable, k)
		}
	}
	
	peer.Addresses = nil
	peer.DNS = nil
	peer.AllowedIPs = nil

	wnet.StagedPeers[peer.PubKey] = *peer
	delete(wnet.ApprovedPeers, peer.PubKey)
	delete(wnet.ApprovalExpirations, peer.PubKey)
}
