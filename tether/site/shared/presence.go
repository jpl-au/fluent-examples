package shared

import (
	"time"

	tether "github.com/jpl-au/tether"
)

// Presence holds the shared reactive primitives for tracking online
// sessions: an observable count and an optional activity bus for
// broadcasting join/leave events to all connected sessions.
type Presence struct {
	// OnlineCount is a reactive counter that tracks connected sessions.
	OnlineCount *tether.Value[int]
	// ActivityBus broadcasts join/leave events to all subscribers.
	// Nil when the section only needs a count (see NewPresenceCountOnly).
	ActivityBus *tether.Bus[ActivityItem]
	// MaxActivity caps the per-session activity feed length.
	MaxActivity int
}

// NewPresence creates a presence tracker with an activity bus and
// a default cap of 20 items in the activity feed.
func NewPresence() *Presence {
	return &Presence{
		OnlineCount: tether.NewValue(0),
		ActivityBus: tether.NewBus[ActivityItem](),
		MaxActivity: 20,
	}
}

// NewPresenceCountOnly creates a presence tracker without an activity
// bus - suitable for sections that only need an online count badge.
func NewPresenceCountOnly() *Presence {
	return &Presence{
		OnlineCount: tether.NewValue(0),
	}
}

// PublishActivity broadcasts an activity item to all subscribed
// sessions. The first six characters of the session ID are used as
// a short display name. No-op if the activity bus is nil.
func (p *Presence) PublishActivity(sessID, action string) {
	if p.ActivityBus == nil {
		return
	}
	p.ActivityBus.Publish(ActivityItem{
		ID: sessID, User: "User " + sessID[:6],
		Action: action, Timestamp: time.Now(),
	})
}

// Watchers returns declarative subscriptions for StatefulConfig.Watchers.
// setCount maps the current online count into the session state.
// addActivity prepends an activity item; pass nil if the section
// does not use an activity feed.
func Watchers[S any](
	p *Presence,
	setCount func(int, S) S,
	addActivity func(ActivityItem, S) S,
) []tether.Watcher[S] {
	w := []tether.Watcher[S]{
		tether.WatchValue(p.OnlineCount, func(n int, s S) S {
			return setCount(n, s)
		}),
	}
	if p.ActivityBus != nil && addActivity != nil {
		w = append(w, tether.WatchBus(p.ActivityBus, func(item ActivityItem, s S) S {
			return addActivity(item, s)
		}))
	}
	return w
}

// TrackPresence performs the imperative parts of presence tracking:
// increments the online count and publishes a "joined" event. Call
// this in OnConnect after Watchers have already subscribed the session
// to the reactive sources via StatefulConfig.Watchers.
func TrackPresence(p *Presence, sessID string) {
	p.OnlineCount.Update(func(n int) int { return n + 1 })
	p.PublishActivity(sessID, "joined")
}

// UntrackPresence decrements the online count and publishes a "left"
// event. The guard against n < 0 handles the edge case where a
// disconnect fires before the session was fully tracked.
func UntrackPresence(p *Presence, sessID string) {
	p.OnlineCount.Update(func(n int) int {
		if n > 0 {
			return n - 1
		}
		return 0
	})

	p.PublishActivity(sessID, "left")
}
