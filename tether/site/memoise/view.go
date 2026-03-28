package memoise

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/table"
	"github.com/jpl-au/fluent/html5/tbody"
	"github.com/jpl-au/fluent/html5/td"
	"github.com/jpl-au/fluent/html5/th"
	"github.com/jpl-au/fluent/html5/thead"
	"github.com/jpl-au/fluent/html5/tr"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the memoisation demo page with two regions:
//   - A memoised table that only re-renders when items change
//   - A plain counter that re-renders on every event
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"What is Memoisation?",
			"Every state change normally re-renders the entire page and "+
				"diffs every Dynamic region. For pages with expensive "+
				"regions (large tables, complex charts), this wastes work "+
				"when only a small part changed. Memoisation lets you mark "+
				"expensive regions so the framework skips them entirely "+
				"when their data has not changed. The render closure never "+
				"runs and no HTML is produced for skipped regions. Cheap "+
				"regions use plain Dynamic keys and always re-render. You "+
				"choose per region.",
			"node.Memoise · tether.Versioned · StatefulConfig.Memoise", panel.WS|panel.SSE,
		),

		panel.Card(
			"Memoised Table vs Plain Counter",
			"The table below is wrapped in node.Memoise with a "+
				"tether.Versioned key. The counter uses a plain Dynamic "+
				"key. Incrementing the counter leaves Items.Version() "+
				"unchanged, so the Memoiser skips the table entirely. "+
				"Adding an item calls With(), which increments the "+
				"version and triggers a re-render of the table.",
			"tether.Versioned.With · node.Memoise · .Dynamic", panel.WS|panel.SSE,
			layout.Stack(
				hint.Text("Click Increment - the counter updates but "+
					"the table is skipped (memoiser hit). Click Add Item "+
					"- the table re-renders with the new row (memoiser miss)."),

				layout.Row(
					button.PrimaryAction("Increment Counter", "memoise.increment"),
					button.SecondaryAction("Add Item", "memoise.add-item"),
				),

				// Counter - plain Dynamic, always re-renders.
				div.New(
					span.Text("Counter: "),
					span.Text(strconv.Itoa(s.Count)).ID("memoise-count"),
				).Dynamic("counter").Class("result-block"),

				// Table - memoised. The closure only runs when
				// Items.Version() changes (i.e. when With was called).
				div.New(
					node.Memoise(s.Items.Version(), func() node.Node {
						return renderTable(s.Items.Val)
					}),
				).Dynamic("items"),
			),
		),
	)
}

// renderTable builds the HTML table from the items slice. This is
// the "expensive" render function that the memoiser skips when the version
// matches.
func renderTable(items []Item) node.Node {
	rows := make([]node.Node, len(items))
	for i, item := range items {
		rows[i] = tr.New(
			td.Text(strconv.Itoa(item.ID)),
			td.Text(item.Name),
		).ID("row-" + strconv.Itoa(item.ID))
	}
	return table.New(
		thead.New(tr.New(th.Text("ID"), th.Text("Name"))),
		tbody.New(rows...),
	).Class("demo-table").ID("memoise-table")
}
