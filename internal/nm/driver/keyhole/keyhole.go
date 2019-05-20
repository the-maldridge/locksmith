package keyhole

import (
	"github.com/the-maldridge/locksmith/internal/keyhole"
	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm/driver"
)

// Driver is an implementation of the driver.Driver interface.
type Driver struct {
	*keyhole.Client
}

func init() {
	driver.Register("keyhole", new)
}

func new() (driver.Driver, error) {
	return &Driver{}, nil
}

// Configure passes off data to keyhole to configure it.
func (d *Driver) Configure(name string, state models.NetState) error {
	return nil
}
