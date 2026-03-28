package memoise

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
	"github.com/jpl-au/tether/sse"
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

var realtimePresence = shared.NewPresenceCountOnly()

// NewRealtime creates a memoised handler for the real-time dashboard.
// Same system monitor as the realtime demo, but with Memoise: true and
// each chart region wrapped in node.Memoise with Versioned keys.
func NewRealtime(app tether.App, assets *tether.Asset) *tether.Handler[RealtimeState] {
	return tether.Stateful(app, tether.StatefulConfig[RealtimeState]{
		Name:     "memoise/realtime",
		Upgrade:  wsupgrade.Upgrade(),
		Fallback: sse.Upgrade(),
		Memoise:  true,

		InitialState: func(_ *http.Request) RealtimeState {
			return RealtimeState{OnlineCount: realtimePresence.OnlineCount.Load()}
		},
		Render: func(s RealtimeState) node.Node {
			return layout.Shell(layout.SectionLive, "/memoise/realtime/", s.OnlineCount, RenderRealtime(s))
		},
		Handle: func(_ tether.Session, s RealtimeState, _ tether.Event) RealtimeState { return s },

		Layout: func(_ RealtimeState, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Memoised Real-time Dashboard"),
					assets.Stylesheet("app.css"),
					assets.Script("echarts.min.js"),
				),
				body.New(content, assets.Script("hooks.js")),
			).Lang("en")
		},

		Watchers: shared.Watchers[RealtimeState](realtimePresence,
			func(n int, s RealtimeState) RealtimeState { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[RealtimeState]) {
			slog.Info("memoise/realtime: connected", "id", sess.ID())
			shared.TrackPresence(realtimePresence, sess.ID())
			startMonitor(sess)
		},
		OnDisconnect: func(sess *tether.StatefulSession[RealtimeState]) {
			slog.Info("memoise/realtime: disconnected", "id", sess.ID())
			shared.UntrackPresence(realtimePresence, sess.ID())
		},
	})
}
