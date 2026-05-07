package timer

// State is the per-session state for the timer demo.
type State struct {
	// OnlineCount tracks connected sessions for the header badge.
	OnlineCount int
	// ElapsedRunning tracks whether the elapsed timer is ticking.
	ElapsedRunning bool
	// CountdownRunning tracks whether the countdown timer is ticking.
	CountdownRunning bool
	// StopwatchRunning tracks whether the stopwatch is ticking.
	StopwatchRunning bool
}
