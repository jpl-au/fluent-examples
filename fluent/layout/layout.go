// Package layout provides the HTML shell for the contact manager.
// Every page is wrapped in a consistent document structure with a
// header bar, footer, stylesheet, and viewport meta tag.
package layout

import (
	"net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/link"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent/components/composite/footer"
	"github.com/jpl-au/fluent-examples/fluent/components/composite/page"
)

// Page renders a full HTML page writing the result to w.
func Page(w http.ResponseWriter, t string, actions node.Node, content ...node.Node) {
	doc(t, actions, nil, content...).Render(w)
}

// PageWithScripts renders a full HTML page with additional script
// tags loaded before the closing body tag.
func PageWithScripts(w http.ResponseWriter, t string, actions node.Node, scripts []string, content ...node.Node) {
	doc(t, actions, scripts, content...).Render(w)
}

// doc builds the full HTML document node tree.
func doc(t string, actions node.Node, scripts []string, content ...node.Node) node.Node {
	children := []node.Node{Shell(t, actions, content...)}
	for _, src := range scripts {
		children = append(children, script.JavaScript(src))
	}
	return html.New(
		head.New(
			meta.UTF8(),
			meta.Viewport("width=device-width, initial-scale=1.0"),
			title.Text("Fluent - "+t),
			link.Stylesheet("/static/app.css"),
		),
		body.New(children...),
	)
}

// Shell wraps content in the header, content area, and footer.
func Shell(t string, actions node.Node, content ...node.Node) node.Node {
	header := []node.Node{
		div.New(
			span.Static("Fluent").Class("header-brand"),
			span.Text(t).Class("header-title"),
		).Class("header-left"),
	}
	if actions != nil {
		header = append(header, div.New(actions).Class("header-actions"))
	}
	return div.New(
		div.New(header...).Class("header"),
		div.New(page.New(content...)).Class("content"),
		footer.New(),
	).Class("shell")
}
