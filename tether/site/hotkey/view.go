package hotkey

import (
	"github.com/jpl-au/fluent/html5/div"
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
			"Press Ctrl+K anywhere on this page. The server receives the event "+
				"and updates the result below. bind.Hotkey fires regardless of which "+
				"element has focus.",
			"bind.Hotkey", panel.WS|panel.SSE,
			layout.Stack(
				hint.Text("Press Ctrl+K, Escape, Ctrl+/, or Shift+? to trigger a hotkey."),
				lastCombo(s.LastCombo),
			),
		),
		div.New(
			bind.Apply(div.New(), bind.Hotkey("ctrl+k", "hotkey.triggered")),
			bind.Apply(div.New(), bind.Hotkey("escape", "hotkey.triggered")),
			bind.Apply(div.New(), bind.Hotkey("ctrl+/", "hotkey.triggered")),
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
