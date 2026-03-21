// Package notelist provides the note list and item components used
// on the contact detail page.
package notelist

import (
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-jit/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent-jit/store"
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
// form.
func Item(contactID string, n store.Note) node.Node {
	return div.New(
		p.Text(n.Content).Class("note-content"),
		span.Text(n.Created.Format("2 Jan 2006, 15:04")).Class("note-time"),
		form.Post(
			"/contacts/"+contactID+"/notes/"+n.ID+"/delete",
			button.Submit("Delete").Class("btn btn-danger btn-sm"),
		),
	).Class("note-item")
}
