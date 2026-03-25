package handler

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h1"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether-app/components/composite/board"
	ccard "github.com/jpl-au/fluent-examples/tether-app/components/composite/card"
	"github.com/jpl-au/fluent-examples/tether-app/components/composite/column"
	"github.com/jpl-au/fluent-examples/tether-app/components/composite/detail"
	"github.com/jpl-au/fluent-examples/tether-app/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether-app/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether-app/layout"
	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// Render returns the top-level render function. It closes over the
// board store so the view always reads the latest shared state.
func Render(b *store.Board, viewers *viewers) func(State) node.Node {
	return func(s State) node.Node {
		if s.Name == "" {
			return landing()
		}

		var content node.Node
		switch s.View {
		case "detail":
			if s.SelectedID == "" {
				content = detail.New(store.Card{})
			} else if c, ok := b.Card(s.SelectedID); ok {
				content = detail.New(c)
			} else {
				content = boardView(b, viewers, s.SessionID)
			}
		default:
			content = boardView(b, viewers, s.SessionID)
		}

		return layout.Shell(s.Name, s.OnlineCount, addButton(), content)
	}
}

// landing renders the name entry page shown on first visit.
func landing() node.Node {
	return div.New(
		div.New(
			h1.Text("Kanban Board").Class("landing-title"),
			p.Text("A collaborative board powered by Tether. Enter your name to get started.").Class("landing-desc"),
			bind.Apply(
				field.Inline(
					field.Text("name", "Your name"),
					button.Submit("Sign In"),
				),
				bind.OnSubmit("name.set"),
				bind.AutoFocus(),
			),
		).Class("landing-card"),
		// Hidden marker so the DnD extension JS loads on initial render.
		bind.Apply(div.New().Class("sr-only"), bind.Draggable()),
	).Class("landing").Dynamic("landing")
}

// boardView renders the three-column kanban grid, or an empty state
// prompt when all columns are empty.
func boardView(b *store.Board, viewers *viewers, sessionID string) node.Node {
	empty := true
	var cols []node.Node
	for _, col := range store.Columns() {
		cards := b.Cards(col)
		if len(cards) > 0 {
			empty = false
		}
		var cardNodes []node.Node
		for _, c := range cards {
			cv := ccard.CardViewers{
				Viewing: viewers.ViewingCard(c.ID, sessionID),
				Typing:  viewers.TypingOnCard(c.ID, sessionID),
			}
			cardNodes = append(cardNodes, ccard.New(c, cv))
		}
		cols = append(cols, columnView(col, cardNodes))
	}
	if empty {
		return div.New(
			p.Text("No cards yet. Click Add Card to get started.").Class("empty-board"),
		).Class("empty-state").Dynamic("board")
	}
	return board.New(cols...)
}

// columnView wraps a column component as a sortable drop zone.
func columnView(col store.Column, cards []node.Node) node.Node {
	var content node.Node
	if len(cards) == 0 {
		content = column.New(col.String(), 0, column.Empty())
	} else {
		content = column.New(col.String(), len(cards), cards...)
	}

	return bind.Apply(
		div.New(content).Class("drop-zone"),
		bind.Sortable("card.move"),
		bind.EventData("column", strconv.Itoa(int(col))),
	).Dynamic("col-" + strconv.Itoa(int(col)))
}

// addButton renders the header action to create a new card.
func addButton() node.Node {
	return button.PrimaryAction("Add Card", "card.new")
}
