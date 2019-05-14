package driver

import "errors"

var (
	// ErrUnknownDriver is returned if a driver is requested that
	// is not currently registered to the system.
	ErrUnknownDriver = errors.New("The specified driver is unknown")
)
