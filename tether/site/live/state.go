package live

import "github.com/jpl-au/fluent-examples/tether/site/shared"

// State is the per-session state for the live updates demo.
type State struct {
	// OnlineCount is the latest snapshot from the shared Presence
	// value, updated reactively via StatefulConfig.Watchers.
	OnlineCount int
	// Activity is the rolling feed of join/leave/broadcast events
	// from the shared Presence bus.
	Activity []shared.ActivityItem
	// LastBroadcast is the most recent message this session sent.
	LastBroadcast string
	// Announcement is the last message pushed to all sessions via
	// Group.Broadcast.
	Announcement string
	// Notification is the last message received from another
	// session via Group.BroadcastOthers.
	Notification string
	// SessionIDs is a snapshot of connected session IDs from
	// Group.All(), refreshed on demand.
	SessionIDs []string
}
