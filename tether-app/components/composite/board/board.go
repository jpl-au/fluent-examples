// Package board provides the top-level kanban board grid layout.
package board

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New renders the three-column board grid from the given column nodes.
func New(columns ...node.Node) node.Node {
	return div.New(columns...).Class("board").Dynamic("board")
}
