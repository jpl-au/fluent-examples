package valuestore

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	btn "github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the value store page with three demo cards:
// read-modify-write via Update, direct set via Store, and local
// versus shared state comparison.
func Render(s State) node.Node {
	shared := strconv.Itoa(s.SharedCount)
	local := strconv.Itoa(s.LocalCount)
	return page.New(
		panel.Card(
			"Update (Read-Modify-Write)",
			"Update takes a function that reads the current value and returns "+
				"a new one, serialised under a mutex. All connected sessions see "+
				"the change via WatchValue.",
			"tether.Value · Update", panel.WS|panel.SSE,
			layout.Row(
				btn.IncrementAction("value.increment"),
				btn.DecrementAction("value.decrement"),
				span.Text("Count: "+shared).Dynamic("update-count"),
			),
		),

		panel.Card(
			"Store (Direct Set)",
			"Store writes a new value directly without reading the old one. "+
				"Use it when the new value does not depend on the previous state.",
			"tether.Value · Store", panel.WS|panel.SSE,
			bind.Apply(field.Inline(
				field.TextWithID("set-value-input", "value", "Enter a number…"),
				btn.Submit("Set Value"),
				btn.DangerAction("Reset to 0", "value.reset"),
			),
				bind.OnSubmit("value.set"),
				bind.Reset(),
			),
			span.Text("Count: "+shared).Dynamic("store-count"),
		),

		panel.Card(
			"Local vs Shared",
			"Normal state fields are per-session. tether.Value is for "+
				"cross-session coordination - changes propagate to every "+
				"connected client via WatchValue.",
			"tether.WatchValue", panel.WS|panel.SSE,
			layout.Row(
				btn.PrimaryAction("Increment Local", "value.local-inc"),
				span.Text("Shared: "+shared).Dynamic("shared-count"),
				span.Text("Local: "+local).Dynamic("local-count"),
			),
			hint.Text(
				"Open two tabs: the shared counter stays in sync, "+
					"but each tab has its own local counter."),
		),
	)
}
