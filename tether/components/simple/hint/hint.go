// Package hint provides muted helper text for descriptions,
// empty-state prompts, and contextual guidance.
package hint

import (
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
)

// Text renders a block-level hint paragraph.
func Text(s string) *p.Element {
	return p.Text(s).Class("hint")
}

// Textf renders a block-level hint with formatted text.
func Textf(format string, args ...any) *p.Element {
	return p.Textf(format, args...).Class("hint")
}

// Success renders a block-level hint with success (green) styling.
func Success(s string) *p.Element {
	return p.Text(s).Class("hint result-success")
}

// Span renders an inline hint for use alongside other inline elements.
func Span(s string) *span.Element {
	return span.Text(s).Class("hint")
}
