package signals

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// RenderSSE builds the signals page for the SSE variant: a subset
// of the WS demos that work over Server-Sent Events.
func RenderSSE(s State) node.Node {
	return page.New(
		panel.Card(
			"Text",
			"Click the button to increment a server-side counter. The server pushes the new value as a signal over the SSE connection, and the client updates the text without re-rendering the page. Identical to the WebSocket version - only the transport differs.",
			"bind.Text", panel.WS|panel.SSE,
			layout.Row(
				button.PrimaryAction("Increment Server Counter", "signals.increment"),
				bind.Apply(span.Text(strconv.Itoa(s.Counter)),
					bind.Text("signals.counter"),
				),
			),
		),

		panel.Card(
			"Show / Hide",
			"Click Toggle Visibility - the server pushes a boolean signal via SSE that shows or hides elements instantly on the client. No page re-render happens; the client toggles CSS display directly.",
			"bind.Show · bind.Hide", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Toggle Visibility", "signals.toggle-panel"),
				bind.Apply(layout.Container(
					panel.SignalText("This panel is visible when the signal is true."),
				), bind.Show("signals.panel_visible")),
				bind.Apply(span.Text("The panel is hidden."),
					bind.Hide("signals.panel_visible"),
				),
			),
		),

		panel.Card(
			"SetSignal (Client-Side)",
			"Click any colour button - the text updates instantly with no server round-trip. SetSignal writes directly to the client-side signal store. No POST request is sent.",
			"bind.SetSignal", panel.AllTransports,
			layout.Row(
				button.Primary("Set to Red", bind.SetSignal("signals.colour", "red")),
				button.Primary("Set to Blue", bind.SetSignal("signals.colour", "blue")),
				button.Primary("Set to Green", bind.SetSignal("signals.colour", "green")),
				bind.Apply(span.Text("none").Dynamic("colour-display"), bind.Text("signals.colour")),
			),
		),

		panel.Card(
			"ToggleSignal (Client-Side)",
			"Click the button to flip a boolean signal on the client. Combined with Show, the panel appears and disappears with no server round-trip.",
			"bind.ToggleSignal · bind.Show", panel.AllTransports,
			layout.Stack(
				button.Primary("Toggle Panel", bind.ToggleSignal("signals.toggle_demo")),
				bind.Apply(panel.SignalText("Toggled on!"), bind.Show("signals.toggle_demo")),
			),
		),

		panel.Card(
			"Optimistic Updates",
			"Click Like - the 'Liked!' text appears instantly before the server's POST response arrives. The client sets the signal optimistically on click for immediate feedback.",
			"bind.Optimistic", panel.AllTransports,
			layout.Row(
				button.PrimaryAction("Like", "signals.like",
					bind.Optimistic("signals.liked", "true"),
				),
				bind.Apply(span.Text("Liked!"), bind.Show("signals.liked")),
			),
		),
	)
}
