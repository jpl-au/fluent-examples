// Package detail renders the card detail/edit view. Used for both
// new and existing cards - the same form, the same layout. An empty
// Card ID means new; a populated one means edit.
package detail

import (
	"time"

	security "github.com/jpl-au/fluent-security"
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether-app/components/simple/badge"
	"github.com/jpl-au/fluent-examples/tether-app/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether-app/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New renders the detail view for a card. When c.ID is empty, the
// form creates a new card. When populated, it updates the existing
// one. Same component either way. Existing cards additionally render
// a read-only preview of the description through the cleaner, so
// users can see how their rich-text input will be sanitised before
// it is shown to other viewers on the board.
func New(c store.Card, cleaner *security.Cleaner, menuOpen bool) node.Node {
	isNew := c.ID == ""

	back := bind.Apply(
		a.New().Class("detail-back").Add(span.Text("\u2190 Back to Board")),
		bind.OnClick("card.back"),
	)

	title := "New Card"
	if !isNew {
		title = c.Title
	}

	header := div.New(
		back,
		div.New(
			span.Text(title).Class("detail-title"),
			columnBadge(c.Column, isNew),
		).Class("detail-title-row"),
		overflow(c, isNew, menuOpen),
	).Class("detail-header")

	f := bind.Apply(
		form.New(
			input.Hidden("id", c.ID),
			div.New(
				field.Label("Title"),
				field.TextValue("title", c.Title, "Card title"),
			).Class("form-group"),
			div.New(
				field.Label("Description"),
				field.Area("description", "Add a description. Basic HTML (strong, em, a, code) survives; scripts are stripped.", c.Description),
			).Class("form-group"),
			div.New(
				button.Submit("Save"),
			).Class("detail-actions"),
		).Class("detail-form"),
		bind.OnSubmit("card.save"),
		bind.OnInput("card.typing"),
		bind.Debounce(500*time.Millisecond),
	)

	return bind.Apply(
		div.New(header, f, preview(c, cleaner, isNew), activity(c.Activity)).Class("detail"),
		bind.Hotkey("escape", "card.back"),
	).Dynamic("detail")
}

// preview renders a read-only view of the description as it will
// appear to other users, with fluent-security's UGC policy applied.
// Hidden for new cards (no stored content yet) and for existing
// cards with empty descriptions.
func preview(c store.Card, cleaner *security.Cleaner, isNew bool) node.Node {
	if isNew || c.Description == "" {
		return nil
	}
	return div.New(
		span.Text("Description preview").Class("preview-title"),
		div.New(cleaner.Clean(c.Description)).Class("preview-body"),
	).Class("preview-section")
}

// activity renders the card's event log. Hidden for new cards.
func activity(events []store.Event) node.Node {
	if len(events) == 0 {
		return nil
	}

	items := make([]node.Node, len(events))
	for i := len(events) - 1; i >= 0; i-- {
		ev := events[i]
		items[len(events)-1-i] = div.New(
			span.Text(ev.User).Class("activity-user"),
			span.Text(ev.Action).Class("activity-action"),
			span.Text(ev.Created.Format("15:04")).Class("activity-time"),
		).Class("activity-item")
	}

	return div.New(
		span.Text("Activity").Class("activity-title"),
		div.New(items...).Class("activity-list"),
	).Class("activity-section")
}

// overflow renders the three-dot menu for existing cards. The menu is
// server-state driven (State.MenuOpen) so it can dismiss on any click
// outside it via bind.Outside, which listens only while the menu is
// actually rendered.
func overflow(c store.Card, isNew, menuOpen bool) node.Node {
	if isNew {
		return nil
	}

	// bind.Stop keeps the trigger click from bubbling to the open menu's
	// document-level Outside listener, so toggling closed sends a single
	// card.menu.toggle rather than also firing card.menu.close.
	toggle := bind.Apply(
		span.Text("\u22EF").Class("overflow-trigger"),
		bind.OnClick("card.menu.toggle"),
		bind.Stop(),
	)

	if !menuOpen {
		return div.New(toggle).Class("overflow")
	}

	menu := bind.Apply(
		div.New(
			bind.Apply(
				span.Text("Delete card").Class("overflow-item overflow-danger"),
				bind.OnClick("card.delete"),
				bind.EventData("id", c.ID),
				bind.Confirm("Delete this card?"),
			),
		).Class("overflow-menu"),
		bind.OnClick("card.menu.close"),
		bind.Outside(),
	)

	return div.New(toggle, menu).Class("overflow")
}

// columnBadge renders the column indicator for existing cards.
func columnBadge(col store.Column, isNew bool) node.Node {
	if isNew {
		return nil
	}
	switch col {
	case store.InProgress:
		return badge.Progress(col.String())
	case store.Done:
		return badge.Done(col.String())
	default:
		return badge.Todo(col.String())
	}
}
