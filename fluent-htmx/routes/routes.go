// Package routes registers all URL patterns for the contact manager.
package routes

import (
	"github.com/jpl-au/chain"

	"github.com/jpl-au/fluent-examples/fluent-htmx/handler"
)

// Register adds all contact and note routes to the given mux.
func Register(mux *chain.Mux) {
	// Contact list (also the home page).
	mux.HandleFunc("GET /{$}", handler.ListContacts)

	// New contact form and submission.
	mux.HandleFunc("GET /contacts/new", handler.NewContactForm)
	mux.HandleFunc("POST /contacts", handler.CreateContact)

	// Contact detail, edit, update, and delete.
	mux.HandleFunc("GET /contacts/{id}", handler.ShowContact)
	mux.HandleFunc("GET /contacts/{id}/edit", handler.EditContactForm)
	mux.HandleFunc("POST /contacts/{id}", handler.UpdateContact)
	mux.HandleFunc("POST /contacts/{id}/delete", handler.DeleteContact)

	// Notes.
	mux.HandleFunc("POST /contacts/{id}/notes", handler.CreateNote)
	mux.HandleFunc("POST /contacts/{id}/notes/{noteID}/delete", handler.DeleteNote)

	// Live log demos.
	mux.HandleFunc("GET /ws", handler.WSPage)
	mux.HandleFunc("GET /ws/feed", handler.WSFeed)
	mux.HandleFunc("GET /sse", handler.SSEPage)
	mux.HandleFunc("GET /sse/feed", handler.SSEFeed)
}
