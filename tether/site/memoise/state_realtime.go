package memoise

import tether "github.com/jpl-au/tether"

// RealtimeState holds rolling metric samples for the system monitor
// charts. Each metric is wrapped in Versioned so the memoisation engine
// can skip unchanged charts during full renders (page load,
// reconnect). Targeted updates via sess.Patch bypass the full
// render entirely.
type RealtimeState struct {
	HeapMB      tether.Versioned[[]float64]
	Goroutines  tether.Versioned[[]int]
	CPUPercent  tether.Versioned[[]float64]
	OnlineCount int
}
