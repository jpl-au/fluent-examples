// Package display provides prominent value elements for counters,
// scores, and other large numeric displays.
package display

import (
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Value renders a prominent display for a string value.
func Value(s string) node.Node {
	return span.New().Class("counter-display").Text(s)
}

// Valuef renders a prominent display with formatted text.
func Valuef(format string, args ...any) node.Node {
	return span.New().Class("counter-display").Textf(format, args...)
}
