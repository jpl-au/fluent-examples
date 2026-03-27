//go:build !windows

package memo

import (
	"log/slog"
	"syscall"
)

// cpuTime returns the total process CPU time (user + system) in
// seconds using syscall.Getrusage.
func cpuTime() float64 {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		slog.Debug("failed to read CPU usage", "error", err)
		return 0
	}
	user := float64(ru.Utime.Sec) + float64(ru.Utime.Usec)/1e6
	sys := float64(ru.Stime.Sec) + float64(ru.Stime.Usec)/1e6
	return user + sys
}
