package handler

import "sync"

// Viewers tracks which card each session is currently viewing.
// Thread-safe for concurrent access across sessions.
type Viewers struct {
	mu    sync.RWMutex
	cards map[string]string // sessionID → cardID
	names map[string]string // sessionID → user name
}

// NewViewers creates an empty viewer tracker.
func NewViewers() *Viewers {
	return &Viewers{
		cards: make(map[string]string),
		names: make(map[string]string),
	}
}

// Set marks a session as viewing a card.
func (v *Viewers) Set(sessionID, cardID, name string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.cards[sessionID] = cardID
	v.names[sessionID] = name
}

// Clear removes a session's viewing state.
func (v *Viewers) Clear(sessionID string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	delete(v.cards, sessionID)
	delete(v.names, sessionID)
}

// For returns the names of users currently viewing the given card,
// excluding the specified session (so you don't see your own name).
func (v *Viewers) For(cardID, excludeSession string) []string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	var out []string
	for sid, cid := range v.cards {
		if cid == cardID && sid != excludeSession {
			if name := v.names[sid]; name != "" {
				out = append(out, name)
			}
		}
	}
	return out
}
