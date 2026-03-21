package signals

import tether "github.com/jpl-au/tether"

// Handle processes events on the signals page, pushing updated signal
// values to the client after each state change. Shared by both the WS
// and SSE handlers - the SSE view simply omits demos for actions it
// does not render.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "signals.increment":
		s.Counter++
		sess.Signal("signals.counter", s.Counter)
	case "signals.toggle-panel":
		s.PanelVisible = !s.PanelVisible
		sess.Signal("signals.panel_visible", s.PanelVisible)
	case "signals.toggle-lock":
		s.Locked = !s.Locked
		sess.Signal("signals.input_locked", s.Locked)
	case "signals.prefill":
		sess.Signal("signals.prefill_value", "hello@example.com")
	case "signals.like":
		sess.Signal("signals.liked", true)
	case "signals.favourite":
		s.Favourited = !s.Favourited
		sess.Signal("signals.favourited", s.Favourited)
	case "signals.toggle-transition":
		s.TransitionVisible = !s.TransitionVisible
	case "signals.reset-all":
		s.Counter = 0
		s.PanelVisible = false
		s.Locked = false
		s.Favourited = false
		s.TransitionVisible = false
		sess.Signals(map[string]any{
			"signals.counter":       0,
			"signals.panel_visible": false,
			"signals.input_locked":  false,
			"signals.liked":         false,
			"signals.favourited":    false,
			"signals.prefill_value": "",
			"signals.toggle_demo":   false,
			"signals.highlight":     false,
			"signals.colour":        "",
		})
	}
	return s
}
