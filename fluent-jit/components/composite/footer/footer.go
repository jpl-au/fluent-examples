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
		a.New().Href("/").Text("Home").Class("footer-link"),
		span.Static("\u00b7").Class("footer-sep"),
		a.New().Href("/ws").Text("WebSocket Demo").Class("footer-link"),
		span.Static("\u00b7").Class("footer-sep"),
		a.New().Href("/sse").Text("SSE Demo").Class("footer-link"),
	).Class("footer")
}
