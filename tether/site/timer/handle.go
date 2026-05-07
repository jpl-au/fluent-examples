package timer

import tether "github.com/jpl-au/tether"

// Handle processes timer control events. The server never ticks the
// timers directly - it only pushes signal changes that tell the
// client-side timer to start, pause, or reset. Each timer has a
// "name.running" boolean signal and a "name" value signal.
//
// Toggle actions flip the running state and push the signal. Reset
// actions clear both the value and running signals in a single batch
// via sess.Signals so the client sees a consistent state.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {

	// Elapsed timer: count up from zero with default 1s precision.
	case "timer.elapsed.toggle":
		s.ElapsedRunning = !s.ElapsedRunning
		sess.Signal("elapsed.running", s.ElapsedRunning)
	case "timer.elapsed.reset":
		s.ElapsedRunning = false
		sess.Signals(map[string]any{
			"elapsed":         0,
			"elapsed.running": false,
		})

	// Countdown timer: count down from 30s, fires "timer.countdown.expired"
	// when the client-side timer reaches zero via TimerOnComplete.
	case "timer.countdown.toggle":
		s.CountdownRunning = !s.CountdownRunning
		sess.Signal("countdown.running", s.CountdownRunning)
	case "timer.countdown.reset":
		s.CountdownRunning = false
		sess.Signals(map[string]any{
			"countdown":         30,
			"countdown.running": false,
		})
	case "timer.countdown.expired":
		// Fired by the client when the countdown reaches zero.
		// The client has already stopped the timer - the server
		// just needs to update its own state and react.
		s.CountdownRunning = false
		sess.Toast("Countdown finished!")

	// Stopwatch: count up with 100ms precision, mm:ss.S format.
	case "timer.stopwatch.toggle":
		s.StopwatchRunning = !s.StopwatchRunning
		sess.Signal("stopwatch.running", s.StopwatchRunning)
	case "timer.stopwatch.reset":
		s.StopwatchRunning = false
		sess.Signals(map[string]any{
			"stopwatch":         0,
			"stopwatch.running": false,
		})
	}
	return s
}
