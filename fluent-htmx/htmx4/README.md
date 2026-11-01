# Fluent-HTMX Example - Contact Manager (htmx 4)

The same contact manager as the [Fluent example](../../fluent/), refactored to use the htmx 4 binding of [Fluent-HTMX](https://github.com/jpl-au/fluent-htmx) for partial page updates. Navigation and form submissions swap the content area via HTMX instead of triggering full page reloads. The [htmx2 example](../htmx2/) is the same application on htmx 2.

This is a **teaching example** showing both HTMX handler patterns so you
can see each approach clearly. A real application would typically settle on
one pattern rather than mixing both.

## What's different from the Fluent version

- **Navigation** adds `hx-get` + `hx-push-url` to links that keep their `href`, so they still work without JavaScript
- **Forms** return HTML partials instead of redirecting - HTMX swaps the content in place; without JavaScript they post and redirect as before
- **GET handlers** check `htmx.HxRequestType(r) == "partial"` to decide between full page and partial
- **POST handlers** use `htmx.Handle(r, func(){})` to return updated partials with URL push
- **The header and tab title** update with each partial, which carries a `<title>` and `<hx-partial>` blocks built with `htmx.Partial`, each routed to its own target
- **Deletes** ask first with `HxConfirm`, and forms show a "Saving..." indicator while their request runs

## What's different from the htmx2 example

- **One script.** The page loads `htmax.js`, the htmx 4 bundle that carries the SSE and WebSocket extensions, so the live log pages load no extension scripts and set no `hx-ext`.
- **Inheritance.** The shell sets `hx-target="#content"` and `hx-swap="innerHTML"` once with `htmx.Inherited`, so links, buttons and forms carry only `hx-get` or `hx-post`. The live log containers set their own target.
- **GET handlers** check `htmx.HxRequestType(r) == "partial"`. A history restore asks for the full page with `HX-Request-Type: full`, so it gets one whatever other headers it carries. The htmx2 example gets the same from its htmx config.
- **No config.** htmx 4 swaps 4xx responses by default, so "Contact not found" shows without the response handling rule the htmx2 example needs.
- **SSE** uses `hx-sse:connect` on the log container. The server sends unnamed events with `sse.Swap`, and htmx swaps each one into the container with its own `hx-swap`.
- **WebSocket** uses `hx-ws:connect` on a wrapper around the send form and the log, since the form has to sit inside the connecting element. Plain HTML messages are swapped into the log; the first message is a JSON frame built with `htmx.WSResponse` that targets the status span instead. Lines sent with the `hx-ws:send` form are read with `htmx.ParseWSMessage`.

## Run

```bash
go run .
# Visit http://localhost:8080
```

## Handler Patterns

Two patterns are demonstrated:

### Pattern 1: `htmx.HxRequestType(r)` if-check (GET handlers)

```go
func ListContacts(w http.ResponseWriter, r *http.Request) {
    content := contactlist.New(store.All())
    actions := button.Primary("+ New", "/contacts/new")

    if htmx.HxRequestType(r) == "partial" {
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
