package handler

import (
	"time"

	tether "github.com/jpl-au/tether"
)

// ViewInfo tracks what a session is doing on the board.
type ViewInfo struct {
	CardID string
	Name   string
	Typing time.Time // zero means not typing
}

// typingTimeout is how long a typing indicator persists after the
// last keystroke before it's considered stale.
const typingTimeout = 3 * time.Second

// viewers wraps tether.Presence[ViewInfo] with helper methods for
// querying viewing and typing state per card.
type viewers struct {
	*tether.Presence[ViewInfo]
}

// newViewers creates a viewer tracker backed by tether.Presence.
func newViewers() *viewers {
	return &viewers{tether.NewPresence[ViewInfo]()}
}

// View marks a session as viewing a card.
func (v *viewers) View(sessionID, cardID, name string) {
	v.Set(sessionID, ViewInfo{CardID: cardID, Name: name})
}

// SetTyping marks a session as actively typing on their current card.
func (v *viewers) SetTyping(sessionID string) {
	if info, ok := v.Get(sessionID); ok {
		info.Typing = time.Now()
		v.Set(sessionID, info)
	}
}

// ViewingCard returns the names of users viewing a card, excluding
// the given session.
func (v *viewers) ViewingCard(cardID, exclude string) []string {
	var out []string
	v.Each(exclude, func(_ string, info ViewInfo) {
		if info.CardID == cardID {
			out = append(out, info.Name)
		}
	})
	return out
}

// TypingOnCard returns the names of users actively typing on a card,
// excluding the given session.
func (v *viewers) TypingOnCard(cardID, exclude string) []string {
	now := time.Now()
	var out []string
	v.Each(exclude, func(_ string, info ViewInfo) {
		if info.CardID == cardID && !info.Typing.IsZero() && now.Sub(info.Typing) < typingTimeout {
			out = append(out, info.Name)
		}
	})
	return out
}
