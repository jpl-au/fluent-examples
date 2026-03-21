// Package columns provides a two-column side-by-side layout for
// comparing related content within demo cards.
package columns

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// New creates a two-column grid layout.
func New(children ...node.Node) node.Node {
	return div.New(children...).Class("demo-columns")
}
