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
				toc.Item(toc.Link("/events/", "Events & Forms"), "Click, Input, Submit, Change, KeyDown, Focus, Confirm, Throttle, DebounceLeading, Outside, EventData, Indicator, Event, Reset, Viewport"),
				toc.Item(toc.Link("/rendering/", "State & Rendering"), "Dynamic keys, tether.Component, tether.RouteTyped"),
				toc.Item(toc.Link("/morph/", "Full-Page Morph"), "Rendering without Dynamic keys, idiomorph fallback"),
				toc.Item(toc.Link("/errors/", "Error Boundaries"), "tether.Catch, panic recovery"),
				toc.Item(toc.Link("/navigation/", "Navigation"), "bind.Link, Params, ReplaceURL, Navigate, Prefetch"),
				toc.Item(toc.Link("/view-transitions/", "View Transitions"), "tether.Client{ViewTransitions}, view-transition-name"),
				toc.Item(toc.Link("/middleware/", "Middleware"), "Chain, Guard, Timing, custom middleware"),
				toc.Item(toc.Link("/html-wire/", "HTML Wire Format"), "WireFormat: wire.HTML, sess.Morph fragments, CacheControl"),
				toc.Item(toc.Link("/client-actions/", "Client-Side Actions"), "CopyToClipboard, FlashText, FlashClass"),
				toc.Item(toc.Link("/selection/", "Multi-Select"), "Selectable, CollectSelected"),
				toc.Item(toc.Link("/touch/", "Touch Gestures"), "OnSwipe, OnLongPress"),
				toc.Item(toc.Link("/security/", "Security"), "security.HTML sanitisation, security.Nonce with CSP"),
			),
		),

		panel.Card(
			"Signals & Directives",
			"Client-side reactivity powered by server-pushed signals. Elements bind to signal values and update instantly without a full re-render.",
			"bind.Text · bind.Show · bind.SetSignal", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/signals/ws/", "WebSocket"), "Text, Show, SetSignal, ToggleSignal, Computed, ShowWhen, ToggleTarget, Class, Attr, Value, Optimistic, Hook, Transition, FocusTrap, Cloak, Permanent"),
				toc.Item(toc.Link("/signals/sse/", "SSE"), "Text, Show, SetSignal, ToggleSignal, Optimistic"),
				toc.Item(toc.Link("/timer/", "Client-Side Timers"), "Timer, Countdown, TimerPrecision, TimerFormat, TimerOnComplete"),
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
				toc.Item(toc.Link("/hotkey/", "Hotkeys"), "bind.Hotkey, the platform-aware mod modifier"),
				toc.Item(toc.Link("/dragdrop/", "Drag and Drop"), "bind.Draggable, bind.Sortable, cross-session sync"),
				toc.Item(toc.Link("/scroll/", "Scroll"), "bind.ScrollTo, sess.ScrollTo, bind.PreserveScroll"),
			),
		),

		panel.Card(
			"Performance",
			"Rendering optimisations for expensive pages: skip unchanged subtrees, re-render single regions, and keep the DOM small for large datasets.",
			"jit.Memoise · sess.Patch", panel.AllTransports,
			toc.List(
				toc.Item(toc.Link("/memoise/", "Memoisation"), "jit.Memoise, tether.Versioned, StatefulConfig.Memoise"),
				toc.Item(toc.Link("/memoise/realtime/", "Memoised Dashboard"), "jit.Memoise and sess.Patch together on live charts"),
				toc.Item(toc.Link("/windowing/", "Windowing"), "Virtual scrolling, slice rendering for large lists"),
				toc.Item(toc.Link("/patch/", "Targeted Updates"), "sess.Patch, single-key re-renders"),
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
				toc.Item(toc.Link("/sw/lifecycle", "PWA Lifecycle"), "Online/offline connectivity, appinstalled events"),
			),
		),
	)
}
