package dragdrop

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/columns"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the drag-and-drop demo page.
func Render(items []Item) node.Node {
	return page.New(
		panel.Card(
			"Drag and Drop",
			"Drag items between the two zones. All connected tabs update "+
				"instantly via Group.Broadcast. The tether-drag-and-drop.js "+
				"extension is auto-included when any element uses bind.Draggable.",
			"bind.Draggable · bind.DropTarget", panel.WS,
			columns.New(
				zone("left", "Zone A", items),
				zone("right", "Zone B", items),
			),
		),
	)
}

// zone renders a drop target with the items that belong to it.
func zone(id, title string, items []Item) node.Node {
	var children []node.Node
	children = append(children, span.Text(title).Class("demo-title"))

	for _, item := range items {
		if item.Zone != id {
			continue
		}
		children = append(children, draggable(item))
	}

	if len(children) == 1 {
		children = append(children, span.Text("Drop items here").Class("hint"))
	}

	return bind.Apply(
		div.New(children...).Class("layout-stack drop-zone"),
		bind.DropTarget("item.drop"),
		bind.EventData("zone", id),
	).Dynamic("zone-" + id)
}

// draggable renders a single draggable item.
func draggable(item Item) node.Node {
	return bind.Apply(
		div.New(span.Text(item.Name)).Class("signal-panel"),
		bind.Draggable(),
		bind.EventData("id", item.ID),
	).Dynamic("item-" + item.ID)
}
