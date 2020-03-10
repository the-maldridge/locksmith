package caller

import "errors"

var (
	// ErrConfUpdate is returned when the WireGuard client fails
	// to update the Configurator
	ErrConfUpdate = errors.New("failure to update mobileconfig")

	// ErrConfUninstalled is returned if an attempt is made to uninstall
	// a Configurator that is non-existent.
	ErrConfUninstalled = errors.New("configuration profile missing")

	// ErrConfInstalled is returned when an attempt is made to install a
	// Configurator when one already exists.
	ErrConfInstalled = errors.New("configuration profile already installed")

	// ErrKeyError is returned when a request is made regarding Keys that
	// fails.
	ErrKeyError = errors.New("error retrieving requested key")
)
