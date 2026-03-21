// Package layout provides the HTML shell for the contact manager.
// The full page render uses jit.Tune to learn optimal buffer sizes
// over repeated renders - no code change needed, just a wrapper
// around the existing render call.
package layout

import (
	"net/http"

	jit "github.com/jpl-au/fluent-jit"
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

	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/footer"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/page"
)

// Page renders a full HTML page with the given title, header actions,
// and content sections, writing the result to the ResponseWriter.
//
// This uses the Global Tune API - jit.Tune("page", doc, w) - because
// the page shell is rendered on every request but the content varies in
// size across different pages (list vs detail vs form). Tune learns the
// optimal buffer size over repeated renders without analysing tree
// structure, so it adapts naturally as traffic patterns change.
//
// For pages where you need per-handler tuning or want to isolate buffer
// statistics, the Instance API (jit.NewTuner()) gives fine-grained
// control - each handler would own its own Tuner instance instead of
// sharing a single global key.
//
// Handlers that need a different JIT strategy (Compile or Flatten) can
// call Document() to get the raw node tree, then apply their own
// strategy directly.
func Page(w http.ResponseWriter, pageTitle string, headerActions node.Node, content ...node.Node) {
	doc := document(pageTitle, headerActions, nil, content...)

	// Global Tune API: the string key "page" identifies this template
	// in the global registry. Tune tracks render sizes and adapts the
	// buffer allocation over time - ideal for content that varies in
	// length across requests.
	jit.Tune("page", doc, w)
}

// PageWithScripts renders a full HTML page with additional script
// tags loaded before the closing body tag.
func PageWithScripts(w http.ResponseWriter, pageTitle string, headerActions node.Node, scripts []string, content ...node.Node) {
	doc := document(pageTitle, headerActions, scripts, content...)
	jit.Tune("page", doc, w)
}

// Document builds the full HTML document node without rendering it.
// This is the building block for handlers that want to apply their own
// JIT strategy (Compile or Flatten) instead of the default Tune used
// by Page. The returned node can be passed directly to jit.Compile,
// jit.Flatten, or a Compiler/Flattener instance.
func Document(pageTitle string, headerActions node.Node, content ...node.Node) node.Node {
	return document(pageTitle, headerActions, nil, content...)
}

// document builds the full HTML document node tree.
func document(pageTitle string, headerActions node.Node, scripts []string, content ...node.Node) node.Node {
	children := []node.Node{Shell(pageTitle, headerActions, content...)}
	for _, src := range scripts {
		children = append(children, script.JavaScript(src))
	}
	return html.New(
		head.New(
			meta.UTF8(),
			meta.Viewport("width=device-width, initial-scale=1.0"),
			title.Text("Fluent-JIT - "+pageTitle),
			link.Stylesheet("/static/app.css"),
		),
		body.New(children...),
	)
}

// Shell wraps content in the header and content structure without
// the outer html/head elements. Each content node becomes a direct
// child of the .page container so the flex gap spaces them evenly.
func Shell(pageTitle string, headerActions node.Node, content ...node.Node) node.Node {
	headerNodes := []node.Node{
		headerLeftNode(pageTitle),
	}
	if headerActions != nil {
		headerNodes = append(headerNodes, div.New(headerActions).Class("header-actions"))
	}

	return div.New(
		div.New(headerNodes...).Class("header"),
		div.New(page.New(content...)).Class("content"),
		footer.New(),
	).Class("shell")
}

// headerLeftNode builds the left side of the header bar - the brand
// label is static (never changes), the page title is dynamic.
func headerLeftNode(pageTitle string) node.Node {
	return div.New(
		span.Static("Fluent-JIT").Class("header-brand"),
		span.Text(pageTitle).Class("header-title"),
	).Class("header-left")
}
