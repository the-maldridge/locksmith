package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/the-maldridge/locksmith/internal/models"
	"github.com/the-maldridge/locksmith/internal/nm"
	"github.com/the-maldridge/locksmith/internal/nm/state"
)

func init() {
	state.Register("json", new)
}

// Store implements the nm.Store interface for keeping a persistent
// record of networks.
type Store struct {
	root string
}

func new() (state.Store, error) {
	x := Store{}
	x.root = filepath.Join(viper.GetString("core.home"), "netstore")

	if err := os.MkdirAll(x.root, 0755); err != nil {
		log.Println(err)
		return &Store{}, nm.ErrInternalError
	}
	return &x, nil
}

// GetState fetches network state and returns it to the caller.
func (s *Store) GetState(id string) (models.NetState, error) {
	in, err := ioutil.ReadFile(filepath.Join(s.root, fmt.Sprintf("%s.json", id)))
	if err != nil {
		if os.IsNotExist(err) {
			nstate := models.NetState{}
			nstate.Initialize()
			return nstate, nil
		}
		return models.NetState{}, nm.ErrInternalError
	}

	state := models.NetState{}
	if err := json.Unmarshal(in, &state); err != nil {
		log.Println(err)
		return models.NetState{}, nm.ErrInternalError
	}
	return state, nil
}

// PutState stores network state for later retrieval
func (s *Store) PutState(id string, st models.NetState) error {
	blob, err := json.Marshal(st)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(s.root, fmt.Sprintf("%s.json", id)), blob, 0640)
}
