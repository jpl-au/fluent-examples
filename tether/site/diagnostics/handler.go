package diagnostics

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

// Entry is a single diagnostic event captured from one or more
// handler Diagnostics buses for display in the UI.
type Entry struct {
	Kind      string
	SessionID string
	Detail    string
}

// diagnosticBus aggregates diagnostic events from all handlers so
// the diagnostics page can observe them via a single WatchBus.
var diagnosticBus = tether.NewBus[Entry]()

// Subscribe bridges a handler's Diagnostics bus to the shared
// diagnosticBus. Call once per handler from main after creation.
// Events are logged at Warn level and republished for UI display.
func Subscribe(ctx context.Context, name string, bus *tether.Bus[tether.Diagnostic]) {
	bus.SubscribeAsync(ctx, func(d tether.Diagnostic) {
		slog.Warn("diagnostic",
			"handler", name, "kind", d.Kind,
			"session", d.SessionID, "error", d.Err,
			"detail", d.Detail)
		diagnosticBus.Publish(Entry{
			Kind:      string(d.Kind),
			SessionID: d.SessionID,
			Detail:    d.Detail,
		})
	})
}

// State is the per-session state for the diagnostics page.
type State struct {
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
	// Events holds the most recent diagnostic events, newest first.
	Events []Entry
}

var presence = shared.NewPresenceCountOnly()

// New creates a WebSocket handler for the diagnostics page.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "diagnostics",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: presence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/diagnostics/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Diagnostics"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Watchers: append(shared.Watchers[State](presence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),
			// Push diagnostic events into session state so the
			// feed updates in real time.
			tether.WatchBus(diagnosticBus, func(e Entry, s State) State {
				s.Events = append([]Entry{e}, s.Events...)
				if len(s.Events) > 20 {
					s.Events = s.Events[:20]
				}
				return s
			}),
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("diagnostics: connected", "id", sess.ID())
			shared.TrackPresence(presence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("diagnostics: disconnected", "id", sess.ID())
			shared.UntrackPresence(presence, sess.ID())
		},
	})
}

// Handle processes events on the diagnostics page. The trigger
// actions deliberately cause framework-level failures so the
// diagnostic pipeline is visible in the UI.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "diag.trigger-panic":
		// The framework recovers this panic and emits a
		// HandlerPanic diagnostic - the event appears in the
		// feed without crashing the session.
		panic("deliberate panic to demonstrate HandlerPanic diagnostic")
	}
	return s
}
