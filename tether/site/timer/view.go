package timer

import (
	"time"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the client-side timer demo page. Each card
// demonstrates a different timer configuration, all controlled by
// the server pushing signals while the client handles the ticking.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Elapsed Timer",
			"Click Start to begin counting up from zero. The timer ticks entirely in the browser "+
				"with no server messages per tick - contrast this with the Live Updates page where "+
				"the uptime ticker uses a server-side goroutine that pushes a signal every second. "+
				"Here the server sends a single signal to start, and the client runs independently "+
				"until told to pause or reset. The display format adapts automatically: seconds only "+
				"below a minute, mm:ss below an hour, hh:mm:ss beyond that.",
			"bind.Timer", panel.WS,
			layout.Stack(
				layout.Row(
					toggleButton("Start", "Pause", s.ElapsedRunning, "timer.elapsed.toggle"),
					button.SecondaryAction("Reset", "timer.elapsed.reset"),
				),
				bind.Apply(
					span.Text("0").Class("timer-display"),
					bind.Timer("elapsed"),
				),
			),
		),

		panel.Card(
			"Countdown Timer",
			"A 30-second countdown that fires an event back to the server when it reaches zero. "+
				"The server handles the completion by showing a toast notification. This pattern is "+
				"useful for quiz timers, session timeouts, or auction clocks where the server needs "+
				"to know when time expires but does not need to track every tick. No polling, no "+
				"background goroutine - the client counts down autonomously and notifies the server "+
				"only at the moment of completion via TimerOnComplete.",
			"bind.Timer · bind.Countdown · bind.TimerOnComplete", panel.WS,
			layout.Stack(
				layout.Row(
					toggleButton("Start", "Pause", s.CountdownRunning, "timer.countdown.toggle"),
					button.SecondaryAction("Reset", "timer.countdown.reset"),
				),
				bind.Apply(
					span.Text("30").Class("timer-display"),
					bind.Timer("countdown"),
					bind.Countdown(30*time.Second),
					bind.TimerOnComplete("timer.countdown.expired"),
				),
			),
		),

		panel.Card(
			"Stopwatch",
			"A sub-second timer with 100ms precision and an explicit mm:ss.S format. "+
				"Doing this with a server-side goroutine would mean 10 WebSocket messages per second "+
				"per session - wasteful for something purely presentational. With a client-side timer "+
				"the server load is identical to the one-second elapsed timer above: one signal to "+
				"start, one to stop. TimerPrecision controls the tick interval and TimerFormat sets "+
				"the display pattern.",
			"bind.Timer · bind.TimerPrecision · bind.TimerFormat", panel.WS,
			layout.Stack(
				layout.Row(
					toggleButton("Start", "Pause", s.StopwatchRunning, "timer.stopwatch.toggle"),
					button.SecondaryAction("Reset", "timer.stopwatch.reset"),
				),
				bind.Apply(
					span.Text("00:00.0").Class("timer-display"),
					bind.Timer("stopwatch"),
					bind.TimerPrecision(100*time.Millisecond),
					bind.TimerFormat("mm:ss.S"),
				),
			),
		),
	)
}

// toggleButton renders a primary action button whose label reflects
// the current running state.
func toggleButton(startLabel, pauseLabel string, running bool, action string) node.Node {
	label := startLabel
	if running {
		label = pauseLabel
	}
	return button.PrimaryAction(label, action)
}
