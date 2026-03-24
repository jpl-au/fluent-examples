// Package badge provides pill-shaped status indicators for the
// example application. Functions return the concrete span element
// so callers can chain bind.Apply or other element methods.
package badge

import "github.com/jpl-au/fluent/html5/span"

// Green renders a green status badge.
func Green(text string) *span.Element {
	return span.Text(text).Class("badge badge-green")
}

// GreenDynamic renders a green status badge with a Tether Dynamic key
// so the differ can track it across re-renders.
func GreenDynamic(key, text string) *span.Element {
	return span.Text(text).Class("badge badge-green").Dynamic(key)
}

// Indigo renders an indigo status badge.
func Indigo(text string) *span.Element {
	return span.Text(text).Class("badge badge-indigo")
}
