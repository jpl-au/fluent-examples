// Package card provides a bordered container for grouping related
// content into visually distinct sections.
package card

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h2"
	"github.com/jpl-au/fluent/node"
)

// New creates a card with an optional title and child content.
func New(title string, children ...node.Node) node.Node {
	return div.New(
		node.When(title != "", h2.Text(title).Class("card-title")),
	).Add(children...).Class("card")
}
