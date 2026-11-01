package handler

import (
	"log/slog"
	"net/http"
	"time"

	htmx "github.com/jpl-au/fluent-htmx/htmx2"
	"github.com/jpl-au/fluent-htmx/htmx2/event"
	"github.com/jpl-au/fluent-htmx/htmx2/swap"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/text"

	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/composite/logentry"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/generate"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/layout"
)

// SSEPage renders the SSE live log page. The page loads the htmx-ext-sse
// extension which opens an EventSource to /sse/feed. Incoming HTML
// fragments are swapped into the log container by the extension.
func SSEPage(w http.ResponseWriter, _ *http.Request) {
	// The status span listens for "status" SSE events so htmx can
	// swap the text from "connecting..." to "connected" when the
	// server sends the first event.
	statusSpan := span.Text("connecting...").Class("status-connecting").ID("sse-status")
	htmx.New(statusSpan).SSESwap("status").HxSwap(swap.InnerHTML)

	// The log container listens for "message" SSE events (the
	// default event name). Each incoming entry's HTML is prepended
	// via afterbegin so the newest entry is always at the top.
	logDiv := div.New().Class("log-feed").ID("sse-log")
	htmx.New(logDiv).SSESwap("message").HxSwap(swap.AfterBegin)

	sseCard := card.New("SSE Live Log",
		div.New(
			span.Static("Status: ").Class("label"),
			statusSpan,
		).Class("layout-row"),
		div.New(
			span.Text(
				"Connected to the server via Server-Sent Events. Log entries are "+
					"pushed as HTML fragments and swapped in by the htmx SSE extension. "+
					"SSE is unidirectional (server to client) and reconnects "+
					"automatically if the connection drops.",
			).Class("hint"),
		),
		logDiv,
	)

	// The sse-connect div wraps the card so htmx opens the
	// EventSource connection to the server. The extension fires its
	// lifecycle events on this element, so HxOn binds the status
	// updates here. After an error it reconnects by itself; it fires
	// close only when the stream ends for good.
	content := div.New(sseCard)
	htmx.New(content).
		HxExt("sse").
		SSEConnect("/sse/feed").
		HxOn(event.SSEOpen, setStatus("sse-status", "connected", "status-connected")).
		HxOn(event.SSEError, setStatus("sse-status", "reconnecting...", "status-reconnecting")).
		HxOn(event.SSEClose, setStatus("sse-status", "closed", "status-disconnected"))

	layout.PageWithScripts(w, "SSE", nil,
		[]string{"/static/htmx-ext-sse.js"}, content)
}

// SSEFeed streams rendered HTML log entries as Server-Sent Events.
// The stream runs until the client disconnects. Each entry is an HTML
// fragment that htmx swaps into the page.
func SSEFeed(w http.ResponseWriter, r *http.Request) {
	sse, err := htmx.NewSSE(w)
	if err != nil {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	slog.Info("sse: client connected", "remote", r.RemoteAddr)

	// Send a status update as the first event so the page shows "connected"
	// instead of "connecting...". The event name "status" matches the
	// sse-swap="status" on the status span.
	if err := sse.Send("status", text.Text("connected")); err != nil {
		slog.Warn("sse: failed to send status", "remote", r.RemoteAddr, "error", err)
		return
	}

	// Push log entries until the client disconnects. The request
	// context is cancelled when the client closes the connection.
	ctx := r.Context()
	for {
		// The "message" event matches the sse-swap="message" attribute on
		// the log container.
		if err := sse.Send("message", logentry.New(generate.LogEntry())); err != nil {
			slog.Warn("sse: failed to send log entry", "remote", r.RemoteAddr, "error", err)
			return
		}

		select {
		case <-ctx.Done():
			slog.Info("sse: client disconnected", "remote", r.RemoteAddr)
			return
		case <-time.After(generate.Jitter()):
		}
	}
}
