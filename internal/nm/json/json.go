package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/nm"
)

func init() {
	nm.RegisterStore("json", new)
}

// Store implements the nm.Store interface for keeping a persistent
// record of networks.
type Store struct {
	root string
}

func new() (nm.Store, error) {
	x := Store{}
	x.root = filepath.Join(viper.GetString("core.home"), "netstore")

	if err := os.MkdirAll(x.root, 0755); err != nil {
		log.Println(err)
		return &Store{}, nm.ErrInternalError
	}
	return &x, nil
}

// GetNetwork fetches a network and returns it to the caller.
func (s *Store) GetNetwork(id string) (nm.Network, error) {
	in, err := ioutil.ReadFile(filepath.Join(s.root, fmt.Sprintf("%s.json", id)))
	if err != nil {
		if os.IsNotExist(err) {
			return nm.Network{}, nm.ErrUnknownNetwork
		}
	}

	net := nm.Network{}
	if err := json.Unmarshal(in, &net); err != nil {
		log.Println(err)
		return nm.Network{}, nm.ErrInternalError
	}

	return net, nil
}

// PutNetwork stores a network for later retrieval
func (s *Store) PutNetwork(n nm.Network) error {
	blob, err := json.Marshal(n)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(s.root, fmt.Sprintf("%s.json", n.ID)), blob, 0640)
}
