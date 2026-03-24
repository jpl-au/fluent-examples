package handler

import (
	"log/slog"
	"net/http"
	"strings"
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
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New creates the kanban board handler. A single handler serves all
// connected browsers; board state is shared via the store and
// synchronised across sessions with Group.Broadcast.
func New(board *store.Board, assets *tether.Asset) *tether.Handler[State] {
	group := tether.NewGroup[State]()
	viewers := NewViewers()

	return tether.Stateful(tether.App{
		DevMode: true,
		Assets:  []*tether.Asset{assets},
	}, tether.StatefulConfig[State]{
		Name:     "kanban",
		Mode:     mode.Both,
		Upgrade:  wsupgrade.Upgrade(),
		Fallback: sse.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{View: "board", OnlineCount: group.Len()}
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

		// Generous idle timeout for a demo app.
		Timeouts: tether.Timeouts{Idle: 10 * time.Minute},

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("connected", "id", sess.ID()[:8], "members", group.Len())
			sess.Signal("online_count", group.Len())
			sess.Update(func(s State) State {
				s.SessionID = sess.ID()
				return s
			})
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("disconnected", "id", sess.ID()[:8])
			viewers.Clear(sess.ID())
		},
	})
}

// navigate handles URL-driven state. When the browser navigates to
// /card/<id>, the detail view opens. When it navigates to /, the
// board view shows. This runs on initial page load and on
// back/forward navigation.
func navigate(board *store.Board) func(tether.Session, State, tether.Params) State {
	return func(_ tether.Session, s State, p tether.Params) State {
		path := p.Path
		if after, ok := strings.CutPrefix(path, "/card/"); ok {
			id := after
			if _, ok := board.Card(id); ok {
				s.View = "detail"
				s.SelectedID = id
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
