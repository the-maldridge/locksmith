package keyhole

import (
	"net/rpc"
)

// NewClient returns a client ready to go.
func NewClient(conn string) (*Client, error) {
	return &Client{conn}, nil
}

// InterfaceInfo returns information for a given interface.
func (c *Client) InterfaceInfo(device string) (InterfaceInfo, error) {
	cl, err := c.connect()
	if err != nil {
		return InterfaceInfo{}, err
	}

	var reply InterfaceInfo
	return reply, cl.Call("Keyhole.DeviceInfo", device, &reply)
}

// ConfigureDevice unsurprisingly configures a device.  The
// configuration supplied should be complete for the desired state,
// and computations will be done on the server to achieve this state.
func (c *Client) ConfigureDevice(dc InterfaceConfig) (string, error) {
	cl, err := c.connect()
	if err != nil {
		return "", err
	}

	var reply string
	return reply, cl.Call("Keyhole.ConfigureDevice", dc, &reply)
}

func (c *Client) connect() (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", c.server)
}
