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

// CardViewers holds presence information for a card on the board.
type CardViewers struct {
	Viewing []string
	Typing  []string
}

// New renders a draggable kanban card with optional presence info.
func New(c store.Card, v ...CardViewers) node.Node {
	var cv CardViewers
	if len(v) > 0 {
		cv = v[0]
	}
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
					presence(cv),
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

// presence renders typing and viewing indicators on the card.
// Users who are typing are excluded from the viewing list to avoid
// showing both "editing" and "viewing" for the same person.
func presence(cv CardViewers) node.Node {
	// Build a set of typing names for exclusion.
	typingSet := make(map[string]bool, len(cv.Typing))
	for _, n := range cv.Typing {
		typingSet[n] = true
	}

	// Viewers who are NOT typing.
	var viewOnly []string
	for _, n := range cv.Viewing {
		if !typingSet[n] {
			viewOnly = append(viewOnly, n)
		}
	}

	if len(cv.Typing) == 0 && len(viewOnly) == 0 {
		return nil
	}

	var nodes []node.Node
	if len(cv.Typing) > 0 {
		if len(cv.Typing) == 1 {
			nodes = append(nodes, span.Text(cv.Typing[0]+" is editing...").Class("card-typing"))
		} else {
			nodes = append(nodes, span.Text(strings.Join(cv.Typing, ", ")+" are editing...").Class("card-typing"))
		}
	}
	if len(viewOnly) > 0 {
		if len(viewOnly) == 1 {
			nodes = append(nodes, span.Text(viewOnly[0]+" is viewing this").Class("card-viewing"))
		} else {
			nodes = append(nodes, span.Text(strings.Join(viewOnly, ", ")+" are viewing this").Class("card-viewing"))
		}
	}

	return div.New(nodes...).Class("card-presence")
}
