// Package list provides a styled list view with bordered items,
// consistent with the shadcn-inspired design system.
package list

import (
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
)

// New creates a styled list from the given items.
func New(items ...node.Node) node.Node {
	return ul.New(items...).Class("item-list")
}

// Empty renders a placeholder message when a list has no items.
func Empty(text string) node.Node {
	return p.Text(text).Class("hint")
}
