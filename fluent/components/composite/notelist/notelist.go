// Package notelist provides the note list and item components used
// on the contact detail page.
package notelist

import (
	security "github.com/jpl-au/fluent-security"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent/store"
)

// New creates a styled note list. If notes is empty a hint is shown
// instead.
func New(contactID string, notes []store.Note) node.Node {
	if len(notes) == 0 {
		return text.Hint("No notes yet.")
	}
	items := make([]node.Node, len(notes))
	for i, n := range notes {
		items[i] = Item(contactID, n)
	}
	return div.New(items...).Class("note-list")
}

// Item renders a single note - content, timestamp, and a delete
// form. Note content is treated as untrusted HTML and passed through
// fluent-security's UGC policy before rendering, so formatting tags
// survive while scripts, event handlers, and javascript: URIs are
// stripped. See the seeded attack fixtures in store.init for what
// this catches in practice.
func Item(contactID string, n store.Note) node.Node {
	return div.New(
		div.New(security.HTML(n.Content)).Class("note-content"),
		span.Text(n.Created.Format("2 Jan 2006, 15:04")).Class("note-time"),
		form.Post(
			"/contacts/"+contactID+"/notes/"+n.ID+"/delete",
			button.Submit("Delete").Class("btn btn-danger btn-sm"),
		),
	).Class("note-item")
}
