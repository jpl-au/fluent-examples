// Package app wires together the tether example application. It
// creates all handlers, bus subscribers, and routes, returning a
// ready-to-serve HTTP handler and the list of drainables for graceful
// shutdown. Both main.go and the playwright integration tests use
// this package to ensure they exercise the same server configuration.
package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jpl-au/chain"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/middleware"
	"github.com/jpl-au/fluent-examples/tether/site/broadcasting"
	"github.com/jpl-au/fluent-examples/tether/site/chat"
	"github.com/jpl-au/fluent-examples/tether/site/clipboard"
	"github.com/jpl-au/fluent-examples/tether/site/components"
	"github.com/jpl-au/fluent-examples/tether/site/configuration"
	"github.com/jpl-au/fluent-examples/tether/site/diagnostics"
	"github.com/jpl-au/fluent-examples/tether/site/dragdrop"
	"github.com/jpl-au/fluent-examples/tether/site/errors"
	"github.com/jpl-au/fluent-examples/tether/site/events"
	"github.com/jpl-au/fluent-examples/tether/site/freeze"
	"github.com/jpl-au/fluent-examples/tether/site/groups"
	"github.com/jpl-au/fluent-examples/tether/site/hotkey"
	httpsite "github.com/jpl-au/fluent-examples/tether/site/http"
	"github.com/jpl-au/fluent-examples/tether/site/live"
	"github.com/jpl-au/fluent-examples/tether/site/memo"
	"github.com/jpl-au/fluent-examples/tether/site/morph"
	mwsite "github.com/jpl-au/fluent-examples/tether/site/mw"
	"github.com/jpl-au/fluent-examples/tether/site/navigation"
	"github.com/jpl-au/fluent-examples/tether/site/notifications"
	"github.com/jpl-au/fluent-examples/tether/site/realtime"
	"github.com/jpl-au/fluent-examples/tether/site/rendering"
	"github.com/jpl-au/fluent-examples/tether/site/scroll"
	"github.com/jpl-au/fluent-examples/tether/site/selection"
	"github.com/jpl-au/fluent-examples/tether/site/signals"
	swsite "github.com/jpl-au/fluent-examples/tether/site/sw"
	swhandler "github.com/jpl-au/fluent-examples/tether/site/sw/handler"
	"github.com/jpl-au/fluent-examples/tether/site/touch"
	"github.com/jpl-au/fluent-examples/tether/site/uploads"
	filteredupload "github.com/jpl-au/fluent-examples/tether/site/uploads/filtered"
	"github.com/jpl-au/fluent-examples/tether/site/valuestore"
)

