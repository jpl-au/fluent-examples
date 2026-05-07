package timer

import (
	"log/slog"
	"net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

var timerPresence = shared.NewPresenceCountOnly()

// New creates a WebSocket handler demonstrating client-side timers.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name: "timer",
		Mode: mode.WebSocket,

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: timerPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/timer/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Client-Side Timers"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("timer: connected", "id", sess.ID())
			shared.TrackPresence(timerPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("timer: disconnected", "id", sess.ID())
			shared.UntrackPresence(timerPresence, sess.ID())
		},

		Watchers: shared.Watchers[State](timerPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),
	})
}
