package notifications

import (
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/columns"
	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/components/simple/result"
	"github.com/jpl-au/fluent-examples/tether/components/simple/spinner"
)

// Render builds the notifications page, demonstrating Toast, Flash,
// Announce, and Signal side effects.
func Render(_ State) node.Node {
	return page.New(
		panel.Card(
			"Toast",
			"Server-initiated notification pushed to the browser. "+
				"The server calls sess.Toast() at any point - in an event handler, a background goroutine, or an upload callback - "+
				"and the client displays a notification in the bottom-right corner that fades out automatically.",
			"sess.Toast", panel.WS|panel.SSE,
			button.PrimaryAction("Show Toast", "notify.toast"),
		),

		panel.Card(
			"Flash",
			"Temporary text replacement pushed from the server. "+
				"The server targets a DOM element by CSS selector and replaces its content for a few seconds, then reverts. "+
				"Unlike re-rendering the page, Flash only touches a single element - useful for inline feedback (e.g. 'Saved!') without a full state change.",
			"sess.Flash", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Flash Message", "notify.flash"),
				p.Text("Waiting for flash...").ID("flash-target"),
			),
		),

		panel.Card(
			"Announce",
			"Server-pushed accessibility announcement. "+
				"The server calls sess.Announce() to post a message to a hidden ARIA live region - "+
				"screen readers speak it aloud without any visual change to the page. "+
				"Essential for dynamic updates (toasts, live feeds, state changes) that are invisible to assistive technology unless explicitly announced. "+
				"This demo also fires a toast and echoes the text below so you can verify it without a screen reader.",
			"sess.Announce", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Announce", "notify.announce"),
				bind.Apply(
					p.Text(""),
					bind.BindText("notify.announced"),
					bind.BindShow("notify.announced"),
				),
			),
		),

		panel.Card(
			"Flash vs Signal",
			"Two ways to show temporary inline feedback, side by side. "+
				"The left button uses sess.Flash - a one-liner that targets a DOM element by CSS selector and auto-reverts. "+
				"The right button uses sess.Signal with bind.BindShow - no selector, fully decoupled, and the server clears the signal after the same delay. "+
				"Flash is faster to wire up; signals are better for reusable components where the server shouldn't know about DOM IDs.",
			"sess.Flash · sess.Signal · bind.BindShow", panel.WS|panel.SSE,
			columns.New(
				layout.Stack(
					result.Label("Selector approach (sess.Flash)"),
					button.PrimaryAction("Flash", "notify.flash-compare"),
					p.Text("Waiting for flash...").ID("flash-compare-target"),
				),
				layout.Stack(
					result.Label("Signal approach (sess.Signal)"),
					button.PrimaryAction("Signal", "notify.signal-flash"),
					bind.Apply(p.Text("Saved!"),
						bind.BindShow("notify.saved"),
					),
					bind.Apply(p.Text("Waiting for signal..."),
						bind.BindHide("notify.saved"),
					),
				),
			),
		),

		panel.Card(
			"Loading Indicator vs Signal",
			"Two ways to show a loading spinner during a slow operation (3 seconds simulated). "+
				"The left button uses bind.Indicator - it targets a spinner element by CSS selector and toggles it automatically while the request is in flight. "+
				"The right button uses bind.Optimistic to flip a signal immediately (showing the spinner), then the server clears it when done. "+
				"Indicator is simpler; signals avoid selector coupling.",
			"bind.Indicator · bind.Optimistic · bind.BindShow", panel.AllTransports,
			columns.New(
				layout.Stack(
					result.Label("Selector approach (bind.Indicator)"),
					button.PrimaryAction("Load", "notify.indicator",
						bind.Indicator("#notify-spinner"),
					),
					spinner.New().ID("notify-spinner"),
				),
				layout.Stack(
					result.Label("Signal approach (bind.Optimistic)"),
					button.PrimaryAction("Load", "notify.signal-indicator",
						bind.Optimistic("notify.loading", "true"),
					),
					bind.Apply(spinner.New(), bind.BindShow("notify.loading")),
				),
			),
		),
	)
}
