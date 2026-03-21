package chat

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

	"github.com/jpl-au/fluent-examples/tether/component/shoutbox"
	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// State is the per-session state for the chat demo.
type State struct {
	// Shoutbox is a cross-session chat component wired via
	// StatefulConfig.Components. Bus subscription delivers messages from
	// other sessions via WatchBus.
	Shoutbox shoutbox.Shoutbox
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
}

var chatPresence = shared.NewPresenceCountOnly()

// New creates a handler demonstrating a real-time chat room. The
// Shoutbox component handles all rendering and event dispatch; a
// WatchBus watcher delivers messages from other sessions.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "chat",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: chatPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/chat/", s.OnlineCount, Render(s))
		},
		Handle: func(_ tether.Session, s State, _ tether.Event) State { return s },

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Chat"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Components: []tether.ComponentMount[State]{
			tether.Mount("shoutbox",
				func(s State) shoutbox.Shoutbox { return s.Shoutbox },
				func(s State, c shoutbox.Shoutbox) State { s.Shoutbox = c; return s },
			),
		},

		Watchers: append(shared.Watchers[State](chatPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),
			tether.WatchBus(shoutbox.Bus, func(m shoutbox.Shout, s State) State {
				s.Shoutbox.Messages = append([]shoutbox.Shout{m}, s.Shoutbox.Messages...)
				if len(s.Shoutbox.Messages) > 50 {
					s.Shoutbox.Messages = s.Shoutbox.Messages[:50]
				}
				return s
			}),
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("chat: connected", "id", sess.ID())
			shared.TrackPresence(chatPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("chat: disconnected", "id", sess.ID())
			shared.UntrackPresence(chatPresence, sess.ID())
		},
	})
}
