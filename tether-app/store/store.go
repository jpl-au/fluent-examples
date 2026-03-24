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
	CreatedAt   time.Time
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
	now := time.Now()
	b.addAt(Todo, "Set up CI pipeline", "Configure GitHub Actions for automated builds and tests.", "System", now.Add(-2*time.Hour))
	b.addAt(Todo, "Write API documentation", "Document all public endpoints with request and response examples.", "System", now.Add(-90*time.Minute))
	b.addAt(InProgress, "Design landing page", "Create mockups for the marketing site hero section.", "System", now.Add(-5*time.Hour))
	b.addAt(InProgress, "Implement user auth", "Add session-based authentication with login and registration flows.", "System", now.Add(-4*time.Hour))
	b.addAt(Done, "Project kickoff", "Align on goals, timeline, and responsibilities.", "System", now.Add(-24*time.Hour))
	b.addAt(Done, "Choose tech stack", "Evaluate options and commit to Go, Tether, and Fluent.", "System", now.Add(-20*time.Hour))
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
	now := time.Now()
	id := newID()
	c := &Card{
		ID:          id,
		Title:       title,
		Description: desc,
		Column:      Todo,
		CreatedBy:   user,
		CreatedAt:   now,
		Activity: []Event{
			{User: user, Action: "created this card", Created: now},
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

// MoveAt relocates a card to a column at a specific index. An index
// of -1 appends to the end. Also handles within-column reordering.
func (b *Board) MoveAt(id string, col Column, idx int, user string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, ok := b.cards[id]
	if !ok {
		return false
	}
	old := c.Column
	b.order[old] = remove(b.order[old], id)
	c.Column = col

	order := b.order[col]
	switch {
	case idx < 0 || idx >= len(order):
		b.order[col] = append(order, id)
	default:
		b.order[col] = append(order[:idx], append([]string{id}, order[idx:]...)...)
	}

	action := fmt.Sprintf("moved from %s to %s", old, col)
	if old == col {
		action = fmt.Sprintf("reordered in %s", col)
	}
	c.Activity = append(c.Activity, Event{
		User:    user,
		Action:  action,
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

// Claim assigns all cards with no real owner to the given user.
// Called when the first person joins so the seed data feels like
// theirs rather than belonging to "System".
func (b *Board) Claim(user string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, c := range b.cards {
		if c.CreatedBy == "System" {
			c.CreatedBy = user
			c.Activity = append(c.Activity, Event{
				User:    user,
				Action:  "claimed this card",
				Created: time.Now(),
			})
		}
	}
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

// TimeAgo returns a human-readable relative time string.
func TimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

// newID generates a short random identifier for cards.
func newID() string {
	return rand.Text()[:12]
}

// addAt inserts a seeded card without locking (used during construction).
func (b *Board) addAt(col Column, title, desc, user string, at time.Time) {
	id := newID()
	b.cards[id] = &Card{
		ID:          id,
		Title:       title,
		Description: desc,
		Column:      col,
		CreatedBy:   user,
		CreatedAt:   at,
		Activity: []Event{
			{User: user, Action: "created this card", Created: at},
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
