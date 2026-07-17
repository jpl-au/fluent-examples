package playwright_test

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
)

// The indicator suite proves the loading indicator is ref-counted: two
// overlapping in-flight events point their bind.Indicator at the same
// spinner, and the spinner must stay lit until the LAST of them settles.
// Without ref-counting, the first event to finish would clear the
// spinner while the second is still in flight.
//
// Neither action changes state, so the echo carries no DOM morph - the
// spinner's tether-pending class is governed solely by the ref-count.

type indicatorState struct{}

func renderIndicator(s indicatorState) node.Node {
	return div.New(
		// The shared loading indicator both buttons target.
		div.New(span.Text("spinner")).ID("spinner"),

		bind.Apply(button.Text("fast"),
			bind.OnClick("fast"),
			bind.Indicator("#spinner"),
			bind.Disable("..."),
		).ID("fast"),

		bind.Apply(button.Text("slow"),
			bind.OnClick("slow"),
			bind.Indicator("#spinner"),
			bind.Disable("..."),
		).ID("slow"),
	)
}

func handleIndicator(_ tether.Session, s indicatorState, ev tether.Event) indicatorState {
	// Serialised per session: the fast action settles well before the
	// slow one, opening a window where only the slow event is in flight.
	switch ev.Action {
	case "fast":
		time.Sleep(300 * time.Millisecond)
	case "slow":
		time.Sleep(2500 * time.Millisecond)
	}
	return s
}

func startIndicatorServer(t *testing.T) string {
	t.Helper()
	h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[indicatorState]{
		InitialState: func(*http.Request) indicatorState { return indicatorState{} },
		Render:       renderIndicator,
		Handle:       handleIndicator,
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

// TestIndicatorRefCounted fires two overlapping events that share one
// indicator and asserts the spinner stays lit until the last one
// settles.
func TestIndicatorRefCounted(t *testing.T) {
	srv := startIndicatorServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	lit := page.Locator("#spinner.tether-pending")

	// Fire both events; the fast one settles first, the slow one keeps
	// running. Clicking dispatches immediately (well before either echo),
	// so both are in flight together.
	if err := page.Locator("#fast").Click(); err != nil {
		t.Fatalf("click fast: %v", err)
	}
	if err := page.Locator("#slow").Click(); err != nil {
		t.Fatalf("click slow: %v", err)
	}

	// Both in flight: the spinner is lit.
	if err := expect(lit).ToHaveCount(1); err != nil {
		t.Fatalf("spinner should be lit while events are in flight: %v", err)
	}

	// Wait for the fast event to settle (its button re-enables).
	if err := expect(page.Locator("#fast")).Not().ToBeDisabled(); err != nil {
		t.Fatalf("fast event did not settle: %v", err)
	}

	// The slow event is still in flight, so ref-counting must keep the
	// spinner lit. Without ref-counting it would have cleared here.
	if err := expect(lit).ToHaveCount(1); err != nil {
		t.Errorf("spinner cleared while the slow event was still in flight (not ref-counted): %v", err)
	}

	// Once the slow event settles, the spinner clears.
	if err := expect(page.Locator("#slow")).Not().ToBeDisabled(); err != nil {
		t.Fatalf("slow event did not settle: %v", err)
	}
	if err := expect(lit).ToHaveCount(0); err != nil {
		t.Errorf("spinner should clear after the last event settles: %v", err)
	}
}
