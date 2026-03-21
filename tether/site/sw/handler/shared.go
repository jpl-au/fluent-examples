// Package handler provides session lifecycle hooks, push notification
// configuration, and shared application state for the Service Worker
// section of the feature explorer.
package handler

import (
	"sync"

	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/push"

	"github.com/jpl-au/fluent-examples/tether/site/shared"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// Group tracks all connected sessions in the SW section. Sessions
// are added and removed automatically via StatefulConfig.Groups in server.go.
var Group = tether.NewGroup[state.State]()

// Presence tracks online sessions. The SW section uses count-only
// presence - no activity bus - because push notification demos
// don't need a live activity feed.
var Presence = shared.NewPresenceCountOnly()

// Watchers returns declarative subscriptions for StatefulConfig.Watchers:
// presence (online count only, no activity feed).
func Watchers() []tether.Watcher[state.State] {
	return shared.Watchers[state.State](Presence,
		func(n int, s state.State) state.State {
			s.OnlineCount = n
			return s
		},
		nil,
	)
}

// Push notification state. SetupPush must be called once at startup
// (from main) to generate VAPID keys and initialise pushSender.
// pushEnabled gates signal values sent to new sessions so the UI
// can hide the subscribe button when keys are unavailable.
//
// subscriptions maps session ID → push subscription. The map is
// mutex-protected because OnSubscribe runs on the session loop
// goroutine while cleanupPushSubscription runs on disconnect.
var (
	pushSender  *push.Sender
	pushEnabled bool

	subscriptions   = make(map[string]push.Subscription)
	subscriptionsMu sync.Mutex
)
