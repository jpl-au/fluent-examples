package viewtransitions

import (
	"net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/layout"
)

// State holds per-request state for the view transitions demo.
type State struct {
	// Swapped tracks which slot the box occupies.
	Swapped bool
}

// New creates a stateless handler with View Transitions switched on.
// The switch is scoped to this handler's own App copy so the rest of
// the site keeps instant, un-animated updates - view transitions are a
// per-application setting, and this page is the one that opts in.
func New(app tether.App, assets *tether.Asset) http.Handler {
	app.Client.ViewTransitions = true

	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State { return State{} },
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/view-transitions/", 0, Render(s))
		},
		Handle: Handle,
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - View Transitions"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes events on the view transitions page.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	if ev.Action == "vt.toggle" {
		s.Swapped = !s.Swapped
	}
	return s
}
