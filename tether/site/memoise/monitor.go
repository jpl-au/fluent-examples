package memoise

import (
	"context"
	"runtime"
	"time"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
)

// maxDataPoints is the number of metric samples kept for the charts.
const maxDataPoints = 60

// startMonitor launches a background goroutine that reads Go runtime
// metrics every second and pushes each chart via sess.Patch. Each
// chart is a targeted update - only the changed chart is re-rendered
// and sent, not the full page.
func startMonitor(sess *tether.StatefulSession[RealtimeState]) {
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

				// Each chart is a targeted Patch. Only the chart div
				// is re-rendered and diffed - the card layout, description,
				// badges, and other charts are untouched.
				sess.Patch("chart-cpu", func(s RealtimeState) (RealtimeState, node.Node) {
					s.CPUPercent = s.CPUPercent.With(appendMetric(s.CPUPercent.Val, cpuPct))
					return s, div.New(
						chartDiv("memoise-cpu", "CPU (%)", "#ee6666", toLineData(s.CPUPercent.Val)),
					).Dynamic("chart-cpu")
				})

				sess.Patch("chart-heap", func(s RealtimeState) (RealtimeState, node.Node) {
					s.HeapMB = s.HeapMB.With(appendMetric(s.HeapMB.Val, heap))
					return s, div.New(
						chartDiv("memoise-heap", "Heap (MB)", "#5470c6", toLineData(s.HeapMB.Val)),
					).Dynamic("chart-heap")
				})

				sess.Patch("chart-goroutines", func(s RealtimeState) (RealtimeState, node.Node) {
					s.Goroutines = s.Goroutines.With(appendMetricInt(s.Goroutines.Val, goroutines))
					return s, div.New(
						chartDiv("memoise-goroutines", "Goroutines", "#91cc75", intsToLineData(s.Goroutines.Val)),
					).Dynamic("chart-goroutines")
				})
			}
		}
	})
}

func appendMetric(data []float64, v float64) []float64 {
	data = append(data, v)
	if len(data) > maxDataPoints {
		data = data[len(data)-maxDataPoints:]
	}
	return data
}

func appendMetricInt(data []int, v int) []int {
	data = append(data, v)
	if len(data) > maxDataPoints {
		data = data[len(data)-maxDataPoints:]
	}
	return data
}
