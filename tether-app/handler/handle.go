package handler

import (
	"fmt"
	"time"

	"github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// Handle processes all kanban board events. The handler demonstrates
// three update strategies depending on what changed:
//
//   - Signals for presence events (card.typing, card.select,
//     card.back). Only text content changes, so we push signal
//     values directly - no render cycle at all.
//   - Full refresh for board mutations (card.save, card.move,
//     card.delete, card.new). The board structure changed, so all
//     sessions re-render. Memoise ensures unchanged columns are
//     skipped.
//   - The online count signal is pushed alongside every refresh
//     so the header badge stays current.
func Handle(board *store.Board, group *tether.Group[State], viewers *viewers) func(tether.Session, State, tether.Event) State {
	return func(sess tether.Session, s State, ev tether.Event) State {
		switch ev.Action {
		case "name.set":
			name, _ := ev.Get("name")
			if name != "" {
				s.Name = name
				s.View = "board"
				sess.ReplaceURL("/")
				board.Claim(name)
				s.BoardVersion++
				refresh(group)
			}

		case "card.new":
			s.View = "detail"
			s.SelectedID = ""
			viewers.Presence.Clear(sess.ID())
			sess.ReplaceURL("/new")

		// Signals: card.typing only updates presence indicators.
		// No board data changed, so a full render would waste work.
		// Instead we push signal values to all sessions - the client
		// updates the bound text elements directly.
		//
		// A delayed re-push runs after the typing timeout so the
		// indicator transitions from "editing" back to "viewing"
		// when the user stops typing. The TypingOnCard function
		// already filters by timestamp, so the re-push naturally
		// clears stale typing indicators.
		case "card.typing":
			viewers.SetTyping(sess.ID())
			if info, ok := viewers.Get(sess.ID()); ok {
				cardID := info.CardID
				pushPresenceSignals(group, viewers, cardID)
				go func() {
					time.Sleep(typingTimeout + 500*time.Millisecond)
					pushPresenceSignals(group, viewers, cardID)
				}()
			}
			return s

		case "card.save":
			id, _ := ev.Get("id")
			title, _ := ev.Get("title")
			desc, _ := ev.Get("description")
			if title == "" {
				return s
			}
			if id == "" {
				c := board.Create(title, desc, s.Name)
				sess.ReplaceURL("/card/" + c.ID)
				s.View = "detail"
				s.SelectedID = c.ID
				notify(group, sess, fmt.Sprintf("%s created \"%s\"", s.Name, title))
			} else {
				board.Update(id, title, desc, s.Name)
				sess.ReplaceURL("/")
				s.View = "board"
				s.SelectedID = ""
				notify(group, sess, fmt.Sprintf("%s updated \"%s\"", s.Name, title))
			}
			// Board data changed - bump version so Memoise re-renders
			// affected columns on every session.
			s.BoardVersion++
			refresh(group)

		case "card.move":
			id, _ := ev.Get("id")
			col, _ := ev.Int("column")
			idx, idxErr := ev.Int("index")
			if idxErr != nil {
				idx = -1
			}
			if c, ok := board.Card(id); ok {
				board.MoveAt(id, store.Column(col), idx, s.Name)
				notify(group, sess, fmt.Sprintf("%s moved \"%s\" to %s", s.Name, c.Title, store.Column(col)))
			}
			s.BoardVersion++
			refresh(group)

		// Signals: card.select updates presence indicators for the
		// card being viewed. Only the local session's state changes
		// (view mode, selected card); other sessions just see updated
		// presence text via signals.
		case "card.select":
			id, _ := ev.Get("id")
			s.View = "detail"
			s.SelectedID = id
			viewers.View(sess.ID(), id, s.Name)
			sess.ReplaceURL("/card/" + id)
			pushPresenceSignals(group, viewers, id)

		// Signals: card.back clears presence for this session.
		// Other sessions see updated presence via signals.
		case "card.back":
			oldID := s.SelectedID
			s.View = "board"
			s.SelectedID = ""
			viewers.Presence.Clear(sess.ID())
			sess.ReplaceURL("/")
			if oldID != "" {
				pushPresenceSignals(group, viewers, oldID)
			}

		case "card.delete":
			id, _ := ev.Get("id")
			if c, ok := board.Card(id); ok {
				board.Delete(id)
				notify(group, sess, fmt.Sprintf("%s deleted \"%s\"", s.Name, c.Title))
			}
			s.View = "board"
			s.SelectedID = ""
			sess.ReplaceURL("/")
			s.BoardVersion++
			refresh(group)
		}
		return s
	}
}

// refresh triggers a re-render on every connected session. Used for
// board mutations where the board structure has changed. Increments
// BoardVersion so the Memoiser cache misses and columns re-render.
// The online count updates separately via WatchValue.
func refresh(group *tether.Group[State]) {
	group.Broadcast(func(_ *tether.StatefulSession[State], s State) State {
		s.BoardVersion++
		return s
	})
}

// notify sends a toast to every session except the one that caused
// the action. Named users see what others are doing in real time.
func notify(group *tether.Group[State], sender tether.Session, msg string) {
	group.BroadcastOthers(sender, func(sess *tether.StatefulSession[State], s State) State {
		sess.Toast(msg)
		return s
	})
}
