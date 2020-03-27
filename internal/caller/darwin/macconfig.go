package caller

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"howett.net/plist"
	"github.com/satori/go.uuid"
)

var (
	tempFile             = "./config_temp.mobileConfig"
	filename             = "./wireguard.mobileConfig"
	identifierTemplate   = "com.%s.wireguard.locksmith"
	badEncoding          = "&#x[1-9A-Z];"
	profiles             = "/usr/bin/profiles"
	plIdentifierTemplate = "com.%s.wireguard.%s-tunnel.tunnel"
	plUUIDTemplate       = "%s: Generated UUID"
)

// vendorConfig represents a mobileConfig vendor configuration.
type VendorConfig struct {
	WgQuickConfig string
}

// vpnConfig represents a mobileconfig VPN configuration.
type VpnConfig struct {
	RemoteAddress        string
	AuthenticationMethod string
}

// payloadContent represents a mobileConfig payload.
type PayloadContent struct {
	PayloadDisplayName string
	PayloadType        string
	PayloadVersion     int
	PayloadIdentifier  string
	PayloadUUID        string
	UserDefinedName    string
	VPNType            string
	VPNSubType         string
	VendorConfig       VendorConfig
	VPN                VpnConfig
}

// mobileConfig represents a WireGuard tunnel mobile configuration.
type MobileConfig struct {
	PayloadDisplayName string
	PayloadType        string
	PayloadVersion     int
	PayloadIdentifier  string
	PayloadUUID        string
	PayloadContent     PayloadContent
}

// MacConfig represents a MacOS X WireGuard configuration.
type MacConfig struct {
	localDevice wgtypes.device
	Config      []byte
	data        MobileConfig
	vpn         VpnConfig
	vendor      VendorConfig
	payload     PayloadContent
	config      MobileConfig
	org         string
}

// init initializes the AppConfigurator
func init() {
	AppConfigurator = NewConfig()
}

// NewConfig creates a new WireGuard mobile configuration.
func NewConfig() *MacConfig {
	return new(MacConfig)
}


func (m *MacConfig) RecompileMobileConfig() {
	defaultVPNConfig := VpnConfig{
		RemoteAddress:        "",
		AuthenticationMethod: "Password",
	}
	defaultConfig := VendorConfig{
		WgQuickConfig: tunnelTemplate,
	}

	defaultPayload := PayloadContent{
		PayloadDisplayName: "VPN",
		PayloadType:        "com.apple.vpn.managed",
		PayloadVersion:     1,
		PayloadIdentifier:  "",
		PayloadUUID:        "Generated UUID",
		UserDefinedName:    "Main Tunnel",
		VPNType:            "VPN",
		VPNSubType:         "com.wireguard.macos",
		VendorConfig:       defaultConfig,
		vpn:                defaultVPNConfig,
	}
	m.payload = defaultPayload

	configTemplate := MobileConfig{
		PayloadDisplayName: "Tunnel",
		PayloadType:        "Configuration",
		PayloadVersion:     1,
		PayloadIdentifier:  "",
		PayloadUUID:        "",
		PayloadContent:     defaultPayload,
	}
	m.config = configTemplate
	wgTunnelConfig, err := plist.Marshal(configTemplate, plist.XMLFormat)
	if err != nil {
		return MacConfig{}
	}

	newDevice := wgtypes.Device{
		Name:       "MacOS Configuration",
		DeviceType: wgtypes.Userspace,
	}
	newMacConfig := MacConfig{
		localDevice: newDevice,
		tunnel:      wgTunnelConfig,
		data:        configTemplate,
	}
}

// SetOrg sets the name of the MacConfig organization.
func (m *MacConfig) SetOrg(orgName string) {
	m.org = orgName
}

// UpdateVpnConfig updates the vpn config for the MacConfig.
func (m *MacConfig) UpdateVpnConfig(newAddr string) {
	newVpnConfig := VpnConfig{
		RemoteAddress:        newAddr,
		AuthenticationMethod: "Password",
	}
	m.vpn = newVpnConfig
}

