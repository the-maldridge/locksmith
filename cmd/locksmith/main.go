package main

import (
	"log"
	"time"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/http"
	_ "github.com/the-maldridge/locksmith/internal/http/auth/dummy"

	"github.com/the-maldridge/locksmith/internal/nm"
	_ "github.com/the-maldridge/locksmith/internal/nm/driver/keyhole"
	_ "github.com/the-maldridge/locksmith/internal/nm/ipam/dummy"
	_ "github.com/the-maldridge/locksmith/internal/nm/ipam/linearv4"
	_ "github.com/the-maldridge/locksmith/internal/nm/state/json"
)

func init() {
	viper.SetDefault("http.token.lifetime", time.Hour*12)
	viper.SetDefault("core.home", ".")
}

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
