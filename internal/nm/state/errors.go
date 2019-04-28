package state

import (
	"errors"
)

var (
	// ErrUnknownStore is returned if the requested store isn't
	// known to the system.
	ErrUnknownStore = errors.New("No store with that name is known")
)
