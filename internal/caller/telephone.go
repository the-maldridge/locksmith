package caller

import (
	"fmt"
	"net/http"
)

var (
	infoHeader = "Authorization: Bearer %s"
)

// Type telephone is a struct that handles all the communication between
// the Caller Client and locksmith.
type Telephone struct {
	locksmithServer string
	token           string
	remoteAddr      string
}

// This method constructs and returns a new Telephone.
func NewTelephone(locksmithLoc, newToken, newAddr string) Telephone {
	newTelephone := Telephone{
		locksmithServer: locksmithLoc,
		token:           newToken,
		remoteAddr:      newAddr,
	}
	return newTelephone
}

// This method returns the remote address for the Telephone.
func (t *Telephone) GetAddress() string {
	return t.remoteAddr
}

// This method sets the remote address for the Telephone.
func (t *Telephone) SetAddress(newAddr string) {
	t.remoteAddr = newAddr
}

// This method returns the token for the telephone.
func (t *Telephone) getToken() string {
	return t.token
}

// This method sends a request to locksmith for network information.
// TODO: This method is WIP
func (t *Telephone) GetNetworkInfo() (*http.Response, error) {
	req, err := http.NewRequest("GET", t.locksmithServer, nil)
	if err != nil {
		return &http.Response{}, err
	}
	authHeader := fmt.Sprintf(infoHeader, t.getToken())
	req.Header.Set("Authorization", authHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	defer resp.Body.Close()

	return resp, nil
}

// This method sends a request to locksmith to add a peer.
func (t *Telephone) AddPeer() {

}
