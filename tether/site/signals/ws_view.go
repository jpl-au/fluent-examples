package signals

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// RenderWS builds the full signals and directives page for the
// WebSocket variant: BindText, BindShow/BindHide, SetSignal,
// ToggleSignal, ToggleTarget, BindClass, BindAttr, BindValue,
// Optimistic, OptimisticToggle, Batch Signals, Cloak, Permanent,
// Hook, Transition, and FocusTrap.
func RenderWS(s State) node.Node {
	return page.New(
		panel.Card(
			"BindText",
			"Click the button to increment a server-side counter. The server pushes the new value as a signal, and the client updates the text without re-rendering the page.",
			"bind.BindText", panel.WS|panel.SSE,
			layout.Row(
				button.PrimaryAction("Increment Server Counter", "signals.increment"),
				bind.Apply(span.Text(strconv.Itoa(s.Counter)),
					bind.BindText("signals.counter"),
				),
			),
		),

		panel.Card(
			"BindShow / BindHide",
			"Click Toggle Visibility - the server pushes a boolean signal that shows or hides elements instantly on the client. No page re-render happens; the client toggles CSS display directly.",
			"bind.BindShow · bind.BindHide", panel.WS|panel.SSE,
			layout.Stack(
				button.PrimaryAction("Toggle Visibility", "signals.toggle-panel"),
				bind.Apply(layout.Container(
					panel.SignalText("This panel is visible when the signal is true."),
				), bind.BindShow("signals.panel_visible")),
				bind.Apply(span.Text("The panel is hidden."),
					bind.BindHide("signals.panel_visible"),
				),
			),
		),

		panel.Card(
			"SetSignal (Client-Side)",
			"Click any colour button - the text updates instantly with no server round-trip. SetSignal writes directly to the client-side signal store.",
			"bind.SetSignal", panel.AllTransports,
			layout.Row(
				button.Primary("Set to Red", bind.SetSignal("signals.colour", "red")),
				button.Primary("Set to Blue", bind.SetSignal("signals.colour", "blue")),
				button.Primary("Set to Green", bind.SetSignal("signals.colour", "green")),
				bind.Apply(span.Text("none").Dynamic("colour-display"), bind.BindText("signals.colour")),
			),
		),

		panel.Card(
			"ToggleSignal (Client-Side)",
			"Click the button to flip a boolean signal on the client. Combined with BindShow, the panel appears and disappears with no server round-trip.",
			"bind.ToggleSignal · bind.BindShow", panel.AllTransports,
			layout.Stack(
				button.Primary("Toggle Panel", bind.ToggleSignal("signals.toggle_demo")),
				bind.Apply(panel.SignalText("Toggled on!"), bind.BindShow("signals.toggle_demo")),
			),
		),

		panel.Card(
			"ToggleTarget",
			"Click the button - it toggles the hidden attribute on the panel below with no server round-trip. ToggleTarget redirects a toggle to a different element identified by CSS selector.",
			"bind.ToggleTarget · bind.ToggleAttr", panel.AllTransports,
			layout.Stack(
				button.Primary("Toggle Panel",
					bind.ToggleTarget("#toggle-target-panel"),
					bind.ToggleAttr("hidden"),
				),
				layout.Stack(
					panel.SignalText("This panel is controlled by the button above."),
				).ID("toggle-target-panel"),
			),
		),

		panel.Card(
			"BindClass",
			"Click Toggle Highlight - the box below gains or loses the 'highlighted' CSS class based on a boolean signal.",
			"bind.BindClass", panel.AllTransports,
			layout.Stack(
				button.Primary("Toggle Highlight", bind.ToggleSignal("signals.highlight")),
				bind.Apply(panel.ToggleDemo(
					p.New().Text("This box gets the 'highlighted' class when the signal is true."),
				), bind.BindClass("highlighted", "signals.highlight")),
			),
		),

		panel.Card(
			"BindAttr",
			"Click Toggle Lock - the server pushes a boolean signal that adds or removes the disabled attribute on the input. BindAttr drives any HTML attribute from a signal.",
			"bind.BindAttr", panel.WS|panel.SSE,
			layout.Row(
				button.PrimaryAction("Toggle Lock", "signals.toggle-lock"),
				bind.Apply(
					field.Text("attr-demo", "Lock me with the button above"),
					bind.BindAttr("disabled", "signals.input_locked"),
				),
			),
		),

		panel.Card(
			"BindValue",
			"Click Pre-fill - the server pushes a signal that sets the value of the form field directly. BindValue is the pattern for server-driven defaults.",
			"bind.BindValue", panel.WS|panel.SSE,
			layout.Row(
				button.PrimaryAction("Pre-fill from Server", "signals.prefill"),
				bind.Apply(
					field.Text("value-demo", "Server will fill this in"),
					bind.BindValue("signals.prefill_value"),
				),
			),
		),

		panel.Card(
			"Optimistic Updates",
			"Click Like - the 'Liked!' text appears instantly before the server even responds. The client sets the signal optimistically on click.",
			"bind.Optimistic", panel.AllTransports,
			layout.Row(
				button.PrimaryAction("Like", "signals.like",
					bind.Optimistic("signals.liked", "true"),
				),
				bind.Apply(span.Text("Liked!"), bind.BindShow("signals.liked")),
			),
		),

		panel.Card(
			"Optimistic Toggle",
			"Click the button - the signal flips instantly on the client before the event reaches the server. OptimisticToggle is the boolean counterpart of Optimistic.",
			"bind.OptimisticToggle", panel.AllTransports,
			layout.Row(
				button.PrimaryAction("Toggle Favourite", "signals.favourite",
					bind.OptimisticToggle("signals.favourited"),
				),
				bind.Apply(span.Text("Favourited!"), bind.BindShow("signals.favourited")),
				bind.Apply(span.Text("Not favourited"), bind.BindHide("signals.favourited")),
			),
		),

		panel.Card(
			"Batch Signals",
			"Click Reset All - the server pushes counter, visibility, and lock state in a single Signals() call. More efficient than individual Signal() calls when updating multiple values at once.",
			"sess.Signals", panel.WS|panel.SSE,
			button.DangerAction("Reset All Signals", "signals.reset-all"),
		),

		panel.Card(
			"Cloak",
			"The element below is server-rendered with the text 'Loading...' but hidden by a data-tether-cloak attribute. "+
				"When the client JavaScript initialises, the cloak is removed and the BindText signal replaces the placeholder with the counter value. "+
				"On localhost the swap is instant - on a slow connection, Cloak prevents a flash of stale placeholder content (FOUC). "+
				"View source to see the hidden element and its data-tether-cloak attribute.",
			"bind.Cloak", panel.AllTransports,
			layout.Stack(
				hint.Text("The element below is cloaked until JS initialises:"),
				bind.Apply(panel.SignalText("Loading...").Dynamic("cloaked"),
					bind.Cloak(),
					bind.BindText("signals.counter"),
				),
			),
		),

		panel.Card(
			"Permanent",
			"This element is marked permanent - even when the server re-renders the page and the differ patches the DOM, this subtree is never touched. Essential for embedded third-party widgets.",
			"bind.Permanent", panel.AllTransports,
			bind.Apply(panel.Signal(
				p.New().Text("This element is never replaced by the differ."),
			), bind.Permanent()),
		),

		panel.Card(
			"Hook",
			"This element has a JavaScript hook named 'tooltip' attached. When the element is mounted into the DOM, the hook's mounted() callback runs. Hooks also fire on updated() and destroyed().",
			"bind.Hook", panel.AllTransports,
			bind.Apply(panel.HookTarget(
				p.New().Text("Hover over this element to see the tooltip."),
			).SetData("tooltip", "Hello from the tooltip hook!"), bind.Hook("tooltip")),
		),

		panel.Card(
			"Transition",
			"Click to toggle the panel. When it appears, a CSS enter class is applied; when it disappears, a leave class runs. The framework manages the class timing automatically.",
			"bind.Transition", panel.AllTransports,
			layout.Stack(
				button.PrimaryAction("Toggle Transition", "signals.toggle-transition"),
				transitionPanel(s.TransitionVisible),
			),
		),

		panel.Card(
			"FocusTrap",
			"Click the first button below, then press Tab repeatedly. "+
				"Focus moves through the three buttons and wraps back to the first - it never escapes the box. "+
				"Shift+Tab wraps in the other direction. "+
				"Essential for modal dialogs and drawers where focus must stay within the overlay.",
			"bind.FocusTrap", panel.AllTransports,
			bind.Apply(panel.SignalFocusTrap(
				button.Primary("Focusable 1"),
				button.Primary("Focusable 2"),
				button.Primary("Focusable 3"),
				hint.Text("Tab cycles within this box - focus never escapes."),
			), bind.FocusTrap()),
		),
	)
}

// transitionPanel renders the CSS transition demo element. Both
// branches share a Dynamic key so the differ can track the change
// when the panel appears or disappears.
func transitionPanel(visible bool) node.Node {
	if visible {
		return bind.Apply(panel.SignalText("Visible with transition."), bind.Transition("fade")).Dynamic("transition-panel")
	}
	return layout.Container().Dynamic("transition-panel")
}
