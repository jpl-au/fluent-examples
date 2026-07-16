// Package handler holds the request handlers for the demo. The same
// page is served twice: / is wired to backend endpoints written with
// the fluent-datastar SSE generator, and /sdk to endpoints written with
// the official Datastar Go SDK (see sdk.go). Both server halves speak
// the same wire protocol, so one fds-decorated client page drives
// either; the only difference between the two pages is which endpoints
// their buttons call.
package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h1"
	"github.com/jpl-au/fluent/html5/h2"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	fds "github.com/jpl-au/fluent-datastar"

	"github.com/jpl-au/fluent-examples/fluent-datastar/layout"
)

// countdownFrom is how many seconds the streaming endpoints count down,
// emitting one element patch and one signal patch per second.
const countdownFrom = 5

// Home renders the demo page wired to the fluent-datastar server half.
func Home(w http.ResponseWriter, _ *http.Request) {
	render(w, "Fluent generator",
		"This page drives endpoints built with the fluent-datastar SSE generator.",
		"/increment", "/clock",
		"/sdk", "the same app driven by the official Datastar Go SDK")
}

// render builds the shared demo page. The two pages differ only in the
// endpoint paths their buttons call and the cross link to the other
// implementation.
func render(w http.ResponseWriter, half string, lede string, incrementPath string, clockPath string, otherHref string, otherLabel string) {
	page := div.New(
		h1.Text("Fluent-Datastar"),
		p.Text("Both halves of the bindings in one small app: reactive signals on the client, and a server patching elements and signals back over SSE.").Class("lede"),
		p.New(
			span.Text(lede+" See "),
			a.Text(otherLabel).Href(otherHref),
			span.Text("."),
		).Class("lede"),
		clientCounterCard(),
		roundTripCard(incrementPath),
		streamingCard(clockPath),
	).Class("wrap")

	// data-signals on the root declares the whole signal store up front.
	// count is client-only; the others are written by the server.
	fds.New(page).Signals(`{count: 0, serverCount: 0, saving: false, tick: 0}`)

	layout.Page(w, half, page)
}

// clientCounterCard is the client-only reactive section. No request is
// made: On("click", ...) mutates the count signal directly, Text binds
// the display to it, and Show toggles a note purely from signal state.
func clientCounterCard() node.Node {
	display := span.Text("0").ID("count-display").Class("counter")
	fds.New(display).Text("$count")

	dec := button.Text("-").Class("btn")
	fds.New(dec).On("click", "$count--")

	inc := button.Text("+").Class("btn")
	fds.New(inc).On("click", "$count++")

	// Show keeps this note in the DOM but toggles its visibility from
	// the signal expression, with no round trip.
	note := p.Text("Count is above zero.").Class("note")
	fds.New(note).Show("$count > 0")

	return card("1. Client-only signals",
		p.Text("count lives entirely in the browser. The buttons write it, the number reads it, and the note below shows only when it is positive."),
		div.New(dec, display, inc).Class("row"),
		note,
	)
}

// roundTripCard is the server round trip. The button issues a POST via
// the Datastar @post action; Datastar sends the whole signal store as
// the request body. Indicator flips the saving signal true while the
// request is in flight, which the spinner reads through Show.
func roundTripCard(incrementPath string) node.Node {
	sync := button.Text("Sync count to server").Class("btn btn-primary")
	fds.New(sync).On("click", fds.Post(incrementPath)).Indicator("saving")

	spinner := span.Text("saving...").Class("spinner")
	fds.New(spinner).Show("$saving")

	serverCount := span.Text("0").Class("counter")
	fds.New(serverCount).Text("$serverCount")

	// The server replaces this element by its id (outer mode).
	result := div.New(
		span.Text("Not synced yet.").Class("muted"),
	).ID("server-result").Class("result")

	return card("2. Server round trip",
		p.Text("The button POSTs the current signals. The handler reads them, then streams back one element patch (the panel below) and one signal patch (the mirrored count) in a single response."),
		div.New(sync, spinner).Class("row"),
		div.New(
			span.Text("Server-side count: ").Class("muted"),
			serverCount,
		).Class("row"),
		result,
	)
}

// streamingCard starts the SSE stream. The GET action opens one
// response that the server holds open, emitting several events over a
// few seconds. tick is written by each signal patch on the stream.
func streamingCard(clockPath string) node.Node {
	start := button.Text("Start countdown").Class("btn btn-primary")
	fds.New(start).On("click", fds.Get(clockPath))

	clock := div.New(
		span.Text("idle").Class("muted"),
	).ID("clock").Class("clock")

	tick := span.Text("0").Class("counter")
	fds.New(tick).Text("$tick")

	return card("3. Streaming SSE",
		p.Text("One GET opens a stream the server keeps open, sending an element patch and a signal patch every second until liftoff. Several events, one response."),
		div.New(start).Class("row"),
		clock,
		div.New(
			span.Text("Last tick signal: ").Class("muted"),
			tick,
		).Class("row"),
	)
}

// Increment is the fluent-half round-trip endpoint. It reads the posted
// signals, then opens an SSE stream and patches back both a rendered
// element and the mirrored signal in one response.
func Increment(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Count int `json:"count"`
	}
	if err := fds.ReadSignals(r, &in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("increment: signals received", "source", "fluent", "count", in.Count)

	gen, err := fds.NewGenerator(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Patch the #server-result element (outer mode, matched by id) with
	// a freshly rendered fluent node.
	if err := gen.PatchElements(resultPanel(in.Count)); err != nil {
		slog.Error("patch elements", "error", err)
		return
	}

	// Mirror the value back into a separate signal so the client shows
	// what the server saw.
	if err := gen.MarshalSignals(map[string]int{"serverCount": in.Count}); err != nil {
		slog.Error("patch signals", "error", err)
	}
}

// Clock is the fluent-half streaming endpoint. It counts down from
// countdownFrom, emitting one element patch and one signal patch per
// second, then a final liftoff patch. It stops early if the client
// disconnects.
func Clock(w http.ResponseWriter, r *http.Request) {
	gen, err := fds.NewGenerator(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("clock: countdown started", "source", "fluent", "remote", r.RemoteAddr)

	ctx := r.Context()
	for n := countdownFrom; n >= 0; n-- {
		if err := gen.PatchElements(clockPanel(n)); err != nil {
			return
		}
		if err := gen.MarshalSignals(map[string]int{"tick": n}); err != nil {
			return
		}

		if n == 0 {
			slog.Info("clock: liftoff", "source", "fluent")
			return
		}
		select {
		case <-ctx.Done():
			slog.Info("clock: client disconnected", "source", "fluent", "remote", r.RemoteAddr)
			return
		case <-time.After(time.Second):
		}
	}
}

// resultPanel renders the round-trip result element both halves patch
// back by its id.
func resultPanel(count int) node.Node {
	stamp := time.Now().Format("15:04:05")
	return div.New(
		span.Textf("Received count %d at %s.", count, stamp).Class("ok"),
	).ID("server-result").Class("result")
}

// clockPanel renders one countdown frame both halves patch back by its id.
func clockPanel(n int) node.Node {
	label := fmt.Sprintf("T-minus %d", n)
	if n == 0 {
		label = "Liftoff!"
	}
	return div.New(
		span.Text(label).Class("ok"),
	).ID("clock").Class("clock")
}

// card wraps a titled section in consistent markup.
func card(heading string, content ...node.Node) node.Node {
	children := append([]node.Node{h2.Text(heading)}, content...)
	return div.New(children...).Class("card")
}
