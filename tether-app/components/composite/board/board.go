// Package board provides the top-level kanban board grid layout.
package board

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// Columns renders the column grid from the given column nodes. The
// caller is responsible for the outer Dynamic wrapper - this allows
// the grid to be used inside a node.Memoise closure where the
// Dynamic key lives on the parent.
func Columns(columns ...node.Node) node.Node {
	return div.New(columns...).Class("board")
}
