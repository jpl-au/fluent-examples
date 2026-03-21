package groups

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/feed"
	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/badge"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the groups page with four demo cards demonstrating
// Group.Add/Remove, Broadcast/BroadcastOthers, and OnJoin/OnLeave.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Join a Room",
			"Group.Add registers a session with a named group. Each session can be in multiple groups simultaneously, but this demo enforces one room at a time to keep things clear.",
			"tether.Group · Add", panel.WS|panel.SSE,
			layout.Stack(
				layout.Row(
					button.PrimaryAction("Join Alpha", "groups.join-alpha"),
					button.PrimaryAction("Join Bravo", "groups.join-bravo"),
				),
				roomStatus(s.CurrentRoom, s.RoomMembers),
			),
		),

		panel.Card(
			"Leave Room",
			"Group.Remove unregisters the session. If the session is not a member, Remove is a no-op.",
			"tether.Group · Remove", panel.WS|panel.SSE,
			layout.Stack(
				button.SecondaryAction("Leave Room", "groups.leave"),
				leaveHint(s.CurrentRoom),
			),
		),

		panel.Card(
			"Broadcast",
			"Broadcast applies a function to every session in the group. BroadcastOthers excludes the sender - the typical pattern when the sender already has the update in their own state.",
			"tether.Group · Broadcast · BroadcastOthers", panel.WS|panel.SSE,
			layout.Stack(
				layout.Row(
					bind.Apply(field.TextWithID("message-input", "message-input", "Type a message…"),
						bind.OnKeyDown("groups.broadcast"), bind.FilterKey("Enter"),
						bind.Collect("#message-input"), bind.Reset(),
					),
					button.PrimaryAction("Send to Room", "groups.broadcast",
						bind.Collect("#message-input"),
						bind.Reset(),
					),
					button.PrimaryAction("Send to Others", "groups.broadcast-others",
						bind.Collect("#message-input"),
						bind.Reset(),
					),
				),
				messageResult(s.RoomMessage),
			),
		),

		panel.Card(
			"OnJoin / OnLeave",
			"Callbacks fire after Add or Remove completes, outside the write lock. Use them for audit logging, notifications, or updating member counts.",
			"tether.Group · OnJoin · OnLeave", panel.WS|panel.SSE,
			activityLog(s.Activity),
		),
	)
}

// roomStatus displays the current room name and member count, or a
// hint when the session is not in any room.
func roomStatus(room string, members int) node.Node {
	if room == "" {
		return layout.Container(
			hint.Text("Not in any room. Click a button above to join."),
		).Dynamic("room-status")
	}
	return layout.Container(
		badge.Green(room),
		span.Text(strconv.Itoa(members)+" member(s)"),
	).Dynamic("room-status")
}

// leaveHint shows contextual guidance depending on whether the
// session is currently in a room.
func leaveHint(room string) node.Node {
	if room == "" {
		return layout.Container(
			hint.Text("You are not in a room."),
		).Dynamic("leave-hint")
	}
	return layout.Container(
		hint.Textf("Currently in %s. Click Leave Room to exit.", room),
	).Dynamic("leave-hint")
}

// messageResult renders the last received room message or a
// placeholder when no message has been received yet.
func messageResult(msg string) node.Node {
	if msg == "" {
		return layout.Container(
			hint.Text("No messages yet. Join a room and send one."),
		).Dynamic("room-message")
	}
	return layout.Container(
		p.Text(msg),
	).Dynamic("room-message")
}

// activityLog renders the join/leave activity feed or a placeholder.
func activityLog(items []string) node.Node {
	if len(items) == 0 {
		return layout.Container(
			hint.Text("No activity yet. Join a room to see OnJoin/OnLeave events."),
		).Dynamic("activity-log")
	}
	nodes := make([]node.Node, len(items))
	for i, item := range items {
		nodes[i] = feed.ActivityText(item)
	}
	return layout.Container(
		feed.Activity(nodes...),
	).Dynamic("activity-log")
}
