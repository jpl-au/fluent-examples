package live

import (
	"context"
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
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// wsGroup tracks all connected WebSocket sessions.
var wsGroup = func() *tether.Group[State] {
	g := tether.NewGroup[State]()
	g.OnJoin = func(sess *tether.StatefulSession[State]) {
		slog.Info("live-ws group: session joined", "id", sess.ID()[:6])
	}
	g.OnLeave = func(sess *tether.StatefulSession[State]) {
		slog.Info("live-ws group: session left", "id", sess.ID()[:6])
	}
	return g
}()

// wsPresence tracks online sessions and broadcasts activity events
// for the WebSocket variant.
var wsPresence = shared.NewPresence()

// NewWS creates a WebSocket handler demonstrating the full set of
// live updates features: uptime ticker, activity feed, online count,
// Group operations, SetTitle, State() in Go(), and Close().
func NewWS(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "live/ws",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: wsPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/live/ws/", s.OnlineCount, RenderWS(s))
		},
		Handle: WSHandle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Live Updates (WebSocket)"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("live-ws: connected", "id", sess.ID())
			shared.TrackPresence(wsPresence, sess.ID())
			sess.Signal("online_count", wsPresence.OnlineCount.Load())
			shared.StartUptimeTicker(sess)
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("live-ws: disconnected", "id", sess.ID())
			shared.UntrackPresence(wsPresence, sess.ID())
		},

		Groups: []*tether.Group[State]{wsGroup},
		Watchers: shared.Watchers(wsPresence,
			func(n int, s State) State {
				s.OnlineCount = n
				return s
			},
			func(item shared.ActivityItem, s State) State {
				s.Activity = append([]shared.ActivityItem{item}, s.Activity...)
				if len(s.Activity) > wsPresence.MaxActivity {
					s.Activity = s.Activity[:wsPresence.MaxActivity]
				}
				return s
			},
		),
	})
}

// WSHandle processes events on the WS live updates page.
func WSHandle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "live.broadcast":
		msg := fmt.Sprintf("broadcast at %s", time.Now().Format("15:04:05"))
		wsPresence.PublishActivity(sess.ID(), msg)
		s.LastBroadcast = msg

	case "live.announce":
		msg := fmt.Sprintf("Announced at %s by User %s", time.Now().Format("15:04:05"), sess.ID()[:6])
		wsGroup.Broadcast(func(_ *tether.StatefulSession[State], s State) State {
			s.Announcement = msg
			return s
		})

	case "live.notify-others":
		msg := fmt.Sprintf("Notification from User %s at %s", sess.ID()[:6], time.Now().Format("15:04:05"))
		wsGroup.BroadcastOthers(sess, func(_ *tether.StatefulSession[State], s State) State {
			s.Notification = msg
			return s
		})
		s.Notification = "You sent the notification - other sessions received it."

	case "live.list-sessions":
		var ids []string
		for member := range wsGroup.All() {
			ids = append(ids, member.ID()[:8])
		}
		s.SessionIDs = ids

	case "live.set-title":
		sess.SetTitle(fmt.Sprintf("tether - Title set at %s", time.Now().Format("15:04:05")))

	case "live.read-state":
		if live, ok := sess.(*tether.StatefulSession[State]); ok {
			live.Go(func(ctx context.Context) {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
				}
				current := live.State()
				announcement := current.Announcement
				if announcement == "" {
					announcement = "(none)"
				}
				live.Toast(fmt.Sprintf("State read at %s - announcement: %s", time.Now().Format("15:04:05"), announcement))
			})
		}

	case "live.close":
		sess.Close()
	}
	return s
}
