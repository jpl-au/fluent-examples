package realtime

import (
	"context"
	"runtime"
	"time"

	tether "github.com/jpl-au/tether"
)

// maxDataPoints is the number of metric samples kept for the charts.
const maxDataPoints = 60

// startMonitor launches a background goroutine that reads Go runtime
// metrics every second and pushes them into session state. The charts
// re-render automatically via the normal diff pipeline.
func startMonitor(sess *tether.StatefulSession[State]) {
	sess.Go(func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		var mem runtime.MemStats
		prevCPU := cpuTime()
		prevWall := time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runtime.ReadMemStats(&mem)
				heap := float64(mem.HeapAlloc) / (1024 * 1024)
				goroutines := runtime.NumGoroutine()

				now := time.Now()
				curCPU := cpuTime()
				wallDelta := now.Sub(prevWall).Seconds()
				cpuPct := 0.0
				if wallDelta > 0 {
					cpuPct = (curCPU - prevCPU) / wallDelta * 100
				}
				prevCPU = curCPU
				prevWall = now

				sess.Update(func(s State) State {
					s.HeapMB = appendMetric(s.HeapMB, heap)
					s.Goroutines = appendMetricInt(s.Goroutines, goroutines)
					s.CPUPercent = appendMetric(s.CPUPercent, cpuPct)
					return s
				})
			}
		}
	})
}

// appendMetric adds a sample to a rolling metric slice, dropping the
// oldest values when the slice exceeds maxDataPoints.
func appendMetric(data []float64, v float64) []float64 {
	data = append(data, v)
	if len(data) > maxDataPoints {
		data = data[len(data)-maxDataPoints:]
	}
	return data
}

// appendMetricInt is the int-typed variant of appendMetric.
func appendMetricInt(data []int, v int) []int {
	data = append(data, v)
	if len(data) > maxDataPoints {
		data = data[len(data)-maxDataPoints:]
	}
	return data
}
