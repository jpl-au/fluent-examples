package hotkey

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the hotkey demo page. The hotkey bindings are
// attached to a wrapper div; the JS runtime scans for
// data-tether-hotkey-* anywhere in the DOM.
func Render(s State) node.Node {
	return page.New(
		panel.Card(
			"Global Hotkeys",
			"Press Mod+K anywhere on this page - that is Cmd+K on macOS and "+
				"Ctrl+K everywhere else. The \"mod\" modifier matches the "+
				"platform's primary command key, so one binding feels native on "+
				"every OS. Ctrl and Meta (Cmd) are distinct modifiers: a ctrl "+
				"combo never swallows Cmd+C on a Mac.",
			"bind.Hotkey", panel.WS|panel.SSE,
			layout.Stack(
				hint.Text("Press Mod+K, Escape, Mod+/, or Shift+? to trigger a hotkey."),
				lastCombo(s.LastCombo),
			),
		),
		panel.Card(
			"Typing Is Not a Shortcut",
			"Hotkeys without Ctrl/Meta/Alt do not fire while an editable "+
				"element has focus. Focus the field below and press Shift+? - it "+
				"types a question mark instead of firing the hotkey. Blur the "+
				"field and press it again to trigger it.",
			"", panel.WS|panel.SSE,
			input.Text("hotkey-probe", "").
				Placeholder("type ? or k in here - no hotkey fires").
				Class("text-input"),
		),
		div.New(
			bind.Apply(div.New(), bind.Hotkey("mod+k", "hotkey.triggered")),
			bind.Apply(div.New(), bind.Hotkey("escape", "hotkey.triggered")),
			bind.Apply(div.New(), bind.Hotkey("mod+/", "hotkey.triggered")),
			bind.Apply(div.New(), bind.Hotkey("shift+?", "hotkey.triggered")),
		).Class("sr-only"),
	)
}

// lastCombo renders the last triggered hotkey combo.
func lastCombo(combo string) node.Node {
	if combo == "" {
		return span.Text("No hotkey triggered yet.").Class("hint").Dynamic("combo")
	}
	return span.Text("Last hotkey: " + combo).Class("result-block").Dynamic("combo")
}
