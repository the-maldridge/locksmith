package state

import (
	"github.com/the-maldridge/locksmith/internal/models"
)

// Store implements the components needed to store network state and
// get them back after a system reload.
type Store interface {
	GetState(string) (models.NetState, error)
	PutState(string, models.NetState) error
}

// Factory returns a store ready to use.
type Factory func() (Store, error)

var (
	stores map[string]Factory
)

func init() {
	stores = make(map[string]Factory)
}

// Register registers a store to the list for possible initialization.
func Register(name string, f Factory) {
	if _, ok := stores[name]; ok {
		// Already registered
		return
	}
	stores[name] = f
}

// Initialize initializes the named store.
func Initialize(name string) (Store, error) {
	s, ok := stores[name]
	if !ok {
		return nil, ErrUnknownStore
	}
	return s()
}
