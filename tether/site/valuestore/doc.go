// Package valuestore demonstrates tether.Value - a thread-safe
// reactive container for shared state. Sessions observe the value
// via tether.WatchValue and see updates in real time. The demo
// contrasts shared state (Value.Store, Value.Update) with
// per-session local state to show when each is appropriate.
package valuestore
