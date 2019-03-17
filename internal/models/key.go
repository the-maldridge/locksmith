package models

// Client is an entity of some description that is requesting a key be
// added to the server.
type Client struct {
	// Owner is the owner associated with this key.
	Owner string

	// PubKey represents the key that is being requested to be
	// installed.
	PubKey string
}
