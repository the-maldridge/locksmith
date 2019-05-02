package nm

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("nm.expiry.interval", "5m")
}

// expirationTimer is meant to be launched as a goroutine that won't
// ever return and handles expiration events.
func (nm *NetworkManager) expirationTimer() {
	log.Printf("Launching expiration timer with an interval of %s",
		viper.GetDuration("nm.expiry.interval"))

	ticker := time.NewTicker(viper.GetDuration("nm.expiry.interval"))
	for range ticker.C {
		nm.ProcessExpirations()
	}
}

// ProcessExpirations handles expiration times that have passed and
// moves keys around as necessary.
func (nm *NetworkManager) ProcessExpirations() {
	for _, net := range nm.networks {
		if net.ActivateExpiry > 0 {
			nm.doActivationExpirations(net.ID)
		}
		if net.ApproveExpiry > 0 {
			nm.doApprovalExpirations(net.ID)
		}
	}
}

func (nm *NetworkManager) doActivationExpirations(id string) error {
	net, err := nm.GetNet(id)
	if err != nil {
		return err
	}
	for key, expiration := range net.ActivationExpirations {
		if time.Now().After(expiration) {
			if err := nm.DeactivatePeer(id, key); err != nil {
				log.Println("Error deactivating key:", err)
			}
			delete(net.ActivationExpirations, key)
		}
	}
	return nil
}

func (nm *NetworkManager) doApprovalExpirations(id string) error {
	net, err := nm.GetNet(id)
	if err != nil {
		return err
	}
	for key, expiration := range net.ApprovalExpirations {
		if time.Now().After(expiration) {
			if err := nm.DisapprovePeer(id, key); err != nil {
				log.Println("Error dissaproving peer:", err)
			}
			delete(net.ApprovalExpirations, key)
		}
	}
	return nil
}
