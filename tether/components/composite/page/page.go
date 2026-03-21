// Package page provides the top-level content wrapper that every
// view's Render function uses - consistent max-width and spacing.
package page

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New wraps demo sections in the standard page container with
// consistent spacing and max-width.
func New(children ...node.Node) node.Node {
	return div.New(children...).Class("page").Dynamic("page")
}
