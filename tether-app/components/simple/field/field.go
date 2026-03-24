// Package field provides styled form input components for the
// kanban board application.
package field

import (
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/label"
	"github.com/jpl-au/fluent/html5/textarea"
	"github.com/jpl-au/fluent/node"
)

// Text creates a text input with a name and placeholder.
func Text(name, placeholder string) *input.Element {
	return input.Text(name, "").Class("input").Placeholder(placeholder)
}

// TextValue creates a text input pre-filled with a value.
func TextValue(name, value, placeholder string) *input.Element {
	return input.Text(name, value).Class("input").Placeholder(placeholder)
}

// Area creates a textarea with a name, placeholder, and content.
func Area(name, placeholder, content string) node.Node {
	return textarea.New().Name(name).Class("input textarea").Placeholder(placeholder).Text(content)
}

// Label creates a form label.
func Label(text string) node.Node {
	return label.New().Class("label").Text(text)
}

// Inline wraps children in a horizontal inline form layout. Returns
// the concrete form element so callers can chain bind.Apply.
func Inline(children ...node.Node) *form.Element {
	return form.New(children...).Class("form-inline")
}
