package caller

import (
	"github.com/the-maldridge/locksmith/pkg/telephone"
)

// Configurators for the Client.
var (
	NativeConfigurator Configurator
	AppConfigurator    Configurator
)

// Client represents the application client.
type Client struct {
	telephone.Telephone
	Configurator
}

// New returns a new WireGuard Client.
func New(t telephone.Telephone) Client {
	newClient := Client{Telephone: t}
	return newClient
}
