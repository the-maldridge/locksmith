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

		ip, err := nm.ipam[h].Assign(wnet.NetState, peer)
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
