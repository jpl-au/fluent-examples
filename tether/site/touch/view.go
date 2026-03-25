package touch

import (
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/components/simple/result"
)

// Render builds the touch gestures demo page.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Swipe",
			"Swipe left, right, up, or down on the area below (touch device required). "+
				"The direction is sent to the server. bind.OnSwipe uses a 30px minimum "+
				"distance within 500ms to distinguish swipes from taps.",
			"bind.OnSwipe", panel.AllTransports,
			layout.Stack(
				bind.Apply(
					panel.Signal(hint.Text("Swipe here on a touch device")),
					bind.OnSwipe("touch.swipe"),
				),
				swipeResult(s.SwipeResult),
			),
		),

		panel.Card(
			"Long Press",
			"Press and hold the area below for 500ms (touch device required). "+
				"Cancelled if the finger moves more than 10px. Common mobile "+
				"alternative to right-click.",
			"bind.OnLongPress", panel.AllTransports,
			layout.Stack(
				bind.Apply(
					panel.Signal(hint.Text("Long-press here on a touch device")),
					bind.OnLongPress("touch.longpress"),
				),
				longPressResult(s.LongPressResult),
			),
		),
	)
}

// swipeResult renders the swipe demo outcome.
func swipeResult(val string) node.Node {
	if val == "" {
		return hint.Text("Swipe to see the direction.")
	}
	return result.Success(val)
}

// longPressResult renders the long-press demo outcome.
func longPressResult(val string) node.Node {
	if val == "" {
		return hint.Text("Press and hold to trigger.")
	}
	return result.Success(val)
}
