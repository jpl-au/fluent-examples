// Package store provides thread-safe in-memory storage for contacts
// and their associated notes. It is seeded with example data on
// initialisation.
package store

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Contact holds a person's details and any attached notes.
type Contact struct {
	ID    string
	Name  string
	Email string
	Phone string
	Notes []Note
}

// Note is a timestamped text entry attached to a contact.
type Note struct {
	ID      string
	Content string
	Created time.Time
}

// Package-level state. The contacts map is guarded by mu (read-lock
// for queries, write-lock for mutations). The atomic counters are
// independent of the mutex and safe to call without holding it.
var (
	mu       sync.RWMutex
	contacts = map[string]Contact{}

	contactSeq atomic.Int64
	noteSeq    atomic.Int64
)

// nextContactID returns a unique contact identifier by atomically
// incrementing the contact sequence counter.
func nextContactID() string {
	return fmt.Sprintf("c%d", contactSeq.Add(1))
}

// nextNoteID returns a unique note identifier by atomically
// incrementing the note sequence counter.
func nextNoteID() string {
	return fmt.Sprintf("n%d", noteSeq.Add(1))
}

// init seeds the store with demo contacts and notes so the
// application has data to display immediately on first launch.
func init() {
	alice := Contact{
		ID:    nextContactID(),
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Phone: "+61 400 111 222",
		Notes: []Note{
			{ID: nextNoteID(), Content: "Met at the Go meetup last Friday.", Created: time.Now().Add(-48 * time.Hour)},
			{ID: nextNoteID(), Content: "Interested in contributing to the fluent project.", Created: time.Now().Add(-24 * time.Hour)},
		},
	}
	bob := Contact{
		ID:    nextContactID(),
		Name:  "Bob Smith",
		Email: "bob@example.com",
		Phone: "+61 400 333 444",
		Notes: []Note{
			{ID: nextNoteID(), Content: "Prefers email over phone calls.", Created: time.Now().Add(-72 * time.Hour)},
		},
	}
	carol := Contact{
		ID:    nextContactID(),
		Name:  "Carol Williams",
		Email: "carol@example.com",
		Phone: "+61 400 555 666",
	}

	contacts[alice.ID] = alice
	contacts[bob.ID] = bob
	contacts[carol.ID] = carol
}

// All returns every contact sorted alphabetically by name.
func All() []Contact {
	mu.RLock()
	defer mu.RUnlock()

	out := make([]Contact, 0, len(contacts))
	for _, c := range contacts {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

// Get returns the contact with the given ID and true, or an empty
// Contact and false if the ID does not exist.
func Get(id string) (Contact, bool) {
	mu.RLock()
	defer mu.RUnlock()

	c, ok := contacts[id]
	return c, ok
}

// Create adds a new contact and returns it.
func Create(name, email, phone string) Contact {
	mu.Lock()
	defer mu.Unlock()

	c := Contact{
		ID:    nextContactID(),
		Name:  name,
		Email: email,
		Phone: phone,
	}
	contacts[c.ID] = c
	return c
}

// Update modifies an existing contact's details. It returns false if
// the contact does not exist.
func Update(id, name, email, phone string) bool {
	mu.Lock()
	defer mu.Unlock()

	c, ok := contacts[id]
	if !ok {
		return false
	}
	c.Name = name
	c.Email = email
	c.Phone = phone
	contacts[id] = c
	return true
}

// Delete removes a contact by ID. It returns false if the contact
// does not exist.
func Delete(id string) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := contacts[id]; !ok {
		return false
	}
	delete(contacts, id)
	return true
}

// AddNote appends a note to the specified contact. It returns false
// if the contact does not exist.
func AddNote(contactID, content string) bool {
	mu.Lock()
	defer mu.Unlock()

	c, ok := contacts[contactID]
	if !ok {
		return false
	}
	c.Notes = append(c.Notes, Note{
		ID:      nextNoteID(),
		Content: content,
		Created: time.Now(),
	})
	contacts[contactID] = c
	return true
}

// DeleteNote removes a specific note from a contact. It returns
// false if either the contact or the note does not exist.
func DeleteNote(contactID, noteID string) bool {
	mu.Lock()
	defer mu.Unlock()

	c, ok := contacts[contactID]
	if !ok {
		return false
	}
	for i, n := range c.Notes {
		if n.ID == noteID {
			c.Notes = append(c.Notes[:i], c.Notes[i+1:]...)
			contacts[contactID] = c
			return true
		}
	}
	return false
}
