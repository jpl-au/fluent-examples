// Package contactlist provides the contact list and item components.
// Links use HTMX to swap the content area without a full page reload.
package contactlist

import (
	"fmt"

	htmx "github.com/jpl-au/fluent-htmx"
	"github.com/jpl-au/fluent-htmx/swap"
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-htmx/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent-htmx/store"
)

// New creates a styled contact list. If contacts is empty a hint is
// shown instead.
func New(contacts []store.Contact) node.Node {
	if len(contacts) == 0 {
		return text.Hint("No contacts yet.")
	}
	items := make([]node.Node, len(contacts))
	for i, c := range contacts {
		items[i] = Item(c)
	}
	return div.New(items...).Class("contact-list")
}

// Item renders a single contact row. The name link uses HTMX to
// swap the content area and push the URL to browser history.
func Item(c store.Contact) node.Node {
	noteLabel := fmt.Sprintf("%d notes", len(c.Notes))
	if len(c.Notes) == 1 {
		noteLabel = "1 note"
	}

	href := "/contacts/" + c.ID
	nameLink := a.Text(c.Name).Href(href).Class("contact-name")
	htmx.New(nameLink).HxGet(href).HxTarget("#content").HxPushURL(href).HxSwap(swap.InnerHTML)

	return div.New(
		nameLink,
		span.Text(c.Email).Class("contact-email"),
		span.Text(noteLabel).Class("contact-notes"),
	).Class("contact-item")
}
