// Package card provides styled content cards for grouping related
// information into visually distinct sections.
package card

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h2"
	"github.com/jpl-au/fluent/node"
)

// New creates a card with an optional title and content children.
// If title is empty, the card contains only the children.
func New(title string, children ...node.Node) node.Node {
	nodes := make([]node.Node, 0, len(children)+1)
	if title != "" {
		nodes = append(nodes, h2.Text(title).Class("card-title"))
	}
	nodes = append(nodes, children...)
	return div.New(nodes...).Class("card")
}

// NewWithAction creates a card with a title and an action element
// (e.g. an Edit link) aligned to the right of the title row.
func NewWithAction(title string, action node.Node, children ...node.Node) node.Node {
	header := div.New(
		h2.Text(title).Class("card-title"),
		action,
	).Class("card-header")

	nodes := make([]node.Node, 0, len(children)+1)
	nodes = append(nodes, header)
	nodes = append(nodes, children...)
	return div.New(nodes...).Class("card")
}
