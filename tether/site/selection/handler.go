package selection

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/layout"
)

// State holds per-request state for the selection demo.
type State struct {
	Result string
}

// New creates a stateless handler demonstrating multi-select.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State { return State{} },
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/selection", 0, Render(s))
		},
		Handle: Handle,
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Selection"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes selection demo events.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "selection.action":
		sel, _ := ev.Get("selected")
		if sel == "" {
			s.Result = "No items selected"
		} else {
			ids := strings.Split(sel, ",")
			s.Result = fmt.Sprintf("Selected %d items: %s", len(ids), sel)
		}
	}
	return s
}
