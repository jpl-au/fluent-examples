package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/lxzan/gws"

	htmx "github.com/jpl-au/fluent-htmx/htmx4"
	"github.com/jpl-au/fluent-htmx/htmx4/event"
	"github.com/jpl-au/fluent-htmx/htmx4/swap"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/span"

	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/composite/logentry"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/composite/row"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/simple/button"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/components/simple/field"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/generate"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx4/layout"
)

// wsUpgrader is shared across requests. The default configuration
// accepts all origins, which is fine for a local demo app.
var wsUpgrader = gws.NewUpgrader(&wsHandler{}, &gws.ServerOption{})

// WSPage renders the WebSocket live log page. A wrapper around the
// form and the log opens the connection with hx-ws:connect, which the
// htmax.js bundle provides; the form has to sit inside it to send over
// that connection. A plain HTML message is swapped into the log with
// afterbegin; a JSON message names its own target. The close and
// error events fire on the connecting element, so HxOn binds the
// status updates there. The connected state comes from the server's
// first message instead.
func WSPage(w http.ResponseWriter, _ *http.Request) {
	// The form sends its fields over the open connection as JSON
	// instead of making a request, and clears itself once sent.
	send := form.New(
		row.New(field.Text("message", "A line to send"), button.Submit("Send")),
	)
	htmx.New(send).WSSend().HxOn(event.WSAfterMessageOutgoing, "this.reset()")

	logDiv := div.New().Class("log-feed").ID("ws-log")

	connection := div.New(send, logDiv)
	htmx.New(connection).
		WSConnect("/ws/feed").
		HxTarget("#ws-log").
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
					"WebSocket extension. The connection is full-duplex, so a line "+
					"sent with the form travels over the same connection and comes "+
					"back as an entry.",
			).Class("hint"),
		),
		connection,
	)

	layout.Page(w, "WebSocket", nil, content)
}

// WSFeed upgrades the connection to a WebSocket and pushes rendered
// HTML log entries at random intervals. The first message is a JSON
// frame that replaces the status span, and the entries follow as
// plain HTML that htmx swaps into the log.
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

// wsHandler satisfies the gws.Event interface. It answers the lines
// the page's form sends; the rest of the connection lifecycle is the
// built-in one.
type wsHandler struct {
	gws.BuiltinEventHandler
}

// OnMessage answers a line sent from the page's form. ParseWSMessage
// splits the hx-ws:send frame into the form's values and the htmx
// request headers. The reply is a log entry carrying the line, as
// plain HTML on the same connection as the feed; gws locks each
// write, so the two never interleave.
func (h *wsHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()

	m, err := htmx.ParseWSMessage(message.Bytes())
	if err != nil {
		slog.Warn("ws: unreadable message", "error", err)
		return
	}
	sent, _ := m.Values["message"].(string)
	line := strings.TrimSpace(sent)
	if line == "" {
		return
	}
	if err := socket.WriteMessage(gws.OpcodeText, logentry.New(generate.Sent(line)).RenderBytes()); err != nil {
		slog.Info("ws: client disconnected on reply", "error", err)
	}
}
