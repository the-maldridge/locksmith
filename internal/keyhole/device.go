package keyhole

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"golang.zx2c4.com/wireguard/wgctrl"
)

// New returns a new keyhole which can accept keys for certain
// devices.
func New() (*Keyhole, error) {
	c, err := wgctrl.New()
	if err != nil {
		return nil, err
	}

	devs, err := c.Devices()
	if err != nil {
		return nil, err
	}
	devList := []string{}
	for _, d := range devs {
		devList = append(devList, d.Name)
	}

	return &Keyhole{c, devList}, nil
}

// DeviceNames returns the device names served by this keyhole server.
func (k *Keyhole) DeviceNames() []string {
	return k.devices
}

// DeviceInfo is an RCP method that provides information about the
// requested device.
func (k *Keyhole) DeviceInfo(name string, reply *InterfaceInfo) error {
	if !k.isKnownDevice(name) {
		return errors.New("Unknown device")
	}
	log.Println("Information requested on", name)

	dev, err := k.wg.Device(name)
	if err != nil {
		return err
	}

	reply.Name = dev.Name
	reply.PublicKey = dev.PublicKey.String()

	for i := range dev.Peers {
		reply.ActivePeers = append(reply.ActivePeers, dev.Peers[i].PublicKey.String())
	}
	return nil
}

// ConfigureDevice configues the device with the provided peerlist by
// computing the difference between the peers that are known and those
// that must be added or removed.
func (k *Keyhole) ConfigureDevice(dc InterfaceConfig, reply *string) error {
	if !k.isKnownDevice(dc.Name) {
		return errors.New("Unknown device")
	}

	cfg, err := k.generateConfig(dc)
	if err != nil {
		return err
	}
	return k.wg.ConfigureDevice(dc.Name, *cfg)
}

// Serve serves the keyhole service on a given port and bind.
func (k *Keyhole) Serve(bind string, port int) error {
	rpc.Register(k)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bind, port))
	if err != nil {
		return err
	}
	return http.Serve(l, nil)
}

func (k *Keyhole) isKnownDevice(name string) bool {
	for i := range k.devices {
		if k.devices[i] == name {
			return true
		}
	}
	return false
}
