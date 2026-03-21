// Package layout provides structural containers for arranging
// children. Each function returns a styled div that owns its
// spacing and direction. View files compose with layout.Row,
// layout.Stack, etc. instead of raw div elements.
package layout

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// Row arranges children horizontally with consistent spacing.
func Row(children ...node.Node) *div.Element {
	return div.New(children...).Class("layout-row")
}

// Stack arranges children vertically with consistent spacing.
func Stack(children ...node.Node) *div.Element {
	return div.New(children...).Class("layout-stack")
}

// Container wraps children in a plain div with no styling. Use it
// in view files instead of importing html5/div directly.
func Container(children ...node.Node) *div.Element {
	return div.New(children...)
}
