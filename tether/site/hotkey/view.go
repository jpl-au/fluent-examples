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
			"bind.Hotkey", panel.WS,
			layout.Stack(
				hint.Text("Press Ctrl+K or Escape to trigger a hotkey."),
				lastCombo(s.LastCombo),
			),
		),
		bind.Apply(
			div.New().Class("sr-only"),
			bind.Hotkey("ctrl+k", "hotkey.triggered"),
			bind.Hotkey("escape", "hotkey.triggered"),
		),
	)
}

// lastCombo renders the last triggered hotkey combo.
func lastCombo(combo string) node.Node {
	if combo == "" {
		return span.New().Class("hint").Text("No hotkey triggered yet.").Dynamic("combo")
	}
	return span.New().Class("result-block").Text("Last hotkey: " + combo).Dynamic("combo")
}
