// Package row provides a horizontal flex layout for controls.
package row

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New arranges children horizontally in a flex row with consistent
// spacing between them.
func New(children ...node.Node) node.Node {
	return div.New(children...).Class("row")
}
