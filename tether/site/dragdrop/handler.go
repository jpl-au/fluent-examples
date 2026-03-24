package dragdrop

import (
	"log/slog"
	"net/http"
	"sync"

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

// Item is a draggable item in the demo.
type Item struct {
	ID   string
	Name string
	Zone string // "left" or "right"
}

// board holds the shared item state across all sessions.
var board = struct {
	mu    sync.RWMutex
	items []Item
}{
	items: []Item{
		{ID: "1", Name: "Alpha", Zone: "left"},
		{ID: "2", Name: "Bravo", Zone: "left"},
		{ID: "3", Name: "Charlie", Zone: "right"},
	},
}

// snapshot returns a copy of the current items.
func snapshot() []Item {
	board.mu.RLock()
	defer board.mu.RUnlock()
	out := make([]Item, len(board.items))
	copy(out, board.items)
	return out
}

// moveItem moves an item to the target zone at the given index.
// An index of -1 appends to the end.
func moveItem(id, zone string, idx int) {
	board.mu.Lock()
	defer board.mu.Unlock()

	// Remove from current position.
	var item Item
	found := false
	for i := range board.items {
		if board.items[i].ID == id {
			item = board.items[i]
			board.items = append(board.items[:i], board.items[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return
	}
	item.Zone = zone

	// Insert at the target index among items in the same zone.
	pos := 0
	zoneIdx := 0
	for pos < len(board.items) {
		if board.items[pos].Zone == zone {
			if zoneIdx == idx {
				break
			}
			zoneIdx++
		}
		pos++
	}
	board.items = append(board.items[:pos], append([]Item{item}, board.items[pos:]...)...)
}

// State is the per-session state.
type State struct {
	OnlineCount int
}

var ddPresence = shared.NewPresenceCountOnly()
var ddGroup = tether.NewGroup[State]()

// New creates a stateful handler demonstrating drag and drop.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "dragdrop",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: ddPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/dragdrop/", s.OnlineCount, Render(snapshot()))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Drag and Drop"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Groups: []*tether.Group[State]{ddGroup},
		Watchers: shared.Watchers[State](ddPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("dragdrop: connected", "id", sess.ID())
			shared.TrackPresence(ddPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("dragdrop: disconnected", "id", sess.ID())
			shared.UntrackPresence(ddPresence, sess.ID())
		},
	})
}

// Handle processes drag-and-drop events.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "item.drop":
		id, _ := ev.Get("id")
		zone, _ := ev.Get("zone")
		idx, idxErr := ev.Int("index")
		if idxErr != nil {
			idx = -1
		}
		if id != "" && zone != "" {
			moveItem(id, zone, idx)
			ddGroup.Broadcast(func(_ *tether.StatefulSession[State], s State) State {
				return s
			})
		}
	}
	return s
}
