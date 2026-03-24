// Package card renders an individual kanban card within a column.
// Each card is draggable (for moving between columns) and clickable
// (for opening the detail view).
package card

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New renders a draggable kanban card. Clicking the card body opens
// the detail view. The card carries its ID via EventData so the
// drop target knows which card was moved.
func New(c store.Card) node.Node {
	return bind.Apply(
		div.New(
			bind.Apply(
				div.New(
					span.New().Class("card-title").Text(c.Title),
					desc(c.Description),
				).Class("card-body"),
				bind.OnClick("card.select"),
				bind.EventData("id", c.ID),
			),
		).Class("card"),
		bind.Draggable(),
		bind.EventData("id", c.ID),
	).Dynamic(c.ID)
}

// desc renders a truncated description, or nil if empty.
func desc(s string) node.Node {
	if s == "" {
		return nil
	}
	if len(s) > 80 {
		s = s[:77] + "..."
	}
	return p.New().Class("card-desc").Text(s)
}
