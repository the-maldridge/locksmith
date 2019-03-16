package main

import (
	"log"

	"github.com/the-maldridge/locksmith/internal/http"
)

func main() {
	log.Println("Keyhole is initializing")

	s, err := http.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Serve())
}
