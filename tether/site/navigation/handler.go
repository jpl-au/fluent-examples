package navigation

import (
	"log/slog"
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

// State holds per-request state for the navigation demo.
type State struct {
	// Page is the current path within this handler.
	Page string
	// Tab is the selected tab from ?tab=<value>.
	Tab string
	// Search is the search term from ?q=<value>.
	Search string
	// QueryPage is the pagination offset from ?page=<value>.
	QueryPage int
	// Tags is a multi-value param from repeated ?tag=<value> keys.
	Tags []string
	// Active is a boolean param from ?active=<true|false>.
	Active bool
	// Price is a float param from ?price=<value>.
	Price float64
	// Quantities is a multi-value int param from repeated ?qty=<value> keys.
	Quantities []int
	// Prices is a multi-value float param from repeated ?price=<value> keys.
	Prices []float64
}

// New creates a stateless page handler for the navigation demo.
// Includes OnNavigate to demonstrate query parameter extraction
// and a small internal routing switch for bind.Link targets.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State {
			return State{Page: "/navigation/", QueryPage: 1}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/navigation/", 0, Render(s))
		},
		Handle: Handle,
		OnNavigate: func(_ tether.Session, s State, p tether.Params) State {
			s.Page = p.Path
			s.Tab = p.Get("tab")
			s.Search = p.Get("q")
			s.QueryPage = p.IntDefault("page", 1)
			s.Tags = p.Strings("tag")
			s.Active = p.BoolDefault("active", false)
			s.Price = p.Float64Default("price", 0)
			if qtys, err := p.Ints("qty"); err != nil {
				slog.Debug("navigation: invalid qty param", "err", err)
			} else {
				s.Quantities = qtys
			}
			if prices, err := p.Float64s("price"); err != nil {
				slog.Debug("navigation: invalid price param", "err", err)
			} else {
				s.Prices = prices
			}
			return s
		},
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Navigation"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes events on the navigation page.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "nav.replace-a":
		sess.ReplaceURL("/navigation/?tab=a")
	case "nav.replace-b":
		sess.ReplaceURL("/navigation/?tab=b")
	case "nav.goto-target":
		sess.Navigate("/navigation/target/")
	}
	return s
}
