package memo

import tether "github.com/jpl-au/tether"

// RealtimeState holds rolling metric samples for the memoised
// system monitor charts. Each metric is wrapped in Versioned so
// the chart regions are only re-rendered when the data changes.
type RealtimeState struct {
	HeapMB      tether.Versioned[[]float64]
	Goroutines  tether.Versioned[[]int]
	CPUPercent  tether.Versioned[[]float64]
	OnlineCount int
}
