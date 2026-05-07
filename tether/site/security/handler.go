package security

import (
	"context"
	"fmt"
	"net/http"

	fsec "github.com/jpl-au/fluent-security"
	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/layout"
)

// State carries the per-request CSP nonce so the render layer can
// stamp it on inline elements and the view can quote it to the user.
type State struct {
	Nonce string
}

// ctxKey is the context key used to pass the per-request nonce from
// the outer handler (which sets the CSP header) to the tether state
// constructor (which reads the nonce into State). Threading via
// context keeps the header and the rendered <script> guaranteed to
// carry the same value.
type ctxKey struct{}

// New creates a stateless page handler for the security demo. Each
// request gets a single fresh nonce, stamped on the Content-Security-
// Policy header, then threaded through the request context into the
// tether State and rendered onto an inline <script>. Because both
// values come from the same call to [security.Nonce], the browser's
// CSP check passes and the inline script runs.
func New(app tether.App, assets *tether.Asset) http.Handler {
	h := tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(r *http.Request) State {
			n, _ := r.Context().Value(ctxKey{}).(string)
			return State{Nonce: n}
		},
		// The security page is read-only - no user events to process -
		// but tether.Stateless requires Handle. Return state unchanged.
		Handle: func(_ tether.Session, s State, _ tether.Event) State {
			return s
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/security/", 0, Render(s))
		},
		Layout: func(s State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Security"),
					assets.Stylesheet("app.css"),
					// Inline script stamped with the per-request nonce.
					// The browser executes it only because the matching
					// 'nonce-...' token appears in the CSP header set
					// by the outer handler on the same response.
					script.RawText("console.log('fluent-security demo: inline script permitted by CSP nonce');").Nonce(s.Nonce),
				),
				body.New(content),
			).Lang("en")
		},
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := fsec.Nonce()
		w.Header().Set("Content-Security-Policy", fmt.Sprintf(
			"default-src 'self'; script-src 'self' 'nonce-%s'; style-src 'self' 'unsafe-inline'",
			n,
		))
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, n)))
	})
}
