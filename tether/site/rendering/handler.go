package rendering

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/component/counter"
	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/middleware"
)

// State holds per-request state for the state and rendering demo.
type State struct {
	// Counter is a plain integer incremented by a button click.
	Counter int
	// Items is a dynamic list grown by the "Add Item" button.
	Items []string
	// Counter2 is a component-managed counter demonstrating
	// tether.RouteTyped event dispatch.
	Counter2 counter.Counter
	// Group is a nested component group demonstrating parent/child
	// event routing via tether.RouteTyped.
	Group counter.Group
}

// New creates a stateless page handler for the rendering demo.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State {
			return State{
				Counter2: counter.New("counter"),
				Group:    counter.NewGroup("group"),
			}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/rendering/", 0, Render(s))
		},
		Handle:     Handle,
		Middleware: []tether.Middleware[State]{middleware.Logging[State]},
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - State & Rendering"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes events on the state and rendering page. Because
// the HTTP section is stateless, counters and item lists are
// round-tripped via bind.EventData on each request.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "rendering.increment":
		s.Counter = readCount(ev, "count") + 1
	case "rendering.decrement":
		if c := readCount(ev, "count"); c > 0 {
			s.Counter = c - 1
		}
	case "rendering.add-item":
		items := decodeItems(ev)
		n := len(items) + 1
		s.Items = append(items, "Item "+strconv.Itoa(n))
	case "rendering.remove-item":
		items := decodeItems(ev)
		if len(items) > 0 {
			s.Items = items[:len(items)-1]
		}
	default:
		s.Counter2 = tether.RouteTyped(s.Counter2, "counter", sess, ev)
		s.Group = tether.RouteTyped(s.Group, "group", sess, ev)
	}
	return s
}

// readCount extracts an integer from event data, defaulting to zero
// - used by the stateless page to round-trip counter values.
func readCount(ev tether.Event, key string) int {
	s, _ := ev.Get(key)
	n, _ := strconv.Atoi(s)
	return n
}

// decodeItems is the inverse of encodeItems - splits the pipe-
// delimited string back into a slice for state reconstruction.
func decodeItems(ev tether.Event) []string {
	s, _ := ev.Get("items")
	if s == "" {
		return nil
	}
	return strings.Split(s, "|")
}
