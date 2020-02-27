package caller

const (
	Preapproved LocksmithStatus = "preapproved"
	Staged      LocksmithStatus = "staged"
	Active      LocksmithStatus = "active"
	Approved    LocksmithStatus = "approved"
)

// A peerStatus type represents what state the Client is in with relation
// to locksmith.
type LocksmithStatus string
