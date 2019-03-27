package nm

// Store implements the components needed to store networks and get
// them back after a system reload.
type Store interface {
	GetNetwork(string) (Network, error)
	PutNetwork(Network) error
}

// StoreFactory returns a store ready to use.
type StoreFactory func() (Store, error)

var (
	stores map[string]StoreFactory
)

func init() {
	stores = make(map[string]StoreFactory)
}

// RegisterStore registers a store to the list for possible initialization.
func RegisterStore(name string, f StoreFactory) {
	if _, ok := stores[name]; ok {
		// Already registered
		return
	}
	stores[name] = f
}

// InitializeStore initializes the named store.
func InitializeStore(name string) (Store, error) {
	s, ok := stores[name]
	if !ok {
		return nil, ErrUnknownStore
	}
	return s()
}
