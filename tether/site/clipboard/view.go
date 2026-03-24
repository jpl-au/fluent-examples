package clipboard

import (
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the clipboard demo page.
func Render(_ State) node.Node {
	return page.New(
		panel.Card(
			"Copy to Clipboard",
			"Click the Copy button to copy the text below to your clipboard. "+
				"bind.CopyToClipboard runs entirely on the client - no server round-trip.",
			"bind.CopyToClipboard", panel.HTTP,
			layout.Stack(
				span.Text("tether-secret-key-abc123").ID("copy-source").Class("result-block"),
				button.Primary("Copy", bind.CopyToClipboard("#copy-source")),
			),
		),
	)
}
