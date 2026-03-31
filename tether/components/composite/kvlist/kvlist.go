// Package kvlist provides a key-value display for showing labelled
// rows of data.
package kvlist

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New creates a key-value list from the given rows.
func New(rows ...node.Node) node.Node {
	return div.New(rows...).Class("config-table")
}
