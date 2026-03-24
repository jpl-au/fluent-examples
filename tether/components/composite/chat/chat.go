// Package chat provides a chat bubble layout for message feeds  -
// right-aligned bubbles for the current user's messages, left-aligned
// bubbles for other participants.
package chat

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Feed creates a scrollable chat feed container. The feed uses CSS
// column-reverse so the newest message (first in the slice) appears
// at the bottom with scroll pinned there automatically.
func Feed(items ...node.Node) *div.Element {
	return div.New(items...).Class("chat-feed")
}

// Bubble renders a single chat message as a styled bubble. When own
// is true the bubble is right-aligned with a rose tint; otherwise it
// is left-aligned with a purple tint.
func Bubble(user, text, timestamp string, own bool) node.Node {
	class := "chat-bubble chat-bubble-other"
	if own {
		class = "chat-bubble chat-bubble-own"
	}
	return div.New(
		div.New(
			span.Text(user).Class("chat-bubble-user"),
			span.Text(timestamp).Class("chat-bubble-time"),
		).Class("chat-bubble-header"),
		div.New().Class("chat-bubble-body").Text(text),
	).Class(class)
}

// Empty renders a centred placeholder when the feed has no messages.
func Empty(hint string) *div.Element {
	return div.New(
		p.Text(hint).Class("hint"),
	).Class("chat-empty")
}
