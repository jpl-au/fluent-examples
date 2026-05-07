// Package footer provides the sticky footer bar linking to the
// WebSocket and SSE live demo pages.
package footer

import (
	"github.com/jpl-au/fluent/html5/a"
	el "github.com/jpl-au/fluent/html5/footer"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// New builds the bottom navigation bar with links to the live demos.
func New() node.Node {
	return el.New(
		a.Text("Home").Href("/").Class("footer-link"),
		span.Static("\u00b7").Class("footer-sep"), // middle dot separator
		a.Text("WebSocket Demo").Href("/ws").Class("footer-link"),
		span.Static("\u00b7").Class("footer-sep"), // middle dot separator
		a.Text("SSE Demo").Href("/sse").Class("footer-link"),
	).Class("footer")
}
