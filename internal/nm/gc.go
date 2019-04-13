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
	for i := range nm.networks {
		for key, expiration := range nm.networks[i].ActivationExpirations {
			if time.Now().After(expiration) {
				if err := nm.DeactivatePeer(nm.networks[i].ID, key); err != nil {
					log.Println("Error deactivating key:", err)
				}
				delete(nm.networks[i].ActivationExpirations, key)
			}
		}
		for key, expiration := range nm.networks[i].ApprovalExpirations {
			if time.Now().After(expiration) {
				if err := nm.DisapprovePeer(nm.networks[i].ID, key); err != nil {
					log.Println("Error dissaproving peer:", err)
				}
				delete(nm.networks[i].ApprovalExpirations, key)
			}
		}
		nm.s.PutNetwork(nm.networks[i])
	}
}
