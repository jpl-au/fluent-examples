package live

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

// sseGroup tracks all connected SSE sessions.
var sseGroup = tether.NewGroup[State]()

// ssePresence tracks online sessions and broadcasts activity events
// for the SSE variant.
var ssePresence = shared.NewPresence()

// NewSSE creates an SSE handler demonstrating a subset of live updates
// features: uptime ticker, activity feed, online count, and broadcast.
func NewSSE(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:     "live/sse",
		Mode:     mode.ServerSentEvents,
		Fallback: sse.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: ssePresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/live/sse/", s.OnlineCount, RenderSSE(s))
		},
		Handle: SSEHandle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Live Updates (SSE)"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("live-sse: connected", "id", sess.ID())
			shared.TrackPresence(ssePresence, sess.ID())
			sess.Signal("online_count", ssePresence.OnlineCount.Load())
			shared.StartUptimeTicker(sess)
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("live-sse: disconnected", "id", sess.ID())
			shared.UntrackPresence(ssePresence, sess.ID())
		},

		Groups: []*tether.Group[State]{sseGroup},
		Watchers: shared.Watchers(ssePresence,
			func(n int, s State) State {
				s.OnlineCount = n
				return s
			},
			func(item shared.ActivityItem, s State) State {
				s.Activity = append([]shared.ActivityItem{item}, s.Activity...)
				if len(s.Activity) > ssePresence.MaxActivity {
					s.Activity = s.Activity[:ssePresence.MaxActivity]
				}
				return s
			},
		),
	})
}

// SSEHandle processes events on the SSE live updates page.
func SSEHandle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "live.broadcast":
		msg := fmt.Sprintf("broadcast at %s", time.Now().Format("15:04:05"))
		ssePresence.PublishActivity(sess.ID(), msg)
		s.LastBroadcast = msg
	}
	return s
}
