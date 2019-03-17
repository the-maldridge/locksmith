package nm

import (
	"log"

	"github.com/the-maldridge/locksmith/internal/models"
)

var (
	preApproveHooks map[string]PreApproveHook
)

func init() {
	preApproveHooks = make(map[string]PreApproveHook)
}

// RegisterPreApproveHook adds a hook for pre-approval to the list
// available for the system to use later on.
func RegisterPreApproveHook(name string, h PreApproveHook) {
	if _, ok := preApproveHooks[name]; ok {
		// Already registered
		return
	}
	log.Println("PreApprove hook", name, "is now registered.")
	preApproveHooks[name] = h
}

// RunPreApproveHook runs the named hook and figures out of the client
// is pre-approved.
func (nm *NetworkManager) RunPreApproveHook(hook, net string, client models.Client) error {
	h, ok := preApproveHooks[hook]
	if !ok {
		return ErrUnknownHook
	}
	return h(net, client)
}
