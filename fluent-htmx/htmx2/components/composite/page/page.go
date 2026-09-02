// Package page provides the content wrapper that every view uses.
package page

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New wraps content in the standard page container with consistent
// spacing and max-width.
func New(children ...node.Node) node.Node {
	return div.New(children...).Class("page")
}
