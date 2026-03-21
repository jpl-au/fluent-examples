// patterns_test.go demonstrates advanced harness patterns: SendInput,
// SendSubmit, SendEvent for different event types, StatefulConfig.Middleware
// for chaining, typed event data parsing (ev.Int, ev.Float64, ev.Bool,
// ev.Bind), Bus.Emit for cross-session messaging, and component
// testing with NewComponent.
package tethertest_test

import (
	"testing"

	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/event"
	"github.com/jpl-au/tether/tethertest"

	"github.com/jpl-au/fluent-examples/tether/component/counter"
	"github.com/jpl-au/fluent-examples/tether/site/broadcasting"
	"github.com/jpl-au/fluent-examples/tether/site/events"
	"github.com/jpl-au/fluent-examples/tether/site/groups"
	"github.com/jpl-au/fluent-examples/tether/site/mw"
	"github.com/jpl-au/fluent-examples/tether/site/valuestore"
)

// ---------- SendInput / SendSubmit / SendEvent ----------

// TestSendInput demonstrates SendInput - fires an input event with
// a value that the handler reads via ev.Value().
func TestSendInput(t *testing.T) {
	h := tethertest.New(tethertest.Config[events.State]{
		Handle: events.Handle,
	})

	h.SendInput("events.input", "hello world")
	if h.State().InputValue != "hello world" {
		t.Fatalf("InputValue = %q", h.State().InputValue)
	}
}

// TestSendSubmit demonstrates SendSubmit - fires a submit event with
// named form fields.
func TestSendSubmit(t *testing.T) {
	h := tethertest.New(tethertest.Config[events.State]{
		Handle: events.Handle,
	})

	h.SendSubmit("events.submit", map[string]string{"name": "Alice"})
	if h.State().SubmitResult != "Hello, Alice!" {
		t.Fatalf("SubmitResult = %q", h.State().SubmitResult)
	}
}

// TestSendEvent demonstrates SendEvent for custom event types - here
// a Change event from a dropdown.
func TestSendEvent(t *testing.T) {
	h := tethertest.New(tethertest.Config[events.State]{
		Handle: events.Handle,
	})

	h.SendEvent(tether.Event{
		Type:   event.Change,
		Action: "events.change",
		Data:   map[string]string{"value": "option-b"},
	})
	if h.State().ChangeValue != "option-b" {
		t.Fatalf("ChangeValue = %q", h.State().ChangeValue)
	}
}

// ---------- Typed Event Data ----------

// TestTypedEventData demonstrates ev.Int, ev.Float64, and ev.Bool
// for parsing typed values from form submissions.
func TestTypedEventData(t *testing.T) {
	h := tethertest.New(tethertest.Config[events.State]{
		Handle: events.Handle,
	})

	h.SendSubmit("events.typed", map[string]string{
		"qty":    "5",
		"price":  "19.99",
		"urgent": "true",
	})
	want := "quantity=5, price=19.99, urgent=true"
	if h.State().TypedResult != want {
		t.Fatalf("TypedResult = %q, want %q", h.State().TypedResult, want)
	}
}

// TestEventBind demonstrates ev.Bind for struct-tag-based form
// binding - the handler maps form fields to struct fields by tag.
func TestEventBind(t *testing.T) {
	h := tethertest.New(tethertest.Config[events.State]{
		Handle: events.Handle,
	})

	h.SendSubmit("events.bind", map[string]string{
		"name":  "Bob",
		"email": "bob@example.com",
	})
	want := `name="Bob", email="bob@example.com"`
	if h.State().BindResult != want {
		t.Fatalf("BindResult = %q, want %q", h.State().BindResult, want)
	}
}

// ---------- Middleware ----------

// TestMiddlewareChain demonstrates StatefulConfig.Middleware - two ordered
// wrappers show the onion execution order (Outer→Inner→Handle→Inner→Outer).
func TestMiddlewareChain(t *testing.T) {
	logging := func(name string) tether.Middleware[mw.State] {
		return func(next tether.HandleFunc[mw.State]) tether.HandleFunc[mw.State] {
			return func(sess tether.Session, s mw.State, ev tether.Event) mw.State {
				s.ChainLog += name + " → "
				s = next(sess, s, ev)
				s.ChainLog += name + " ← "
				return s
			}
		}
	}

	h := tethertest.New(tethertest.Config[mw.State]{
		Handle: mw.Handle,
		Middleware: []tether.Middleware[mw.State]{
			logging("Outer"),
			logging("Inner"),
		},
	})

	h.Send("mw.ping")
	want := "Outer → Inner → Inner ← Outer ← "
	if h.State().ChainLog != want {
		t.Fatalf("ChainLog = %q, want %q", h.State().ChainLog, want)
	}
}

