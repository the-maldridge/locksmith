package nm

import (
	"errors"
)

var (
	// ErrUnknownHook is returned if a hook is requested by
	// configuration and isn't actually known or installed in the
	// system.
	ErrUnknownHook = errors.New("No hook with that name is known")

	// ErrUnknownNetwork is returns if a network with an unknown
	// ID is requested.
	ErrUnknownNetwork = errors.New("No network with that ID exists")
)
