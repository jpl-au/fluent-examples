package signals

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
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

var wsPresence = shared.NewPresenceCountOnly()

// NewWS creates a WebSocket handler demonstrating the full set of
// signal bindings and directives.
func NewWS(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "signals/ws",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: wsPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/signals/ws/", s.OnlineCount, RenderWS(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Signals (WebSocket)"),
					assets.Stylesheet("app.css"),
				),
				body.New(content, assets.Script("hooks.js")),
			).Lang("en")
		},

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("signals-ws: connected", "id", sess.ID())
			shared.TrackPresence(wsPresence, sess.ID())
			// Push initial signal values so BindShow/BindHide
			// elements display correctly before any user interaction.
			sess.Signals(map[string]any{
				"signals.counter":       0,
				"signals.panel_visible": false,
				"signals.input_locked":  false,
				"signals.favourited":    false,
				"signals.liked":         false,
				"signals.toggle_demo":   false,
				"signals.highlight":     false,
			})
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("signals-ws: disconnected", "id", sess.ID())
			shared.UntrackPresence(wsPresence, sess.ID())
		},

		Watchers: shared.Watchers[State](wsPresence,
			func(n int, s State) State {
				s.OnlineCount = n
				return s
			},
			nil,
		),
	})
}