// TestMiddlewareGuard demonstrates a guard middleware that
// short-circuits the chain - Handle never runs for blocked actions.
func TestMiddlewareGuard(t *testing.T) {
	guard := func(next tether.HandleFunc[mw.State]) tether.HandleFunc[mw.State] {
		return func(sess tether.Session, s mw.State, ev tether.Event) mw.State {
			if ev.Action == "mw.blocked" {
				s.BlockedResult = "Blocked"
				return s
			}
			return next(sess, s, ev)
		}
	}

	h := tethertest.New(tethertest.Config[mw.State]{
		Handle:     mw.Handle,
		Middleware: []tether.Middleware[mw.State]{guard},
	})

	h.Send("mw.blocked")
	if h.State().BlockedResult != "Blocked" {
		t.Fatal("guard should have blocked the action")
	}
	if h.State().LastAction != "" {
		t.Fatal("Handle should not have run")
	}
}

// ---------- Groups & Broadcasting ----------

// TestGroupStateTransitions demonstrates testing a handler that uses
// tether.Group - the state logic (room tracking, messages) works
// with any Session, even though Group.Add/Remove require StatefulSession.
func TestGroupStateTransitions(t *testing.T) {
	h := tethertest.New(tethertest.Config[groups.State]{
		Handle: groups.Handle,
	})

	h.Send("groups.join-alpha")
	if h.State().CurrentRoom != "alpha" {
		t.Fatalf("CurrentRoom = %q", h.State().CurrentRoom)
	}

	h.Send("groups.leave")
	if h.State().CurrentRoom != "" {
		t.Fatalf("CurrentRoom = %q, want empty", h.State().CurrentRoom)
	}
}

// TestBroadcastSenderState demonstrates that the sender sees their
// own broadcast message immediately in state.
func TestBroadcastSenderState(t *testing.T) {
	h := tethertest.New(tethertest.Config[broadcasting.State]{
		Handle: broadcasting.Handle,
	})

	h.SendSubmit("broadcast.send", map[string]string{"broadcast-input": "hello"})
	if len(h.State().Messages) != 1 || h.State().Messages[0].Text != "hello" {
		t.Fatalf("Messages = %v", h.State().Messages)
	}
}

// ---------- Value Store ----------

// TestValueStoreLocalVsShared demonstrates that per-session local
// state and shared tether.Value state remain independent.
func TestValueStoreLocalVsShared(t *testing.T) {
	h := tethertest.New(tethertest.Config[valuestore.State]{
		Handle: valuestore.Handle,
	})

	h.Send("value.local-inc")
	h.Send("value.local-inc")
	h.Send("value.increment") // shared - doesn't affect local

	if h.State().LocalCount != 2 {
		t.Fatalf("LocalCount = %d, want 2", h.State().LocalCount)
	}
}

// ---------- Component Testing ----------

// TestComponent demonstrates NewComponent - the component harness
// tests a tether.Component directly without parent state wiring.
func TestComponent(t *testing.T) {
	h := tethertest.NewComponent(counter.New("c"))

	h.Send("increment")
	h.Send("increment")
	if h.Component().Count != 2 {
		t.Fatalf("Count = %d, want 2", h.Component().Count)
	}

	h.Send("reset")
	if h.Component().Count != 0 {
		t.Fatal("Count should be 0 after reset")
	}
	if !h.HasToast("Counter reset") {
		t.Fatalf("Toast = %q", h.Toast())
	}
}

// TestComponentMount demonstrates the Mounter lifecycle hook  -
// Mount() fires the component's one-time setup.
func TestComponentMount(t *testing.T) {
	h := tethertest.NewComponent(counter.New("likes"))
	h.Mount()

	if !h.HasToast("likes counter ready") {
		t.Fatalf("Toast = %q", h.Toast())
	}
}

// TestComponentHTML demonstrates HTML() on the component harness.
func TestComponentHTML(t *testing.T) {
	h := tethertest.NewComponent(counter.New("c"))
	h.Send("increment")

	html := h.HTML()
	if html == "" {
		t.Fatal("HTML() returned empty string")
	}
}
