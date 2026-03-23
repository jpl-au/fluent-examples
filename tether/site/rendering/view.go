package rendering

import (
	"strconv"
	"strings"

	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/list"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the state and rendering page, demonstrating how
// Render produces a node tree from state: counters, dynamic lists,
// component routing, and nested components.
func Render(s State) node.Node {
	counterStr := strconv.Itoa(s.Counter)

	return cpage.New(
		panel.Card("Dynamic Keys",
			"Click + and - to change the counter. Each element has a unique key (via Dynamic) so the framework updates only the text that changed rather than replacing everything on the page. This keeps transitions smooth and avoids unnecessary DOM work.",
			"Dynamic", panel.AllTransports,
			layout.Row(
				button.DecrementAction("rendering.decrement",
					bind.EventData("count", counterStr),
				),
				button.IncrementAction("rendering.increment",
					bind.EventData("count", counterStr),
				),
				span.Text("Count: "+counterStr).Dynamic("rendering-counter"),
			),
		),

		panel.Card("Dynamic List",
			"Click Add Item to grow the list and Remove Last to shrink it. Each list item has a unique key, so the framework inserts or removes just that item instead of rebuilding the entire list every time.",
			"Dynamic", panel.AllTransports,
			layout.Stack(
				layout.Row(
					button.PrimaryAction("Add", "rendering.add-item",
						bind.EventData("items", encodeItems(s.Items)),
					),
					button.SecondaryAction("Remove", "rendering.remove-item",
						bind.EventData("items", encodeItems(s.Items)),
					),
				),
				layout.Container(itemList(s.Items)).Dynamic("item-list"),
			),
		),

		panel.Card("Component",
			"A self-contained counter that implements tether.Component. The component owns its own Render and Handle methods - the page handler delegates to it via tether.RouteTyped, which strips the prefix and preserves the concrete type. The component has no knowledge of the page's state type.",
			"tether.Component · tether.RouteTyped", panel.AllTransports,
			bind.Apply(layout.Container(s.Counter2.Render()), bind.Prefix("counter")).Dynamic("component-counter"),
		),

		panel.Card("Nested Components",
			"Two independent counters nested inside a parent group. Each counter manages "+
				"its own state - click + and - on either one independently. The parent's "+
				"Reset All button reaches into both children and zeros them, showing how a "+
				"parent component can intercept actions and modify child state. The combined "+
				"total demonstrates the parent observing child state during render.",
			"tether.RouteTyped · nesting", panel.AllTransports,
			bind.Apply(layout.Container(s.Group.Render()), bind.Prefix("group")).Dynamic("nested-group"),
		),
	)
}

// itemList renders the dynamic list demo - each item gets a stable
// Dynamic key so the differ can track additions and removals.
func itemList(items []string) node.Node {
	if len(items) == 0 {
		return hint.Text("No items yet. Click Add to create one.")
	}
	nodes := make([]node.Node, len(items))
	for i, item := range items {
		nodes[i] = li.New().Text(item).Dynamic("item-" + strconv.Itoa(i))
	}
	return list.New(nodes...)
}

// encodeItems serialises the item list into a pipe-delimited string
// for round-tripping through bind.EventData on the stateless page.
func encodeItems(items []string) string {
	return strings.Join(items, "|")
}
