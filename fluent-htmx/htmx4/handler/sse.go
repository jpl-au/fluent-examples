package handler

import (
	"log/slog"
	"net/http"
	"time"

	htmx "github.com/jpl-au/fluent-htmx/htmx4"
	"github.com/jpl-au/fluent-htmx/htmx4/event"
	"github.com/jpl-au/fluent-htmx/htmx4/swap"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"

	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/composite/logentry"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/generate"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/layout"
)

// SSEPage renders the SSE live log page. The log container opens the
// stream with hx-sse:connect, which the htmax.js bundle provides.
// Each unnamed event from the server is swapped into the container
// with afterbegin, so the newest entry is at the top. The lifecycle
// events fire on the connecting element, so HxOn binds the status
// updates there.
func SSEPage(w http.ResponseWriter, _ *http.Request) {
	logDiv := div.New().Class("log-feed").ID("sse-log")
	htmx.New(logDiv).
		SSEConnect("/sse/feed").
		HxSwap(swap.AfterBegin).
		HxOn(event.SSEAfterConnection, setStatus("sse-status", "connected", "status-connected")).
		HxOn(event.SSEClose, setStatus("sse-status", "reconnecting...", "status-reconnecting")).
		HxOn(event.SSEError, setStatus("sse-status", "disconnected", "status-disconnected"))

	content := card.New("SSE Live Log",
		div.New(
			span.Static("Status: ").Class("label"),
			span.Text("connecting...").Class("status-connecting").ID("sse-status"),
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

	layout.Page(w, "SSE", nil, content)
}

// SSEFeed streams rendered HTML log entries as Server-Sent Events.
// The stream runs until the client disconnects. Each entry goes out
// as an unnamed event, which htmx swaps into the connecting element
// with its own hx-swap style.
func SSEFeed(w http.ResponseWriter, r *http.Request) {
	sse, err := htmx.NewSSE(w)
	if err != nil {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	slog.Info("sse: client connected", "remote", r.RemoteAddr)

	// Push log entries until the client disconnects. The request
	// context is cancelled when the client closes the connection.
	ctx := r.Context()
	for {
		if err := sse.Swap(logentry.New(generate.LogEntry())); err != nil {
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
