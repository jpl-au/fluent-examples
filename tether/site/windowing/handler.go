package windowing

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

const (
	totalItems = 1000
	pageSize   = 30
)

// Item represents a row in the windowed table.
type Item struct {
	ID   int
	Name string
}

// State holds per-session state for the windowing demo.
type State struct {
	Items        []Item
	ScrollOffset int
	OnlineCount  int
}

var windowPresence = shared.NewPresenceCountOnly()

// seedItems creates the full dataset.
func seedItems() []Item {
	items := make([]Item, totalItems)
	for i := range items {
		items[i] = Item{ID: i + 1, Name: "Item " + strconv.Itoa(i+1)}
	}
	return items
}

// New creates a stateful handler demonstrating windowed rendering.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:     "windowing",
		Upgrade:  wsupgrade.Upgrade(),
		Fallback: sse.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{
				Items:       seedItems(),
				OnlineCount: windowPresence.OnlineCount.Load(),
			}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/windowing/", s.OnlineCount, Render(s))
		},
		Handle: Handle,
		OnNavigate: func(_ tether.Session, s State, p tether.Params) State {
			s.ScrollOffset = (p.IntDefault("page", 1) - 1) * pageSize
			if s.ScrollOffset < 0 {
				s.ScrollOffset = 0
			}
			if s.ScrollOffset > len(s.Items)-pageSize {
				s.ScrollOffset = len(s.Items) - pageSize
			}
			return s
		},

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Windowing"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](windowPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("windowing: connected", "id", sess.ID())
			shared.TrackPresence(windowPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("windowing: disconnected", "id", sess.ID())
			shared.UntrackPresence(windowPresence, sess.ID())
		},
	})
}

// Handle processes windowing demo events. After each page change,
// the URL is updated so refreshing or sharing the link lands on
// the same page.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "window.next":
		s.ScrollOffset += pageSize
		if s.ScrollOffset > len(s.Items)-pageSize {
			s.ScrollOffset = len(s.Items) - pageSize
		}
	case "window.prev":
		s.ScrollOffset -= pageSize
		if s.ScrollOffset < 0 {
			s.ScrollOffset = 0
		}
	case "window.first":
		s.ScrollOffset = 0
	case "window.last":
		s.ScrollOffset = len(s.Items) - pageSize
	default:
		return s
	}
	page := s.ScrollOffset/pageSize + 1
	sess.ReplaceURL("/windowing/?page=" + strconv.Itoa(page))
	return s
}
