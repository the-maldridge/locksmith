package models

// Profile represents a locksmith configuration profile
type Profile struct {
	name       string
	privateKey string
	PublicKey  string
	serverAddr string
	tunnel     InterfaceConfig
}

// NewProfile returns a new uninitialized Profile
func NewProfile() *Profile {
	return new(Profile)
}

// GetName returns the Profile name as a string
func (p *Profile) GetName() string {
	return p.name
}

// SetName receives a string and stores it as the name for the Profile
func (p *Profile) SetName(n string) {
	p.name = n
}

// GetPubkey returns the PublicKey member string
func (p *Profile) GetPubkey() string {
	return p.PublicKey
}

// SetPubkey receives a string and sets it as the PublicKey of the Profile
func (p *Profile) SetPubkey(k string) {
	p.PublicKey = k
}

// SetPrivateKey receives a string and sets it as the privateKey for the
// Profile
func (p *Profile) SetPrivateKey(k string) {
    p.privateKey = k
}

// SetAddr receives a string and stores it as the locksmith server address
// for the Profile
func (p *Profile) SetAddr(a string) {
	p.serverAddr = a
}

// SetInterfaceConfig receives an InterfaceConfig and stores it in the profile
func (p *Profile) SetInterfaceConfig(i InterfaceConfig) {
	p.tunnel = i
}
