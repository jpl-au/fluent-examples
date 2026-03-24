package scroll

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the scroll demo page.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Client-Side ScrollTo",
			"Click the button to scroll the target element into view. "+
				"bind.ScrollTo runs entirely on the client - no server round-trip.",
			"bind.ScrollTo", panel.WS,
			layout.Stack(
				button.Primary("Scroll to Target", bind.ScrollTo("#scroll-target")),
				scrollSpacer(),
				span.Text("Scroll target reached!").ID("scroll-target").Class("result-block"),
			),
		),

		panel.Card(
			"Server-Side ScrollTo",
			"Click the button - the server calls sess.ScrollTo to "+
				"scroll the same target element into view. The scroll "+
				"command travels over the WebSocket.",
			"sess.ScrollTo", panel.WS,
			button.PrimaryAction("Server Scroll", "scroll.server-scroll"),
		),

		panel.Card(
			"PreserveScroll",
			"Scroll inside the list below, then click Add Items. "+
				"The list re-renders with more items but the scroll "+
				"position is preserved because bind.PreserveScroll "+
				"saves and restores scrollTop across morphs.",
			"bind.PreserveScroll", panel.WS,
			layout.Stack(
				button.PrimaryAction("Add 5 Items", "scroll.add"),
				preserveList(s.Items),
			),
		),
	)
}

// scrollSpacer creates vertical space so the target is off-screen.
func scrollSpacer() node.Node {
	return div.New().SetData("testid", "spacer").
		Class("demo-description").
		Style("height:20rem")
}

// preserveList renders a scrollable list with PreserveScroll. The
// scroll position survives re-renders when items are added.
func preserveList(n int) node.Node {
	items := make([]node.Node, n)
	for i := range n {
		items[i] = div.New(
			span.Text(fmt.Sprintf("Item %d", i+1)),
		).ID(fmt.Sprintf("item-%d", i+1)).Class("list-item")
	}
	return bind.Apply(
		div.New(items...).Class("item-list viewport-list").Style("max-height:12rem;overflow-y:auto"),
		bind.PreserveScroll(),
	).Dynamic("preserve-list")
}
