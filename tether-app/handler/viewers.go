package handler

import (
	"sync"
	"time"
)

// Viewers tracks which card each session is currently viewing and
// whether they are actively typing.
type Viewers struct {
	mu     sync.RWMutex
	cards  map[string]string    // sessionID → cardID
	names  map[string]string    // sessionID → user name
	typing map[string]time.Time // sessionID → last typing timestamp
}

// NewViewers creates an empty viewer tracker.
func NewViewers() *Viewers {
	return &Viewers{
		cards:  make(map[string]string),
		names:  make(map[string]string),
		typing: make(map[string]time.Time),
	}
}

// Set marks a session as viewing a card.
func (v *Viewers) Set(sessionID, cardID, name string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.cards[sessionID] = cardID
	v.names[sessionID] = name
	delete(v.typing, sessionID)
}

// Clear removes a session's viewing state.
func (v *Viewers) Clear(sessionID string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	delete(v.cards, sessionID)
	delete(v.names, sessionID)
	delete(v.typing, sessionID)
}

// SetTyping marks a session as actively typing on their current card.
// Expires after typingTimeout so stale indicators clear automatically.
func (v *Viewers) SetTyping(sessionID string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.typing[sessionID] = time.Now()
}

// typingTimeout is how long a typing indicator persists after the
// last keystroke before it's considered stale.
const typingTimeout = 3 * time.Second

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

// Typing returns the names of users actively typing on the given
// card, excluding the specified session.
func (v *Viewers) Typing(cardID, excludeSession string) []string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	now := time.Now()
	var out []string
	for sid, cid := range v.cards {
		if cid == cardID && sid != excludeSession {
			if t, ok := v.typing[sid]; ok && now.Sub(t) < typingTimeout {
				if name := v.names[sid]; name != "" {
					out = append(out, name)
				}
			}
		}
	}
	return out
}
