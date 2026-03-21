// Package mw demonstrates tether.Middleware - wrapping a HandleFunc
// with cross-cutting behaviour. The demo chains five middleware:
// timing (measures handler duration), guard (short-circuits blocked
// actions), counting (increments an event counter), and two ordered
// wrappers that log entry/exit to visualise the onion execution
// order.
package mw
