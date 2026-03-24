// Package field provides styled form field elements - text inputs,
// file inputs, labels, form groups, and inline form layouts.
package field

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/label"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/node"
)

// Text creates a styled text input with the given name and placeholder.
func Text(name, placeholder string) *input.Element {
	return input.Text(name, "").Class("input").Placeholder(placeholder)
}

// TextValue creates a styled text input with a pre-filled value.
func TextValue(name, value, placeholder string) *input.Element {
	return input.Text(name, value).Class("input").Placeholder(placeholder)
}

// TextWithID creates a styled text input with an explicit ID.
func TextWithID(id, name, placeholder string) *input.Element {
	return input.Text(name, "").ID(id).Class("input").Placeholder(placeholder)
}

// File creates a styled file input with an explicit ID.
func File(id, name string) *input.Element {
	return input.File(name).ID(id).Class("input")
}

// Label creates a styled form label.
func Label(forID, text string) node.Node {
	return label.For(forID, text).Class("label")
}

// Group wraps children in a form group container.
func Group(children ...node.Node) node.Node {
	return div.New(children...).Class("form-group")
}

// Inline wraps children in a horizontal inline form layout. Returns
// the concrete form element so callers can chain bind.Apply.
func Inline(children ...node.Node) *form.Element {
	return form.New(children...).Class("form-inline")
}

// Error renders a validation error message below a form field.
func Error(s string) node.Node {
	return p.Text(s).Class("form-error")
}
