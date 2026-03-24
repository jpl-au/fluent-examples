package handler

import (
	"fmt"
	"net/http"

	htmx "github.com/jpl-au/fluent-htmx"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/contactlist"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/menu"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/notelist"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/composite/row"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/simple/button"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/simple/field"
	"github.com/jpl-au/fluent-examples/fluent-htmx/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent-htmx/layout"
	"github.com/jpl-au/fluent-examples/fluent-htmx/store"
)

// ListContacts renders the contact list page. HTMX requests receive
// just the contact list; standard requests receive the full page.
// Demonstrates the htmx.HxRequest(r) if-check pattern.
func ListContacts(w http.ResponseWriter, r *http.Request) {
	contacts := store.All()
	content := contactlist.New(contacts)
	actions := button.Primary("+ New", "/contacts/new")

	if htmx.HxRequest(r) {
		layout.Partial(w, "Contacts", actions, content)
		return
	}
	layout.Page(w, "Contacts", actions, content)
}

// NewContactForm renders the "add contact" form.
// Demonstrates the htmx.HxRequest(r) if-check pattern.
func NewContactForm(w http.ResponseWriter, r *http.Request) {
	content := card.New("New Contact",
		form.Post("/contacts",
			field.Group(field.Label("name", "Name"), field.Text("name", "Full name")),
			field.Group(field.Label("email", "Email"), field.Text("email", "email@example.com")),
			field.Group(field.Label("phone", "Phone"), field.Text("phone", "+61 400 000 000")),
			row.New(button.Submit("Create Contact")),
		),
	)
	actions := button.Back("/")

	if htmx.HxRequest(r) {
		layout.Partial(w, "New Contact", actions, content)
		return
	}
	layout.Page(w, "New Contact", actions, content)
}

// CreateContact processes the form and returns the updated contact
// list. Demonstrates the htmx.Handle(r, func(){}) closure pattern  -
// the closure renders the partial and sets the URL; the fallback
// uses a standard redirect.
func CreateContact(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	store.Create(r.FormValue("name"), r.FormValue("email"), r.FormValue("phone"))

	if htmx.Handle(r, func() {
		htmx.HxPushURL(w, "/")
		layout.Partial(w, "Contacts", button.Primary("+ New", "/contacts/new"), contactlist.New(store.All()))
	}) {
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ShowContact renders the contact detail page with notes.
// Demonstrates the htmx.HxRequest(r) if-check pattern.
func ShowContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, ok := store.Get(id)
	if !ok {
		notFound(w, r)
		return
	}

	details, notes := contactDetail(c)
	actions := button.Back("/")

	if htmx.HxRequest(r) {
		layout.Partial(w, c.Name, actions, details, notes)
		return
	}
	layout.Page(w, c.Name, actions, details, notes)
}

// EditContactForm renders the edit form.
// Demonstrates the htmx.HxRequest(r) if-check pattern.
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
	actions := button.Back("/contacts/" + c.ID)

	if htmx.HxRequest(r) {
		layout.Partial(w, "Edit Contact", actions, content)
		return
	}
	layout.Page(w, "Edit Contact", actions, content)
}

// UpdateContact processes the edit form. Uses htmx.Handle to return
// the updated detail partial.
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	store.Update(id, r.FormValue("name"), r.FormValue("email"), r.FormValue("phone"))

	if htmx.Handle(r, func() {
		c, ok := store.Get(id)
		if !ok {
			return
		}
		htmx.HxPushURL(w, "/contacts/"+id)
		details, notes := contactDetail(c)
		layout.Partial(w, c.Name, button.Back("/"), details, notes)
	}) {
		return
	}
	http.Redirect(w, r, "/contacts/"+id, http.StatusSeeOther)
}

// DeleteContact removes a contact. Uses htmx.Handle to return the
// updated contact list.
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	store.Delete(id)

	if htmx.Handle(r, func() {
		htmx.HxPushURL(w, "/")
		layout.Partial(w, "Contacts", button.Primary("+ New", "/contacts/new"), contactlist.New(store.All()))
	}) {
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// contactDetail builds the detail and notes cards for a contact.
// Shared between ShowContact, UpdateContact, and note handlers.
func contactDetail(c store.Contact) (node.Node, node.Node) {
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

	return details, notes
}

// detailRow renders a label–value pair. The label is static (never
// changes); the value is dynamic user input escaped via Text().
func detailRow(label, value string) node.Node {
	return div.New(
		span.Static(label).Class("detail-label"),
		span.Text(value).Class("detail-value"),
	).Class("detail-row")
}

// notFound renders a "Contact not found" page.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	content := card.New("",
		text.Hint("Contact not found."),
		row.New(button.Back("/")),
	)

	if htmx.HxRequest(r) {
		layout.Partial(w, "Not Found", nil, content)
		return
	}
	layout.Page(w, "Not Found", nil, content)
}
