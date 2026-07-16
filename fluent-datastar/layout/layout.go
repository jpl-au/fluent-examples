// Package layout provides the HTML shell for the example. The full
// page loads the vendored Datastar client as an ES module and renders
// the demo content inside a simple card layout.
package layout

import (
	"net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/link"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
)

// Page renders a full HTML page with the Datastar client loaded. The
// client is served locally from /static/datastar.js and loaded as an
// ES module, which is how Datastar v1 ships.
func Page(w http.ResponseWriter, pageTitle string, content ...node.Node) {
	doc := html.New(
		head.New(
			meta.UTF8(),
			meta.Viewport("width=device-width, initial-scale=1.0"),
			title.Text("Fluent-Datastar - "+pageTitle),
			link.Stylesheet("/static/app.css"),
			script.Module("/static/datastar.js"),
		),
		body.New(content...),
	)
	doc.Render(w)
}
