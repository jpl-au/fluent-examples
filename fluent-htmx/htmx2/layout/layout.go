// Package layout provides the HTML shell for the contact manager.
// The full page loads the HTMX library, carries the htmx config and
// gives the content area a target ID. HTMX requests receive partials
// from the Partial function, which sends the page title and the
// header sections as <hx-partial> blocks so they update alongside
// the content.
package layout

import (
	"log/slog"
	"net/http"

	htmx "github.com/jpl-au/fluent-htmx/htmx2"
	"github.com/jpl-au/fluent-htmx/htmx2/swap"
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

	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/composite/footer"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/composite/page"
)

// htmxConfig is the htmx config every page carries in its head.
//
// The handlers answer HX-Request with a partial, so a history restore
// must not carry it: htmx 2 sends HX-Request on a history cache miss
// by default and would swap the partial in as the whole page.
//
// htmx 2 swaps no 4xx response by default, so "Contact not found"
// would never show. The 404 rule swaps it like a normal response;
// htmx replaces the whole list, so its default rules are restated
// around it. htmx 4 swaps 4xx responses already.
var htmxConfig = func() string {
	cfg, err := htmx.Config().
		HistoryRestoreAsHxRequest(false).
		ResponseHandling([]map[string]any{
			{"code": "204", "swap": false},
			{"code": "[23]..", "swap": true},
			{"code": "404", "swap": true},
			{"code": "[45]..", "swap": false, "error": true},
		}).
		ToJSON()
	if err != nil {
		// The settings are fixed, so this fails only if they are mistyped.
		panic(err)
	}
	return cfg
}()

// Page renders a full HTML page with the HTMX library loaded. The
// content area has id="content" which HTMX requests target for
// partial swaps.
func Page(w http.ResponseWriter, pageTitle string, headerActions node.Node, content ...node.Node) {
	slog.Debug("render", "mode", "full page", "title", pageTitle)
	doc := html.New(
		headNode(pageTitle),
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
		headNode(pageTitle),
		body.New(children...),
	)
	doc.Render(w)
}

// Partial renders the page content for an HTMX swap, plus the page
// title and the two header sections as <hx-partial> blocks. htmx
// swaps the main content into the target of the element that made
// the request, sets the tab title from the <title>, and routes each
// block to its own target, so the header stays in sync with the
// content.
func Partial(w http.ResponseWriter, pageTitle string, headerActions node.Node, content ...node.Node) {
	slog.Debug("render", "mode", "htmx partial", "title", pageTitle)

	title.Text("Fluent-HTMX - " + pageTitle).Render(w)
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

// headNode builds the document head shared by every full page.
func headNode(pageTitle string) node.Node {
	return head.New(
		meta.UTF8(),
		meta.Viewport("width=device-width, initial-scale=1.0"),
		meta.Name("htmx-config", htmxConfig),
		title.Text("Fluent-HTMX - "+pageTitle),
		link.Stylesheet("/static/app.css"),
		script.JavaScript("/static/htmx.min.js"),
	)
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
