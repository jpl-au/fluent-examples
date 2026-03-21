package notifications

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// State is the per-session state for the notifications demo.
type State struct {
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
}

var notifyPresence = shared.NewPresenceCountOnly()

// New creates a handler demonstrating server-initiated notifications.
// Requires a persistent connection because the server pushes Toast,
// Flash, Announce, and Signal to the client at any time.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "notifications",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: notifyPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/notifications/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Notifications"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](notifyPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("notifications: connected", "id", sess.ID())
			shared.TrackPresence(notifyPresence, sess.ID())
			sess.Signals(map[string]any{
				"notify.saved":     false,
				"notify.loading":   false,
				"notify.announced": false,
			})
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("notifications: disconnected", "id", sess.ID())
			shared.UntrackPresence(notifyPresence, sess.ID())
		},
	})
}

// simulatedDelay is used by several demos to make server-side work
// visible in the UI. In a real application the delay would come from
// database queries, API calls, or other I/O. We use an artificial
// pause so the loading spinners and indicators are clearly visible
// when demonstrating the feature.
const simulatedDelay = 3 * time.Second

// Handle processes events on the notifications page, firing the
// appropriate side effect for each demo button.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "notify.toast":
		sess.Toast("This is a toast notification!")

	case "notify.flash":
		sess.Flash("#flash-target", "Flashed!")

	case "notify.announce":
		text := "New item added to your feed"
		sess.Announce(text)
		// Echo the announcement visually so sighted users can verify
		// it was sent. In production you would not normally do this  -
		// Announce is meant for screen readers only.
		sess.Signal("notify.announced", "Screen readers heard: \""+text+"\"")
		sess.Toast("Announcement sent to ARIA live region")

	case "notify.flash-compare":
		sess.Flash("#flash-compare-target", "Saved!")

	case "notify.signal-flash":
		sess.Signal("notify.saved", true)
		// Auto-hide after the same duration as Flash so the two
		// approaches behave identically from the user's perspective.
		sess.Go(func(_ context.Context) {
			time.Sleep(5 * time.Second)
			sess.Signal("notify.saved", false)
		})

	case "notify.indicator":
		sess.Go(func(_ context.Context) {
			time.Sleep(simulatedDelay)
			sess.Toast("bind.Indicator - loading complete")
		})

	case "notify.signal-indicator":
		sess.Go(func(_ context.Context) {
			time.Sleep(simulatedDelay)
			sess.Signal("notify.loading", false)
			sess.Toast("bind.Optimistic - loading complete")
		})
	}
	return s
}
