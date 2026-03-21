// Package http provides the stateless HTTP section of the feature
// explorer, demonstrating features that work without a persistent
// transport connection (no WebSocket, no SSE).
package http

import (
	gohttp "net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/http/state"
)

// New creates a stateless page handler for the HTTP section. No
// WebSocket, no SSE - each request is independent.
func New(app tether.App, assets *tether.Asset) gohttp.Handler {
	r := newRouter()

	return tether.Stateless(app, tether.StatelessConfig[state.State]{
		InitialState: func(_ *gohttp.Request) state.State {
			return state.State{}
		},
		Render: func(s state.State) node.Node {
			return layout.Shell(layout.SectionHTTP, s.Page, 0, r.Render(s))
		},
		Handle: r.Handle,
		OnNavigate: r.OnNavigate(func(s *state.State, p tether.Params) {
			s.Page = p.Path
		}),

		Layout: func(_ state.State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - HTTP"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}
