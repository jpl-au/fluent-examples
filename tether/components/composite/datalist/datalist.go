// Package datalist provides a styled row-based data list for
// displaying tabular data without HTML tables. Each row shows a
// label and value pair. Designed for use with the tether/window
// package for virtualised lists.
package datalist

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New creates a data list container from the given rows.
func New(rows ...node.Node) node.Node {
	return div.New(rows...).Class("data-list")
}
