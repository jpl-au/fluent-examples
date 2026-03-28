package memoise

import (
	"log/slog"
	"syscall"
)

// cpuTime returns the total process CPU time (user + kernel) in
// seconds using GetProcessTimes.
func cpuTime() float64 {
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		slog.Debug("failed to get current process handle", "error", err)
		return 0
	}
	var creationTime, exitTime, kernelTime, userTime syscall.Filetime
	if err := syscall.GetProcessTimes(h, &creationTime, &exitTime, &kernelTime, &userTime); err != nil {
		slog.Debug("failed to read CPU usage", "error", err)
		return 0
	}
	user := float64(userTime.Nanoseconds()) / 1e9
	kernel := float64(kernelTime.Nanoseconds()) / 1e9
	return user + kernel
}
