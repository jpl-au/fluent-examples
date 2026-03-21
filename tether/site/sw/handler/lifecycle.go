package handler

import (
	"log/slog"

	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/site/shared"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// OnConnect performs imperative setup for a newly connected session.
// Reactive subscriptions are handled declaratively via StatefulConfig.Watchers.
func OnConnect(sess *tether.StatefulSession[state.State]) {
	slog.Info("connected", "id", sess.ID())
	shared.TrackPresence(Presence, sess.ID())
	pushInitialSignals(sess)
}

// OnDisconnect tears down session state and removes any push
// subscription so stale subscriptions don't accumulate. Called when
// the transport closes permanently (not on temporary disconnects).
func OnDisconnect(sess *tether.StatefulSession[state.State]) {
	slog.Info("disconnected", "id", sess.ID())
	cleanupPushSubscription(sess.ID())
	shared.UntrackPresence(Presence, sess.ID())
}
