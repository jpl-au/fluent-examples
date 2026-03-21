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
	"github.com/jpl-au/tether/sse"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

var ssePresence = shared.NewPresenceCountOnly()

// NewSSE creates an SSE handler demonstrating a subset of signal
// bindings over Server-Sent Events: BindText, BindShow/BindHide,
// SetSignal, ToggleSignal, and Optimistic.
func NewSSE(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:     "signals/sse",
		Mode:     mode.ServerSentEvents,
		Fallback: sse.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: ssePresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/signals/sse/", s.OnlineCount, RenderSSE(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Signals (SSE)"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("signals-sse: connected", "id", sess.ID())
			shared.TrackPresence(ssePresence, sess.ID())
			sess.Signals(map[string]any{
				"signals.counter":       0,
				"signals.panel_visible": false,
				"signals.liked":         false,
				"signals.toggle_demo":   false,
			})
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("signals-sse: disconnected", "id", sess.ID())
			shared.UntrackPresence(ssePresence, sess.ID())
		},

		Watchers: shared.Watchers[State](ssePresence,
			func(n int, s State) State {
				s.OnlineCount = n
				return s
			},
			nil,
		),
	})
}
