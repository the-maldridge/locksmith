package main

import (
	"log"

	"github.com/the-maldridge/locksmith/internal/keyhole"
)

func main() {
	c, err := keyhole.NewClient("localhost:1234")
	if err != nil {
		log.Fatal(err)
	}

	reply, err := c.InterfaceInfo("wg0")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(reply)

	cfg := keyhole.InterfaceConfig{
		Name: "wg0",
		Peers: []keyhole.Peer{
			keyhole.Peer{
				Pubkey: "F3sSyEZ/VSHVurBN3oAAL+Vt5+6/zlnbeiUJwy4kbwU=",
				AllowedIPs: []string{
					"192.168.0.2/32",
				},
			},
			keyhole.Peer{
				Pubkey: "F3sSyEZ/VSHVurBN3oAAL+Vt4+6/zlnbeiUJwy4kbwU=",
				AllowedIPs: []string{
					"192.168.0.3/32",
				},
			},
		},
	}

	rep, err := c.ConfigureDevice(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rep)
}
