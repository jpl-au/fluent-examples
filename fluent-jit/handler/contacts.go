// Package handler provides HTTP request handlers for the contact
// manager. Each handler follows the PRG (Post-Redirect-Get) pattern:
// GET handlers render pages, POST handlers mutate state and redirect.
//
// This file demonstrates all three JIT strategies across different
// handlers, using both the Global API (string-keyed registry) and the
// Instance API (fine-grained control) to show the full range of
// fluent-jit usage.
package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	jit "github.com/jpl-au/fluent-jit"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/form"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/card"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/contactlist"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/menu"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/notelist"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/composite/row"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/simple/button"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/simple/field"
	"github.com/jpl-au/fluent-examples/fluent-jit/components/simple/text"
	"github.com/jpl-au/fluent-examples/fluent-jit/layout"
	"github.com/jpl-au/fluent-examples/fluent-jit/store"
)

// showContactCompiler is a package-level Compiler instance demonstrating
// the Instance API. It gives fine-grained control over the compilation
// lifecycle - useful when you want to isolate a specific template's
// execution plan from the global registry, or when you need access to
// methods like Validate() for testing.
var showContactCompiler = jit.NewCompiler()

// newContactFlattener holds the pre-rendered "new contact" form.
// Flatten only works when every node in the tree uses Static() content.
// In practice, components like card, field, and button use Text()
// internally (for HTML escaping), so NewFlattener returns an error here.
// This demonstrates the Instance API's explicit error handling - the
// Global API (jit.Flatten) would silently fall back instead.
var newContactFlattener *jit.Flattener

func init() {
	content := newContactContent()
	var err error
	newContactFlattener, err = jit.NewFlattener(content)
	if err != nil {
		// Expected: the form components use Text() internally for labels and
		// button text, which marks them as dynamic even though the values
		// never change. This is the tradeoff - Flatten requires Static()
		// throughout, while Text() provides HTML escaping safety.
		slog.Debug("new contact form contains dynamic content, using standard render",
			"error", err)
	}
}

// ListContacts renders the contact list page.
//
// Uses the Global Compile API because the list structure is always the
// same (a div containing contact-item children) - only the text content
// (names, emails, note counts) changes between renders. Compile freezes
// the static HTML skeleton after the first render and re-evaluates only
// the dynamic Text() nodes on subsequent calls, avoiding repeated
// serialisation of all the surrounding divs, classes, and anchors.
func ListContacts(w http.ResponseWriter, r *http.Request) {
	contacts := store.All()
	content := contactlist.New(contacts)
	doc := layout.Document("Contacts", button.Primary("+ New", "/contacts/new"), content)
	jit.Compile("contacts-list", doc, w)
}

// NewContactForm renders the "add contact" form.
//
// Demonstrates the Instance Flatten API - jit.NewFlattener() pre-renders
// a fully static tree to raw bytes once, then serves those bytes on every
// subsequent call. If the tree contains any dynamic content (Text/Textf
// nodes), NewFlattener returns an error and the handler falls back to
// standard rendering. This shows the explicit error path that the Instance
// API provides, compared to the Global API which falls back silently.
func NewContactForm(w http.ResponseWriter, r *http.Request) {
	// If flattening succeeded at init, serve the pre-rendered bytes
	// directly - zero tree traversal, zero serialisation.
	if newContactFlattener != nil {
		newContactFlattener.Render(w)
		return
	}
	// Fallback: render the tree normally. This path is taken when the
	// form contains dynamic content (Text() nodes in components).
	newContactContent().Render(w)
}

// newContactContent builds the static "new contact" form tree. This is
// extracted so both init (for flattening) and the fallback path can
// share the same structure.
func newContactContent() node.Node {
	content := card.New("New Contact",
		form.Post("/contacts",
			field.Group(field.Label("name", "Name"), field.Text("name", "Full name")),
			field.Group(field.Label("email", "Email"), field.Text("email", "email@example.com")),
			field.Group(field.Label("phone", "Phone"), field.Text("phone", "+61 400 000 000")),
			row.New(button.Submit("Create Contact")),
		),
	)
	return layout.Shell("New Contact", button.Back("/"), content)
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
//
// Uses the Instance Compile API (showContactCompiler) because the
// detail page has a fixed structure - a detail card with label/value
// rows and a notes card with a form - but the text content changes per
// contact. The Instance API gives fine-grained control: the compiler
// lives at package scope, its execution plan is isolated from the
// global registry, and Validate() can be used in tests to assert
// structural compatibility.
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

	// Instance Compile API: the compiler builds an execution plan on
	// its first call and reuses it for every subsequent render. Static
	// HTML (divs, classes, labels) is frozen; dynamic Text() nodes are
	// re-evaluated from the fresh tree each time.
	doc := layout.Document(c.Name, button.Back("/"), details, notes)
	showContactCompiler.Render(doc, w)
}

// EditContactForm renders the "edit contact" form pre-filled with
// the current values.
//
// Uses the Global Tune API because the form content varies in size
// depending on the pre-filled values (short vs long names, emails).
// Tune adapts the buffer size over time without analysing tree
// structure, so it handles this natural variation gracefully.
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
//
// Demonstrates the Global Flatten API - jit.Flatten("not-found", ...)
// attempts to pre-render the tree to raw bytes on first call and cache
// them for subsequent calls. If the tree contains dynamic content (as
// it does here - components use Text() internally), the global API
// falls back to standard rendering silently. This is the key difference
// from the Instance API: no error to handle, just graceful degradation.
// For truly static content where every node uses Static(), Flatten
// caches the output and serves it as a direct byte copy.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	content := card.New("",
		text.Hint("Contact not found."),
		row.New(button.Back("/")),
	)
	shell := layout.Shell("Not Found", nil, content)
	jit.Flatten("not-found", shell, w)
}
