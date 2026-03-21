package live

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/badge"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// RenderSSE builds the live updates page for the SSE variant: uptime
// ticker, activity feed, online count, broadcast, and group info.
func RenderSSE(s State) node.Node {
	return page.New(
		panel.Card(
			"Uptime Ticker",
			"Watch the counter increment every second. A background goroutine started with sess.Go ticks on the server and pushes the new value as a signal over the SSE connection.",
			"sess.Go · sess.Signal", panel.WS|panel.SSE,
			layout.Row(
				bind.Apply(span.Text("0"), bind.BindText("live.uptime")),
				hint.Span(" seconds since connect"),
			),
		),

		panel.Card(
			"Activity Feed",
			"Open this page in a second browser tab, then close it. Join and leave events appear in the feed below. A tether.Bus broadcasts these events to every connected session via SSE messages.",
			"tether.Bus · tether.WatchBus", panel.WS|panel.SSE,
			sseActivityFeed(s.Activity),
		),

		panel.Card(
			"Online Count (tether.Value)",
			"This badge shows how many sessions are connected right now. Open or close tabs to see it change. The count is a shared tether.Value that pushes updates as signals via SSE whenever it changes.",
			"tether.NewValue · tether.WatchValue", panel.WS|panel.SSE,
			layout.Stack(
				bind.Apply(badge.GreenDynamic("online-count-live", fmt.Sprintf("%d online", s.OnlineCount)), bind.BindText("online_count")),
				hint.Text("Same as the header badge - both bound to the online_count signal."),
			),
		),

		panel.Card(
			"Broadcast Message",
			"Click Send Broadcast - the client sends a POST to the server, which publishes the message to the tether.Bus. Every connected session receives it over SSE and sees it in the activity feed above.",
			"tether.Bus.Publish", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Send Broadcast", "live.broadcast"),
				sseBroadcastResult(s.LastBroadcast),
			),
		),

		panel.Card(
			"Group",
			"The session count below reflects the number of active SSE connections in this handler's Group. Open more tabs to see it increase.",
			"tether.NewGroup · Group.Len", panel.WS|panel.SSE,
			layout.Stack(
				hint.Textf("Group has %d session(s)", sseGroup.Len()),
			),
		),
	)
}

// sseActivityFeed renders the activity list for the SSE variant.
func sseActivityFeed(items []shared.ActivityItem) node.Node {
	// Reuse the WS variant - identical rendering.
	return activityFeed(items)
}

// sseBroadcastResult shows the last broadcast message for the SSE variant.
func sseBroadcastResult(msg string) node.Node {
	return broadcastResult(msg)
}
