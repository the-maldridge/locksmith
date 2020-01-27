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
		wnet.AddressTable[peerAddr.String()] = peer.PubKey
		log.Printf("Peer '%s' has been issued address '%s' on net '%s'", peer.PubKey, ip, wnet.ID)
	}
}

func (nm *NetworkManager) deconfigurePeer(wnet *Network, peer models.Peer) {
	// Remove address
	for _, h := range wnet.IPAM {
		addr, err := nm.ipam[h].Release(peer)
		if err != nil {
			log.Printf("Error releasing address on net '%s': '%s'", wnet.ID, err)
			continue
		}
		log.Printf("Address '%s' on net '%s' has been released", addr, wnet.ID)
	}

	// Dump from the address table
	for k, v := range wnet.AddressTable {
		if v == peer.PubKey {
			delete(wnet.AddressTable, k)
		}
	}

	wnet.StagedPeers[peer.PubKey] = peer
	delete(wnet.ApprovedPeers, peer.PubKey)
	delete(wnet.ApprovalExpirations, peer.PubKey)
}
