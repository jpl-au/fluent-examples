// Package errors demonstrates tether.Catch - the error boundary that
// recovers from panics in child components and renders a fallback
// node instead of crashing the page. The demo deliberately panics
// inside a nested component to show how the boundary contains the
// failure and keeps the rest of the page functional.
package errors
