package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/lxzan/gws"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/html5/span"

	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/logentry"
	"github.com/jpl-au/fluent-examples/fluent-htmx/generate"
	"github.com/jpl-au/fluent-examples/fluent-htmx/layout"
)

// wsUpgrader is shared across requests. The default configuration
// accepts all origins, which is fine for a local demo app.
var wsUpgrader = gws.NewUpgrader(&wsHandler{}, &gws.ServerOption{})

// WSPage renders the WebSocket live log page. The page loads the
// htmx-ext-ws extension which connects to /ws/feed and swaps HTML
// fragments into the log container as they arrive.
func WSPage(w http.ResponseWriter, _ *http.Request) {
	wsCard := card.New("WebSocket Live Log",
		div.New(
			span.Static("Status: ").Class("label"),
			span.Text("connecting...").Class("status-connecting").ID("ws-status"),
		).Class("layout-row"),
		div.New(
			span.Static(
				"Connected to the server via a WebSocket. Log entries are pushed "+
					"from the server as HTML fragments and swapped in by the htmx "+
					"WebSocket extension. The connection is full-duplex so entries "+
					"arrive with minimal latency.",
			).Class("hint"),
		),
		div.New().Class("log-feed").ID("ws-log"),
	)

	// Listen for htmx WebSocket lifecycle events to update the
	// connection status. The extensions fire these on the element
	// that carries ws-connect.
	statusScript := script.Static(`
document.body.addEventListener("htmx:wsOpen", function() {
  var s = document.getElementById("ws-status");
  s.textContent = "connected";
  s.className = "status-connected";
});
document.body.addEventListener("htmx:wsClose", function() {
  var s = document.getElementById("ws-status");
  s.textContent = "reconnecting...";
  s.className = "status-reconnecting";
});
document.body.addEventListener("htmx:wsError", function() {
  var s = document.getElementById("ws-status");
  s.textContent = "disconnected";
  s.className = "status-disconnected";
});`)

	// The ws-connect div wraps the card so htmx opens the WebSocket
	// and processes incoming HTML fragments as OOB swaps.
	content := div.New(wsCard, statusScript)
	content.SetAttribute("hx-ext", "ws")
	content.SetAttribute("ws-connect", "/ws/feed")

	layout.PageWithScripts(w, "WebSocket", nil,
		[]string{"/static/htmx-ext-ws.js"}, content)
}

// WSFeed upgrades the connection to a WebSocket and pushes rendered
// HTML log entries at random intervals. Each message is an OOB swap
// fragment that htmx appends to the #ws-log container.
func WSFeed(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r)
	if err != nil {
		slog.Error("ws upgrade failed", "error", err)
		return
	}

	// ReadLoop processes control frames (ping/pong/close) so the
	// connection stays alive and closes cleanly.
	go conn.ReadLoop()

	slog.Info("ws: client connected", "remote", r.RemoteAddr)

	// Send a status update so the page shows "connected" instead of
	// "connecting...". The span carries hx-swap-oob so htmx replaces
	// the existing status element by ID.
	status := span.Text("connected").Class("status-connected").ID("ws-status")
	status.SetAttribute("hx-swap-oob", "true")
	if err := conn.WriteMessage(gws.OpcodeText, status.Render()); err != nil {
		slog.Info("ws: client disconnected on status", "remote", r.RemoteAddr)
		return
	}

	// Push log entries until the connection drops.
	for {
		entry := generate.LogEntry()

		// Wrap the entry in a div targeting #ws-log with afterbegin
		// so htmx prepends (newest at top).
		wrapper := div.New(logentry.New(entry)).ID("ws-log")
		wrapper.SetAttribute("hx-swap-oob", "afterbegin")

		if err := conn.WriteMessage(gws.OpcodeText, wrapper.Render()); err != nil {
			slog.Info("ws: client disconnected", "remote", r.RemoteAddr)
			return
		}
		time.Sleep(generate.Jitter())
	}
}

// wsHandler satisfies the gws.Event interface. We only need the
// connection lifecycle for the log feed, not incoming messages.
type wsHandler struct {
	gws.BuiltinEventHandler
}
