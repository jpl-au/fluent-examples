package datalist

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Row creates a single data row with a label and value.
func Row(label, value string) node.Node {
	return div.New(
		span.Text(label).Class("data-label"),
		span.Text(value).Class("data-value"),
	).Class("data-row")
}

// RowWithID creates a data row with an HTML ID for test targeting.
func RowWithID(id, label, value string) node.Node {
	return div.New(
		span.Text(label).Class("data-label"),
		span.Text(value).Class("data-value"),
	).Class("data-row").ID(id)
}
