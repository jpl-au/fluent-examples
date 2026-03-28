package components

import (
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the component demo page. The two counters (Likes and
// Stars) are wired via StatefulConfig.Components on the handler - the
// framework dispatches their events automatically. This page only
// calls Render; Handle is never involved.
func Render(s State) node.Node {
	return cpage.New(
		panel.Card("StatefulConfig.Components",
			"These two counters are wired declaratively via StatefulConfig.Components and tether.Mount. "+
				"The framework intercepts events matching each mount's prefix and dispatches them "+
				"to the component's Handle method - the page's Handle function never sees them. "+
				"Each counter is an independent instance of the same component.Counter type.",
			"tether.Mount · StatefulConfig.Components", panel.AllTransports,
			layout.Stack(
				hint.Text("Likes counter:"),
				bind.Apply(layout.Container(s.Likes.Render()), bind.Prefix("likes")).Dynamic("likes-section"),
				hint.Text("Stars counter:"),
				bind.Apply(layout.Container(s.Stars.Render()), bind.Prefix("stars")).Dynamic("stars-section"),
			),
		),

		panel.Card("Component Interface",
			"Every tether.Component implements two methods: Render() builds the UI tree, "+
				"Handle() processes events and returns the updated component. Components are value "+
				"types - Handle returns a new value, the receiver is never mutated. The component "+
				"has no knowledge of the parent's state type.",
			"tether.Component", panel.AllTransports,
			layout.Stack(
				hint.Text("The counter component's Handle method calls sess.Toast(\"Counter reset\") when reset is clicked - session side-effects work inside components just like they do in the page handler."),
			),
		),

		panel.Card("Event.Target",
			"When StatefulConfig.Components dispatches an event, the framework sets Event.Target to "+
				"the mount's prefix (e.g. \"likes\" or \"stars\"). Middleware and logging can inspect "+
				"this field to identify which component handled the event without parsing the action string.",
			"Event.Target", panel.AllTransports,
			layout.Stack(
				hint.Text("Open the server console and click a counter button - the middleware log line includes the component target."),
			),
		),
	)
}
