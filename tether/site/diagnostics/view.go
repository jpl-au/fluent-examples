package diagnostics

import (
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/diagnostic"
	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the diagnostics page: trigger buttons, a live event
// feed, and a reference of all diagnostic kinds.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Trigger Diagnostic Events",
			"Click a button to deliberately cause a framework-level failure. "+
				"The framework recovers the error, emits a Diagnostic event on the handler's "+
				"Diagnostics bus, and the event appears in the feed below in real time. "+
				"Your session stays alive - handler panics are caught, not fatal.",
			"Handler.Diagnostics · HandlerPanic", panel.WS|panel.SSE,
			layout.Stack(
				button.DangerAction("Trigger Handler Panic", "diag.trigger-panic"),
				hint.Text(
					"You can also trigger diagnostics from other pages: "+
						"upload a wrong file type on Filtered Uploads (UploadRejected), "+
						"or disconnect your network briefly on any live page (TransportError).",
				),
			),
		),

		panel.Card(
			"Diagnostic Event Feed",
			"This feed subscribes to Handler.Diagnostics via WatchBus and renders the 20 most recent events. "+
				"Events are aggregated from multiple handlers across the application.",
			"WatchBus · DiagnosticKind · Bus.Publish", panel.WS|panel.SSE,
			eventFeed(s.Events),
		),

		panel.Card(
			"DiagnosticKind Reference",
			"Every diagnostic event carries a Kind that identifies its category. "+
				"Subscribe to Handler.Diagnostics and switch on Kind to route events to "+
				"metrics, alerting, or logging. All diagnostics are non-fatal.",
			"DiagnosticKind · Diagnostic", panel.AllTransports,
			kindReference(),
		),
	)
}

// eventFeed renders the diagnostic event list or a placeholder.
func eventFeed(entries []Entry) node.Node {
	if len(entries) == 0 {
		return hint.Text(
			"No diagnostic events yet. Click the Trigger Handler Panic button above to generate one.",
		)
	}
	items := make([]node.Node, len(entries))
	for i, e := range entries {
		t := e.Kind
		if e.Detail != "" {
			t += " - " + e.Detail
		}
		if e.SessionID != "" && len(e.SessionID) >= 6 {
			t += " (session " + e.SessionID[:6] + ")"
		}
		items[i] = diagnostic.Event(t)
	}
	return diagnostic.EventList(items...).Dynamic("diagnostics")
}

// kindEntry describes a single DiagnosticKind for the reference.
type kindEntry struct {
	kind    string
	desc    string
	trigger string // plain-text trigger instruction (fallback)
	link    string // optional URL to open in a new tab
	action  string // optional Tether action to fire on click
	linkTxt string // button/link label
}

// kinds lists all DiagnosticKind values with descriptions and how
// to trigger them in this demo application.
var kinds = []kindEntry{
	{
		kind:    "handler_panic",
		desc:    "Recovered panic inside Handle, Update, or a command callback. The session stays alive - the framework catches the panic and continues.",
		action:  "diag.trigger-panic",
		linkTxt: "Trigger Handler Panic",
	},
	{
		kind:    "upload_rejected",
		desc:    "Upload rejected because its MIME type did not match UploadConfig.Accept. The Detail field contains the rejected content type.",
		trigger: "Upload a non-image file (e.g. a .txt) on",
		link:    "/uploads/filtered/",
		linkTxt: "Filtered Uploads",
	},
	{
		kind:    "transport_error",
		desc:    "Failure reading from or writing to the WebSocket or SSE transport. Normal disconnects (EOF) are not emitted - only genuine failures.",
		trigger: "Disconnect your network briefly while on any live page.",
	},
	{
		kind:    "buffer_overflow",
		desc:    "Session command channel was full - a goroutine was spawned to deliver the command. Sustained overflow indicates a slow handler or broadcast storm.",
		trigger: "Requires high-frequency broadcasts exceeding CmdBufferSize.",
	},
	{
		kind:    "encode_error",
		desc:    "Failure encoding a wire update (JSON serialisation). Usually indicates an unencodable type in state or render output.",
		trigger: "Not easily triggered from the UI.",
	},
	{
		kind:    "upload_error",
		desc:    "Failure inside an UploadConfig.Handle callback (e.g. disk full, permission denied).",
		trigger: "Not easily triggered from the UI.",
	},
	{
		kind:    "command_dropped",
		desc:    "Command discarded because both the buffer and overflow goroutine cap were exhausted. Data was lost - unlike buffer_overflow, this is unrecoverable.",
		trigger: "Requires extreme overload beyond buffer_overflow.",
	},
	{
		kind:    "session_binding_failed",
		desc:    "Client attempted to reconnect with a different User-Agent than the original session. May indicate a stolen session ID.",
		trigger: "Requires a spoofed User-Agent on reconnect.",
	},
	{
		kind:    "store_error",
		desc:    "Failure saving, loading, or deleting differ snapshots from the DiffStore. Non-fatal - the framework falls back to in-memory behaviour.",
		trigger: "Requires a filesystem failure on the diff store directory.",
	},
	{
		kind:    "session_store_error",
		desc:    "Failure saving, loading, or deleting session state from the SessionStore. Non-fatal - continues with in-memory state.",
		trigger: "Requires a filesystem failure on the session store directory.",
	},
}

// kindReference renders the full DiagnosticKind reference as a list
// of cards, each showing the kind name, description, and how to
// trigger it in the demo application.
func kindReference() node.Node {
	items := make([]node.Node, len(kinds))
	for i, k := range kinds {
		items[i] = diagnostic.Item(k.kind, k.desc, triggerNode(k))
	}
	return diagnostic.List(items...)
}

// triggerNode renders a button for the diagnostic kind: either a
// Tether action button, a link to another demo page (new tab), or
// plain text when neither is available.
func triggerNode(k kindEntry) node.Node {
	if k.action != "" {
		return button.SmallAction(k.linkTxt, k.action)
	}
	if k.link != "" {
		return button.Link(k.linkTxt, k.link)
	}
	return diagnostic.Trigger(k.trigger)
}
