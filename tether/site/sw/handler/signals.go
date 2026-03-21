package handler

import (
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
	tether "github.com/jpl-au/tether"
)

// pushInitialSignals sends the current online count and push
// availability to a newly connected session. The header badge needs
// the count before any tether.Observe fires, and the push subscribe
// button must know immediately whether VAPID keys are available so
// it can show or hide itself via bind.BindShow/BindHide.
func pushInitialSignals(sess *tether.StatefulSession[state.State]) {
	sess.Signals(map[string]any{
		"online_count":     Presence.OnlineCount.Load(),
		"push.available":   pushEnabled,
		"push.unavailable": !pushEnabled,
		"sw.online":        true,
	})
}
