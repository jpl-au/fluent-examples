// Package card renders an individual kanban card within a column.
// Each card is draggable (for moving between columns) and clickable
// (for opening the detail view). Shows the title, description
// snippet, creator, timestamp, and who is currently viewing.
//
// Presence indicators (typing and viewing) use Signals rather than
// server-rendered HTML. The handler pushes signal values like
// "typing-{id}" and "viewing-{id}" directly to the client, and
// bind.BindText updates the text in place with no render cycle.
// This is the right tool for high-frequency, text-only updates
// where the DOM structure never changes - only the content.
package card

import (
	security "github.com/jpl-au/fluent-security"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// New renders a draggable kanban card with signal-bound presence
// indicators. Presence text is pushed via sess.Signal() - see
// handler/viewers.go for the signal push logic.
func New(c store.Card) node.Node {
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
					presence(c.ID),
				).Class("card-body"),
				bind.OnClick("card.select"),
				bind.EventData("id", c.ID),
			),
		).Class("card"),
		bind.Draggable(),
		bind.EventData("id", c.ID),
	).Dynamic(c.ID)
}

// desc renders a truncated plain-text snippet of the description, or
// nil if empty. The card grid is not the place for rich formatting -
// even if the description contains <strong> or <a>, the snippet
// should read as a single line of preview text. security.PlainText
// strips every tag and returns plain text; it is the correct preset
// for this use case. The full rich HTML is shown on the detail view,
// where the hoisted cleaner enforces the UGC policy.
//
// Truncation is rune-based rather than byte-based so multi-byte
// characters (emoji, CJK, accented letters) are never split in the
// middle of a codepoint, which would produce invalid UTF-8 in the
// rendered HTML.
func desc(s string) node.Node {
	if s == "" {
		return nil
	}
	plain := string(security.PlainText(s).Render())
	runes := []rune(plain)
	if len(runes) > 80 {
		plain = string(runes[:77]) + "..."
	}
	return p.Text(plain).Class("card-desc")
}

// presence renders signal-bound typing and viewing indicators.
// The server pushes text values via sess.Signal("typing-{id}", ...)
// and the client updates these elements directly - no render/diff
// cycle needed. BindShow hides the element when the signal value
// is empty, so the presence section collapses automatically.
func presence(cardID string) node.Node {
	return div.New(
		bind.Apply(span.New().Class("card-typing"),
			bind.BindText("typing-"+cardID),
			bind.BindShow("typing-"+cardID),
		),
		bind.Apply(span.New().Class("card-viewing"),
			bind.BindText("viewing-"+cardID),
			bind.BindShow("viewing-"+cardID),
		),
	).Class("card-presence")
}
