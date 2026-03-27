package patch

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/datalist"
	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the patch demo page.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"When to Use Targeted Updates",
			"Use sess.Patch when you know exactly which Dynamic region "+
				"changed and want to skip the full render pipeline. A "+
				"full Update re-renders every Dynamic region and diffs "+
				"them all. Patch re-renders only the targeted key - over "+
				"1,000x faster for a single key out of many. Use it "+
				"from timers, broadcast callbacks, and Go goroutines "+
				"where the mutation is known and targeted.",
			"sess.Patch · DiffKey", panel.WS|panel.SSE,
		),

		panel.Card(
			"Live Targeted Updates",
			"A background timer increments one counter every 500ms "+
				"using sess.Patch. Only the active row is re-rendered "+
				"and sent to the client. The other "+
				strconv.Itoa(numCounters-1)+" rows are untouched. "+
				"Watch the render logs to see patch durations vs full "+
				"render durations.",
			"sess.Patch · sess.Go · timer", panel.WS|panel.SSE,
			layout.Stack(
				hint.Text("Watch the counters increment one at a time. "+
					"Each update patches only the active row. Click "+
					"Reset to zero all counters via a full Update."),

				button.PrimaryAction("Reset All", "patch.reset"),

				renderCounterList(s),
			),
		),
	)
}

// renderCounterList builds the full list of counters with Dynamic
// keys for targeted patching.
func renderCounterList(s State) node.Node {
	rows := make([]node.Node, numCounters)
	for i := range numCounters {
		rows[i] = RenderCounter(i, s.Counters[i])
	}
	return div.New(rows...).Class("data-list")
}

// RenderCounter renders a single counter row. Exported so the Patch
// closure can call it to produce the targeted subtree.
func RenderCounter(index, value int) node.Node {
	id := strconv.Itoa(index)
	return datalist.RowDynamic(
		"counter-"+id,
		"Counter "+id,
		strconv.Itoa(value),
	)
}
