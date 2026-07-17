package playwright_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
)

// The input-preservation suite proves ignoreActiveValue: a server patch
// that would rewrite the value of the input the user is actively typing
// in must not clobber the in-progress text. The input carries a
// server-controlled value attribute that changes on every render, so
// without ignoreActiveValue the morph would overwrite what the user
// typed.

type preserveState struct {
	Rewrites int
}

func renderPreserve(s preserveState) node.Node {
	// The input's value attribute is server-controlled. It starts empty
	// so the user types into a clean field; once the server has seen an
	// input it renders value="server-N", which a morph would push into
	// the field - ignoreActiveValue is what keeps the user's text.
	serverVal := ""
	if s.Rewrites > 0 {
		serverVal = fmt.Sprintf("server-%d", s.Rewrites)
	}
	field := bind.Apply(
		input.Text("field", serverVal),
		bind.OnInput("rewrite"),
		bind.Debounce(50*time.Millisecond),
	).ID("field")

	return div.New(
		span.Text(fmt.Sprintf("rewrites:%d", s.Rewrites)).ID("count"),
		field,
	).Dynamic("box")
}

func handlePreserve(_ tether.Session, s preserveState, ev tether.Event) preserveState {
	if ev.Action == "rewrite" {
		s.Rewrites++
	}
	return s
}

func startPreserveServer(t *testing.T) string {
	t.Helper()
	h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[preserveState]{
		InitialState: func(*http.Request) preserveState { return preserveState{} },
		Render:       renderPreserve,
		Handle:       handlePreserve,
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := &http.Server{Handler: h}
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close() })
	return "http://" + ln.Addr().String()
}

// TestInputValueSurvivesMorph types into an input, triggers a server
// render whose morph would rewrite the input's value attribute, and
// asserts the user's in-progress text survives (ignoreActiveValue).
func TestInputValueSurvivesMorph(t *testing.T) {
	srv := startPreserveServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	field := page.Locator("#field")
	if err := field.Focus(); err != nil {
		t.Fatalf("focus field: %v", err)
	}
	// Type while the field is focused. Each keystroke fires a debounced
	// input event; the server increments Rewrites and morphs the box,
	// which would reset the value attribute to "server-1".
	if err := field.PressSequentially("hello world"); err != nil {
		t.Fatalf("type: %v", err)
	}

	// Wait for the server render to land so a morph has definitely run.
	// One or more debounce windows may fire; any non-zero count proves a
	// morph carrying a fresh value attribute reached the input.
	if err := expect(page.Locator("#count")).Not().ToHaveText("rewrites:0"); err != nil {
		text, _ := page.Locator("#count").TextContent()
		t.Fatalf("server did not render after input: %q (%v)", text, err)
	}

	// The morph carried value="server-N", but the field was the active
	// element, so the typed text must be intact.
	if err := expect(field).ToHaveValue("hello world"); err != nil {
		val, _ := field.InputValue()
		t.Errorf("morph clobbered the active input: value = %q, want %q", val, "hello world")
	}
}
