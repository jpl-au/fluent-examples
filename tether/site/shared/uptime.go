package shared

import (
	"context"
	"time"

	tether "github.com/jpl-au/tether"
)

// StartUptimeTicker launches a background goroutine that pushes an
// incrementing counter as a signal every second. The counter is
// delivered exclusively via sess.Signal - no sess.Update call, no
// render-diff cycle. The bound element on the live page updates
// directly through bind.BindText.
//
// The goroutine is tied to the session lifetime via sess.Go and
// stops automatically when the session is destroyed.
func StartUptimeTicker[S any](sess *tether.StatefulSession[S]) {
	sess.Go(func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		var uptime int
		for {
			select {
			case <-ticker.C:
				uptime++
				sess.Signal("live.uptime", uptime)
			case <-ctx.Done():
				return
			}
		}
	})
}
