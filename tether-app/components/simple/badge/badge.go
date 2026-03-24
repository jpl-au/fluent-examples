// Package badge provides status indicator components for the kanban
// board. Each variant applies a colour matching the column's meaning.
package badge

import (
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Todo renders a muted badge for the To Do column.
func Todo(text string) node.Node {
	return span.New().Class("badge badge-muted").Text(text)
}

// Progress renders a blue badge for the In Progress column.
func Progress(text string) node.Node {
	return span.New().Class("badge badge-blue").Text(text)
}

// Done renders a green badge for the Done column.
func Done(text string) node.Node {
	return span.New().Class("badge badge-green").Text(text)
}

// Count renders a small count indicator.
func Count(text string) node.Node {
	return span.New().Class("badge badge-count").Text(text)
}
