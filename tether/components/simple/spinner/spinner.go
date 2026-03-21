// Package spinner provides a loading indicator element for use with
// bind.Indicator to show progress during server requests.
package spinner

import "github.com/jpl-au/fluent/html5/span"

// New creates a loading spinner. Returns the concrete element so
// callers can chain .ID() for bind.Indicator targeting.
func New() *span.Element {
	return span.New().Class("spinner")
}
