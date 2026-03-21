// Package button provides styled button and link elements for the
// contact manager. Anchor variants navigate via GET; Submit creates
// a form submission button.
package button

import (
	"github.com/jpl-au/fluent/html5/a"
	el "github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/node"
)

// Primary creates a high-contrast link styled as a button.
func Primary(label, href string) node.Node {
	return a.Text(label).Href(href).Class("btn btn-primary")
}

// Secondary creates a medium-contrast link styled as a button.
func Secondary(label, href string) node.Node {
	return a.Text(label).Href(href).Class("btn btn-secondary")
}

// Danger creates a destructive-action link styled as a button.
func Danger(label, href string) node.Node {
	return a.Text(label).Href(href).Class("btn btn-danger")
}

// Link creates a plain styled navigation link.
func Link(label, href string) node.Node {
	return a.Text(label).Href(href).Class("btn btn-link")
}

// Back creates a secondary button with a left arrow for navigation.
// The label is always "← Back" - static content, no escaping needed.
func Back(href string) node.Node {
	return a.Static("← Back").Href(href).Class("btn btn-secondary")
}

// Submit creates a form submit button.
func Submit(label string) node.Node {
	return el.Submit(label).Class("btn btn-primary")
}

// DangerSubmit creates a destructive form submit button.
func DangerSubmit(label string) node.Node {
	return el.Submit(label).Class("btn btn-danger btn-sm")
}
