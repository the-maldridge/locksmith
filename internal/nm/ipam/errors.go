package ipam

import (
	"errors"
)

var (
	// ErrUnkownAddresser is returned if the reqeusted addresser
	// isn't known to the system.
	ErrUnknownAddresser = errors.New("No addresser with that name is known")
)
