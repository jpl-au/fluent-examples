package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/lxzan/gws"

	htmx "github.com/jpl-au/fluent-htmx/htmx2"
	"github.com/jpl-au/fluent-htmx/htmx2/event"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/span"

	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/composite/logentry"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/composite/row"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/simple/button"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/simple/field"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/generate"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/layout"
)

// wsUpgrader is shared across requests. The default configuration
// accepts all origins, which is fine for a local demo app.
var wsUpgrader = gws.NewUpgrader(&wsHandler{}, &gws.ServerOption{})

// WSPage renders the WebSocket live log page. The page loads the
// htmx-ext-ws extension which connects to /ws/feed and swaps HTML
// fragments into the log container as they arrive.
func WSPage(w http.ResponseWriter, _ *http.Request) {
	// The form sends its fields over the open connection as JSON
	// instead of making a request, and clears itself once sent.
	send := form.New(
		row.New(field.Text("message", "A line to send"), button.Submit("Send")),
	)
	htmx.New(send).WSSend().HxOn(event.WSAfterSend, "this.reset()")

	wsCard := card.New("WebSocket Live Log",
		div.New(
			span.Static("Status: ").Class("label"),
			span.Text("connecting...").Class("status-connecting").ID("ws-status"),
		).Class("layout-row"),
		div.New(
			span.Text(
				"Connected to the server via a WebSocket. Log entries are pushed "+
					"from the server as HTML fragments and swapped in by the htmx "+
					"WebSocket extension. The connection is full-duplex, so a line "+
					"sent with the form travels over the same connection and comes "+
					"back as an entry.",
			).Class("hint"),
		),
		send,
		div.New().Class("log-feed").ID("ws-log"),
	)

	// The ws-connect div wraps the card so htmx opens the WebSocket,
	// sends the form over it and processes incoming HTML fragments as
	// OOB swaps. The extension fires its lifecycle events on this
	// element, so HxOn binds the status updates here.
	content := div.New(wsCard)
	htmx.New(content).
		HxExt("ws").
		WSConnect("/ws/feed").
		HxOn(event.WSOpen, setStatus("ws-status", "connected", "status-connected")).
		HxOn(event.WSClose, setStatus("ws-status", "reconnecting...", "status-reconnecting")).
		HxOn(event.WSError, setStatus("ws-status", "disconnected", "status-disconnected"))

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

	// ReadLoop processes control frames (ping/pong/close) and the
	// lines the page sends, so the connection stays alive and closes
	// cleanly.
	go conn.ReadLoop()

	slog.Info("ws: client connected", "remote", r.RemoteAddr)

	// Send a status update so the page shows "connected" instead of
	// "connecting...". The span is swapped out of band, so htmx
	// replaces the existing status element by ID.
	status := span.Text("connected").Class("status-connected").ID("ws-status")
	htmx.New(status).HxSwapOOB("true")
	if err := conn.WriteMessage(gws.OpcodeText, status.RenderBytes()); err != nil {
		slog.Info("ws: client disconnected on status", "remote", r.RemoteAddr)
		return
	}

	// Push log entries until the connection drops.
	for {
		if err := conn.WriteMessage(gws.OpcodeText, wsEntry(generate.LogEntry())); err != nil {
			slog.Info("ws: client disconnected", "remote", r.RemoteAddr)
			return
		}
		time.Sleep(generate.Jitter())
	}
}

// wsEntry renders one log entry as a WebSocket message: a div with
// the log's ID, swapped out of band with afterbegin so htmx prepends
// the entry (newest at top).
func wsEntry(e generate.Entry) []byte {
	wrapper := div.New(logentry.New(e)).ID("ws-log")
	htmx.New(wrapper).HxSwapOOB("afterbegin")
	return wrapper.RenderBytes()
}

// wsHandler satisfies the gws.Event interface. It answers the lines
// the page's form sends; the rest of the connection lifecycle is the
// built-in one.
type wsHandler struct {
	gws.BuiltinEventHandler
}

// OnMessage answers a line sent from the page's form. htmx-ext-ws
// sends the form's fields as a JSON object, with the htmx request
// headers under HEADERS. The reply is a log entry carrying the line,
// written to the same connection as the feed; gws locks each write,
// so the two never interleave.
func (h *wsHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()

	var sent struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(message.Bytes(), &sent); err != nil {
		slog.Warn("ws: unreadable message", "error", err)
		return
	}
	line := strings.TrimSpace(sent.Message)
	if line == "" {
		return
	}
	if err := socket.WriteMessage(gws.OpcodeText, wsEntry(generate.Sent(line))); err != nil {
		slog.Info("ws: client disconnected on reply", "error", err)
	}
}
