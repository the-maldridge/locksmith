package models

import (
	"github.com/mdlayher/wireguardctrl/tree/master/wgtypes"
)

// Client is an entity of some description that is requesting a key be
// added to the server.
type Client struct {
	// Owner is the owner associated with this key.
	Owner string

	// PubKey represents the key that is being requested to be
	// installed.
	PubKey wgtypes.Key
}
