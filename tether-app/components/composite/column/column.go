// Package column renders a single kanban swimlane with a header
// (title and card count) and a list of card nodes.
package column

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h2"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// New renders a column with its title, card count, and card children.
func New(title string, count int, cards ...node.Node) node.Node {
	hdr := div.New(
		h2.Text(title).Class("column-title"),
		span.Text(fmt.Sprintf("%d", count)).Class("badge badge-count"),
	).Class("column-header")

	body := div.New(cards...).Class("column-body")

	return div.New(hdr, body).Class("column")
}

// Empty renders a placeholder when a column has no cards.
func Empty() node.Node {
	return div.New(
		span.Text("No cards").Class("column-empty"),
	).Class("column-placeholder")
}
