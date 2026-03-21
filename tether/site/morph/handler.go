package morph

import (
	"net/http"
	"strconv"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/layout"
)

// State holds per-request state for the morph demo.
type State struct {
	// Counter is a plain integer incremented by button clicks.
	Counter int
}

// New creates a stateless page handler for the full-page morph demo.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State { return State{} },
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/morph", 0, Render(s))
		},
		Handle: Handle,
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Full-Page Morph"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes events for the morph demo.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "morph.increment":
		s.Counter = readCount(ev) + 1
	case "morph.decrement":
		if c := readCount(ev); c > 0 {
			s.Counter = c - 1
		}
	}
	return s
}

func readCount(ev tether.Event) int {
	raw, _ := ev.Get("count")
	n, _ := strconv.Atoi(raw)
	return n
}
