package handler

import (
	"net/http"

	htmx "github.com/jpl-au/fluent-htmx/htmx2"

	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/components/simple/button"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/layout"
	"github.com/jpl-au/fluent-examples/fluent-htmx/htmx2/store"
)

// CreateNote adds a note and returns the updated detail partial.
// Uses htmx.Handle to avoid a full page redirect.
func CreateNote(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	store.AddNote(contactID, r.FormValue("content"))

	if htmx.Handle(r, func() {
		c, ok := store.Get(contactID)
		if !ok {
			notFound(w, r)
			return
		}
		details, notes := contactDetail(c)
		layout.Partial(w, c.Name, button.Back("/"), details, notes)
	}) {
		return
	}
	http.Redirect(w, r, "/contacts/"+contactID, http.StatusSeeOther)
}

// DeleteNote removes a note and returns the updated detail partial.
// Uses htmx.Handle to swap the content in place.
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("id")
	noteID := r.PathValue("noteID")
	store.DeleteNote(contactID, noteID)

	if htmx.Handle(r, func() {
		c, ok := store.Get(contactID)
		if !ok {
			notFound(w, r)
			return
		}
		details, notes := contactDetail(c)
		layout.Partial(w, c.Name, button.Back("/"), details, notes)
	}) {
		return
	}
	http.Redirect(w, r, "/contacts/"+contactID, http.StatusSeeOther)
}
