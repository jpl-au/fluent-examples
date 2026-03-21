package shoutbox

import (
	"time"

	"github.com/jpl-au/fluent-examples/tether/components/composite/chat"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/field"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
)

// Bus routes shoutbox messages between sessions. Components
// subscribe via StatefulConfig.Watchers; pages emit via Bus.Emit.
var Bus = tether.NewBus[Shout]()

// Shout is a single message in the shoutbox feed.
type Shout struct {
	User string
	Text string
	Time time.Time
}

// maxShouts is the number of messages kept in the feed.
const maxShouts = 50

// Shoutbox is a cross-session message component. It implements
// [tether.Component] and [tether.Mounter]. Mount performs one-time
// setup when the component is first added to a session. Messages
// from other sessions arrive via a WatchBus watcher declared in the
// handler package.
type Shoutbox struct {
	// Messages is the rolling feed of shouts, newest first.
	Messages []Shout
}

// Mount is called once during session startup by the framework.
func (s Shoutbox) Mount(sess tether.Session) tether.Component {
	sess.Toast("Shoutbox connected")
	return s
}

// Render builds the shoutbox UI: a message feed and input form.
func (s Shoutbox) Render() node.Node {
	return div.New(
		shoutFeed(s.Messages),
		div.New(
			bind.Apply(field.TextWithID("shout-input", "shout-input", "Type a message…"),
				bind.OnKeyDown("send"), bind.FilterKey("Enter"),
				bind.Collect("#shout-input"), bind.Reset(),
			),
			button.SmallPrimaryAction("Send", "send", bind.Collect("#shout-input"), bind.Reset()),
		).Class("shoutbox-form"),
	).Class("shoutbox")
}

// Handle processes shoutbox events. The prefix is already stripped
// by StatefulConfig.Components.
func (s Shoutbox) Handle(sess tether.Session, ev tether.Event) tether.Component {
	switch ev.Action {
	case "send":
		text := ev.Data["shout-input"]
		if text == "" {
			return s
		}
		// Emit to all other sessions via the bus.
		Bus.Emit(sess, Shout{
			User: sess.ID()[:6],
			Text: text,
			Time: time.Now(),
		})
		// Add to the sender's own feed immediately.
		s.Messages = appendShout(s.Messages, Shout{
			User: "you",
			Text: text,
			Time: time.Now(),
		})
	}
	return s
}

// shoutFeed renders the message list or an empty-state placeholder.
// Both branches share Dynamic("shout-feed") so the diff engine can
// detect the transition when messages arrive.
func shoutFeed(messages []Shout) node.Node {
	if len(messages) == 0 {
		return chat.Empty("No messages yet. Open another tab and start chatting.").Dynamic("shout-feed")
	}
	bubbles := make([]node.Node, len(messages))
	for i, m := range messages {
		bubbles[i] = chat.Bubble(m.User, m.Text, m.Time.Format("15:04"), m.User == "you")
	}
	return chat.Feed(bubbles...).Dynamic("shout-feed")
}

// appendShout prepends a message to the feed and trims it to
// maxShouts so the list doesn't grow without bound.
func appendShout(msgs []Shout, m Shout) []Shout {
	msgs = append([]Shout{m}, msgs...)
	if len(msgs) > maxShouts {
		msgs = msgs[:maxShouts]
	}
	return msgs
}
