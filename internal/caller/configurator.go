package caller

// Configurator represents a wireguard configuration implementation for a
// platform.
type Configurator interface {
	InstallConfig() error
	UninstallConfig() error
}
