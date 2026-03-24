package handler

// State is the per-session state for the kanban board. The board
// data itself lives in the shared store; this struct tracks only
// the session's view state and reactive values.
type State struct {
	// Name is the user's display name, set on the landing page.
	// When empty, the landing page is shown instead of the board.
	Name string
	// View is "board" or "detail".
	View string
	// SelectedID is the card being viewed in detail mode.
	SelectedID string
	// OnlineCount tracks connected sessions for the header badge.
	OnlineCount int
}
