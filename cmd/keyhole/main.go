package main

import (
	"log"

	"github.com/the-maldridge/locksmith/internal/keyhole"
)

func main() {
	log.Println("Keyhole is initializing")

	kh, err := keyhole.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Managing:")
	for _, d := range kh.DeviceNames() {
		log.Printf("  %s", d)
	}

	kh.Serve("", 1234)
}
