package counter

import (
	"testing"

	tether "github.com/jpl-au/tether"
)

func TestLiveCounterUsesServerState(t *testing.T) {
	sess := &tether.CaptureSession{SessionID: "live"}
	c := New("count")
	for range 3 {
		c = c.Handle(sess, tether.Event{Action: "increment", Data: map[string]string{"count": "0"}}).(Counter)
	}
	if c.Count != 3 {
		t.Fatalf("rapid clicks restored stale browser state: count=%d", c.Count)
	}
}

func TestLiveGroupKeepsSiblingState(t *testing.T) {
	g := NewGroup("group")
	g.Left.Count, g.Right.Count = 2, 5
	g = g.Handle(&tether.CaptureSession{SessionID: "live"}, tether.Event{Action: "left.increment", Data: map[string]string{"count": "0", "sibling": "0"}}).(Group)
	if g.Left.Count != 3 || g.Right.Count != 5 {
		t.Fatalf("stale browser state overwrote group: %+v", g)
	}
}

func TestStatelessCounterRestoresEventState(t *testing.T) {
	c := New("count").Handle(&tether.CaptureSession{}, tether.Event{Action: "increment", Data: map[string]string{"count": "7"}}).(Counter)
	if c.Count != 8 {
		t.Fatalf("stateless count=%d", c.Count)
	}
}

func TestStatelessGroupRestoresBothCounters(t *testing.T) {
	for _, tc := range []struct {
		action string
		left   int
		right  int
	}{
		{action: "left.increment", left: 8, right: 4},
		{action: "right.increment", left: 4, right: 8},
	} {
		t.Run(tc.action, func(t *testing.T) {
			g := NewGroup("group").Handle(&tether.CaptureSession{}, tether.Event{
				Action: tc.action,
				Data:   map[string]string{"count": "7", "sibling": "4"},
			}).(Group)
			if g.Left.Count != tc.left || g.Right.Count != tc.right {
				t.Fatalf("stateless group counts=(%d, %d), want (%d, %d)", g.Left.Count, g.Right.Count, tc.left, tc.right)
			}
		})
	}
}
