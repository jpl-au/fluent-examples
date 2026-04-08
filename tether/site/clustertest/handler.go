// Package clustertest provides a hidden test-only handler for
// verifying cluster integration in playwright tests. It is not
// registered in the main application routes or navigation.
package clustertest

import (
	"net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
	"github.com/jpl-au/tether/mode"
)

// Message is a message sent between sessions via the cluster bus.
type Message struct {
	Text string
}

// State is the per-session state for the cluster test page.
type State struct {
	Messages []Message
}

// New creates a handler for the cluster test page. The bus parameter
// allows tests to inject a bus wired to a specific cluster, avoiding
// package-level state that would be shared between servers in the
// same process.
func New(app tether.App, bus *tether.Bus[Message]) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name: "cluster-test",
		Mode: mode.WebSocket,

		InitialState: func(_ *http.Request) State {
			return State{}
		},

		Render: func(s State) node.Node {
			return render(s)
		},

		Handle: func(sess tether.Session, s State, ev tether.Event) State {
			switch ev.Action {
			case "cluster.send":
				text := ev.Data["cluster-input"]
				if text == "" {
					return s
				}
				msg := Message{Text: text}

				// Append to sender's state directly so they see
				// their own message without waiting for the bus.
				s.Messages = append(s.Messages, msg)

				// Emit delivers to all other sessions. When a
				// cluster is configured, the bus also publishes
				// to the cluster for cross-node delivery.
				bus.Emit(sess, msg)
			}
			return s
		},

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Cluster Test"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: []tether.Watcher[State]{
			tether.WatchBus(bus, func(m Message, s State) State {
				s.Messages = append(s.Messages, m)
				return s
			}),
		},
	})
}

// render builds the minimal UI for the cluster test page.
func render(s State) node.Node {
	return div.New(
		div.New(
			bind.Apply(
				input.Text("cluster-input", "").
					ID("cluster-input").
					Placeholder("Type a message..."),
				bind.OnKeyDown("cluster.send"),
				bind.FilterKey("Enter"),
				bind.Collect("#cluster-input"),
				bind.Reset(),
			),
			bind.Apply(
				button.Text("Send"),
				bind.OnClick("cluster.send"),
				bind.Collect("#cluster-input"),
				bind.Reset(),
			),
		),
		messageList(s.Messages),
	)
}

// messageList renders received messages or an empty placeholder.
func messageList(msgs []Message) node.Node {
	if len(msgs) == 0 {
		return div.Text("No messages yet.").Dynamic("cluster-messages")
	}
	nodes := make([]node.Node, len(msgs))
	for i, m := range msgs {
		nodes[i] = div.Text(m.Text)
	}
	return div.New(nodes...).Dynamic("cluster-messages")
}
