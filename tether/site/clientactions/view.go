package clientactions

import (
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the client-side actions demo page.
func Render(_ State) node.Node {
	return page.New(
		panel.Card(
			"Copy to Clipboard with FlashText",
			"Click Copy to copy the API key to your clipboard. "+
				"bind.CopyToClipboard handles the copy entirely on the "+
				"client. bind.FlashText swaps the button label to "+
				"\"Copied!\" for two seconds as confirmation, then "+
				"reverts. No server round-trip for either action.",
			"bind.CopyToClipboard · bind.FlashText", panel.AllTransports,
			layout.Stack(
				span.Text("sk_live_abc123def456").ID("copy-source").Class("result-block"),
				button.Primary("Copy",
					bind.CopyToClipboard("#copy-source"),
					bind.FlashText("Copied!"),
				),
			),
		),

		panel.Card(
			"Copy to Clipboard with FlashClass",
			"FlashClass temporarily adds a CSS class to the element "+
				"after a client-side action succeeds. Use it for richer "+
				"feedback like colour changes, icon swaps, or "+
				"animations. Here the button turns green on copy. "+
				"The class is removed after the flash duration.",
			"bind.CopyToClipboard · bind.FlashClass", panel.AllTransports,
			layout.Stack(
				span.Text("https://example.com/share/abc123").ID("link-source").Class("result-block"),
				button.Primary("Copy Link",
					bind.CopyToClipboard("#link-source"),
					bind.FlashClass("btn-flashed"),
				),
			),
		),
	)
}
