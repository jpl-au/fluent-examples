package counter

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
)

// Counter is a self-contained component that manages its own count.
// It implements [tether.Component], [tether.Mounter], and
// [tether.EqualComponent] - the Render method builds the UI, Handle
// processes events, Mount performs one-time setup, and EqualComponent
// lets the framework skip re-rendering when the count is unchanged.
//
// The name field provides unique Dynamic keys so multiple Counter
// instances on the same page don't collide in the differ. Event
// namespacing is handled by [bind.Prefix] on the container element
// - the component uses bare action names.
//
//	counter.New("likes")
//	counter.New("stars")
type Counter struct {
	name string
	// Count is the current counter value, updated by Handle.
	Count int
}

// New creates a Counter. The name is used for Dynamic keys
// (e.g. "likes-count") so multiple counters coexist on one page.
// Event namespacing is handled by bind.Prefix on the container.
func New(name string) Counter {
	return Counter{name: name}
}

// EqualComponent lets the framework skip re-rendering this component
// when the count has not changed - an optimisation for parent handlers
// that update other state without affecting the counter.
func (c Counter) EqualComponent(other tether.Component) bool {
	o, ok := other.(Counter)
	return ok && c.Count == o.Count
}

// Render builds the counter UI: increment/decrement buttons and the
// current count. Actions are bare names - bind.Prefix on the
// container element handles namespacing.
func (c Counter) Render() node.Node {
	countStr := strconv.Itoa(c.Count)
	return layout.Row(
		button.DecrementAction("decrement", bind.EventData("count", countStr)),
		button.IncrementAction("increment", bind.EventData("count", countStr)),
		button.ResetAction("reset"),
		span.Text("Count: "+countStr).Dynamic(c.name+"-count"),
	)
}

// Mount is called once when the component is first added to a session
// via StatefulConfig.Components. It fires a toast to confirm the counter is
// ready - a trivial example of the Mounter lifecycle hook.
func (c Counter) Mount(sess tether.Session) tether.Component {
	sess.Toast(c.name + " counter ready")
	return c
}

// Handle processes counter events. Actions are bare names - the
// framework strips the mount prefix before calling Handle.
func (c Counter) Handle(sess tether.Session, ev tether.Event) tether.Component {
	// Restore count from event data so the component works on stateless
	// pages where state is reconstructed from scratch each request.
	if s, _ := ev.Get("count"); s != "" {
		c.Count, _ = strconv.Atoi(s)
	}
	switch ev.Action {
	case "increment":
		c.Count++
	case "decrement":
		if c.Count > 0 {
			c.Count--
		}
	case "reset":
		c.Count = 0
		sess.Toast("Counter reset")
	}
	return c
}
