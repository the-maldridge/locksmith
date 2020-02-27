package caller

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"howett.net/plist"
)

var (
	tempFile             = "./config_temp.mobileConfig"
	filename             = "./bh-wireguard.mobileConfig"
	identifierTemplate   = "com.%s.wireguard.locksmith"
	badEncoding          = "&#x[1-9A-Z];"
	plIdentifierTemplate = "com.%s.wireguard.%s-tunnel.tunnel"
	plUUIDTemplate       = "%s: Generated UUID"
	tunnelTemplate       = `[Interface]
PrivateKey = %s
Address = 10.10.1.%s/24

[Peer]
PublicKey = %s
Endpoint = %s
AllowedIps = 0.0.0.0/0`
)

// This type represents a mobileConfig vendor configuration.
type vendorConfig struct {
	WgQuickConfig string
}

// This type represents a mobileconfig VPN configuration.
type vpnConfig struct {
	RemoteAddress        string
	AuthenticationMethod string
}

// This type represents a mobileConfig payload.
type payloadContent struct {
	PayloadDisplayName string
	PayloadType        string
	PayloadVersion     int
	PayloadIdentifier  string
	PayloadUUID        string
	UserDefinedName    string
	VPNType            string
	VPNSubType         string
	VendorConfig       vendorConfig
	VPN                vpnConfig
}

// This type represents a WireGuard tunnel mobile configuration.
type WGConfig struct {
	PayloadDisplayName string
	PayloadType        string
	PayloadVersion     int
	PayloadIdentifier  string
	PayloadUUID        string
	PayloadContent     payloadContent
}

// This function creates a new WireGuard mobile configuration.
func NewConfig(pubKey, privKey, addr, endpoint, remoteAddr, company string) (
	WGConfig, []byte, string, error) {
	tunnel := fmt.Sprintf(tunnelTemplate, privKey, addr, pubKey, endpoint)

	defaultVPNConfig := vpnConfig{
		RemoteAddress:        remoteAddr,
		AuthenticationMethod: "Password",
	}
	defaultConfig := vendorConfig{
		WgQuickConfig: tunnel,
	}

	identifier := fmt.Sprintf(identifierTemplate, company)
	defaultPayload := payloadContent{
		PayloadDisplayName: "VPN",
		PayloadType:        "com.apple.vpn.managed",
		PayloadVersion:     1,
		PayloadIdentifier:  identifier,
		PayloadUUID:        "Generated UUID",
		UserDefinedName:    "Main Tunnel",
		VPNType:            "VPN",
		VPNSubType:         "com.wireguard.macos",
		VendorConfig:       defaultConfig,
		VPN:                defaultVPNConfig,
	}

	configIdentifier := fmt.Sprintf(plIdentifierTemplate, company, company)
	configUUID := fmt.Sprintf(plUUIDTemplate, company)
	configTemplate := WGConfig{
		PayloadDisplayName: "Tunnel",
		PayloadType:        "Configuration",
		PayloadVersion:     1,
		PayloadIdentifier:  configIdentifier,
		PayloadUUID:        configUUID,
		PayloadContent:     defaultPayload,
	}

	wgTunnelConfig, err := plist.Marshal(configTemplate, plist.XMLFormat)
	if err != nil {
		return configTemplate, []byte(""), " ", err
	}

	WriteMobileConfig(wgTunnelConfig)
	removeBadCharacters(tempFile, filename)
	configTemplate.RemoveLocalFiles()

	return configTemplate, wgTunnelConfig, filename, nil
}

// This method removes the local files written during the config installation
// process.
func (wg *WGConfig) RemoveLocalFiles() {
	_ = os.Remove(tempFile)
	_ = os.Remove(filename)
}

// This method returns the identifier for the mobileconfig.
func (wg *WGConfig) GetIdentifier() string {
	return wg.PayloadIdentifier
}

// This function writes a new configuration file.
func WriteMobileConfig(mobileConfig []uint8) (string, error) {

	_ = os.Remove(tempFile)

	f, err := os.OpenFile(tempFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
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

// This function removes any newline characters that rendered poorly
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
	nf, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
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
