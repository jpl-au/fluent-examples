// Package sw provides the Service Worker section of the feature
// explorer, demonstrating push notifications, asset caching, and
// offline support via the browser's Service Worker API.
package sw

import (
	"log/slog"
	"time"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
	"github.com/jpl-au/tether/sse"
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/middleware"
	"github.com/jpl-au/fluent-examples/tether/site/sw/handler"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// New creates a tether.Handler for the Service Worker section,
// configured with mode.Both transport, the full service worker,
// push notifications, and the SW router.
func New(app tether.App, assets *tether.Asset) *tether.Handler[state.State] {
	r := newRouter()

	return tether.Stateful(app, tether.StatefulConfig[state.State]{
		Mode:     mode.Both,
		Upgrade:  wsupgrade.Upgrade(),
		Fallback: sse.Upgrade(),
		Worker:   true,
		// Name distinguishes this handler in startup logs - without it the
		// transport label (ws+sse) would be identical to the main ws handler.
		// Worker: true is the real differentiator; Name makes it human-readable.
		Name: "sw",

		InitialState: handler.InitialState,
		Render: func(s state.State) node.Node {
			return layout.Shell(layout.SectionSW, s.Page, s.OnlineCount, r.Render(s))
		},
		Handle:     r.Handle,
		OnNavigate: r.OnNavigate(func(s *state.State, p tether.Params) { s.Page = p.Path }),

		Middleware: []tether.Middleware[state.State]{middleware.Logging[state.State]},

		Layout: func(_ state.State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Service Worker"),
					assets.Stylesheet("app.css"),
				),
				body.New(content, assets.Script("hooks.js")),
			).Lang("en")
		},

		OnConnect:    handler.OnConnect,
		OnDisconnect: handler.OnDisconnect,
		OnStructuralChange: func(sess *tether.StatefulSession[state.State], change tether.StructuralChange) {
			slog.Warn("structural change",
				"session", sess.ID(),
				"added", change.Added,
				"removed", change.Removed,
				"reordered", change.Reordered,
				"bytes", change.Bytes,
			)
		},
		OnNoPatch: func(sess *tether.StatefulSession[state.State], info tether.NoPatch) {
			if info.Source == "update" {
				slog.Debug("no-patch update", "session", sess.ID())
				return
			}
			slog.Warn("no patches produced",
				"session", sess.ID(),
				"source", info.Source,
				"action", info.Action,
			)
		},
		Groups:   []*tether.Group[state.State]{handler.Group},
		Watchers: handler.Watchers(),
		Push:     handler.PushConfig(),

		// Generous timeouts for a demo app where users may leave tabs open.
		Timeouts: tether.Timeouts{Idle: 10 * time.Minute, Reconnect: 30 * time.Second},
		// 128 KB event limit prevents accidental payload bloat during demos.
		Limits: tether.Limits{MaxEventBytes: 128 << 10},
	})
}
