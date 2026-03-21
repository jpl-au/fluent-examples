package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jpl-au/fluent-examples/fluent-jit/generate"
	"github.com/jpl-au/fluent-examples/fluent-jit/layout"
)

// SSEPage renders the SSE live log page. The page loads sse.js which
// opens an EventSource to /sse/feed and appends entries to the log
// container as they arrive.
func SSEPage(w http.ResponseWriter, _ *http.Request) {
	content := logCard(
		"SSE Live Log",
		"Connected to the server via Server-Sent Events. Log entries are "+
			"pushed as SSE data frames and rendered by client-side JavaScript. "+
			"SSE is unidirectional (server to client) and reconnects automatically "+
			"if the connection drops.",
		"sse-status",
		"sse-log",
	)
	layout.PageWithScripts(w, "SSE", nil, []string{"/static/sse.js"}, content)
}

// SSEFeed streams fake log entries as Server-Sent Events. The stream
// runs until the client disconnects. Each entry is a JSON-encoded
// data frame so the client can parse it consistently with the
// WebSocket variant.
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

	// Push log entries until the client disconnects. The request
	// context is cancelled when the client closes the connection.
	ctx := r.Context()
	for {
		entry := generate.LogEntry()
		data, err := json.Marshal(entry)
		if err != nil {
			slog.Warn("sse: marshal failed", "error", err)
			continue
		}

		// SSE data frame format: "data: <payload>\n\n"
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()

		select {
		case <-ctx.Done():
			slog.Info("sse: client disconnected", "remote", r.RemoteAddr)
			return
		case <-time.After(generate.Jitter()):
		}
	}
}
