package playwright_test

import (
	"net/http"
	"testing"

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

func startIndicatorServer(t *testing.T) (base string, finishFast, finishSlow func()) {
	t.Helper()
	fast, slow := make(chan struct{}), make(chan struct{})
	h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[indicatorState]{
		InitialState: func(*http.Request) indicatorState { return indicatorState{} },
		Render:       renderIndicator,
		Handle: func(sess tether.Session, s indicatorState, ev tether.Event) indicatorState {
			var release <-chan struct{}
			switch ev.Action {
			case "fast":
				release = fast
			case "slow":
				release = slow
			default:
				return s
			}
			// This fixture deliberately holds replies until the browser has
			// observed each in-flight state. Shutdown also releases a handler
			// if an assertion fails before the test opens its gate.
			select {
			case <-release:
			case <-sess.Context().Done():
			}
			return s
		},
	})
	return startServer(t, RealHTTP1, h, h.Shutdown), func() { close(fast) }, func() { close(slow) }
}

// TestIndicatorRefCounted fires two overlapping events that share one
// indicator and asserts the spinner stays lit until the last one
// settles.
func TestIndicatorRefCounted(t *testing.T) {
	srv, finishFast, finishSlow := startIndicatorServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	lit := page.Locator("#spinner.tether-pending")

	// Hold both replies until the browser confirms both events are in flight.
	if err := page.Locator("#fast").Click(); err != nil {
		t.Fatalf("click fast: %v", err)
	}
	if err := page.Locator("#slow").Click(); err != nil {
		t.Fatalf("click slow: %v", err)
	}
	for _, id := range []string{"#fast", "#slow"} {
		if err := expect(page.Locator(id)).ToBeDisabled(); err != nil {
			t.Fatalf("%s event is not in flight: %v", id, err)
		}
	}

	// Both in flight: the spinner is lit.
	if err := expect(lit).ToHaveCount(1); err != nil {
		t.Fatalf("spinner should be lit while events are in flight: %v", err)
	}

	// Wait for the fast event to settle (its button re-enables).
	finishFast()
	if err := expect(page.Locator("#fast")).Not().ToBeDisabled(); err != nil {
		t.Fatalf("fast event did not settle: %v", err)
	}

	// The slow event is still in flight, so ref-counting must keep the
	// spinner lit. Without ref-counting it would have cleared here.
	if err := expect(lit).ToHaveCount(1); err != nil {
		t.Errorf("spinner cleared while the slow event was still in flight (not ref-counted): %v", err)
	}

	// Once the slow event settles, the spinner clears.
	finishSlow()
	if err := expect(page.Locator("#slow")).Not().ToBeDisabled(); err != nil {
		t.Fatalf("slow event did not settle: %v", err)
	}
	if err := expect(lit).ToHaveCount(0); err != nil {
		t.Errorf("spinner should clear after the last event settles: %v", err)
	}
}
