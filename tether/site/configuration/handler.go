package configuration

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
	ws "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
	"github.com/jpl-au/fluent-examples/tether/store"
)

// PageView is emitted by OnNavigate on every page load. A global
// subscriber counts these - including views from the initial GET
// (pre-warming), where CaptureSession's synchronous enqueue ensures
// the emission reaches subscribers immediately.
type PageView struct {
	Path string
}

// pageViewBus routes page-view events to global subscribers.
var pageViewBus = tether.NewBus[PageView]()

// pageViewCount is incremented by a global subscriber on every
// page view - including pre-warm views from the initial GET.
var pageViewCount atomic.Int64

// Setup wires the global page-view subscriber. Call from main with
// the root context so the subscription is cancelled on shutdown.
func Setup(ctx context.Context) {
	pageViewBus.Subscribe(ctx, func(pv PageView) {
		n := pageViewCount.Add(1)
		slog.Debug("page view", "path", pv.Path, "total", n)
	})
}

// State is the per-session state for the configuration demo.
type State struct {
	// PageViews is the total page-view count at render time.
	PageViews int64
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
}

var presence = shared.NewPresenceCountOnly()

// Configured values are stored at package level so view.go can
// reference them in the rendered cards.
var (
	configuredCompression = ws.Compression{
		Level:           ws.CompressionFastest,
		Threshold:       1024,
		ContextTakeover: true,
	}
	configuredTimeouts = tether.Timeouts{
		Idle:              5 * time.Minute,
		MaxLifetime:       30 * time.Minute,
		Reconnect:         15 * time.Second,
		Pending:           20 * time.Second,
		ShutdownGrace:     15 * time.Second,
		Heartbeat:         25 * time.Second,
		Retry:             time.Second,
		MaxRetry:          20 * time.Second,
		BackoffMultiplier: 2.0,
		DisableJitter:     false,
	}
	configuredLimits = tether.Limits{
		MaxSessions:   100,
		MaxPending:    64,
		CmdBufferSize: 128,
		MaxEventBytes: 128 << 10, // 128 KB
	}
	configuredSecurity = tether.Security{
		TrustedOrigins: []string{"http://localhost:8080"},
	}

	// configuredSessionStore persists session state to disk so sessions
	// survive server restarts. The framework serialises the State
	// struct via the Codec (default CBOR), saves on disconnect and
	// graceful shutdown, and restores when a reconnecting client
	// reaches a server with no in-memory session.
	configuredSessionStore = store.NewFileSessionStore("/tmp/tether-sessions")

	// configuredDiffStore offloads differ snapshots to disk during the
	// reconnect window, freeing Go memory for disconnected sessions.
	// The framework saves on disconnect and deletes on reconnect
	// (Render re-seeds the differ).
	configuredDiffStore = store.NewFileDiffStore("/tmp/tether-diffs")
)

// New creates a handler demonstrating the Timeouts, Limits, and
// Security configuration fields. The page is informational - the
// handler uses non-default values so each field is visible in the
// rendered cards.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name: "configuration",
		Mode: mode.WebSocket,
		Upgrade: ws.Upgrade(ws.Options{
			ReadLimit:   128 << 10,
			Compression: configuredCompression,
		}),

		Timeouts:     configuredTimeouts,
		Limits:       configuredLimits,
		SessionStore: configuredSessionStore,
		DiffStore:    configuredDiffStore,

		InitialState: func(_ *http.Request) State {
			return State{
				PageViews:   pageViewCount.Load(),
				OnlineCount: presence.OnlineCount.Load(),
			}
		},
		OnNavigate: func(sess tether.Session, s State, params tether.Params) State {
			// Emit fires during pre-warming (initial GET) because
			// CaptureSession.enqueue runs synchronously - the global
			// subscriber increments the counter before the page renders.
			pageViewBus.Emit(sess, PageView{Path: params.Path})
			s.PageViews = pageViewCount.Load()
			return s
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/configuration/", s.OnlineCount, Render(s))
		},
		Handle: func(_ tether.Session, s State, _ tether.Event) State { return s },

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Configuration"),
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
			slog.Info("configuration: connected", "id", sess.ID())
			shared.TrackPresence(presence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("configuration: disconnected", "id", sess.ID())
			shared.UntrackPresence(presence, sess.ID())
		},

		// OnRestore fires instead of OnConnect when a session is
		// recovered from the SessionStore (e.g. after a server restart).
		// Re-establish any runtime resources the session needs - here
		// we rejoin presence tracking and log the recovery.
		OnRestore: func(sess *tether.StatefulSession[State]) {
			slog.Info("configuration: restored from store", "id", sess.ID(), "pageViews", sess.State().PageViews)
			shared.TrackPresence(presence, sess.ID())
		},
	})
}
