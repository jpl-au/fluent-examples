// Package state defines the session state types for the Service Worker
// section of the feature explorer.
package state

// State is the per-session state for the Service Worker section.
// Each browser tab gets its own copy, persistent across reconnects
// until the session is reaped.
type State struct {
	// Page is the current route path, set by OnNavigate.
	Page string
	// OnlineCount is the latest snapshot from the shared Presence
	// value, updated reactively via StatefulConfig.Watchers.
	OnlineCount int
	// Lifecycle holds the PWA lifecycle events page state.
	Lifecycle LifecycleState
}

// LifecycleState tracks PWA lifecycle events detected by the client.
type LifecycleState struct {
	// Installed is set to true when the browser fires the
	// appinstalled event, indicating the user added the PWA
	// to their home screen.
	Installed bool
}
