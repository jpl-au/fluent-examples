package patch

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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

const numCounters = 20

// State holds per-session state for the patch demo.
type State struct {
	Counters    [numCounters]int
	ActiveIndex int
	OnlineCount int
}

var patchPresence = shared.NewPresenceCountOnly()

// New creates a stateful handler demonstrating targeted updates.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:     "patch",
		Upgrade:  wsupgrade.Upgrade(),
		Fallback: sse.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: patchPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/patch/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Targeted Updates"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](patchPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("patch: connected", "id", sess.ID())
			shared.TrackPresence(patchPresence, sess.ID())
			startTicker(sess)
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("patch: disconnected", "id", sess.ID())
			shared.UntrackPresence(patchPresence, sess.ID())
		},
	})
}

// startTicker launches a background goroutine that increments one
// counter every 500ms using sess.Patch. Only the targeted row is
// re-rendered - the other 19 are untouched.
func startTicker(sess *tether.StatefulSession[State]) {
	sess.Go(func(ctx context.Context) {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		idx := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				i := idx % numCounters
				key := "counter-" + strconv.Itoa(i)
				sess.Patch(key, func(s State) (State, node.Node) {
					s.Counters[i]++
					s.ActiveIndex = i
					return s, RenderCounter(i, s.Counters[i])
				})
				idx++
			}
		}
	})
}

// Handle processes patch demo events.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "patch.reset":
		s.Counters = [numCounters]int{}
	}
	return s
}
