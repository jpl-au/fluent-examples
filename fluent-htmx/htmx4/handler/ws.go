package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/lxzan/gws"

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

// wsUpgrader is shared across requests. The default configuration
// accepts all origins, which is fine for a local demo app.
var wsUpgrader = gws.NewUpgrader(&wsHandler{}, &gws.ServerOption{})

// WSPage renders the WebSocket live log page. The log container opens
// the connection with hx-ws:connect, which the htmax.js bundle
// provides. A plain HTML message is swapped into the container with
// afterbegin; a JSON message names its own target. The close and
// error events fire on the connecting element, so HxOn binds the
// status updates there. The connected state comes from the server's
// first message instead.
func WSPage(w http.ResponseWriter, _ *http.Request) {
	logDiv := div.New().Class("log-feed").ID("ws-log")
	htmx.New(logDiv).
		WSConnect("/ws/feed").
		HxSwap(swap.AfterBegin).
		HxOn(event.WSClose, setStatus("ws-status", "reconnecting...", "status-reconnecting")).
		HxOn(event.WSError, setStatus("ws-status", "disconnected", "status-disconnected"))

	content := card.New("WebSocket Live Log",
		div.New(
			span.Static("Status: ").Class("label"),
			span.Text("connecting...").Class("status-connecting").ID("ws-status"),
		).Class("layout-row"),
		div.New(
			span.Text(
				"Connected to the server via a WebSocket. Log entries are pushed "+
					"from the server as HTML fragments and swapped in by the htmx "+
					"WebSocket extension. The connection is full-duplex so entries "+
					"arrive with minimal latency.",
			).Class("hint"),
		),
		logDiv,
	)

	layout.Page(w, "WebSocket", nil, content)
}

// WSFeed upgrades the connection to a WebSocket and pushes rendered
// HTML log entries at random intervals. The first message is a JSON
// frame that replaces the status span, and the entries follow as
// plain HTML that htmx swaps into the connecting element.
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

	// A JSON message carries its own target and swap style, so it
	// can update an element other than the connection's own.
	status, err := htmx.WSResponse{
		Content: span.Text("connected").Class("status-connected").ID("ws-status"),
		Target:  "#ws-status",
		Swap:    swap.OuterHTML,
	}.JSON()
	if err != nil {
		slog.Error("ws: failed to encode status", "error", err)
		return
	}
	if err := conn.WriteMessage(gws.OpcodeText, status); err != nil {
		slog.Info("ws: client disconnected on status", "remote", r.RemoteAddr)
		return
	}

	// Push log entries until the connection drops. Plain HTML is
	// swapped with the connecting element's own hx-target and
	// hx-swap, so each entry lands at the top of the log.
	for {
		entry := logentry.New(generate.LogEntry())
		if err := conn.WriteMessage(gws.OpcodeText, entry.RenderBytes()); err != nil {
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
