# Fluent-HTMX Example - Contact Manager (htmx 4)

The same contact manager as the [Fluent example](../../fluent/), refactored to use the htmx 4 binding of [Fluent-HTMX](https://github.com/jpl-au/fluent-htmx) for partial page updates. Navigation and form submissions swap the content area via HTMX instead of triggering full page reloads. The [htmx2 example](../htmx2/) is the same application on htmx 2.

This is a **teaching example** showing both HTMX handler patterns so you
can see each approach clearly. A real application would typically settle on
one pattern rather than mixing both.

## What's different from the Fluent version

- **Navigation** uses `hx-get` + `hx-target="#content"` + `hx-push-url` instead of plain `<a>` links
- **Forms** return HTML partials instead of redirecting - HTMX swaps the content in place
- **GET handlers** check `htmx.HxRequest(r)` to decide between full page and partial
- **POST handlers** use `htmx.Handle(r, func(){})` to return updated partials with URL push

## What's different from the htmx2 example

- **One script.** The page loads `htmax.js`, the htmx 4 bundle that carries the SSE and WebSocket extensions, so the live log pages load no extension scripts and set no `hx-ext`.
- **Header updates** travel as `<hx-partial>` blocks built with `htmx.Partial`, each routed to its own target, instead of `hx-swap-oob` elements.
- **SSE** uses `hx-sse:connect` on the log container. The server sends unnamed events with `sse.Swap`, and htmx swaps each one into the container with its own `hx-swap`.
- **WebSocket** uses `hx-ws:connect` on the log container. Plain HTML messages are swapped into the container; the first message is a JSON frame built with `htmx.WSResponse` that targets the status span instead.
- **Connection status** is bound on the connecting element with `HxOn` and the constants in the `event` package, rather than a separate script listening on the body.

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

    if htmx.HxRequest(r) {
        layout.Partial(w, content)
        return
    }
    layout.Page(w, "Contacts", headerActions, content)
}
```

### Pattern 2: `htmx.Handle(r, func(){})` closure (POST handlers)

```go
func CreateContact(w http.ResponseWriter, r *http.Request) {
    store.Create(r.FormValue("name"), ...)

    if htmx.Handle(r, func() {
        htmx.HxPushURL(w, "/")
        layout.Partial(w, contactlist.New(store.All()))
    }) {
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

## Structure

Same structure as the Fluent example - see [../../fluent/README.md](../../fluent/README.md).
