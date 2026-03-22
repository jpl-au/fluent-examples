package page

import (
	"github.com/jpl-au/fluent/node"

	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/composite/toc"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/site/http/state"
)

// OverviewRender builds the landing page with a table of contents
// covering all sections.
func OverviewRender(_ state.State) node.Node {
	return cpage.New(
		panel.Card(
			"Welcome",
			"tether is a Go framework for building interactive web applications with server-side rendering, real-time updates, and client-side reactivity.",
			"", panel.AllTransports,
			hint.Text("This feature explorer is split into sections by concern. Each section demonstrates a different aspect of the framework."),
		),

		panel.Card(
			"Stateless - tether.Stateless",
			"Stateless pages served via plain HTTP. No WebSocket, no SSE. State is reconstructed from each request. Events are sent as fetch POST requests and the response carries the HTML update.",
			"tether.Stateless", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/events", "Events & Forms"), "Click, Input, Submit, Change, KeyDown, Focus, Confirm, Throttle, EventData, Indicator, On, Reset, Viewport"),
				toc.Item(toc.Link("/rendering", "State & Rendering"), "Dynamic keys, tether.Catch, tether.Component"),
				toc.Item(toc.Link("/morph", "Full-Page Morph"), "Rendering without Dynamic keys, idiomorph fallback"),
				toc.Item(toc.Link("/errors", "Error Boundaries"), "tether.Catch, panic recovery"),
				toc.Item(toc.Link("/navigation", "Navigation"), "Navigate, ReplaceURL, SetTitle"),
				toc.Item(toc.Link("/middleware", "Middleware"), "Chain, Guard, Timing, custom middleware"),
			),
		),

		panel.Card(
			"Signals & Directives",
			"Client-side reactivity powered by server-pushed signals. Elements bind to signal values and update instantly without a full re-render.",
			"bind.BindText · bind.BindShow · bind.SetSignal", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/signals/ws/", "WebSocket"), "BindText, BindShow, SetSignal, ToggleClass, Optimistic, Hook, Transition, FocusTrap, Cloak, Permanent"),
				toc.Item(toc.Link("/signals/sse/", "SSE"), "BindText, BindShow, SetSignal, ToggleClass, Optimistic"),
			),
		),

		panel.Card(
			"Live Updates",
			"Real-time features that require a persistent connection: groups, broadcasting, shared values, and server-pushed rendering.",
			"tether.Handler", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/live/ws/", "WebSocket"), "Group, Bus, Value, Observe, Go"),
				toc.Item(toc.Link("/live/sse/", "SSE"), "Group, Bus, Value, Observe, Go"),
			),
		),

		panel.Card(
			"Features",
			"Functionality that works with any live transport. Each feature uses the persistent session for real-time feedback.",
			"", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/notifications/", "Notifications"), "Toast, Flash, Announce, Signal"),
				toc.Item(toc.Link("/uploads/", "File Uploads"), "Upload, UploadInput, UploadProgress"),
				toc.Item(toc.Link("/uploads/filtered/", "Filtered Uploads"), "UploadConfig.Accept, MIME filtering"),
				toc.Item(toc.Link("/broadcasting/", "Broadcasting"), "bus.Emit, WatchBus, SubscribeAsync"),
				toc.Item(toc.Link("/components/", "Components"), "tether.Component, Mount, EqualComponent"),
				toc.Item(toc.Link("/chat/", "Chat Room"), "Component, Mounter, Bus, WatchBus"),
				toc.Item(toc.Link("/realtime/", "Real-time Dashboard"), "sess.Go, sess.Update, go-echarts"),
				toc.Item(toc.Link("/configuration/", "Configuration"), "tether.Value, WatchValue, configuration sync"),
				toc.Item(toc.Link("/valuestore/", "Value Store"), "tether.Value, Store, Update, WatchValue"),
				toc.Item(toc.Link("/groups/", "Groups"), "tether.Group, Broadcast, BroadcastOthers, OnJoin, OnLeave"),
				toc.Item(toc.Link("/freeze/", "Freeze & Restore"), "FreezeMode, SessionStore"),
			),
		),

		panel.Card(
			"Observability",
			"Runtime diagnostics and monitoring for debugging and operational visibility.",
			"", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/diagnostics/", "Diagnostics"), "DiagnosticKind, handler_panic, upload_rejected, transport_error"),
			),
		),

		panel.Card(
			"Service Worker",
			"Push notifications, asset caching, and offline support via the browser's Service Worker API. Built on top of mode.Both with Worker: true.",
			"tether.Handler · Worker", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/sw/", "Overview"), "Service Worker architecture"),
				toc.Item(toc.Link("/sw/push", "Push Notifications"), "PushSubscribe, Push, Notification, NotificationAction"),
				toc.Item(toc.Link("/sw/caching", "Caching & Offline"), "Precache, cache-first assets, offline shell fallback"),
			),
		),
	)
}
