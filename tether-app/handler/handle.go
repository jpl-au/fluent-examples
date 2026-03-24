package handler

import (
	"github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// Handle processes all kanban board events. Board mutations are
// applied to the shared store and then broadcast to every session
// via Group.Broadcast so all browsers stay in sync.
func Handle(board *store.Board, group *tether.Group[State]) func(tether.Session, State, tether.Event) State {
	return func(sess tether.Session, s State, ev tether.Event) State {
		switch ev.Action {
		case "name.set":
			name, _ := ev.Get("name")
			if name != "" {
				s.Name = name
				s.View = "board"
				sess.ReplaceURL("/")
			}

		case "card.new":
			s.View = "detail"
			s.SelectedID = ""
			sess.ReplaceURL("/new")

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
			} else {
				board.Update(id, title, desc, s.Name)
				sess.ReplaceURL("/")
				s.View = "board"
				s.SelectedID = ""
			}
			refresh(group)

		case "card.move":
			id, _ := ev.Get("id")
			col, _ := ev.Int("column")
			board.Move(id, store.Column(col), s.Name)
			refresh(group)

		case "card.select":
			id, _ := ev.Get("id")
			s.View = "detail"
			s.SelectedID = id
			sess.ReplaceURL("/card/" + id)

		case "card.back":
			s.View = "board"
			s.SelectedID = ""
			sess.ReplaceURL("/")

		case "card.delete":
			id, _ := ev.Get("id")
			board.Delete(id)
			s.View = "board"
			s.SelectedID = ""
			sess.ReplaceURL("/")
			refresh(group)
		}
		return s
	}
}

// refresh triggers a re-render on every connected session. Because
// the render function reads from the shared store, all sessions pick
// up the latest board state.
func refresh(group *tether.Group[State]) {
	group.Broadcast(func(_ *tether.StatefulSession[State], s State) State {
		return s
	})
}
