package groups

import (
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

// State is the per-session state for the groups demo.
type State struct {
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
	// CurrentRoom is the room the session belongs to: "alpha",
	// "bravo", or "" when not in any room.
	CurrentRoom string
	// RoomMembers is the member count of the current room.
	RoomMembers int
	// RoomMessage is the last message received in the room.
	RoomMessage string
	// Activity holds join/leave log entries pushed by OnJoin
	// and OnLeave callbacks.
	Activity []string
}

var (
	roomAlpha = tether.NewGroup[State]()
	roomBravo = tether.NewGroup[State]()
	presence  = shared.NewPresenceCountOnly()
)

// init wires OnJoin/OnLeave callbacks for both rooms so activity
// logging starts before any session connects.
func init() {
	roomAlpha.OnJoin = func(sess *tether.StatefulSession[State]) {
		broadcastActivity(roomAlpha, "User "+sess.ID()[:6]+" joined Alpha")
	}
	roomAlpha.OnLeave = func(sess *tether.StatefulSession[State]) {
		broadcastActivity(roomAlpha, "User "+sess.ID()[:6]+" left Alpha")
	}
	roomBravo.OnJoin = func(sess *tether.StatefulSession[State]) {
		broadcastActivity(roomBravo, "User "+sess.ID()[:6]+" joined Bravo")
	}
	roomBravo.OnLeave = func(sess *tether.StatefulSession[State]) {
		broadcastActivity(roomBravo, "User "+sess.ID()[:6]+" left Bravo")
	}
}

// broadcastActivity pushes an activity log entry to every session
// in the group and updates each session's member count.
func broadcastActivity(g *tether.Group[State], msg string) {
	g.Broadcast(func(_ *tether.StatefulSession[State], s State) State {
		s.Activity = append([]string{msg}, s.Activity...)
		if len(s.Activity) > 10 {
			s.Activity = s.Activity[:10]
		}
		s.RoomMembers = g.Len()
		return s
	})
}

// New creates a handler demonstrating explicit tether.Group
// membership with Add/Remove, Broadcast/BroadcastOthers, and
// OnJoin/OnLeave callbacks.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "groups",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: presence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/groups/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Groups"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: shared.Watchers[State](presence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("groups: connected", "id", sess.ID())
			shared.TrackPresence(presence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("groups: disconnected", "id", sess.ID())
			shared.UntrackPresence(presence, sess.ID())
			// Remove from whichever room the session was in.
			leaveRoom(sess, sess.State().CurrentRoom)
		},
	})
}

// Handle processes events on the groups page. State logic (room
// tracking, message text) uses the Session interface so it is
// testable with tethertest. Group operations (Add, Remove, Broadcast)
// require a StatefulSession and are guarded by a type assertion.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	// Group operations need the concrete StatefulSession - nil in tests.
	live, _ := sess.(*tether.StatefulSession[State])

	switch ev.Action {
	case "groups.join-alpha":
		leaveRoom(live, s.CurrentRoom)
		if live != nil {
			roomAlpha.Add(live)
		}
		s.CurrentRoom = "alpha"
		s.RoomMembers = roomAlpha.Len()

	case "groups.join-bravo":
		leaveRoom(live, s.CurrentRoom)
		if live != nil {
			roomBravo.Add(live)
		}
		s.CurrentRoom = "bravo"
		s.RoomMembers = roomBravo.Len()

	case "groups.leave":
		leaveRoom(live, s.CurrentRoom)
		s.CurrentRoom = ""
		s.RoomMembers = 0

	case "groups.broadcast":
		msg := ev.Data["message-input"]
		if msg == "" {
			return s
		}
		g := currentGroup(s.CurrentRoom)
		if g != nil {
			g.Broadcast(func(_ *tether.StatefulSession[State], st State) State {
				st.RoomMessage = msg
				return st
			})
		}
		// The sender's state is updated by Broadcast, so set it
		// here as well for the immediate return value.
		s.RoomMessage = msg

	case "groups.broadcast-others":
		msg := ev.Data["message-input"]
		if msg == "" {
			return s
		}
		g := currentGroup(s.CurrentRoom)
		if g != nil {
			g.BroadcastOthers(sess, func(_ *tether.StatefulSession[State], st State) State {
				st.RoomMessage = msg
				return st
			})
		}
		// The sender already has the message; note that they sent it.
		s.RoomMessage = "You sent: " + msg
	}
	return s
}

// leaveRoom removes the session from the named room. The room name
// is passed explicitly (from state) so the function does not need
// to read StatefulSession.State(). Safe to call with a nil session
// (tethertest) or an empty room name.
func leaveRoom(live *tether.StatefulSession[State], room string) {
	if live == nil {
		return
	}
	switch room {
	case "alpha":
		roomAlpha.Remove(live)
	case "bravo":
		roomBravo.Remove(live)
	}
}

// currentGroup returns the group for the given room name, or nil
// if the session is not in a room.
func currentGroup(room string) *tether.Group[State] {
	switch room {
	case "alpha":
		return roomAlpha
	case "bravo":
		return roomBravo
	default:
		return nil
	}
}
