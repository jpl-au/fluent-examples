// Package store provides the shared kanban board state. All methods
// are safe for concurrent use - a single Board instance is shared
// across all tether sessions.
package store

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// Column identifies a kanban swimlane.
type Column int

const (
	Todo       Column = iota // To Do
	InProgress               // In Progress
	Done                     // Done
	columnCount
)

// String returns the display name for a column.
func (c Column) String() string {
	switch c {
	case Todo:
		return "To Do"
	case InProgress:
		return "In Progress"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}
}

// Columns returns all column values in display order.
func Columns() []Column { return []Column{Todo, InProgress, Done} }

// Event records a single action taken on a card.
type Event struct {
	User    string
	Action  string
	Created time.Time
}

// Card is a single kanban card.
type Card struct {
	ID          string
	Title       string
	Description string
	Column      Column
	CreatedBy   string
	Activity    []Event
}

// Board holds the shared kanban state.
type Board struct {
	mu    sync.RWMutex
	cards map[string]*Card
	order [columnCount][]string // card IDs per column, in display order
}

// New creates a board seeded with example cards.
func New() *Board {
	b := &Board{cards: make(map[string]*Card)}
	b.add(Todo, "Set up CI pipeline", "Configure GitHub Actions for automated builds and tests.", "System")
	b.add(Todo, "Write API documentation", "Document all public endpoints with request and response examples.", "System")
	b.add(InProgress, "Design landing page", "Create mockups for the marketing site hero section.", "System")
	b.add(InProgress, "Implement user auth", "Add session-based authentication with login and registration flows.", "System")
	b.add(Done, "Project kickoff", "Align on goals, timeline, and responsibilities.", "System")
	b.add(Done, "Choose tech stack", "Evaluate options and commit to Go, Tether, and Fluent.", "System")
	return b
}

// Cards returns the cards in a column in display order.
func (b *Board) Cards(col Column) []Card {
	b.mu.RLock()
	defer b.mu.RUnlock()
	ids := b.order[col]
	out := make([]Card, 0, len(ids))
	for _, id := range ids {
		if c, ok := b.cards[id]; ok {
			out = append(out, *c)
		}
	}
	return out
}

// Card returns a single card by ID.
func (b *Board) Card(id string) (Card, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	c, ok := b.cards[id]
	if !ok {
		return Card{}, false
	}
	return *c, true
}

// Create adds a new card to the To Do column and returns it.
func (b *Board) Create(title, desc, user string) Card {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := newID()
	c := &Card{
		ID:          id,
		Title:       title,
		Description: desc,
		Column:      Todo,
		CreatedBy:   user,
		Activity: []Event{
			{User: user, Action: "created this card", Created: time.Now()},
		},
	}
	b.cards[id] = c
	b.order[Todo] = append(b.order[Todo], id)
	return *c
}

// Move relocates a card to a different column.
func (b *Board) Move(id string, col Column, user string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, ok := b.cards[id]
	if !ok {
		return false
	}
	old := c.Column
	if old == col {
		return false
	}
	b.order[old] = remove(b.order[old], id)
	c.Column = col
	b.order[col] = append(b.order[col], id)
	c.Activity = append(c.Activity, Event{
		User:    user,
		Action:  fmt.Sprintf("moved from %s to %s", old, col),
		Created: time.Now(),
	})
	return true
}

// Update modifies a card's title and description.
func (b *Board) Update(id, title, desc, user string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, ok := b.cards[id]
	if !ok {
		return false
	}
	c.Title = title
	c.Description = desc
	c.Activity = append(c.Activity, Event{
		User:    user,
		Action:  "updated this card",
		Created: time.Now(),
	})
	return true
}

// Delete removes a card from the board.
func (b *Board) Delete(id string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, ok := b.cards[id]
	if !ok {
		return false
	}
	b.order[c.Column] = remove(b.order[c.Column], id)
	delete(b.cards, id)
	return true
}

// newID generates a short random identifier for cards.
func newID() string {
	return rand.Text()[:12]
}

// add inserts a seeded card without locking (used during construction).
func (b *Board) add(col Column, title, desc, user string) {
	id := newID()
	b.cards[id] = &Card{
		ID:          id,
		Title:       title,
		Description: desc,
		Column:      col,
		CreatedBy:   user,
		Activity: []Event{
			{User: user, Action: "created this card", Created: time.Now()},
		},
	}
	b.order[col] = append(b.order[col], id)
}

// remove filters an ID from a slice, preserving order.
func remove(ids []string, target string) []string {
	out := ids[:0]
	for _, id := range ids {
		if id != target {
			out = append(out, id)
		}
	}
	return out
}