func (m *MacConfig) UpdatePayload() {
	newPayload := PayloadContent{
		PayloadDisplayName: "VPN",
		PayloadType:        "com.apple.vpn.managed",
		PayloadVersion:     1,
		PayloadIdentifier:  "",
		PayloadUUID:        "Generated UUID",
		UserDefinedName:    "Main Tunnel",
		VPNType:            "VPN",
		VPNSubType:         "com.wireguard.macos",
		VendorConfig:       m.vendor,
		VPN:                m.vpn,
	}
	m.payload = newPayload

}

// GenerateKeys generates a public and private key for the configuration.
func (m *MacConfig) GenerateKeys() error {
	newPrivKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return err
	}
	m.localDevice.PrivateKey = newPrivKey
	m.localDevice.PublicKey = privKey.PublicKey()
}

// UpdateTunnel updates the tunnel information for the mobile configuration.
func (m *MacConfig) UpdateTunnel() {
	newTunnel := fmt.Sprintf(tunnelTemplate,
		m.localDevice.PrivateKey,
		m.addr,
		m.localDevice.PublicKey,
		m.allowedIps,
		m.endpoint,
	)
	newVendorConfig := VendorConfig {
		WgQuickConfig: newTunnel,
	}
	m.vendor = newVendorConfig
}

// UpdateMobileConfig updates the mobileConfig plist stored in the
// MacConfig.
func (m *MacConfig) UpdateMobileConfig() {
	configIdentifier := fmt.Sprintf(plIdentifierTemplate, m.org, m.org)
	configUUID := fmt.Sprintf(plUUIDTemplate, m.org)
	configTemplate := MobileConfig{
		PayloadDisplayName: "Tunnel",
		PayloadType:        "Configuration",
		PayloadVersion:     1,
		PayloadIdentifier:  configIdentifier,
		PayloadUUID:        configUUID,
		PayloadContent:     defaultPayload,
	}
	m.config = configTemplate
}

// InstallConfig installs the tunnel configuration onto the Mac.
func (m *MacConfig) InstallConfig() error {
	WriteMobileConfig([]byte(m.tunnel))
	removeBadCharacters(tempFile, filename)
	cmd := exec.Command(profiles, "-I", "-F", filename)
	err := cmd.Run()
	if err != nil {
		return err
	}
	RemoveLocalFiles()
	return nil
}

// UninstallConfig removes the tunnel configuration from the Mac.
func (m *MacConfig) UninstallConfig() error {
	cmd := exec.Command(profiles, "remove", m.GetIdentifier())
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// WriteMobileConfig  writes a new configuration file for installation.
func WriteMobileConfig(mobileConfig []uint8) (string, error) {

	_ = os.Remove(tempFile)

	f, err := os.OpenFile(tempFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = f.Write(mobileConfig)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// removeBadCharacters removes any newline characters that rendered poorly
// in the mobileConfig.
func removeBadCharacters(inFile, outFile string) error {
	garbage := regexp.MustCompile(badEncoding)

	_ = os.Remove(outFile)

	// Open temporary file
	f, err := os.Open(inFile)
	if err != nil {
		return err
	}

	defer f.Close()

	// Open final config file
	nf, err := os.OpenFile(outFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600)
	if err != nil {
		return err
	}

	defer nf.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		newlineByte := []byte("\n")
		nextLine := []byte(scanner.Text())
		regScan := garbage.ReplaceAll(nextLine, newlineByte)
		_, err = nf.Write(regScan)

		if err != nil {
			return err
		}
	}
	_ = os.Remove(inFile)
	return nil
}

// RemoveLocalFiles removes the local files written during the config
// installation process.
func RemoveLocalFiles() {
	_ = os.Remove(tempFile)
	_ = os.Remove(filename)
}

// CreateUUID returns a string UUID or an error.
func CreateUUID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return id, nil
}
