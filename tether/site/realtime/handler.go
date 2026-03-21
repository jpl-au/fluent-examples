package realtime

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

// State holds rolling metric samples for the system monitor charts.
// Each slice grows up to maxDataPoints entries, trimmed on each tick.
type State struct {
	HeapMB      []float64
	Goroutines  []int
	CPUPercent  []float64
	OnlineCount int
}

var realtimePresence = shared.NewPresenceCountOnly()

// New creates a WebSocket handler for the real-time dashboard.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "realtime",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: realtimePresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/realtime/", s.OnlineCount, Render(s))
		},
		Handle: func(_ tether.Session, s State, _ tether.Event) State { return s },

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Real-time Dashboard"),
					assets.Stylesheet("app.css"),
					assets.Script("echarts.min.js"),
				),
				body.New(content, assets.Script("hooks.js")),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](realtimePresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("realtime: connected", "id", sess.ID())
			shared.TrackPresence(realtimePresence, sess.ID())
			startMonitor(sess)
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("realtime: disconnected", "id", sess.ID())
			shared.UntrackPresence(realtimePresence, sess.ID())
		},
	})
}
