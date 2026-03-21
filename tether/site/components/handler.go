package components

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

	"github.com/jpl-au/fluent-examples/tether/component/counter"
	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// State is the per-session state for the components demo.
type State struct {
	// Likes is a component-managed counter wired via StatefulConfig.Components.
	Likes counter.Counter
	// Stars is a second independent counter demonstrating multiple
	// component instances on one page.
	Stars counter.Counter
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
}

var compPresence = shared.NewPresenceCountOnly()

// New creates a handler demonstrating tether.Component with two
// independent counter instances wired via StatefulConfig.Components. The
// framework dispatches prefixed events directly to each component's
// Handle method - the page-level Handle is never involved.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "components",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{
				Likes:       counter.New("likes"),
				Stars:       counter.New("stars"),
				OnlineCount: compPresence.OnlineCount.Load(),
			}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/components/", s.OnlineCount, Render(s))
		},
		Handle: func(_ tether.Session, s State, _ tether.Event) State { return s },

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Components"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Components: []tether.ComponentMount[State]{
			tether.Mount("likes",
				func(s State) counter.Counter { return s.Likes },
				func(s State, c counter.Counter) State { s.Likes = c; return s },
			),
			tether.Mount("stars",
				func(s State) counter.Counter { return s.Stars },
				func(s State, c counter.Counter) State { s.Stars = c; return s },
			),
		},

		Watchers: shared.Watchers[State](compPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("components: connected", "id", sess.ID())
			shared.TrackPresence(compPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("components: disconnected", "id", sess.ID())
			shared.UntrackPresence(compPresence, sess.ID())
		},
	})
}
