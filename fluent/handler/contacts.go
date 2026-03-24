package handler

import (
	"fmt"
	"net/http"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent/components/composite/contactlist"
	"github.com/jpl-au/fluent-examples/fluent/components/composite/menu"
	"github.com/jpl-au/fluent-examples/fluent/components/composite/notelist"
	"github.com/jpl-au/fluent-examples/fluent/components/composite/row"
	"github.com/jpl-au/fluent-examples/fluent/components/simple/button"
	"github.com/jpl-au/fluent-examples/fluent/components/simple/field"
	"github.com/jpl-au/fluent-examples/fluent/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent/layout"
	"github.com/jpl-au/fluent-examples/fluent/store"
)

// ListContacts renders the contact list page.
func ListContacts(w http.ResponseWriter, r *http.Request) {
	contacts := store.All()
	content := contactlist.New(contacts)
	layout.Page(w, "Contacts", button.Primary("+ New", "/contacts/new"), content)
}

// NewContactForm renders the "add contact" form.
func NewContactForm(w http.ResponseWriter, r *http.Request) {
	content := card.New("New Contact",
		form.Post("/contacts",
			field.Group(field.Label("name", "Name"), field.Text("name", "Full name")),
			field.Group(field.Label("email", "Email"), field.Text("email", "email@example.com")),
			field.Group(field.Label("phone", "Phone"), field.Text("phone", "+61 400 000 000")),
			row.New(button.Submit("Create Contact")),
		),
	)
	layout.Page(w, "New Contact", button.Back("/"), content)
}

// CreateContact processes the new-contact form submission and
// redirects to the list.
func CreateContact(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	store.Create(r.FormValue("name"), r.FormValue("email"), r.FormValue("phone"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ShowContact renders the contact detail page with notes.
func ShowContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, ok := store.Get(id)
	if !ok {
		notFound(w, r)
		return
	}

	noteLabel := fmt.Sprintf("%d notes", len(c.Notes))
	if len(c.Notes) == 1 {
		noteLabel = "1 note"
	}

	details := card.NewWithAction(c.Name,
		menu.New(
			menu.Link("Edit", "/contacts/"+c.ID+"/edit"),
			menu.FormAction("Delete", "/contacts/"+c.ID+"/delete"),
		),
		div.New(
			detailRow("Email", c.Email),
			detailRow("Phone", c.Phone),
			detailRow("Notes", noteLabel),
		).Class("detail-grid"),
	)

	notes := card.New("Notes",
		notelist.New(c.ID, c.Notes),
		form.Post("/contacts/"+c.ID+"/notes",
			field.Group(field.TextArea("content", "Add a note...")),
			row.New(button.Submit("Add Note")),
		),
	)

	layout.Page(w, c.Name, button.Back("/"), details, notes)
}

// EditContactForm renders the "edit contact" form pre-filled with
// the current values.
func EditContactForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, ok := store.Get(id)
	if !ok {
		notFound(w, r)
		return
	}

	content := card.New("Edit Contact",
		form.Post("/contacts/"+c.ID,
			field.Group(field.Label("name", "Name"), field.TextValue("name", c.Name, "Full name")),
			field.Group(field.Label("email", "Email"), field.TextValue("email", c.Email, "email@example.com")),
			field.Group(field.Label("phone", "Phone"), field.TextValue("phone", c.Phone, "+61 400 000 000")),
			row.New(button.Submit("Save Changes")),
		),
	)
	layout.Page(w, "Edit Contact", button.Back("/contacts/"+c.ID), content)
}

// UpdateContact processes the edit form and redirects to the detail
// page.
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	store.Update(id, r.FormValue("name"), r.FormValue("email"), r.FormValue("phone"))
	http.Redirect(w, r, "/contacts/"+id, http.StatusSeeOther)
}

// DeleteContact deletes a contact and redirects to the list.
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	store.Delete(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// detailRow renders a label–value pair for the contact detail card.
// The label is treated as static (known at definition time); the
// value is escaped via Text() because it contains user input.
func detailRow(label, value string) node.Node {
	return div.New(
		span.Static(label).Class("detail-label"),
		span.Text(value).Class("detail-value"),
	).Class("detail-row")
}

// notFound renders a simple "Contact not found" page.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	content := card.New("",
		text.Hint("Contact not found."),
		row.New(button.Back("/")),
	)
	layout.Page(w, "Not Found", nil, content)
}
