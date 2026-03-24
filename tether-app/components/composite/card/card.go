// Package card renders an individual kanban card within a column.
// Each card is draggable (for moving between columns) and clickable
// (for opening the detail view). Shows the title, description
// snippet, creator, timestamp, and who is currently viewing.
package card

import (
	"strings"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New renders a draggable kanban card. The viewers parameter lists
// other users currently viewing this card (empty when nobody is).
func New(c store.Card, viewers ...string) node.Node {
	return bind.Apply(
		div.New(
			bind.Apply(
				div.New(
					span.Text(c.Title).Class("card-title"),
					desc(c.Description),
					div.New(
						span.Text(c.CreatedBy).Class("card-author"),
						span.Text(store.TimeAgo(c.CreatedAt)).Class("card-time"),
					).Class("card-meta"),
					viewing(viewers),
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
	return p.Text(s).Class("card-desc")
}

// viewing renders a small indicator showing who is viewing this card.
func viewing(names []string) node.Node {
	if len(names) == 0 {
		return nil
	}
	if len(names) == 1 {
		return span.Text(names[0] + " is viewing this").Class("card-viewing")
	}
	return span.Text(strings.Join(names, ", ") + " are viewing this").Class("card-viewing")
}
