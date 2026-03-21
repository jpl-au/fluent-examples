// Package viewport provides scrollable list elements for infinite
// scroll patterns - a constrained list container, styled items, and
// a sentinel element that triggers loading when it enters the viewport.
package viewport

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
)

// List creates a scrollable viewport list container.
func List(items ...node.Node) node.Node {
	return ul.New(items...).Class("viewport-list")
}

// Item creates a single viewport list item with text.
func Item(s string) node.Node {
	return li.New().Class("viewport-item").Text(s)
}

// Itemf creates a single viewport list item with formatted text.
func Itemf(format string, args ...any) node.Node {
	return li.Textf(format, args...).Class("viewport-item")
}

// Sentinel creates the trailing element that triggers the next page
// load via bind.OnViewport. Returns the concrete element so callers
// can chain .Dynamic() and bind.Apply.
func Sentinel() *div.Element {
	return div.New().Class("viewport-sentinel")
}
