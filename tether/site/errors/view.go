package errors

import (
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the error boundary page demonstrating tether.Catch.
func Render(_ State) node.Node {
	return cpage.New(
		panel.Card("Error Boundary",
			"The component below deliberately panics during render. tether.Catch catches the panic and shows the fallback content instead, so one broken component does not take down the rest of the page.",
			"tether.Catch", panel.AllTransports,
			tether.Catch(func() node.Node {
				return layout.Container(p.New().Text("This component panics during render."), triggerPanic())
			}, panel.SignalSuccess("Caught by error boundary - component recovered gracefully.")),
		),
	)
}

// triggerPanic triggers a panic inside a tether.Catch boundary
// to demonstrate that the error boundary recovers gracefully.
func triggerPanic() node.Node {
	panic("simulated failure to demonstrate tether.Catch error boundary")
}
