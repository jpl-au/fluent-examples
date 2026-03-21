# Fluent Example - Contact Manager

A server-rendered contact manager demonstrating [Fluent](https://github.com/jpl-au/fluent) with plain HTTP - no HTMX, no WebSocket, no framework. Pure Go, pure HTML.

## What it demonstrates

- **Fluent API** - HTML5 elements built with method chaining, no templates
- **Components** - reusable simple and composite components (buttons, fields, cards, lists)
- **Routing** - [chain](https://github.com/jpl-au/chain) middleware router with Go 1.22 patterns
- **PRG pattern** - POST handlers redirect after mutation to prevent duplicate submissions
- **Project structure** - clean separation: `main.go` → `server/` → `routes/` → `handler/`

## Run

```bash
go run .
# Visit http://localhost:8080
```

## Structure

```
main.go              ← entry point
server/server.go     ← chain.Mux, middleware, static assets
routes/routes.go     ← all route registration
handler/
  contacts.go        ← contact CRUD handlers
  notes.go           ← note handlers
store/store.go       ← in-memory storage (seeded with example data)
layout/layout.go     ← HTML shell (head, body, header)
components/
  simple/            ← button, field, text
  composite/         ← page, card, row, contactlist, notelist
static/              ← CSS, fonts
```

## Pages

| Path | Method | Description |
|------|--------|-------------|
| `/` | GET | Contact list |
| `/contacts/new` | GET | Add contact form |
| `/contacts` | POST | Create contact |
| `/contacts/{id}` | GET | Contact detail + notes |
| `/contacts/{id}/edit` | GET | Edit contact form |
| `/contacts/{id}` | POST | Update contact |
| `/contacts/{id}/delete` | POST | Delete contact |
| `/contacts/{id}/notes` | POST | Add note |
| `/contacts/{id}/notes/{noteID}/delete` | POST | Delete note |
