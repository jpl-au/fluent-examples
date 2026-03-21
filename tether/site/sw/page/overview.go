package page

import (
	"github.com/jpl-au/fluent/node"

	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/composite/toc"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// OverviewRender builds the landing page for the Service Worker section,
// explaining the full service worker mode and linking to the push,
// caching, and lifecycle sub-pages.
func OverviewRender(_ state.State) node.Node {
	return cpage.New(
		panel.Card(
			"Service Worker Section",
			"This section demonstrates service worker capabilities: push notifications, asset caching, and offline support. These features are independent of the transport (WebSocket or SSE) - they rely on the browser's Service Worker API.",
			"Worker: true · push.Sender", panel.WS|panel.SSE,
			hint.Text("Setting Worker: true in the handler config registers a full service worker that caches assets, serves offline shells, and enables background sync. Push notifications use the Web Push protocol with VAPID authentication."),
		),

		panel.Card(
			"How it works",
			"The framework embeds a service worker (tether-worker.js) that handles three concerns: asset caching for fast loads and offline access, push event handling for notifications, and background sync for queuing events when offline. All of this is automatic from a single config flag.",
			"", panel.WS|panel.SSE,
		),

		panel.Card(
			"Features in this section",
			"",
			"", panel.WS|panel.SSE,
			featureList(),
		),
	)
}

// featureList builds the table-of-contents links for the SW section.
func featureList() node.Node {
	return toc.List(
		toc.Item(toc.NavLink("/sw/push", "Push Notifications"), "VAPID-authenticated Web Push - even when the tab is closed"),
		toc.Item(toc.NavLink("/sw/caching", "Caching & Offline"), "Cache-first assets, offline shell fallback"),
		toc.Item(toc.NavLink("/sw/lifecycle", "PWA Lifecycle"), "Online/offline connectivity, app installation events"),
	)
}
