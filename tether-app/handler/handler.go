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
	tetherfs "github.com/jpl-au/tether-store/fs"
	"github.com/jpl-au/tether/mode"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New creates the kanban board handler. A single handler serves all
// connected browsers; board state is shared via the store and
// synchronised across sessions with Group.Broadcast.
//
// The handler demonstrates three optimisation strategies:
//
//   - Signals for presence indicators (typing, viewing, online count).
//     High-frequency, text-only updates that skip the render cycle.
//   - Memoise for board columns. Expensive subtrees that only
//     re-render when BoardVersion changes (board mutations).
//   - Patch for targeted card updates after edits. Only the saved
//     card is re-rendered and diffed, not the entire board.
func New(board *store.Board, assets *tether.Asset) *tether.Handler[State] {
	group := tether.NewGroup[State]()
	viewers := newViewers()

	return tether.Stateful(tether.App{
		DevMode: true,
		Assets:  []*tether.Asset{assets},
		// ViewTransitions cross-fades the content region when it swaps
		// between the board and a card detail, so the SPA-style region
		// morph animates instead of snapping. The browser falls back to
		// an instant swap when it lacks the API or the user has
		// prefers-reduced-motion set.
		Client: tether.Client{ViewTransitions: true},
	}, tether.StatefulConfig[State]{
		Name: "kanban",
		Mode: mode.Both,

		// Memoise enables subtree memoisation. Column renders are
		// wrapped in jit.Memoise(s.BoardVersion, ...) so they are
		// skipped entirely when the board hasn't changed. See view.go.
		Memoise: true,

		// Persistence: SessionStore saves session state (State struct)
		// to disk on disconnect and graceful shutdown. DiffStore saves
		// differ snapshots so reconnecting clients receive targeted
		// patches instead of a full morph. Together they provide
		// crash recovery and seamless server restarts.
		SessionStore: tetherfs.NewSessionStore(".tether/sessions"),
		DiffStore:    tetherfs.NewDiffStore(".tether/diffs"),

		InitialState: func(_ *http.Request) State {
			return State{View: "board", OnlineCount: group.Count().Load()}
		},
		Render:     Render(board, descCleaner),
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
			// WatchValue tracks the online count in state for the
			// initial SSR render. The signal push in the callback
			// keeps the badge up to date on subsequent changes
			// without a render cycle.
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

		// OnRestore fires instead of OnConnect when a session is
		// recovered from the SessionStore (server restart, crash
		// recovery). The state has been deserialised - rejoin the
		// group and re-establish presence tracking.
		OnRestore: func(sess *tether.StatefulSession[State]) {
			slog.Info("restored", "id", sess.ID()[:8], "name", sess.State().Name)
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
		s.MenuOpen = false
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
