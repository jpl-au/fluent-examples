package handler

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/card"
)

// logCard builds the card content shared by the WebSocket and SSE
// pages: a status indicator, description, and a scrollable container
// where JavaScript appends log entries as they arrive.
func logCard(title, desc, statusID, logID string) node.Node {
	return card.New(title,
		div.New(
			span.Text("Status: ").Class("label"),
			span.Text("connecting...").Class("status-connecting").ID(statusID),
		).Class("layout-row"),
		div.New(
			span.Text(desc).Class("hint"),
		),
		div.New().Class("log-feed").ID(logID),
	)
}
