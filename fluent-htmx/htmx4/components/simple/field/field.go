// Package field provides styled form field elements for labels,
// text inputs, and text areas.
package field

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/label"
	"github.com/jpl-au/fluent/html5/textarea"
	"github.com/jpl-au/fluent/node"
)

// Text creates a single-line text input with a name and placeholder.
func Text(name, placeholder string) *input.Element {
	return input.Text(name, "").Class("input").Placeholder(placeholder)
}

// TextValue creates a single-line text input pre-filled with a value.
func TextValue(name, value, placeholder string) *input.Element {
	return input.Text(name, value).Class("input").Placeholder(placeholder)
}

// TextArea creates a multi-line text area with a name and placeholder.
func TextArea(name, placeholder string) node.Node {
	return textarea.New().Name(name).Class("input").Placeholder(placeholder)
}

// TextAreaValue creates a multi-line text area pre-filled with a value.
func TextAreaValue(name, value, placeholder string) node.Node {
	return textarea.Text(value).Name(name).Class("input").Placeholder(placeholder)
}

// Label creates a styled label associated with a form element by ID.
func Label(forID, text string) node.Node {
	return label.For(forID, text).Class("label")
}

// Group wraps child nodes in a styled form-group container.
func Group(children ...node.Node) node.Node {
	return div.New(children...).Class("form-group")
}
