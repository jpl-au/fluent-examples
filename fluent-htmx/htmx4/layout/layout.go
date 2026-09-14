// Package layout provides the HTML shell for the contact manager.
// The full page loads the htmax.js bundle and gives the content
// area a target ID. HTMX requests receive partials from the Partial
// function, which sends the header sections as <hx-partial> blocks
// so the title updates alongside the content.
package layout

import (
	"log/slog"
	"net/http"

	htmx "github.com/jpl-au/fluent-htmx/htmx4"
	"github.com/jpl-au/fluent-htmx/htmx4/swap"
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

	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/composite/footer"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/composite/page"
)

// Page renders a full HTML page with the htmax.js bundle loaded. The
// bundle carries the SSE and WebSocket extensions, so the live log
// pages need no extra scripts. The content area has id="content"
// which HTMX requests target for partial swaps.
func Page(w http.ResponseWriter, pageTitle string, headerActions node.Node, content ...node.Node) {
	slog.Debug("render", "mode", "full page", "title", pageTitle)
	doc := html.New(
		head.New(
			meta.UTF8(),
			meta.Viewport("width=device-width, initial-scale=1.0"),
			title.Text("Fluent-HTMX - "+pageTitle),
			link.Stylesheet("/static/app.css"),
			script.JavaScript("/static/htmax.js"),
		),
		body.New(
			Shell(pageTitle, headerActions, content...),
		),
	)
	doc.Render(w)
}

// Partial renders the page content for an HTMX swap, plus the two
// header sections as <hx-partial> blocks. htmx swaps the main content
// into the target of the element that made the request and routes
// each block to its own target, so the title and actions stay in
// sync with the content.
func Partial(w http.ResponseWriter, pageTitle string, headerActions node.Node, content ...node.Node) {
	slog.Debug("render", "mode", "htmx partial", "title", pageTitle)

	page.New(content...).Render(w)

	htmx.Partial("#header-left", headerLeftNode(pageTitle)).HxSwap(swap.OuterHTML).Render(w)
	htmx.Partial("#header-actions", headerActionsNode(headerActions)).HxSwap(swap.OuterHTML).Render(w)
}

// Shell wraps content in the header and content structure.
func Shell(pageTitle string, headerActions node.Node, content ...node.Node) node.Node {
	return div.New(
		div.New(headerLeftNode(pageTitle), headerActionsNode(headerActions)).Class("header"),
		div.New(page.New(content...)).Class("content").ID("content"),
		footer.New(),
	).Class("shell")
}

// headerLeftNode builds the brand + title section of the header with
// a stable ID so an <hx-partial> block can replace it.
func headerLeftNode(pageTitle string) *div.Element {
	return div.New(
		span.Static("Fluent-HTMX").Class("header-brand"),
		span.Text(pageTitle).Class("header-title"),
	).Class("header-left").ID("header-left")
}

// headerActionsNode builds the header actions container with a
// stable ID so an <hx-partial> block can replace it. A nil action
// renders an empty container so the ID is always present.
func headerActionsNode(headerActions node.Node) *div.Element {
	actions := div.New()
	if headerActions != nil {
		actions = div.New(headerActions)
	}
	return actions.Class("header-actions").ID("header-actions")
}
