package selection

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the selection demo page.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Multi-Select",
			"Click an item to select it. Ctrl+click to toggle. "+
				"Shift+click to select a range. Then click Show Selected "+
				"to send the selected IDs to the server. Selection is "+
				"purely client-side until you collect it.",
			"bind.Selectable · bind.CollectSelected", panel.AllTransports,
			layout.Stack(
				bind.Apply(
					div.New(itemList()...).ID("select-list").Class("item-list"),
					bind.Selectable(),
				),
				button.PrimaryAction("Show Selected", "selection.action",
					bind.CollectSelected("#select-list"),
				),
				selectionResult(s.Result),
			),
		),
	)
}

// itemList generates a list of selectable items.
func itemList() []node.Node {
	names := []string{"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf", "Hotel"}
	items := make([]node.Node, len(names))
	for i, name := range names {
		items[i] = bind.Apply(
			div.New(span.Text(name)).Class("list-item"),
			bind.EventData("id", fmt.Sprintf("%d", i+1)),
		)
	}
	return items
}

// selectionResult renders the server's response.
func selectionResult(val string) node.Node {
	if val == "" {
		return hint.Text("Select items and click the button.")
	}
	return span.Text(val).Class("result-block").Dynamic("selection-result")
}
