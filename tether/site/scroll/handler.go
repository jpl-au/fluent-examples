package scroll

import (
	"fmt"
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

// State holds per-session state for the scroll demo.
type State struct {
	OnlineCount int
	Items       int // number of items in the preserve-scroll list
}

var scrollPresence = shared.NewPresenceCountOnly()

// New creates a stateful handler demonstrating scroll features.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "scroll",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: scrollPresence.OnlineCount.Load(), Items: 30}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/scroll/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Scroll"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](scrollPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("scroll: connected", "id", sess.ID())
			shared.TrackPresence(scrollPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("scroll: disconnected", "id", sess.ID())
			shared.UntrackPresence(scrollPresence, sess.ID())
		},
	})
}

// Handle processes scroll demo events.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "scroll.add":
		s.Items += 5
		sess.ScrollTo(fmt.Sprintf("#item-%d", s.Items))
	case "scroll.server-scroll":
		sess.ScrollTo("#scroll-target")
	}
	return s
}
