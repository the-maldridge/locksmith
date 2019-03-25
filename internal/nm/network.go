package nm

import (
	"log"
)

// Sync ensures that the network interface is synchronized with the
// state requested.
func (n *Network) Sync() error {
	log.Println("Synchronizing...")
	log.Println("Net:", n.ID)
	log.Println("Interface:", n.Interface)
	log.Println("Keys:")
	for k := range n.ActivePeers {
		log.Printf("  %s\n", k)
	}
	return nil
}
