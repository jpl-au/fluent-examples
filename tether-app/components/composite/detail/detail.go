// Package detail renders the card detail/edit view. Used for both
// new and existing cards - the same form, the same layout. An empty
// Card ID means new; a populated one means edit.
package detail

import (
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
// one. Same component either way.
func New(c store.Card) node.Node {
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
			span.New().Class("detail-title").Text(title),
			columnBadge(c.Column, isNew),
		).Class("detail-title-row"),
		overflow(c, isNew),
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
				field.Area("description", "Add a description...", c.Description),
			).Class("form-group"),
			div.New(
				button.Submit("Save"),
			).Class("detail-actions"),
		).Class("detail-form"),
		bind.OnSubmit("card.save"),
	)

	return bind.Apply(
		div.New(header, f, activity(c.Activity)).Class("detail"),
		bind.Hotkey("escape", "card.back"),
	).Dynamic("detail")
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
			span.New().Class("activity-user").Text(ev.User),
			span.New().Class("activity-action").Text(ev.Action),
			span.New().Class("activity-time").Text(ev.Created.Format("15:04")),
		).Class("activity-item")
	}

	return div.New(
		span.New().Class("activity-title").Text("Activity"),
		div.New(items...).Class("activity-list"),
	).Class("activity-section")
}

// overflow renders the three-dot menu for existing cards.
func overflow(c store.Card, isNew bool) node.Node {
	if isNew {
		return nil
	}

	toggle := bind.Apply(
		span.New().Class("overflow-trigger").Text("\u22EF"),
		bind.ToggleClass("overflow-open"),
		bind.ToggleTarget(".overflow-menu"),
	)

	menu := div.New(
		bind.Apply(
			span.New().Class("overflow-item overflow-danger").Text("Delete card"),
			bind.OnClick("card.delete"),
			bind.EventData("id", c.ID),
			bind.Confirm("Delete this card?"),
		),
	).Class("overflow-menu")

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
