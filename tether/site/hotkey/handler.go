package hotkey

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

// State is the per-session state for the hotkey demo.
type State struct {
	LastCombo   string
	OnlineCount int
}

var hotkeyPresence = shared.NewPresenceCountOnly()

// New creates a stateful handler demonstrating global hotkeys.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "hotkey",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: hotkeyPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/hotkey/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Hotkeys"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](hotkeyPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("hotkey: connected", "id", sess.ID())
			shared.TrackPresence(hotkeyPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("hotkey: disconnected", "id", sess.ID())
			shared.UntrackPresence(hotkeyPresence, sess.ID())
		},
	})
}

// Handle processes hotkey events.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "hotkey.triggered":
		combo, _ := ev.Get("combo")
		s.LastCombo = combo
	}
	return s
}
