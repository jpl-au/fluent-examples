package htmlwire

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
	"github.com/jpl-au/tether/wire"

	"github.com/jpl-au/fluent-examples/tether/layout"
)

// State holds per-request state for the HTML wire format demo.
type State struct {
	// Counter travels with each event as event data - stateless
	// pages reconstruct state from the request.
	Counter int
}

// New creates a stateless page handler that answers POST events with
// plain HTML instead of the JSON envelope.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		// Responses are plain HTML: morph fragments as the body,
		// side effects in a JSON island. Inspect them with curl.
		WireFormat: wire.HTML,

		// Stateless pages embed no session token, so the initial GET
		// is safe for browsers and CDNs to cache.
		CacheControl: "public, max-age=30",

		InitialState: func(_ *http.Request) State { return State{} },
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/html-wire/", 0, Render(s))
		},
		Handle: Handle,
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - HTML Wire Format"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes events for the HTML wire format demo.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "htmlwire.fragment":
		s.Counter = readCount(ev) + 1
		// Only the "hw-count" fragment travels back: the response
		// body is the single keyed element, and a Tether-Morph:
		// keyed header tells the client to apply it as a targeted
		// morph.
		sess.Morph("hw-count")
	case "htmlwire.full":
		s.Counter = readCount(ev) + 1
		// No Morph call - the response body is the whole rendered
		// page. The toast rides in the JSON effects island appended
		// to the HTML.
		sess.Toast("This update arrived as plain HTML")
	}
	return s
}

func readCount(ev tether.Event) int {
	raw, _ := ev.Get("count")
	n, _ := strconv.Atoi(raw)
	return n
}
