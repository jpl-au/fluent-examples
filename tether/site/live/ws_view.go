package live

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/feed"
	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/badge"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// RenderWS builds the full live updates page for the WebSocket variant.
func RenderWS(s State) node.Node {
	return page.New(
		panel.Card(
			"Uptime Ticker",
			"Watch the counter increment every second. A background goroutine started with sess.Go ticks on the server and pushes the new value as a signal. The number updates on the client without re-rendering the page.",
			"sess.Go · sess.Signal", panel.WS|panel.SSE,
			layout.Row(
				bind.Apply(span.Text("0"), bind.BindText("live.uptime")),
				hint.Span(" seconds since connect"),
			),
		),

		panel.Card(
			"Activity Feed",
			"Open this page in a second browser tab, then close it. You will see join and leave events appear in the feed below. A tether.Bus broadcasts these events to every connected session, and tether.WatchBus updates each session's state, triggering a re-render.",
			"tether.Bus · tether.WatchBus", panel.WS|panel.SSE,
			activityFeed(s.Activity),
		),

		panel.Card(
			"Online Count (tether.Value)",
			"This badge shows how many sessions are connected right now. Open or close tabs to see it change. tether.NewValue holds the count on the server; tether.WatchValue pushes the new value as a signal to every session whenever it changes.",
			"tether.NewValue · tether.WatchValue", panel.WS|panel.SSE,
			layout.Stack(
				bind.Apply(badge.GreenDynamic("online-count-live", fmt.Sprintf("%d online", s.OnlineCount)), bind.BindText("online_count")),
				hint.Text("This badge is identical to the one in the header - both bound to the same signal."),
			),
		),

		panel.Card(
			"Broadcast Message",
			"Click Send Broadcast - the message appears in the activity feed above, and in every other connected tab's feed too. The server publishes the message to the tether.Bus, which delivers it to all sessions.",
			"tether.Bus.Publish", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Send Broadcast", "live.broadcast"),
				broadcastResult(s.LastBroadcast),
			),
		),

		panel.Card(
			"Group.Broadcast",
			"Click Announce - the server calls Group.Broadcast to push a state mutation to every connected session simultaneously, including the sender. Open multiple tabs to see all of them update at the same time.",
			"Group.Broadcast", panel.WS|panel.SSE,
			layout.Stack(
				layout.Row(
					button.PrimaryAction("Announce to All", "live.announce"),
					hint.Textf("Group has %d session(s)", wsGroup.Len()),
				),
				announcementResult(s.Announcement),
			),
		),

		panel.Card(
			"Group.BroadcastOthers",
			"Click Notify Others - all other connected sessions receive the notification, but the sender sees a confirmation instead. BroadcastOthers excludes the calling session.",
			"Group.BroadcastOthers", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Notify Others", "live.notify-others"),
				notificationResult(s.Notification),
			),
		),

		panel.Card(
			"Group.All()",
			"Click List Sessions - the server iterates Group.All() to collect every connected session's ID. Group.All() returns an iter.Seq so you can range over it directly.",
			"Group.All", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("List Sessions", "live.list-sessions"),
				sessionList(s.SessionIDs),
			),
		),

		panel.Card(
			"Group.OnJoin / Group.OnLeave",
			"Open and close browser tabs - each connect and disconnect triggers the OnJoin and OnLeave callbacks set directly on the Group. Check the server console to see them log.",
			"Group.OnJoin · Group.OnLeave", panel.WS|panel.SSE,
			hint.Text("Connect or disconnect a tab and check the server console."),
		),

		panel.Card(
			"SetTitle",
			"Click the button - the browser tab title updates immediately. SetTitle is on Session so it works directly in Handle without a type-assert.",
			"sess.SetTitle", panel.WS|panel.SSE,
			button.PrimaryAction("Set Tab Title", "live.set-title"),
		),

		panel.Card(
			"State() in Go()",
			"Click the button - a background goroutine starts, waits one second, then reads the current session state and toasts what it found. Try clicking Announce to All before the toast fires - the goroutine will see the updated announcement.",
			"sess.Go · sess.State()", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Read State in 1s", "live.read-state"),
				hint.Text("A toast will appear after one second with the current state."),
			),
		),

		panel.Card(
			"Close()",
			"Click the button - the session closes and the transport disconnects. The browser reconnects automatically. Use Close() for admin kick, forced logout, or expiring stale sessions.",
			"sess.Close()", panel.WS|panel.SSE,
			layout.Stack(
				button.DangerAction("Disconnect Session", "live.close"),
				hint.Text("The page will reconnect automatically after closing."),
			),
		),
	)
}

// activityFeed renders the activity list or placeholder.
func activityFeed(items []shared.ActivityItem) node.Node {
	if len(items) == 0 {
		return layout.Container(
			hint.Text("No activity yet. Open another tab to see join/leave events."),
		).Dynamic("activity")
	}
	nodes := make([]node.Node, len(items))
	for i, item := range items {
		nodes[i] = feed.ActivityItem(item.User, item.Action, item.Timestamp.Format("15:04:05"))
	}
	return layout.Container(feed.Activity(nodes...)).Dynamic("activity")
}

// announcementResult renders the last Group.Broadcast announcement.
func announcementResult(msg string) node.Node {
	if msg == "" {
		return layout.Container().Dynamic("announcement")
	}
	return layout.Container(
		span.Text(msg),
	).Dynamic("announcement")
}

// notificationResult renders the last notification from BroadcastOthers.
func notificationResult(msg string) node.Node {
	if msg == "" {
		return layout.Container(
			hint.Text("Open another tab and click Notify Others to see a notification here."),
		).Dynamic("notification")
	}
	return layout.Container(
		span.Text(msg),
	).Dynamic("notification")
}

// sessionList renders connected session IDs.
func sessionList(ids []string) node.Node {
	if len(ids) == 0 {
		return layout.Container(
			hint.Text("Click to see connected sessions."),
		).Dynamic("session-list")
	}
	nodes := make([]node.Node, len(ids))
	for i, id := range ids {
		nodes[i] = feed.ActivityText(id)
	}
	return layout.Container(
		feed.Activity(nodes...),
	).Dynamic("session-list")
}

// broadcastResult shows the last broadcast message.
func broadcastResult(msg string) node.Node {
	if msg == "" {
		return layout.Container().Dynamic("broadcast")
	}
	return layout.Container(
		span.Text("Last: " + msg),
	).Dynamic("broadcast")
}