// New creates the complete example application. The context controls
// the lifetime of bus subscribers and background goroutines - cancel
// it to clean up. The assets parameter is the shared [tether.Asset]
// for stylesheets and scripts - the caller creates it from the
// embedded filesystem. Returns the HTTP handler (a chain.Mux with all
// routes) and the list of tether handlers that need draining on
// shutdown.
func New(ctx context.Context, assets *tether.Asset) (http.Handler, []tether.Drainable) {

	app := tether.App{
		DevMode: true,
		Assets:  []*tether.Asset{assets},
	}

	// Wire bus subscribers before creating handlers so subscribers
	// are in place when OnNavigate fires during pre-warming.
	broadcasting.Setup(ctx)
	configuration.Setup(ctx)
	swhandler.SetupPush()

	// Stateless HTTP features (tether.Stateless).
	httpHandler := httpsite.New(app, assets)
	eventsHandler := events.New(app, assets)
	errorsHandler := errors.New(app, assets)
	navigationHandler := navigation.New(app, assets)
	renderingHandler := rendering.New(app, assets)
	morphHandler := morph.New(app, assets)
	middlewareHandler := mwsite.New(app, assets)
	clipboardHandler := clipboard.New(app, assets)
	selectionHandler := selection.New(app, assets)
	touchHandler := touch.New(app, assets)

	// WebSocket features (tether.Handler).
	notificationsHandler := notifications.New(app, assets)
	uploadsHandler := uploads.New(app, assets)
	broadcastingHandler := broadcasting.New(app, assets)
	componentsHandler := components.New(app, assets)
	chatHandler := chat.New(app, assets)
	filteredUploadHandler := filteredupload.New(app, assets)
	configurationHandler := configuration.New(app, assets)
	valuestoreHandler := valuestore.New(app, assets)
	groupsHandler := groups.New(app, assets)

	// Hotkey, drag-and-drop, and scroll demos.
	hotkeyHandler := hotkey.New(app, assets)
	dragdropHandler := dragdrop.New(app, assets)
	scrollHandler := scroll.New(app, assets)

	// Freeze demo - FreezeWithConnect with SessionStore.
	freezeHandler := freeze.New(app, assets)

	// Memo demos - subtree memoisation with Versioned keys.
	memoHandler := memo.New(app, assets)
	memoRealtimeHandler := memo.NewRealtime(app, assets)

	// Signal demos - WS and SSE variants.
	signalsWSHandler := signals.NewWS(app, assets)
	signalsSSEHandler := signals.NewSSE(app, assets)

	// Live updates demos - WS and SSE variants.
	liveWSHandler := live.NewWS(app, assets)
	liveSSEHandler := live.NewSSE(app, assets)

	// Real-time dashboard (go-echarts system monitor).
	realtimeHandler := realtime.New(app, assets)

	// Service Worker section.
	swHandler := swsite.New(app, assets)

	// Diagnostics page - aggregates events from all live handlers.
	diagnosticsHandler := diagnostics.New(app, assets)

	// Subscribe handler diagnostics to the shared diagnostics bus
	// so events from any handler appear on the diagnostics page.
	diagnostics.Subscribe(ctx, "diagnostics", diagnosticsHandler.Diagnostics)
	diagnostics.Subscribe(ctx, "signals/ws", signalsWSHandler.Diagnostics)
	diagnostics.Subscribe(ctx, "configuration", configurationHandler.Diagnostics)
	diagnostics.Subscribe(ctx, "notifications", notificationsHandler.Diagnostics)
	diagnostics.Subscribe(ctx, "uploads", uploadsHandler.Diagnostics)
	diagnostics.Subscribe(ctx, "groups", groupsHandler.Diagnostics)

	mux := chain.New()
	mux.Use(middleware.RequestLogging)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(signalsWSHandler.Health()); err != nil {
			slog.Warn("health: encode failed", "error", err)
		}
	})

	mux.Handle("/events", eventsHandler)
	mux.Handle("/errors", errorsHandler)
	mux.Handle("/navigation/", navigationHandler)
	mux.Handle("/rendering", renderingHandler)
	mux.Handle("/notifications/", notificationsHandler)
	mux.HandleFunc("GET /uploads/files/", uploads.ServeFile)
	mux.Handle("/uploads/", uploadsHandler)
	mux.Handle("/broadcasting/", broadcastingHandler)
	mux.Handle("/components/", componentsHandler)
	mux.Handle("/chat/", chatHandler)
	mux.Handle("/uploads/filtered/", filteredUploadHandler)
	mux.Handle("/configuration/", configurationHandler)
	mux.Handle("/valuestore/", valuestoreHandler)
	mux.Handle("/groups/", groupsHandler)
	mux.Handle("/morph", morphHandler)
	mux.Handle("/middleware", middlewareHandler)
	mux.Handle("/signals/ws/", signalsWSHandler)
	mux.Handle("/signals/sse/", signalsSSEHandler)
	mux.Handle("/live/ws/", liveWSHandler)
	mux.Handle("/live/sse/", liveSSEHandler)
	mux.Handle("/realtime/", realtimeHandler)
	mux.Handle("/diagnostics/", diagnosticsHandler)
	mux.Handle("/freeze/", freezeHandler)
	mux.Handle("/clipboard", clipboardHandler)
	mux.Handle("/selection", selectionHandler)
	mux.Handle("/touch", touchHandler)
	mux.Handle("/hotkey/", hotkeyHandler)
	mux.Handle("/dragdrop/", dragdropHandler)
	mux.Handle("/scroll/", scrollHandler)
	mux.Handle("/memo/realtime/", memoRealtimeHandler)
	mux.Handle("/memo/", memoHandler)
	mux.Handle("/sw/", swHandler)

	// HTTP section as the catch-all (must be registered last).
	mux.Handle("/", httpHandler)

	drainables := []tether.Drainable{
		notificationsHandler, uploadsHandler, broadcastingHandler,
		componentsHandler, chatHandler, filteredUploadHandler,
		configurationHandler, valuestoreHandler, groupsHandler,
		signalsWSHandler, signalsSSEHandler,
		liveWSHandler, liveSSEHandler,
		realtimeHandler, diagnosticsHandler, freezeHandler,
		hotkeyHandler, dragdropHandler, scrollHandler, memoHandler, memoRealtimeHandler, swHandler,
	}

	return mux, drainables
}
