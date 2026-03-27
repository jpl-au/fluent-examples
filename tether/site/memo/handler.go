package memo

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
	"github.com/jpl-au/tether/sse"
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// Item represents a row in the memoised table.
type Item struct {
	ID   int
	Name string
}

// State holds per-session state for the memo demo.
type State struct {
	// Items is wrapped in Versioned so the memo key tracks changes
	// automatically via With().
	Items tether.Versioned[[]Item]
	// Count is a plain field - the counter region uses a standard
	// Dynamic key and re-renders on every cycle.
	Count       int
	OnlineCount int
}

var memoPresence = shared.NewPresenceCountOnly()

// seedItems creates the initial dataset.
func seedItems() []Item {
	items := make([]Item, 20)
	for i := range items {
		items[i] = Item{ID: i + 1, Name: "Item " + strconv.Itoa(i+1)}
	}
	return items
}

// New creates a stateful handler demonstrating memoisation.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:     "memo",
		Upgrade:  wsupgrade.Upgrade(),
		Fallback: sse.Upgrade(),
		Memo:     true,

		InitialState: func(_ *http.Request) State {
			return State{
				Items:       tether.NewVersioned(seedItems()),
				OnlineCount: memoPresence.OnlineCount.Load(),
			}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/memo/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Memoisation"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](memoPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("memo: connected", "id", sess.ID())
			shared.TrackPresence(memoPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("memo: disconnected", "id", sess.ID())
			shared.UntrackPresence(memoPresence, sess.ID())
		},
	})
}

// Handle processes memo demo events.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "memo.increment":
		// Only the counter changes. Items.Version() is unchanged,
		// so the Memoiser skips the table entirely.
		s.Count++
	case "memo.add-item":
		// Items change via With() - version increments automatically.
		// The Memoiser detects the miss and re-renders the table.
		id := len(s.Items.Val) + 1
		s.Items = s.Items.With(append(s.Items.Val, Item{
			ID:   id,
			Name: "Item " + strconv.Itoa(id),
		}))
	}
	return s
}
