package handler

import (
	"context"
	"log/slog"

	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/push"

	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// SetupPush generates ephemeral VAPID keys and initialises the push
// sender. Must be called once from main before the handler starts
// accepting connections. Keys are generated at startup rather than
// loaded from config because this is a demo - production apps would
// persist keys so existing subscriptions survive restarts.
func SetupPush() {
	pub, priv, err := push.GenerateVAPIDKeys()
	if err != nil {
		slog.Warn("push notifications disabled", "error", err)
		return
	}
	pushSender = push.NewSender(push.Config{
		VAPIDPublicKey:  pub,
		VAPIDPrivateKey: priv,
		Subject:         "mailto:demo@example.com",
	})
	pushEnabled = true
}

// cleanupPushSubscription removes a session's push subscription on
// disconnect so stale entries don't accumulate. The browser-side
// subscription remains valid - it simply won't receive server-sent
// pushes until the user reconnects and re-subscribes.
func cleanupPushSubscription(sessID string) {
	subscriptionsMu.Lock()
	delete(subscriptions, sessID)
	subscriptionsMu.Unlock()
}

// PushConfig returns the tether.PushConfig wired to our shared sender
// and subscription map, or nil if VAPID key generation failed.
// Returning nil tells Tether to disable push - the handler won't
// register the service worker or expose the subscribe endpoint.
func PushConfig() *tether.PushConfig[state.State] {
	if !pushEnabled {
		return nil
	}
	return &tether.PushConfig[state.State]{
		Sender: pushSender,
		OnSubscribe: func(ctx context.Context, sess *tether.StatefulSession[state.State], sub push.Subscription) {
			slog.Info("push subscription", "id", sess.ID())
			subscriptionsMu.Lock()
			subscriptions[sess.ID()] = sub
			subscriptionsMu.Unlock()
		},
	}
}
