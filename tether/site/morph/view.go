package morph

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the full-page morph demo. Every element intentionally
// omits .Dynamic() to demonstrate the morph fallback path.
func Render(s State) node.Node {
	counterStr := strconv.Itoa(s.Counter)
	return cpage.New(
		panel.Card("Counter Without Dynamic Keys",
			"Click + and - to change the counter. This page has no Dynamic keys, "+
				"so the differ finds nothing to patch. The framework falls back to a "+
				"full-page morph - it sends the entire rendered HTML and the client-side "+
				"idiomorph library diffs the whole DOM. Check the server logs: you will "+
				"see 'structural=true' instead of targeted patches.",
			"Full morph fallback", panel.AllTransports,
			layout.Row(
				button.DecrementAction("morph.decrement",
					bind.EventData("count", counterStr),
				),
				button.IncrementAction("morph.increment",
					bind.EventData("count", counterStr),
				),
				// Intentionally no .Dynamic() - this forces the full-page
				// morph path instead of targeted patches.
				span.Text("Count: "+counterStr),
			),
		),

		panel.Card("Why This Matters",
			"Compare this page with the State & Rendering page, which uses "+
				"Dynamic keys. Both counters work, but the Dynamic-keyed version "+
				"sends only the changed element (~50 bytes) while this page sends "+
				"the entire rendered page (~5 KB) on every click. For a single "+
				"counter the difference is negligible, but for complex pages with "+
				"many dynamic elements, targeted patches are significantly more "+
				"efficient.",
			"", panel.AllTransports,
			hint.Text("Open the browser DevTools Network tab (WS frames) "+
				"and compare the payload sizes between this page and State & Rendering."),
		),
	)
}
