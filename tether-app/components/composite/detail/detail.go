// Package detail renders the card detail/edit view. When a user
// clicks a card on the board, the board region swaps out entirely
// for this view. Saving or going back swaps it back.
package detail

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/h2"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether-app/components/simple/badge"
	"github.com/jpl-au/fluent-examples/tether-app/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether-app/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New renders the full detail view for a card. The form submits
// title and description back to the server as card.update.
func New(c store.Card) node.Node {
	f := bind.Apply(
		form.New(
			input.Hidden("id", c.ID),
			div.New(
				field.Label("Title"),
				field.TextValue("title", c.Title, "Card title"),
			).Class("form-group"),
			div.New(
				field.Label("Description"),
				field.Area("description", "Card description", c.Description),
			).Class("form-group"),
			div.New(
				div.New(
					button.SecondaryAction("Back", "card.back"),
					button.Submit("Save"),
				).Class("detail-actions-left"),
				button.DangerAction("Delete", "card.delete",
					bind.EventData("id", c.ID),
					bind.Confirm("Delete this card?"),
				),
			).Class("detail-actions"),
		).Class("detail-form"),
		bind.OnSubmit("card.update"),
	)

	return div.New(
		div.New(
			h2.New().Class("detail-title").Text(c.Title),
			badgeFor(c.Column),
		).Class("detail-header"),
		f,
	).Class("detail").Dynamic("detail")
}

// badgeFor returns the appropriate column badge.
func badgeFor(col store.Column) node.Node {
	switch col {
	case store.InProgress:
		return badge.Progress(col.String())
	case store.Done:
		return badge.Done(col.String())
	default:
		return badge.Todo(col.String())
	}
}
