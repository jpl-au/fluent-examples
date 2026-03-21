package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/lxzan/gws"

	"github.com/jpl-au/fluent-examples/fluent-jit/generate"
	"github.com/jpl-au/fluent-examples/fluent-jit/layout"
)

// wsUpgrader is shared across requests. The default configuration
// accepts all origins, which is fine for a local demo app.
var wsUpgrader = gws.NewUpgrader(&wsHandler{}, &gws.ServerOption{})

// WSPage renders the WebSocket live log page. The page loads ws.js
// which connects back to /ws/feed and appends entries to the log
// container as they arrive.
func WSPage(w http.ResponseWriter, _ *http.Request) {
	content := logCard(
		"WebSocket Live Log",
		"Connected to the server via a WebSocket. Log entries are pushed "+
			"from the server as JSON and rendered by client-side JavaScript. "+
			"The connection is full-duplex so entries arrive with minimal latency.",
		"ws-status",
		"ws-log",
	)
	layout.PageWithScripts(w, "WebSocket", nil, []string{"/static/ws.js"}, content)
}

// WSFeed upgrades the connection to a WebSocket and pushes fake log
// entries at random intervals. The feed goroutine starts only after
// the upgrade succeeds and stops when the connection closes.
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

	// Push log entries until the connection drops. The timer only
	// runs while a client is connected.
	for {
		entry := generate.LogEntry()
		data, err := json.Marshal(entry)
		if err != nil {
			slog.Warn("ws: marshal failed", "error", err)
			continue
		}
		if err := conn.WriteMessage(gws.OpcodeText, data); err != nil {
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
