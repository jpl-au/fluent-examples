package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
	"github.com/jpl-au/tether/sse"
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New creates the kanban board handler. A single handler serves all
// connected browsers; board state is shared via the store and
// synchronised across sessions with Group.Broadcast.
func New(board *store.Board, assets *tether.Asset) *tether.Handler[State] {
	group := tether.NewGroup[State]()
	viewers := newViewers()

	return tether.Stateful(tether.App{
		DevMode: true,
		Assets:  []*tether.Asset{assets},
	}, tether.StatefulConfig[State]{
		Name:     "kanban",
		Mode:     mode.Both,
		Upgrade:  wsupgrade.Upgrade(),
		Fallback: sse.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{View: "board", OnlineCount: group.Count().Load()}
		},
		Render:     Render(board, viewers),
		Handle:     Handle(board, group, viewers),
		OnNavigate: navigate(board),

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Kanban Board"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Groups: []*tether.Group[State]{group},
		Watchers: []tether.Watcher[State]{
			tether.WatchValue(group.Count(), func(n int, s State) State {
				s.OnlineCount = n
				return s
			}),
		},

		// No idle timeout - kanban boards are left open indefinitely.

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("connected", "id", sess.ID()[:8])
			sess.Update(func(s State) State {
				s.SessionID = sess.ID()
				return s
			})
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("disconnected", "id", sess.ID()[:8])
			viewers.Presence.Clear(sess.ID())
		},
	})
}

// navigate handles URL-driven state.
func navigate(board *store.Board) func(tether.Session, State, tether.Params) State {
	return func(_ tether.Session, s State, p tether.Params) State {
		path := p.Path
		if after, ok := strings.CutPrefix(path, "/card/"); ok {
			if _, ok := board.Card(after); ok {
				s.View = "detail"
				s.SelectedID = after
				return s
			}
		}
		if path == "/new" {
			s.View = "detail"
			s.SelectedID = ""
			return s
		}
		s.View = "board"
		s.SelectedID = ""
		return s
	}
}
