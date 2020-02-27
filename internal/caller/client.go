package caller

import (
	"fmt"
	"os/exec"
	"runtime"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	CLIENT_NAME        = "Caller"
	LINUX_INSTALL_PATH = "/etc/wireguard"
	PROFILES           = "/usr/bin/profiles"
)

// Type WgClient represents the application client.
type Client struct {
	status            LocksmithStatus
	localDevice       wgtypes.Device
	telephone         Telephone
	configuration     WGConfig
	confInstallStatus bool
	operatingSystem   string
	remoteAddr        string
	company           string
	configData        []byte
}

// This function returns a new WireGuard client.
func New() Client {

	// Initialize Device
	newDevice := wgtypes.Device{
		Name: CLIENT_NAME,
		Type: wgtypes.Userspace,
	}

	// Initialize Client
	newClient := Client{
		localDevice:       newDevice,
		confInstallStatus: false,
		status:            Preapproved,
		operatingSystem:   runtime.GOOS,
		company:           "",
	}
	newClient.generateKeys()
	newClient.UpdateConfiguration()

	return newClient
}

// This method updates the WGConfig for the Client using the currently
// stored keys.
func (c *Client) UpdateConfiguration() error {
	configFile, data, _, err := NewConfig(c.GetPublicKey().String(),
		c.getPrivateKey().String(), "1", "1", c.telephone.GetAddress(),
		c.GetCompany())
	if err != nil {
		return err
	}
	c.SetConfiguration(configFile)
	c.setData(data)
	return nil
}

// This method sets the Telephone for the Client.
func (c *Client) SetAddress(newAddr string) {
	c.remoteAddr = newAddr
}

// This method returns the company name of the Client.
func (c *Client) GetAddress() string {
	return c.remoteAddr
}


// This method sets the Telephone for the Client.
func (c *Client) SetTelephone(newPhone Telephone) {
	c.telephone = newPhone
}

// This method returns the company name of the Client.
func (c *Client) GetTelephone() Telephone {
	return c.telephone
}

// This method sets the company name for the Client.
func (c *Client) SetCompany(newCompany string) {
	c.company = newCompany
}

// This method returns the company name of the Client.
func (c *Client) GetCompany() string {
	return c.company
}

// This method sets the configuration data for the Client.
func (c *Client) setData(newData []byte) {
	c.configData = newData
}

// This method returns the configuration data for the Client.
func (c *Client) GetData() []byte {
	return c.configData
}

// This method sets the configuration for the Client.
func (c *Client) SetConfiguration(newConfig WGConfig) {
	c.configuration = newConfig
}

// This method returns the operating system the client is running on.
func (c *Client) GetOS() string {
	return c.operatingSystem
}

// This method returns the PublicKey for the Client.
func (c *Client) GetPublicKey() wgtypes.Key {
	return c.localDevice.PublicKey
}

// This method sets the PublicKey for the Client.
func (c *Client) setPublicKey(pubkey wgtypes.Key) {
	c.localDevice.PublicKey = pubkey
}

// This method returns the PrivateKey for the Client.
func (c *Client) getPrivateKey() wgtypes.Key {
	return c.localDevice.PrivateKey
}

// This method sets the PrivateKey for the Client.
func (c *Client) setPrivateKey(privKey wgtypes.Key) {
	c.localDevice.PrivateKey = privKey
}

// This method generates Keys and stores them in the Client.
func (c *Client) generateKeys() error {
	newKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return ErrKeyError
	}
	c.setPrivateKey(newKey)
	c.setPublicKey(c.localDevice.PrivateKey.PublicKey())
	return nil
}

// This method returns the locksmith status of the WireGuard client.
func (c *Client) GetStatus() LocksmithStatus {
	return c.status
}

// This method updates the locksmith status of the WireGuard
// client.
func (c *Client) SetStatus(newStatus LocksmithStatus) {
	c.status = newStatus
}

// This method returns a bool describing whether the WireGuard tunnel is
// installed.
func (c *Client) IsInstalled() bool {
	return c.confInstallStatus
}

// This method sets a boolean flag describing the configuration installation
// status of the Client.
func (c *Client) SetIsInstalled(newStatus bool) {
	c.confInstallStatus = newStatus
}

// This method installs the provided mobileConfig onto the localhost machine.
//TODO: Test Installation
func (c *Client) InstallMobileConfig() error {

	if c.IsInstalled() {
		return ErrConfInstalled
	}
	file, err := WriteMobileConfig(c.GetData())
	if err != nil {
		return err
	}

	switch c.GetOS() {
	case "linux":
		cmd := exec.Command("sudo", "mkdir", LINUX_INSTALL_PATH)
		err = cmd.Run()
		if err != nil {
			return err
		}
		wgConf := LINUX_INSTALL_PATH + "wg0.conf"
		cmd = exec.Command("sudo", "cp", file, wgConf)
		err = cmd.Run()
		if err != nil {
			return err
		}
		cmd = exec.Command("sudo", "wg-quick", "up", wgConf)
		err = cmd.Run()
		if err != nil {
			return err
		}
	case "darwin":
		cmd := exec.Command(PROFILES, "-I", "-F", file)
		err = cmd.Run()
		if err != nil {
			return err
		}
		c.configuration.RemoveLocalFiles()
	case "windows":
		//TODO: Windows
		fmt.Println("Work in progress.")
	}
	c.SetIsInstalled(true)
	return nil
}

// This method removes the WireGuard Client mobile config from the localhost
// machine.
//TODO: Test Removal
func (c *Client) RemoveMobileConfig() error {
	if !c.IsInstalled() {
		return nil
	}

	switch c.GetOS() {
	case "linux":
		cmd := exec.Command("sudo", "wg-quick", "down", "wg0")
		err := cmd.Run()
		if err != nil {
			return err
		}
		wgConf := LINUX_INSTALL_PATH + "wg0.conf"
		cmd = exec.Command("sudo", "rm", wgConf)
	case "darwin":
		cmd := exec.Command(PROFILES, "remove",
			c.configuration.GetIdentifier())
		err := cmd.Run()
		if err != nil {
			return err
		}
	//TODO: Make work
	case "windows":
		fmt.Println("lol")
	}
	c.SetIsInstalled(false)
	return nil
}
