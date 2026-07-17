package viewtransitions

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the view transitions demo page.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"View Transitions",
			"This page opts in with tether.Client{ViewTransitions: true}, so every server-driven DOM update is wrapped in a native View Transition and cross-fades instead of snapping. Click Move - the box carries a CSS view-transition-name, so the browser animates it smoothly to its new slot (a shared-element transition) rather than making it jump. It is a single opt-in switch and needs no per-element API. Where the browser lacks the API, or the visitor has prefers-reduced-motion set, the update applies instantly with nothing else changed.",
			"tether.Client · view-transition-name", panel.AllTransports,
			layout.Stack(
				button.PrimaryAction("Move the Box", "vt.toggle"),
				stage(s.Swapped),
				hint.Text("Nothing but the box position changes between clicks - the smooth motion is the browser animating the transition."),
			),
		),
	)
}

// stage renders the two slots the box moves between. The box keeps a
// stable view-transition-name in both states, which is what lets the
// browser animate it across the gap rather than fading one out and
// another in.
func stage(swapped bool) node.Node {
	top := slot()
	bottom := slot()
	if swapped {
		bottom = slot(box())
	} else {
		top = slot(box())
	}
	return layout.Stack(top, bottom).Dynamic("vt-stage")
}

// box is the shared element that animates between slots.
func box() node.Node {
	return div.New(span.Text("Box")).
		Style("view-transition-name: vt-box; background: #4f46e5; color: #fff; " +
			"padding: 1rem 1.5rem; border-radius: 8px; width: fit-content; font-weight: 600")
}

// slot is one landing area for the box, with a fixed height so the two
// slots hold their shape whether or not they contain the box.
func slot(children ...node.Node) node.Node {
	return div.New(children...).
		Style("min-height: 64px; display: flex; align-items: center; " +
			"padding: 0.5rem; border: 1px dashed var(--border, #d1d5db); border-radius: 8px")
}
