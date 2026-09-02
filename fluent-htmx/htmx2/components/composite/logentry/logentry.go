// Package logentry renders a single log entry as an HTML fragment.
// Both the WebSocket and SSE handlers use this to produce server-side
// rendered HTML that htmx swaps directly into the page.
package logentry

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-htmx/generate"
)

// New builds a single log entry div with time, level, message, and
// attributes spans. The returned node can be rendered to HTML bytes
// for sending over WebSocket or SSE.
func New(e generate.Entry) node.Node {
	return div.New(
		span.Text(e.Time).Class("log-time"),
		span.Text(e.Level).Class("log-level log-level-"+e.Level),
		span.Text(e.Message).Class("log-message"),
		span.Text(e.Attrs).Class("log-attrs"),
	).Class("log-entry")
}
