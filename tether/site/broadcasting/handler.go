package broadcasting

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// State is the per-session state for the broadcasting demo.
type State struct {
	// Messages holds broadcast messages from all sessions.
	Messages []Message
	// MessageCount tracks total messages for the counter demo.
	MessageCount int
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
}

// Message is a broadcast message sent between sessions via the bus.
type Message struct {
	// User is a display name derived from the sender's session ID.
	User string
	// Text is the message content.
	Text string
}

// messageBus routes broadcast messages between sessions. Emit
// delivers to everyone except the sender; Subscribe and
// SubscribeAsync receive every message unconditionally.
var messageBus = tether.NewBus[Message]()

// messageCount tracks the total number of messages across all
// sessions. Raw bus subscribers increment this; sessions observe
// it via tether.WatchValue to keep their counter in sync.
var messageCount = tether.NewValue(0)

// Setup wires bus subscribers that must live for the duration of
// the server process. Pass the root context from main so subscriptions
// are cancelled cleanly on shutdown.
func Setup(ctx context.Context) {
	// Synchronous subscriber: increments the shared counter on
	// every message.
	messageBus.Subscribe(ctx, func(_ Message) {
		messageCount.Update(func(n int) int { return n + 1 })
	})

	// Asynchronous subscriber: simulates audit logging. Runs in
	// its own goroutine so I/O never blocks the publisher.
	messageBus.SubscribeAsync(ctx, func(m Message) {
		slog.Info("audit: broadcast message", "user", m.User, "text", m.Text)
	})
}

var broadcastPresence = shared.NewPresenceCountOnly()

// New creates a handler demonstrating cross-session broadcasting
// via tether.Bus, tether.Value, and declarative watchers.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "broadcasting",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: broadcastPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/broadcasting/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Broadcasting"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: append(shared.Watchers[State](broadcastPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),
			tether.WatchBus(messageBus, func(m Message, s State) State {
				s.Messages = append(s.Messages, m)
				return s
			}),
			tether.WatchValue(messageCount, func(n int, s State) State {
				s.MessageCount = n
				return s
			}),
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("broadcasting: connected", "id", sess.ID())
			shared.TrackPresence(broadcastPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("broadcasting: disconnected", "id", sess.ID())
			shared.UntrackPresence(broadcastPresence, sess.ID())
		},
	})
}

// Handle processes events on the broadcasting page, emitting
// messages to the shared bus on user input.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "broadcast.send":
		text := ev.Data["broadcast-input"]
		if text == "" {
			return s
		}
		msg := Message{
			User: "User " + sess.ID()[:6],
			Text: text,
		}

		// Append to the sender's state directly - they see it
		// immediately without waiting for the bus round-trip.
		s.Messages = append(s.Messages, msg)

		// Emit delivers to all other sessions via tether.WatchBus.
		// The sender is filtered out automatically.
		messageBus.Emit(sess, msg)
	}
	return s
}
