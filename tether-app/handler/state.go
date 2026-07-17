package handler

// State is the per-session state for the kanban board. The board
// data itself lives in the shared store; this struct tracks only
// the session's view state and reactive values.
//
// State is serialised to disk via SessionStore (CBOR by default)
// so sessions survive server restarts. All exported fields are
// persisted automatically.
//
// BoardVersion is the memoisation key for column rendering. The
// handler increments it on board mutations (create, save, move,
// delete) so the Memoiser skips unchanged columns. See view.go.
type State struct {
	// SessionID identifies this session for viewer tracking.
	SessionID string
	// Name is the user's display name, set on the landing page.
	// When empty, the landing page is shown instead of the board.
	Name string
	// View is "board" or "detail".
	View string
	// SelectedID is the card being viewed in detail mode.
	SelectedID string
	// OnlineCount tracks connected sessions for the header badge.
	OnlineCount int
	// BoardVersion is the memoisation key for column rendering.
	// Incremented on board mutations so Memoise cache misses and
	// columns re-render.
	BoardVersion int
	// MenuOpen tracks whether the card detail's overflow menu is
	// showing. The menu renders only when open so its click-outside
	// binding (bind.Outside) exists solely while it is dismissable.
	MenuOpen bool
}
