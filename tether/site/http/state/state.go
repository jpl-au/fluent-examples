package state

// State is reconstructed from each HTTP request - there is no
// persistent session.
type State struct {
	// Page is the current route path, set by the router.
	Page string
}
