// Package store provides the shared kanban board state. All methods
// are safe for concurrent use - a single Board instance is shared
// across all tether sessions.
package store

import (
	"fmt"
	"sync"
	"sync/atomic"
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

// Card is a single kanban card.
type Card struct {
	ID          string
	Title       string
	Description string
	Column      Column
}

// Board holds the shared kanban state.
type Board struct {
	mu    sync.RWMutex
	cards map[string]*Card
	order [columnCount][]string // card IDs per column, in display order
	seq   atomic.Int64
}

// New creates a board seeded with example cards.
func New() *Board {
	b := &Board{cards: make(map[string]*Card)}
	b.add(Todo, "Set up CI pipeline", "Configure GitHub Actions for automated builds and tests.")
	b.add(Todo, "Write API documentation", "Document all public endpoints with request and response examples.")
	b.add(InProgress, "Design landing page", "Create mockups for the marketing site hero section.")
	b.add(InProgress, "Implement user auth", "Add session-based authentication with login and registration flows.")
	b.add(Done, "Project kickoff", "Align on goals, timeline, and responsibilities.")
	b.add(Done, "Choose tech stack", "Evaluate options and commit to Go, Tether, and Fluent.")
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
func (b *Board) Create(title, desc string) Card {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := fmt.Sprintf("card-%d", b.seq.Add(1))
	c := &Card{ID: id, Title: title, Description: desc, Column: Todo}
	b.cards[id] = c
	b.order[Todo] = append(b.order[Todo], id)
	return *c
}

// Move relocates a card to a different column.
func (b *Board) Move(id string, col Column) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, ok := b.cards[id]
	if !ok {
		return false
	}
	old := c.Column
	b.order[old] = remove(b.order[old], id)
	c.Column = col
	b.order[col] = append(b.order[col], id)
	return true
}

// Update modifies a card's title and description.
func (b *Board) Update(id, title, desc string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, ok := b.cards[id]
	if !ok {
		return false
	}
	c.Title = title
	c.Description = desc
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

// add inserts a seeded card without locking (used during construction).
func (b *Board) add(col Column, title, desc string) {
	id := fmt.Sprintf("card-%d", b.seq.Add(1))
	b.cards[id] = &Card{ID: id, Title: title, Description: desc, Column: col}
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
