package nm

import (
	"log"
)

// SyncNet is used to push the desired state down to the driver.  It
// handles pulling the correct driver handle and pushing the correct
// config downwards.
func (nm *NetworkManager) SyncNet(wnet Network) error {
	req := wnet.Driver
	if req == "" {
		req = "LOCAL"
	}

	drvr, ok := nm.driver[req]
	if !ok {
		log.Printf("Failure to activate peer on '%s': missing driver '%s'", wnet.ID, req)
		return ErrInternalError
	}
	ident := wnet.ID
	if wnet.Interface != "" {
		ident = wnet.Interface
	}
	return drvr.Configure(ident, wnet.NetState)
}

// GetNet returns the network if known or an error if not known.
func (nm *NetworkManager) GetNet(id string) (Network, error) {
	net := Network{}

	for i := range nm.networks {
		if nm.networks[i].ID == id {
			net.NetConfig = nm.networks[i]
			s, err := nm.GetState(net.ID)
			if err != nil {
				return Network{}, err
			}
			net.NetState = s
			return net, nil
		}
	}
	return Network{}, ErrUnknownNetwork
}

// StoreNet is a convenience function that stores network state.
func (nm *NetworkManager) StoreNet(net Network) error {
	return nm.PutState(net.ID, net.NetState)
}
