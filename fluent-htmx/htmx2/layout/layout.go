// Package layout provides the HTML shell for the contact manager.
// The full page includes the HTMX library and a target ID on the
// content area. HTMX requests receive partials via the Partial
// function which also includes an out-of-band header swap so the
// page title updates without a full reload.
package layout

import (
	"log/slog"
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

	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/footer"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/page"
)

// Page renders a full HTML page with the HTMX library loaded. The
// content area has id="content" which HTMX requests target for
// partial swaps. The header has id="header-left" for out-of-band
// title updates.
func Page(w http.ResponseWriter, pageTitle string, headerActions node.Node, content ...node.Node) {
	slog.Debug("render", "mode", "full page", "title", pageTitle)
	doc := html.New(
		head.New(
			meta.UTF8(),
			meta.Viewport("width=device-width, initial-scale=1.0"),
			title.Text("Fluent-HTMX - "+pageTitle),
			link.Stylesheet("/static/app.css"),
			script.JavaScript("/static/htmx.min.js"),
		),
		body.New(
			Shell(pageTitle, headerActions, content...),
		),
	)
	doc.Render(w)
}

// PageWithScripts renders a full HTML page with additional script
// tags loaded before the closing body tag. Used by the WebSocket and
// SSE pages to load the htmx extensions they need.
func PageWithScripts(w http.ResponseWriter, pageTitle string, headerActions node.Node, scripts []string, content ...node.Node) {
	slog.Debug("render", "mode", "full page with scripts", "title", pageTitle)
	children := []node.Node{Shell(pageTitle, headerActions, content...)}
	for _, src := range scripts {
		children = append(children, script.JavaScript(src))
	}
	doc := html.New(
		head.New(
			meta.UTF8(),
			meta.Viewport("width=device-width, initial-scale=1.0"),
			title.Text("Fluent-HTMX - "+pageTitle),
			link.Stylesheet("/static/app.css"),
			script.JavaScript("/static/htmx.min.js"),
		),
		body.New(children...),
	)
	doc.Render(w)
}

// Partial renders the page content for an HTMX swap, plus an
// out-of-band header update so the title changes alongside the
// content. The OOB element carries hx-swap-oob="true" which tells
// HTMX to swap it by ID independently of the main target.
func Partial(w http.ResponseWriter, pageTitle string, headerActions node.Node, content ...node.Node) {
	slog.Debug("render", "mode", "htmx partial", "title", pageTitle)

	// Main content - swapped into #content by the triggering element.
	page.New(content...).Render(w)

	// Out-of-band header update - HTMX matches by id="header-left"
	// and swaps it independently so the title stays in sync.
	headerLeft := headerLeftNode(pageTitle)
	headerLeft.SetAttribute("hx-swap-oob", "true")
	headerLeft.Render(w)

	// Out-of-band header actions update.
	actionsNode := div.New()
	if headerActions != nil {
		actionsNode = div.New(headerActions)
	}
	actionsNode.Class("header-actions").ID("header-actions")
	actionsNode.SetAttribute("hx-swap-oob", "true")
	actionsNode.Render(w)
}

// Shell wraps content in the header and content structure.
func Shell(pageTitle string, headerActions node.Node, content ...node.Node) node.Node {
	headerNodes := []node.Node{
		headerLeftNode(pageTitle),
	}
	actionsNode := div.New()
	if headerActions != nil {
		actionsNode = div.New(headerActions)
	}
	headerNodes = append(headerNodes, actionsNode.Class("header-actions").ID("header-actions"))

	return div.New(
		div.New(headerNodes...).Class("header"),
		div.New(page.New(content...)).Class("content").ID("content"),
		footer.New(),
	).Class("shell")
}

// headerLeftNode builds the brand + title section of the header with
// a stable ID so it can be swapped out-of-band on HTMX requests.
func headerLeftNode(pageTitle string) *div.Element {
	return div.New(
		span.Static("Fluent-HTMX").Class("header-brand"),
		span.Text(pageTitle).Class("header-title"),
	).Class("header-left").ID("header-left")
}
