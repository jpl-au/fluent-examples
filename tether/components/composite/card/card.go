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
	nodes := make([]node.Node, 0, len(children)+1)
	if title != "" {
		nodes = append(nodes, h2.New().Class("card-title").Text(title))
	}
	nodes = append(nodes, children...)
	return div.New(nodes...).Class("card")
}
