package chat

import (
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the chat room page. The shoutbox component handles
// all rendering and event dispatch - it is wired via StatefulConfig.Components
// on the handler, so this page only calls Render.
func Render(s State) node.Node {
	return cpage.New(
		panel.Card(
			"Chat Room",
			"A real-time chat room powered by tether.Component, Mounter, and Bus. "+
				"Messages are broadcast to every connected tab via ShoutBus.Emit and "+
				"delivered by a WatchBus watcher. The component's Mount method fires "+
				"a toast when the session first connects. Open two browser tabs and "+
				"start chatting.",
			"Component · Mounter · Bus · WatchBus", panel.WS|panel.SSE,
			bind.Apply(layout.Container(s.Shoutbox.Render()), bind.Prefix("shoutbox")).Dynamic("chat-section"),
		),
	)
}
