package driver

import (
	"github.com/the-maldridge/locksmith/internal/models"
)

var (
	impls map[string]Factory
)

// A Driver provides a mechanism to configure a wireguard interface
// either remotely or locally.
type Driver interface {
	Configure(string, models.NetState) error
}

// A Factory provides an initialized driver and possibly an
// initialization error.
type Factory func() (Driver, error)

func init() {
	impls = make(map[string]Factory)
}

// Register allows each driver to register itself to the system.
func Register(name string, f Factory) {
	if _, ok := impls[name]; ok {
		return
	}
	impls[name] = f
}

// Initialize provides the mechanism needed to initialize a driver for
// use.
func Initialize(name string) (Driver, error) {
	d, ok := impls[name]
	if !ok {
		return nil, ErrUnknownDriver
	}
	return d()
}
