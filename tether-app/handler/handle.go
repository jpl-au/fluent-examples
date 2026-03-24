package handler

import (
	"fmt"

	"github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether-app/store"
)

// Handle processes all kanban board events. Board mutations are
// applied to the shared store and then broadcast to every session
// via Group.Broadcast so all browsers stay in sync.
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
				refresh(group)
			}

		case "card.new":
			s.View = "detail"
			s.SelectedID = ""
			viewers.Presence.Clear(sess.ID())
			sess.ReplaceURL("/new")
			refresh(group)

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
			refresh(group)

		case "card.typing":
			viewers.SetTyping(sess.ID())
			refresh(group)

		case "card.select":
			id, _ := ev.Get("id")
			s.View = "detail"
			s.SelectedID = id
			viewers.View(sess.ID(), id, s.Name)
			sess.ReplaceURL("/card/" + id)
			refresh(group)

		case "card.back":
			s.View = "board"
			s.SelectedID = ""
			viewers.Presence.Clear(sess.ID())
			sess.ReplaceURL("/")
			refresh(group)

		case "card.delete":
			id, _ := ev.Get("id")
			if c, ok := board.Card(id); ok {
				board.Delete(id)
				notify(group, sess, fmt.Sprintf("%s deleted \"%s\"", s.Name, c.Title))
			}
			s.View = "board"
			s.SelectedID = ""
			sess.ReplaceURL("/")
			refresh(group)
		}
		return s
	}
}

// refresh triggers a re-render on every connected session.
func refresh(group *tether.Group[State]) {
	group.Broadcast(func(_ *tether.StatefulSession[State], s State) State {
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
