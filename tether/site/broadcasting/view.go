package broadcasting

import (
	"strconv"
	"strings"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/feed"
	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the broadcasting page with three demo cards:
// cross-session events, message counter, and async subscribers.
func Render(s State, who *tether.Presence[string]) node.Node {
	return page.New(
		panel.Card(
			"Who's Here",
			"tether.Presence tracks per-session metadata and makes it available to all sessions. "+
				"Open multiple tabs to see each session appear. Presence.Each excludes the current session "+
				"so you don't see yourself in the list.",
			"tether.Presence · Presence.Each", panel.WS|panel.SSE,
			whoIsHere(who, s.SessionID),
		),

		panel.Card(
			"Cross-Session Events",
			"Type a message and click Send. The sender sees their own message immediately (returned from Handle); other sessions receive it via tether.WatchBus with automatic sender filtering. bus.Emit publishes to everyone except the sender, preventing double-apply.",
			"bus.Emit · tether.WatchBus", panel.WS|panel.SSE,
			layout.Stack(
				layout.Row(
					bind.Apply(field.TextWithID("broadcast-input", "broadcast-input", "Type a message…"),
						bind.OnKeyDown("broadcast.send"), bind.FilterKey("Enter"),
						bind.Collect("#broadcast-input"), bind.Reset(),
					),
					button.PrimaryAction("Send", "broadcast.send",
						bind.Collect("#broadcast-input"),
						bind.Reset(),
					),
				),
				messageList(s.Messages),
			),
		),

		panel.Card(
			"Message Counter",
			"Displays the total number of messages broadcast across all sessions. A raw bus.Subscribe callback increments a tether.Value[int] on every message. Each session observes the value via tether.WatchValue, which delivers the current count immediately on subscription and pushes updates automatically.",
			"bus.Subscribe · tether.Value · tether.WatchValue", panel.WS|panel.SSE,
			span.Text("Total messages: "+strconv.Itoa(s.MessageCount)).Dynamic("message-count"),
		),

		panel.Card(
			"Async Subscribers",
			"bus.SubscribeAsync registers a callback that runs in its own goroutine for every event. Use it for I/O-bound consumers - database writes, HTTP calls, audit logging - that must not block the publisher. This demo's async subscriber logs each message to slog. Check the server console to see it in action.",
			"bus.SubscribeAsync", panel.WS|panel.SSE,
			hint.Text("Async subscriber logs to the server console via slog.Info."),
		),
	)
}

// whoIsHere renders the list of users on this page via Presence.Each.
func whoIsHere(who *tether.Presence[string], sessionID string) node.Node {
	var names []string
	who.Each(sessionID, func(_ string, name string) {
		names = append(names, name)
	})
	if len(names) == 0 {
		return hint.Text("No other users on this page. Open another tab to see presence.")
	}
	return span.Text("Also here: " + strings.Join(names, ", ")).Dynamic("who-here")
}

// messageList renders broadcast messages or a placeholder.
func messageList(msgs []Message) node.Node {
	if len(msgs) == 0 {
		return layout.Container(
			hint.Text("No messages yet. Open another tab and start broadcasting."),
		).Dynamic("messages")
	}
	nodes := make([]node.Node, len(msgs))
	for i, m := range msgs {
		nodes[i] = feed.MessageItem(m.User, m.Text)
	}
	return layout.Container(feed.Messages(nodes...)).Dynamic("messages")
}
