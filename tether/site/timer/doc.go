// Package timer demonstrates client-side timers: elapsed counters,
// countdowns, and sub-second stopwatches that tick entirely in the
// browser. The server controls timers by pushing signals - no
// background goroutines or per-tick WebSocket messages required.
package timer
