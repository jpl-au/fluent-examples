package freeze

import (
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the freeze demo page.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Frozen Counter",
			"Click increment a few times, then disconnect (kill the WebSocket or navigate away briefly). When you reconnect, the counter is restored from the SessionStore - the server released all session memory during the disconnect.",
			"FreezeWithConnect · SessionStore · OnConnect", panel.WS,
			layout.Stack(
				layout.Row(
					button.IncrementAction("freeze.increment"),
					button.DecrementAction("freeze.decrement"),
				),
				span.Textf("Count: %d", s.Count).Dynamic("count"),
				hint.Text("Disconnect and reconnect - the count survives."),
			),
		),
	)
}
