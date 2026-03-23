package shared

import "time"

// ActivityItem represents one entry in the activity feed. Sessions
// subscribe to activity events via tether.On to see join, leave,
// and broadcast messages from other sessions.
type ActivityItem struct {
	// ID is the originating session's unique identifier.
	ID string
	// User is the display name derived from the session ID.
	User string
	// Action describes what happened (e.g. "joined", "left").
	Action string
	// Timestamp records when the activity occurred.
	Timestamp time.Time
}
