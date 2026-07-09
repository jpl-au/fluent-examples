# Fluent Example - Contact Manager

A server-rendered contact manager demonstrating [Fluent](https://github.com/jpl-au/fluent) with plain HTTP - no HTMX, no framework. Pure Go, pure HTML. A pair of live-log demos (`/ws`, `/sse`) also show raw WebSocket and Server-Sent Events, again using only the standard library.

This is a **teaching example** - a deliberately simple application designed
to show how Fluent's API works. The contact manager is a toy domain chosen
because it exercises routing, forms, and CRUD without domain complexity
getting in the way.

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
  ws.go              ← WebSocket live-log demo page + feed
  sse.go             ← Server-Sent Events live-log demo page + feed
  live.go            ← shared live-log page scaffold
store/store.go       ← in-memory storage (seeded with example data)
layout/layout.go     ← HTML shell (head, body, header)
components/
  simple/            ← button, field, text
  composite/         ← page, card, row, contactlist, notelist, footer, menu
static/              ← CSS, fonts, ws.js, sse.js
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
| `/ws` | GET | WebSocket live-log demo page |
| `/ws/feed` | GET | WebSocket feed endpoint (upgrades the connection) |
| `/sse` | GET | Server-Sent Events live-log demo page |
| `/sse/feed` | GET | SSE feed endpoint (streams events) |
