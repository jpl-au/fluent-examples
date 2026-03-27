package windowing

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
	"github.com/jpl-au/fluent-examples/tether/components/simple/result"
)

// Render builds the windowing demo page.
func Render(s State) node.Node {
	end := min(s.ScrollOffset+pageSize, len(s.Items))

	return page.New(
		panel.Card(
			"When to Use Windowing",
			"Windowing reduces render and diff cost by only building "+
				"DOM nodes for the visible portion of a large dataset. "+
				"The full dataset can live in memory, in a database, or "+
				"behind an API - the render function only touches the "+
				"current page. This means fewer nodes to diff, less HTML "+
				"on the wire, and a smaller browser DOM. In production, "+
				"combine windowing with database pagination so the server "+
				"fetches only the current page from storage on each "+
				"navigation.",
			"tether/window · window.New · window.Config", panel.WS|panel.SSE,
		),

		panel.Card(
			"Paginated List",
			"This demo holds "+strconv.Itoa(len(s.Items))+" items in "+
				"memory, but the render function only builds nodes for "+
				"the current page of "+strconv.Itoa(pageSize)+". The "+
				"differ walks "+strconv.Itoa(pageSize)+" nodes instead "+
				"of "+strconv.Itoa(len(s.Items))+". The URL updates on "+
				"each page change so refreshing or sharing the link "+
				"lands on the same page.",
			"Slice rendering · Dynamic · ReplaceURL · OnNavigate", panel.WS|panel.SSE,
			layout.Stack(
				hint.Text("Navigate with the buttons below. Open your "+
					"browser's inspector to confirm only "+
					strconv.Itoa(pageSize)+" rows exist in the DOM."),

				layout.Row(
					button.PrimaryAction("First", "window.first"),
					button.PrimaryAction("Previous", "window.prev"),
					button.PrimaryAction("Next", "window.next"),
					button.PrimaryAction("Last", "window.last"),
				),

				result.BlockDynamic("position",
					"Showing rows "+strconv.Itoa(s.ScrollOffset+1)+
						" - "+strconv.Itoa(end)+
						" of "+strconv.Itoa(len(s.Items))),

				renderVisibleRows(s.Items, s.ScrollOffset, end),
			),
		),
	)
}

// renderVisibleRows builds a datalist from the visible slice.
func renderVisibleRows(items []Item, start, end int) node.Node {
	rows := make([]node.Node, end-start)
	for i := start; i < end; i++ {
		rows[i-start] = datalist.RowWithID(
			"wrow-"+strconv.Itoa(items[i].ID),
			strconv.Itoa(items[i].ID),
			items[i].Name,
		)
	}
	return div.New(rows...).Dynamic("items").Class("data-list")
}
