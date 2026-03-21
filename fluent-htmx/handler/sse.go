package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/script"
	"github.com/jpl-au/fluent/html5/span"

	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/logentry"
	"github.com/jpl-au/fluent-examples/fluent-htmx/generate"
	"github.com/jpl-au/fluent-examples/fluent-htmx/layout"
)

// SSEPage renders the SSE live log page. The page loads the htmx-ext-sse
// extension which opens an EventSource to /sse/feed. Incoming HTML
// fragments are swapped into the log container by the extension.
func SSEPage(w http.ResponseWriter, _ *http.Request) {
	// The status span listens for "status" SSE events so htmx can
	// swap the text from "connecting..." to "connected" when the
	// server sends the first event.
	statusSpan := span.Text("connecting...").Class("status-connecting").ID("sse-status")
	statusSpan.SetAttribute("sse-swap", "status")
	statusSpan.SetAttribute("hx-swap", "innerHTML")

	// The log container listens for "message" SSE events (the
	// default event name). Each incoming entry's HTML is prepended
	// via afterbegin so the newest entry is always at the top.
	logDiv := div.New().Class("log-feed").ID("sse-log")
	logDiv.SetAttribute("sse-swap", "message")
	logDiv.SetAttribute("hx-swap", "afterbegin")

	sseCard := card.New("SSE Live Log",
		div.New(
			span.Static("Status: ").Class("label"),
			statusSpan,
		).Class("layout-row"),
		div.New(
			span.Static(
				"Connected to the server via Server-Sent Events. Log entries are "+
					"pushed as HTML fragments and swapped in by the htmx SSE extension. "+
					"SSE is unidirectional (server to client) and reconnects "+
					"automatically if the connection drops.",
			).Class("hint"),
		),
		logDiv,
	)

	// Listen for htmx SSE lifecycle events to update the
	// connection status.
	statusScript := script.Static(`
document.body.addEventListener("htmx:sseOpen", function() {
  var s = document.getElementById("sse-status");
  s.textContent = "connected";
  s.className = "status-connected";
});
document.body.addEventListener("htmx:sseClose", function() {
  var s = document.getElementById("sse-status");
  s.textContent = "reconnecting...";
  s.className = "status-reconnecting";
});
document.body.addEventListener("htmx:sseError", function() {
  var s = document.getElementById("sse-status");
  s.textContent = "disconnected";
  s.className = "status-disconnected";
});`)

	// The sse-connect div wraps the card so htmx opens the
	// EventSource connection to the server.
	content := div.New(sseCard, statusScript)
	content.SetAttribute("hx-ext", "sse")
	content.SetAttribute("sse-connect", "/sse/feed")

	layout.PageWithScripts(w, "SSE", nil,
		[]string{"/static/htmx-ext-sse.js"}, content)
}

// SSEFeed streams rendered HTML log entries as Server-Sent Events.
// The stream runs until the client disconnects. Each entry is an HTML
// fragment that htmx swaps into the page.
func SSEFeed(w http.ResponseWriter, r *http.Request) {
	// SSE requires these headers to keep the connection open and
	// prevent buffering by proxies or the Go runtime.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	slog.Info("sse: client connected", "remote", r.RemoteAddr)

	// Send a status update as the first event so the page shows
	// "connected" instead of "connecting...". The event name
	// "status" matches the sse-swap="status" on the status span.
	fmt.Fprintf(w, "event: status\ndata: connected\n\n")
	flusher.Flush()

	// Push log entries until the client disconnects. The request
	// context is cancelled when the client closes the connection.
	ctx := r.Context()
	for {
		entry := generate.LogEntry()
		html := logentry.New(entry).Render()

		// SSE data frame format: "data: <payload>\n\n"
		// The default event name is "message" which matches the
		// sse-swap="message" attribute on the log container.
		fmt.Fprintf(w, "data: %s\n\n", html)
		flusher.Flush()

		select {
		case <-ctx.Done():
			slog.Info("sse: client disconnected", "remote", r.RemoteAddr)
			return
		case <-time.After(generate.Jitter()):
		}
	}
}
