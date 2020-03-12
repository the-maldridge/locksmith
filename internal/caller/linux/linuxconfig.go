package caller

import (
	"io/ioutil"
	"os"
)

const (
	installPath = "/etc/wireguard/"
)

// LinuxConfig represents the tunnel configuration for a linux machine.
type LinuxConfig struct {
	DeviceName string
	Config     []byte
}

// NewConfig generates and returns a new LinuxConfig.
func NewConfig() *LinuxConfig {
	return new(LinuxConfig)
}

// InstallConfig installs the wireguard tunnel to the machine.
func (l *LinuxConfig) InstallConfig() error {
	filename := installPath + l.DeviceName

	// Check the existence of the wireguard directory
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		os.Mkdir(installPath, 0700)
	}

	// Check the existence of a wireguard conf
	if _, err := os.Stat(filename); os.IsExist(err) {
		if e := os.Remove(filename); e != nil {
			return e
		}
	}

	// Write file
	if err := ioutil.WriteFile(filename, l.Config, 0600); err != nil {
		return err
	}

	return nil
}

// UninstallConfig removes the Wireguard tunnel from the machine.
func (l *LinuxConfig) UninstallConfig() error {
	filename := installPath + l.DeviceName

	// See if a conf is installed
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	}

	// Remove the unused conf file
	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}

// SetDeviceName sets the name of the LinuxConfig
func (l *LinuxConfig) SetDeviceName(newName string) {
	l.DeviceName = newName
}

// SetConfig receives an InterfaceConfig and stores it in the LinuxConfig.
func (l *LinuxConfig) SetConfig(c []byte) {
	l.Config = c
}
