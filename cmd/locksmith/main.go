package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/http"
	"github.com/the-maldridge/locksmith/internal/nm"
	_ "github.com/the-maldridge/locksmith/internal/nm/driver/keyhole"
	_ "github.com/the-maldridge/locksmith/internal/nm/ipam/dummy"
	_ "github.com/the-maldridge/locksmith/internal/nm/state/json"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/locksmith/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Fatal error config file:", err)
	}

	log.Println("Locksmith is initializing")

	nm, err := nm.New()
	if err != nil {
		log.Fatal(err)
	}

	s, err := http.New(nm)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Serve())
}
