// Package notelist provides the note list and item components. Delete
// forms use HTMX to swap the content area without a full page reload.
package notelist

import (
	htmx "github.com/jpl-au/fluent-htmx"
	"github.com/jpl-au/fluent-htmx/swap"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-htmx/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent-htmx/store"
)

// New creates a styled note list. If notes is empty a hint is shown.
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

// Item renders a single note with content, timestamp, and a delete
// button. The delete form uses HTMX to swap the content in place.
func Item(contactID string, n store.Note) node.Node {
	action := "/contacts/" + contactID + "/notes/" + n.ID + "/delete"
	f := form.Post(action)
	htmx.New(f).HxPost(action).HxTarget("#content").HxSwap(swap.InnerHTML)

	return div.New(
		p.Text(n.Content).Class("note-content"),
		span.Text(n.Created.Format("2 Jan 2006, 15:04")).Class("note-time"),
		f.Add(button.Submit("Delete").Class("btn btn-danger btn-sm")),
	).Class("note-item")
}
