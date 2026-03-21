package freeze

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
	"github.com/jpl-au/fluent-examples/tether/store"
)

// State is the per-session state. Count is the only field - simple
// enough to verify round-trip serialisation through freeze/thaw.
type State struct {
	Count       int
	OnlineCount int
}

var freezePresence = shared.NewPresenceCountOnly()

// New creates a handler with FreezeOnDisconnect enabled. On
// disconnect the session is frozen (state persisted, memory
// released). On reconnect the counter is restored from the store.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	sessionStore := store.NewFileSessionStore("tmp/freeze-sessions")
	diffStore := store.NewFileDiffStore("tmp/freeze-diffs")

	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "freeze",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: freezePresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/freeze/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Freeze Demo"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		FreezeOnDisconnect: true,
		SessionStore:       sessionStore,
		DiffStore:          diffStore,

		Watchers: shared.Watchers[State](freezePresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("freeze: connected", "id", sess.ID())
			shared.TrackPresence(freezePresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("freeze: disconnected", "id", sess.ID())
			shared.UntrackPresence(freezePresence, sess.ID())
		},
	})
}

// Handle processes counter events.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "freeze.increment":
		s.Count++
	case "freeze.decrement":
		s.Count--
	}
	return s
}
