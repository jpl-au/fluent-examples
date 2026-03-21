// Package configtable provides a key-value display for showing
// configuration settings as labelled rows.
package configtable

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New creates a config table from the given rows.
func New(rows ...node.Node) node.Node {
	return div.New(rows...).Class("config-table")
}
