// harness_test.go demonstrates the core tethertest harness: creating a
// harness, sending events, reading state, rendering HTML, and verifying
// that the rendered output reflects state changes.
package tethertest_test

import (
	"strings"
	"testing"

	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/tethertest"

	"github.com/jpl-au/fluent-examples/tether/site/signals"
)

// TestSendAndState demonstrates the simplest use case: send an event,
// read the resulting state. State accumulates across multiple sends.
func TestSendAndState(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Handle: signals.Handle,
	})

	h.Send("signals.increment")
	if h.State().Counter != 1 {
		t.Fatalf("Counter = %d, want 1", h.State().Counter)
	}

	h.Send("signals.increment")
	if h.State().Counter != 2 {
		t.Fatalf("Counter = %d, want 2", h.State().Counter)
	}
}

// TestRenderHTML demonstrates HTML() - the rendered output after each
// Send. Requires Config.Render to be set.
func TestRenderHTML(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Render: signals.RenderWS,
		Handle: signals.Handle,
	})

	h.Send("signals.increment")
	html := h.HTML()
	if !strings.Contains(html, "Increment Server Counter") {
		t.Fatal("HTML should contain the increment button")
	}
}

// TestRender demonstrates Render() - the full rendered HTML for the
// current state without sending any event first.
func TestRender(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Render: signals.RenderWS,
		Handle: signals.Handle,
	})

	html := h.Render()
	if !strings.Contains(html, "BindText") {
		t.Fatal("Render should contain the BindText demo section")
	}
}

// TestRenderNode demonstrates RenderNode() - returns the raw node
// tree for programmatic inspection rather than a string.
func TestRenderNode(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Render: signals.RenderWS,
		Handle: signals.Handle,
	})

	n := h.RenderNode()
	if n == nil {
		t.Fatal("RenderNode should not return nil")
	}
	// Render the node to HTML to verify it produces output.
	html := string(n.Render())
	if html == "" {
		t.Fatal("node.Render() should produce non-empty HTML")
	}
}

// TestRenderOptional demonstrates that Config.Render is optional  -
// tests that only check state and effects don't need a render function.
func TestRenderOptional(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Handle: signals.Handle,
	})

	h.Send("signals.increment")
	if h.State().Counter != 1 {
		t.Fatal("state should work without Render")
	}
	if h.HTML() != "" {
		t.Fatal("HTML should be empty when Render is not configured")
	}
}

// TestEffectsResetBetweenSends demonstrates that each Send produces
// a fresh set of effects - a toast from one Send does not carry over
// to the next.
func TestEffectsResetBetweenSends(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Handle: func(sess tether.Session, s signals.State, ev tether.Event) signals.State {
			if ev.Action == "toast-action" {
				sess.Toast("hello")
			}
			return signals.Handle(sess, s, ev)
		},
	})

	h.Send("toast-action")
	if h.Toast() == "" {
		t.Fatal("expected toast after first send")
	}

	h.Send("signals.increment")
	if h.Toast() != "" {
		t.Fatalf("Toast = %q, should be empty after non-toast send", h.Toast())
	}
}
