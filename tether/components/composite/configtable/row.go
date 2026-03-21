package configtable

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Row creates a single label-value pair inside a config table.
func Row(label, value string) node.Node {
	return div.New(
		span.New().Class("config-label").Text(label),
		span.New().Class("config-value").Text(value),
	).Class("config-row")
}
