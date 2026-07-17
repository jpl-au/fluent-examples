package playwright_test

import (
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
)

// The computed-signals suite drives client-side derived signals through a
// real browser. The server only ever pushes inputs (qty, price); every
// derived value - total, its double, the show flag - is computed on the
// client by the postfix VM, with no eval. The suite also proves the
// conditional bindings still work now that they share that VM, and that a
// dependency cycle is reported rather than crashing the runtime.

type computeState struct {
	Qty int
}

func renderCompute(s computeState) node.Node {
	return div.New(
		// total = qty * price, both server-pushed inputs.
		bind.Apply(span.New(),
			bind.Computed("total", "qty * price"),
			bind.Text("total"),
		).ID("total"),

		// Chained computed: double reads total, another computed's output.
		bind.Apply(span.New(),
			bind.Computed("double", "total + total"),
			bind.Text("double"),
		).ID("double"),

		// A computed drives Show: the flag is derived, never pushed.
		bind.Apply(span.Text("BIG"),
			bind.Computed("big", "total > 50"),
			bind.Show("big"),
		).ID("big"),

		// ShowWhen still works after the collapse onto the shared VM.
		bind.Apply(span.Text("MANY"), bind.ShowWhen("qty", ">=", 3)).ID("many"),

		// A self-referential computed: seeding it forces a cycle, which
		// the runtime must report (not loop or throw).
		bind.Apply(span.New(),
			bind.Computed("loop", "loop + 1"),
			bind.Text("loop"),
		).ID("loop"),

		bind.Apply(button.Text("inc"), bind.OnClick("inc")).ID("inc"),
		bind.Apply(button.Text("cycle"), bind.OnClick("cycle")).ID("cycle"),
	)
}

func handleCompute(sess tether.Session, s computeState, ev tether.Event) computeState {
	switch ev.Action {
	case "inc":
		s.Qty++
		sess.Signal("qty", s.Qty)
		sess.Signal("price", 10)
	case "cycle":
		// Seed the self-referential computed to trigger the cycle guard.
		sess.Signal("loop", 1)
	}
	return s
}

func startComputeServer(t *testing.T) string {
	t.Helper()
	h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[computeState]{
		InitialState: func(*http.Request) computeState { return computeState{} },
		Render:       renderCompute,
		Handle:       handleCompute,
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

// TestComputedCartTotal clicks inc, which pushes qty and price, and
// asserts the client derives total = qty * price live - no total signal
// is ever sent by the server.
func TestComputedCartTotal(t *testing.T) {
	srv := startComputeServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	inc := page.Locator("#inc")
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	// qty=1, price=10 -> total=10.
	if err := expect(page.Locator("#total")).ToHaveText("10"); err != nil {
		t.Errorf("total should be 10 at qty=1: %v", err)
	}

	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	// qty=2 -> total=20.
	if err := expect(page.Locator("#total")).ToHaveText("20"); err != nil {
		t.Errorf("total should be 20 at qty=2: %v", err)
	}
}

// TestComputedChained proves a computed that reads another computed's
// output updates when the underlying input changes: double = total + total.
func TestComputedChained(t *testing.T) {
	srv := startComputeServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#inc").Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	// total=10 -> double=20.
	if err := expect(page.Locator("#double")).ToHaveText("20"); err != nil {
		t.Errorf("double should chain to 20: %v", err)
	}
}

// TestComputedDrivesShow raises the total past the threshold and asserts a
// computed boolean (big = total > 50) toggles a Show binding.
func TestComputedDrivesShow(t *testing.T) {
	srv := startComputeServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	inc := page.Locator("#inc")
	// One click: total=10, not > 50, BIG hidden.
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	if err := expect(page.Locator("#big")).ToBeHidden(); err != nil {
		t.Errorf("BIG should be hidden at total=10: %v", err)
	}
	// Six clicks total: qty=6, total=60 > 50, BIG visible.
	for i := 0; i < 5; i++ {
		if err := inc.Click(); err != nil {
			t.Fatalf("click inc: %v", err)
		}
	}
	if err := expect(page.Locator("#big")).ToBeVisible(); err != nil {
		t.Errorf("BIG should be visible at total=60: %v", err)
	}
}

// TestShowWhenPostCollapse checks the conditional-binding sugar still
// works after compiling to the shared VM: MANY appears once qty >= 3.
func TestShowWhenPostCollapse(t *testing.T) {
	srv := startComputeServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	inc := page.Locator("#inc")
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	if err := expect(page.Locator("#many")).ToBeHidden(); err != nil {
		t.Errorf("MANY should be hidden at qty=1: %v", err)
	}
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	if err := expect(page.Locator("#many")).ToBeVisible(); err != nil {
		t.Errorf("MANY should be visible at qty=3: %v", err)
	}
}

// TestComputedCycleReported seeds a self-referential computed and asserts
// the cycle is reported via Tether.onError without crashing the runtime -
// a subsequent inc still drives total.
func TestComputedCycleReported(t *testing.T) {
	srv := startComputeServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	// Collect client-side error reports.
	if _, err := page.Evaluate(`() => { window.__errs = []; window.Tether.onError = function (e) { window.__errs.push(e && e.message ? e.message : String(e)); }; }`); err != nil {
		t.Fatalf("install onError: %v", err)
	}

	if err := page.Locator("#cycle").Click(); err != nil {
		t.Fatalf("click cycle: %v", err)
	}

	// The click round-trips to the server, which seeds the signal; only then
	// does the cascade run and the guard fire - so wait for the report.
	if _, err := page.WaitForFunction(`() => (window.__errs || []).length > 0`, nil); err != nil {
		t.Fatalf("no cycle report arrived: %v", err)
	}

	// The cycle guard must have reported at least one error mentioning it.
	got, err := page.Evaluate(`() => (window.__errs || []).join("\n")`)
	if err != nil {
		t.Fatalf("read errs: %v", err)
	}
	if msg, _ := got.(string); !strings.Contains(strings.ToLower(msg), "cycle") {
		t.Errorf("expected a cycle report, got: %q", msg)
	}

	// The runtime survived: a normal computed still updates.
	if err := page.Locator("#inc").Click(); err != nil {
		t.Fatalf("click inc after cycle: %v", err)
	}
	if err := expect(page.Locator("#total")).ToHaveText("10"); err != nil {
		t.Errorf("total should still compute after a cycle: %v", err)
	}
}
