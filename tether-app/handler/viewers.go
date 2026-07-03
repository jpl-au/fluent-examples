package handler

import (
	"strings"
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

// pushPresenceSignals pushes signal values for a card's typing and
// viewing indicators to all connected sessions. Each session sees
// different text because it excludes itself from the list.
//
// Uses Group.Each instead of Broadcast because only signal values
// change - the DOM structure is unchanged, so a render cycle would
// be wasted work. Signals push the value directly to the client's
// signal store, and bind.Text updates the element in place.
func pushPresenceSignals(group *tether.Group[State], v *viewers, cardID string) {
	group.Each(func(sess *tether.StatefulSession[State]) {
		typing := v.TypingOnCard(cardID, sess.ID())
		viewing := v.ViewingCard(cardID, sess.ID())
		sess.Signals(map[string]any{
			"typing-" + cardID:  formatTyping(typing),
			"viewing-" + cardID: formatViewing(viewing, typing),
		})
	})
}

// formatTyping returns the display text for typing indicators.
func formatTyping(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0] + " is editing..."
	default:
		return strings.Join(names, ", ") + " are editing..."
	}
}

// formatViewing returns the display text for viewing indicators,
// excluding anyone who is already shown as typing.
func formatViewing(viewing, typing []string) string {
	typingSet := make(map[string]bool, len(typing))
	for _, n := range typing {
		typingSet[n] = true
	}
	var viewOnly []string
	for _, n := range viewing {
		if !typingSet[n] {
			viewOnly = append(viewOnly, n)
		}
	}
	switch len(viewOnly) {
	case 0:
		return ""
	case 1:
		return viewOnly[0] + " is viewing this"
	default:
		return strings.Join(viewOnly, ", ") + " are viewing this"
	}
}
