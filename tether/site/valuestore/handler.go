package valuestore

import (
	"log/slog"
	"net/http"
	"strconv"

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

// State is the per-session state for the value store demo.
type State struct {
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
	// SharedCount mirrors the shared counter observed via WatchValue.
	SharedCount int
	// LocalCount is a per-session counter that does not propagate
	// to other sessions.
	LocalCount int
}

// sharedCounter is a cross-session reactive integer. All connected
// sessions observe it via tether.WatchValue and see updates in real
// time whenever any session calls Store or Update.
var sharedCounter = tether.NewValue(0)

var presence = shared.NewPresenceCountOnly()

// New creates a WebSocket handler demonstrating tether.Value with
// Store (direct set) and Update (read-modify-write) methods.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "valuestore",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{
				OnlineCount: presence.OnlineCount.Load(),
				SharedCount: sharedCounter.Load(),
			}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/valuestore/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Value Store"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: append(shared.Watchers[State](presence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),
			tether.WatchValue(sharedCounter, func(n int, s State) State {
				s.SharedCount = n
				return s
			}),
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("valuestore: connected", "id", sess.ID())
			shared.TrackPresence(presence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("valuestore: disconnected", "id", sess.ID())
			shared.UntrackPresence(presence, sess.ID())
		},
	})
}

// Handle processes events for the value store demo, dispatching
// each action to the appropriate tether.Value operation.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "value.increment":
		sharedCounter.Update(func(n int) int { return n + 1 })

	case "value.decrement":
		sharedCounter.Update(func(n int) int {
			if n > 0 {
				return n - 1
			}
			return 0
		})

	case "value.reset":
		sharedCounter.Store(0)

	case "value.set":
		raw := ev.Value()
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			slog.Warn("valuestore: invalid value for set", "raw", raw, "err", err)
			return s
		}
		sharedCounter.Store(parsed)

	case "value.local-inc":
		s.LocalCount++
	}
	return s
}
