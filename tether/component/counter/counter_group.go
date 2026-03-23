package counter

import (
	"strconv"
	"strings"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/columns"
	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/result"
)

// Group is a parent component containing two independent child
// counters. It demonstrates nested event routing: each child manages
// its own count, while the parent can observe both (total) and act on
// both at once (Reset All). Events route through the chain:
// group → left/right → increment/decrement.
//
// Event namespacing is handled by nested bind.Prefix containers  -
// the component uses bare action names.
type Group struct {
	name string
	// Left is the first child counter.
	Left Counter
	// Right is the second child counter.
	Right Counter
}

// NewGroup creates a Group with two nested child counters.
func NewGroup(name string) Group {
	return Group{
		name:  name,
		Left:  New("left"),
		Right: New("right"),
	}
}

// Render builds the group UI: two child counters side by side, a
// combined total, and a parent-level Reset All button. Each child
// is wrapped with bind.Prefix so the nested prefix chain routes
// events correctly (e.g. "group.left.increment").
func (g Group) Render() node.Node {
	leftStr := strconv.Itoa(g.Left.Count)
	rightStr := strconv.Itoa(g.Right.Count)
	total := g.Left.Count + g.Right.Count

	return div.New(
		columns.New(
			div.New(
				result.Label("Counter A"),
				bind.Apply(div.New(counterButtons(leftStr, rightStr)), bind.Prefix("left")),
				span.Text("Count: "+leftStr).Dynamic(g.name+"-left-count"),
			),
			div.New(
				result.Label("Counter B"),
				bind.Apply(div.New(counterButtons(rightStr, leftStr)), bind.Prefix("right")),
				span.Text("Count: "+rightStr).Dynamic(g.name+"-right-count"),
			),
		),
		layout.Row(
			span.Text("Combined total: "+strconv.Itoa(total)).Dynamic(g.name+"-total"),
			button.SmallDangerAction("Reset All", "reset-all"),
		),
	).Class("nested-group")
}

// counterButtons renders -/+ buttons that carry both their own count
// and the sibling's count so neither is lost on stateless round-trips.
// Actions are bare names - bind.Prefix handles namespacing.
func counterButtons(count, sibling string) node.Node {
	return layout.Row(
		button.DecrementAction("decrement", bind.EventData("count", count), bind.EventData("sibling", sibling)),
		button.IncrementAction("increment", bind.EventData("count", count), bind.EventData("sibling", sibling)),
	)
}

// Handle processes group events and delegates child events via
// RouteTyped. The parent intercepts "reset-all" to zero both
// children; everything else is routed by prefix. The sibling's
// count is restored from event data before routing so both
// children survive stateless page reconstruction.
func (g Group) Handle(sess tether.Session, ev tether.Event) tether.Component {
	// Restore the sibling counter's count from event data so both
	// children survive stateless page reconstruction.
	if s, _ := ev.Get("sibling"); s != "" {
		n, _ := strconv.Atoi(s)
		if strings.HasPrefix(ev.Action, "left.") {
			g.Right.Count = n
		} else if strings.HasPrefix(ev.Action, "right.") {
			g.Left.Count = n
		}
	}

	switch ev.Action {
	case "reset-all":
		g.Left.Count = 0
		g.Right.Count = 0
		sess.Toast("All counters reset")
	default:
		g.Left = tether.RouteTyped(g.Left, "left", sess, ev)
		g.Right = tether.RouteTyped(g.Right, "right", sess, ev)
	}
	return g
}
