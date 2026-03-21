package page

import (
	"github.com/jpl-au/fluent/node"

	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// CachingRender builds the caching and offline page, explaining how
// StatefulConfig.Worker enables cache-first asset serving and offline page
// shells via the service worker.
func CachingRender(_ state.State) node.Node {
	return cpage.New(
		panel.Card(
			"Asset Precaching",
			"When the service worker installs, it pre-caches every file listed in the Asset's Precache slice. These assets are served cache-first on subsequent requests, so CSS, JS, and images load instantly without a network round-trip.",
			"tether.Asset · Precache", panel.AllTransports,
			hint.Text("This app pre-caches app.css and hooks.js. Open DevTools > Application > Cache Storage to see the cached entries."),
		),

		panel.Card(
			"Cache-First Static Assets",
			"All requests matching the asset prefix (/static/) are served cache-first. If the asset is in the cache, the service worker returns it immediately. If not, it fetches from the network and caches the response for next time.",
			"Worker: true", panel.AllTransports,
			hint.Text("Disable your network in DevTools and reload - the page still loads with cached assets."),
		),

		panel.Card(
			"Network-First Navigation",
			"Navigation requests (HTML pages) use a network-first strategy. The service worker tries the network first and falls back to the cached shell when offline. This ensures you always get the latest content when connected.",
			"Worker: true", panel.AllTransports,
			hint.Text("The offline shell is the initial page HTML cached on first load. When offline, it serves as a fallback so the app doesn't show a browser error."),
		),

		panel.Card(
			"Background Sync",
			"In SSE mode, when the network drops, events are queued in IndexedDB and replayed when connectivity returns. The service worker's Background Sync API handles the retry automatically. Enable this with the BackgroundSync client config option.",
			"BackgroundSync: true", panel.SSE,
			hint.Text("This feature is specific to SSE mode where events are sent via POST. WebSocket connections handle reconnection at the transport level."),
		),
	)
}
