// Package broadcasting demonstrates cross-session communication via
// tether.Bus: emit events to all connected sessions, track a shared
// message counter with tether.Value, and observe both reactively with
// tether.WatchBus and tether.WatchValue. Includes an async subscriber
// that logs to slog without blocking the publisher.
package broadcasting
