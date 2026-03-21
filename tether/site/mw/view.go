package mw

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/components/simple/result"
)

// Render builds the middleware demo page, showing four cards that
// demonstrate custom middleware composition: execution order,
// request timing, guard short-circuiting, and post-processing.
func Render(s State) node.Node {
	return page.New(
		panel.Card("Middleware Chain",
			"Middleware wraps the handler in layers. The outermost middleware in the slice runs first on the way in and last on the way out, forming an onion-like call stack.",
			"tether.Middleware", panel.AllTransports,
			layout.Stack(
				button.PrimaryAction("Send Event", "mw.ping",
					bind.EventData("count", fmt.Sprintf("%d", s.EventCount)),
				),
				layout.Container(chainResult(s.ChainLog)).Dynamic("chain-result"),
			),
		),

		panel.Card("Request Timing",
			"A timing middleware records the wall-clock duration of the inner handler. The slow event sleeps briefly to make the difference visible.",
			"tether.Middleware", panel.AllTransports,
			layout.Stack(
				layout.Row(
					button.PrimaryAction("Fast Event", "mw.ping",
						bind.EventData("count", fmt.Sprintf("%d", s.EventCount)),
					),
					button.PrimaryAction("Slow Event", "mw.slow",
						bind.EventData("count", fmt.Sprintf("%d", s.EventCount)),
					),
				),
				layout.Container(timingResult(s.TimingResult)).Dynamic("timing-result"),
			),
		),

		panel.Card("Guard (Short-Circuit)",
			"Guard middleware can short-circuit the chain by returning early without calling the next handler. Use this pattern for authorisation checks, rate limiting, or action filtering.",
			"tether.Middleware", panel.AllTransports,
			layout.Stack(
				layout.Row(
					button.PrimaryAction("Allowed Action", "mw.ping",
						bind.EventData("count", fmt.Sprintf("%d", s.EventCount)),
					),
					button.DangerAction("Blocked Action", "mw.blocked",
						bind.EventData("count", fmt.Sprintf("%d", s.EventCount)),
					),
				),
				layout.Container(guardResult(s.BlockedResult)).Dynamic("guard-result"),
			),
		),

		panel.Card("Event Counting",
			"Post-processing middleware modifies the state after the handler returns. This counter increments on every event regardless of the action.",
			"tether.Middleware", panel.AllTransports,
			layout.Row(
				button.PrimaryAction("Click Me", "mw.ping",
					bind.EventData("count", fmt.Sprintf("%d", s.EventCount)),
				),
				span.Text("Events processed: "+fmt.Sprintf("%d", s.EventCount)),
			),
		),
	)
}

// chainResult renders the middleware execution log when present,
// or a hint paragraph when there is nothing to show yet.
func chainResult(log string) node.Node {
	if log == "" {
		return hint.Text("Click to see middleware execution order")
	}
	return layout.Container(
		result.Label("Execution order"),
		result.Block(log),
	)
}

// timingResult renders the measured handler duration, or a hint
// paragraph when no timing data is available.
func timingResult(r string) node.Node {
	if r == "" {
		return hint.Text("Click a button to see handler timing")
	}
	return layout.Container(
		result.Label("Duration"),
		result.Block(r),
	)
}

// guardResult renders the guard middleware outcome, or a hint
// paragraph when no guard action has been triggered.
func guardResult(r string) node.Node {
	if r == "" {
		return hint.Text("Click Blocked Action to see the guard in action")
	}
	return layout.Container(
		result.Label("Result"),
		result.Block(r),
	)
}
