// Package contactlist provides the contact list and item components
// used on the main listing page.
package contactlist

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-jit/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent-jit/store"
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

// Item renders a single contact row - name, email, note count, and
// a link to the detail page.
func Item(c store.Contact) node.Node {
	noteLabel := fmt.Sprintf("%d notes", len(c.Notes))
	if len(c.Notes) == 1 {
		noteLabel = "1 note"
	}

	return div.New(
		a.Text(c.Name).Href("/contacts/"+c.ID).Class("contact-name"),
		span.Text(c.Email).Class("contact-email"),
		span.Text(noteLabel).Class("contact-notes"),
	).Class("contact-item")
}
