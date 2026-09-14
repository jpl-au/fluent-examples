// Package text provides styled text elements for hints, results,
// and inline feedback.
package text

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
)

// Hint creates a muted paragraph for guidance or empty-state messages.
func Hint(s string) *p.Element {
	return p.Text(s).Class("hint")
}

// Hintf creates a formatted hint paragraph.
func Hintf(format string, args ...any) *p.Element {
	return p.Text(fmt.Sprintf(format, args...)).Class("hint")
}

// Result creates an inline span for displaying a value.
func Result(s string) *span.Element {
	return span.Text(s).Class("result-text")
}

// Muted creates a de-emphasised inline span.
func Muted(s string) *span.Element {
	return span.Text(s).Class("result-text result-muted")
}
