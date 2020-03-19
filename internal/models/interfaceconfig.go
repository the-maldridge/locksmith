package models

import (
	"bytes"
	"text/template"
)

var (
	tmpl           = template.Must(template.New("Tunnel").Parse(tunnelTemplate))
	tunnelTemplate = `[Interface]
PrivateKey = {{.privateKey}}
Address = {{.address}}

{{range .Peers}}
[Peer]
PublicKey = {{.PublicKey}}
AllowedIPs = {{.AllowedIPs}}
Endpoint = {{.Endpoint}}
{{end}}`
)

// InterfaceConfig type represents a WireGuard tunnel configuration.
type InterfaceConfig struct {
	PrivateKey string
	Address    string
	Domain     string
	DNS        []string
	Peers      []InterfacePeer
}

// InterfacePeer type represents a WireGuard peer connection.
type InterfacePeer struct {
	Name       string
	Key        string
	Endpoint   string
	AllowedIPs []string
}

// String returns a string representation of the InterfaceConfig Tunnel.
func (i *InterfaceConfig) String() (string, error) {
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, i)
	if err != nil {
		return "Error:", err
	}
	return tpl.String(), nil
}

// NewInterfaceConfig returns a new interfaceConfig.
func NewInterfaceConfig() *InterfaceConfig {
	return new(InterfaceConfig)
}

// AddPeer receives a new Peer and stores it in the interfaceConfig.
func (i *InterfaceConfig) AddPeer(p InterfacePeer) {
	i.Peers = append(i.Peers, p)
}

// AddDNS adds another entry to the InterfaceConfig's DNS record.
func (i *InterfaceConfig) AddDNS(d string) {
	i.DNS = append(i.DNS, d)
}

// SetPrivateKey receives a string and sets the
// privateKey for the Tunnel.
func (i *InterfaceConfig) SetPrivateKey(pk string) {
	i.PrivateKey = pk
}

// SetAddress receives a string and sets the address for the Tunnel.
func (i *InterfaceConfig) SetAddress(a string) {
	i.Address = a
}

// SetDomain sets the domain of the InterfaceConfig
func (i *InterfaceConfig) SetDomain(d string) {
	i.Domain = d
}
