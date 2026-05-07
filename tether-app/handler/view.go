package handler

import (
	"strconv"

	security "github.com/jpl-au/fluent-security"
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
// board store so the view always reads the latest shared state, and
// over the description cleaner so the detail view can sanitise rich
// descriptions with the hoisted UGC policy. Card previews on the
// board use security.PlainText directly and do not need the cleaner.
//
// Presence indicators are signal-bound (not rendered here) and the
// online count badge uses bind.BindText - see layout.Shell.
//
// The board view is wrapped in node.Memoise keyed on BoardVersion.
// When BoardVersion hasn't changed (e.g. navigation between views),
// the Memoiser skips the entire board subtree - no column renders,
// no card renders, no HTML generated. The closure only runs on a
// cache miss (board mutation incremented the version).
func Render(b *store.Board, cleaner *security.Cleaner) func(State) node.Node {
	return func(s State) node.Node {
		if s.Name == "" {
			return landing()
		}

		var content node.Node
		switch s.View {
		case "detail":
			if s.SelectedID == "" {
				content = detail.New(store.Card{}, cleaner)
			} else if c, ok := b.Card(s.SelectedID); ok {
				content = detail.New(c, cleaner)
			} else {
				content = memoiseBoard(b, s.BoardVersion)
			}
		default:
			content = memoiseBoard(b, s.BoardVersion)
		}

		return layout.Shell(s.Name, s.OnlineCount, addButton(), content)
	}
}

// memoiseBoard wraps the board rendering in node.Memoise so the
// entire column grid is skipped when the board hasn't changed.
// The boardVersion key is incremented by the handler on every
// board mutation (create, save, move, delete).
//
// The Memoise node is a child of the Dynamic div (Pattern 1 from
// fluent-jit docs). On a cache hit, the Memoiser finds the key on
// the Dynamic's child and skips the closure entirely - no column
// renders, no card renders, no HTML generated.
func memoiseBoard(b *store.Board, boardVersion int) node.Node {
	return div.New(
		node.Memoise(boardVersion, func() node.Node {
			return boardColumns(b)
		}),
	).Dynamic("board")
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

// boardColumns renders the column grid contents. Called inside the
// Memoise closure in memoiseBoard, so this only runs on a cache miss
// (board mutation). Presence indicators are signal-bound and update
// without re-rendering.
func boardColumns(b *store.Board) node.Node {
	empty := true
	var cols []node.Node
	for _, col := range store.Columns() {
		cards := b.Cards(col)
		if len(cards) > 0 {
			empty = false
		}
		var cardNodes []node.Node
		for _, c := range cards {
			cardNodes = append(cardNodes, ccard.New(c))
		}
		cols = append(cols, columnView(col, cardNodes))
	}
	if empty {
		return div.New(
			p.Text("No cards yet. Click Add Card to get started.").Class("empty-board"),
		).Class("empty-state")
	}
	return board.Columns(cols...)
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
