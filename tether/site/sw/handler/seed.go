package handler

import (
	"net/http"

	"github.com/jpl-au/fluent-examples/tether/site/sw/state"
)

// InitialState creates the starting state for a new session,
// seeding the online count so the header badge is correct on the
// very first render before any tether.Observe subscription fires.
func InitialState(_ *http.Request) state.State {
	return state.State{OnlineCount: Presence.OnlineCount.Load()}
}
