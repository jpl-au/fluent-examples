package errors

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

// State is empty - the error boundary demo has no interactive state.
type State struct{}

// New creates a stateless page handler for the error boundary demo.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State { return State{} },
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/errors/", 0, Render(s))
		},
		Handle: func(_ tether.Session, s State, _ tether.Event) State { return s },
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Error Boundaries"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}
