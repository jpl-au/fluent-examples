# Fluent-HTMX Example - Contact Manager

The same contact manager as the [Fluent example](../../fluent/), refactored to use [Fluent-HTMX](https://github.com/jpl-au/fluent-htmx) for partial page updates. Navigation and form submissions swap the content area via HTMX instead of triggering full page reloads.

This is a **teaching example** showing both HTMX handler patterns so you
can see each approach clearly. A real application would typically settle on
one pattern rather than mixing both.

## What's different from the Fluent version

- **Navigation** adds `hx-get` + `hx-target="#content"` + `hx-push-url` to links that keep their `href`, so they still work without JavaScript
- **Forms** return HTML partials instead of redirecting - HTMX swaps the content in place; without JavaScript they post and redirect as before
- **GET handlers** check `htmx.HxRequest(r)` to decide between full page and partial
- **POST handlers** use `htmx.Handle(r, func(){})` to return updated partials with URL push
- **The header and tab title** update with each partial, which carries a `<title>` and `<hx-partial>` blocks built with `htmx.Partial`, each routed to its own target (htmx 2.0.11 or later)
- **Deletes** ask first with `HxConfirm`, and forms show a "Saving..." indicator while their request runs

## htmx config

The page head carries an htmx config built with `htmx.Config()`, for two htmx 2 defaults that don't suit a site that answers `HX-Request` with a partial:

- `HistoryRestoreAsHxRequest(false)`: on a history cache miss htmx 2 sends `HX-Request` by default, so the handlers would return a partial and htmx would swap it in as the whole page. Without the header they return the full page.
- `ResponseHandling`: htmx 2 swaps no 4xx response by default, so "Contact not found" would never show. A rule for 404 swaps it like a normal response; htmx replaces the whole list, so its default rules are restated around it.

htmx 4 needs neither: it swaps 4xx responses by default, and a history restore asks for the full page with `HX-Request-Type: full` (see the [htmx4 example](../htmx4/)).

## Live pages

The footer links to two pages that stream log entries from the server:

- **WebSocket** (`/ws`) loads `htmx-ext-ws`. The server pushes each entry as an `hx-swap-oob` fragment, and a line sent with the `ws-send` form comes back on the same connection as an entry.
- **SSE** (`/sse`) loads `htmx-ext-sse`. The server sends named events, which `sse-swap` routes to the status span and the log.

Both bind the connection status on the connecting element with `HxOn` and the constants in the `event` package.

## Run

```bash
go run .
# Visit http://localhost:8080
```

## Handler Patterns

Two patterns are demonstrated:

### Pattern 1: `htmx.HxRequest(r)` if-check (GET handlers)

```go
func ListContacts(w http.ResponseWriter, r *http.Request) {
    content := contactlist.New(store.All())
    actions := button.Primary("+ New", "/contacts/new")

    if htmx.HxRequest(r) {
        layout.Partial(w, "Contacts", actions, content)
        return
    }
    layout.Page(w, "Contacts", actions, content)
}
```

### Pattern 2: `htmx.Handle(r, func(){})` closure (POST handlers)

```go
func CreateContact(w http.ResponseWriter, r *http.Request) {
    store.Create(r.FormValue("name"), ...)

    if htmx.Handle(r, func() {
        htmx.HxPushURL(w, "/")
        layout.Partial(w, "Contacts", button.Primary("+ New", "/contacts/new"), contactlist.New(store.All()))
    }) {
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

## Structure

Same structure as the Fluent example - see [../../fluent/README.md](../../fluent/README.md).
