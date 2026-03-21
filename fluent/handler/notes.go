package handler

import (
	"net/http"

	"github.com/jpl-au/fluent-examples/fluent/store"
)

// CreateNote adds a note to a contact and redirects to the detail
// page.
func CreateNote(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	store.AddNote(contactID, r.FormValue("content"))
	http.Redirect(w, r, "/contacts/"+contactID, http.StatusSeeOther)
}

// DeleteNote removes a note from a contact and redirects to the
// detail page.
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("id")
	noteID := r.PathValue("noteID")
	store.DeleteNote(contactID, noteID)
	http.Redirect(w, r, "/contacts/"+contactID, http.StatusSeeOther)
}
