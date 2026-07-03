package signals

// State is the per-session state for the signals demo.
type State struct {
	// OnlineCount tracks connected sessions for the header badge.
	OnlineCount int
	// Counter is incremented by the Text demo button.
	Counter int
	// PanelVisible controls the Show/Hide demo panel.
	PanelVisible bool
	// TransitionVisible controls the CSS transition demo panel.
	TransitionVisible bool
	// Locked tracks whether the Attr demo input is disabled.
	Locked bool
	// Favourited tracks the OptimisticToggle demo state.
	Favourited bool
}
