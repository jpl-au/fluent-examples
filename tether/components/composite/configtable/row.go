package configtable

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Row creates a single label-value pair inside a config table.
func Row(label, value string) node.Node {
	return div.New(
		span.Text(label).Class("config-label"),
		span.Text(value).Class("config-value"),
	).Class("config-row")
}
