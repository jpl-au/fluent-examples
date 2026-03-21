// Package realtime provides a standalone real-time dashboard
// demonstrating live Go runtime metrics pushed from the server
// every second via Session.Go. CPU usage, heap allocation, and
// goroutine count are rendered as go-echarts line charts that
// update automatically through the normal diff pipeline.
package realtime
