package keyhole

import (
	"log"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/keyhole"
	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm/driver"
)

// Driver is an implementation of the driver.Driver interface.
type Driver struct {
	*keyhole.Client
}

func init() {
	driver.Register("KEYHOLE", new)
}

func new() (driver.Driver, error) {
	return &Driver{}, nil
}

// Configure passes off data to keyhole to configure it.
func (d *Driver) Configure(name string, state models.NetState) error {
	condensed := make(map[string][]string)
	for peer, addresses := range state.AddressTable {
		for addr := range addresses {
			condensed[peer] = append(condensed[peer], addr)
		}
	}

	peers := []keyhole.Peer{}
	for key := range state.ActivePeers {
		p := keyhole.Peer{
			Pubkey:     key,
			AllowedIPs: condensed[key],
		}
		peers = append(peers, p)
	}

	ic := keyhole.InterfaceConfig{
		Name:  name,
		Peers: peers,
	}

	for _, server := range viper.GetStringSlice("keyhole.servers") {
		c, err := keyhole.NewClient(server)
		if err != nil {
			log.Println("Error initializing keyhole:", err)
			continue
		}
		_, err = c.ConfigureDevice(ic)
		if err != nil {
			log.Println("Error configuring interface:", err)
			return err
		}
	}

	return nil
}
