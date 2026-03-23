# Fluent-HTMX Example - Contact Manager

The same contact manager as the [Fluent example](../fluent/), refactored to use [Fluent-HTMX](https://github.com/jpl-au/fluent-htmx) for partial page updates. Navigation and form submissions swap the content area via HTMX instead of triggering full page reloads.

This is a **teaching example** showing both HTMX handler patterns so you
can see each approach clearly. A real application would typically settle on
one pattern rather than mixing both.

## What's different from the Fluent version

- **Navigation** uses `hx-get` + `hx-target="#content"` + `hx-push-url` instead of plain `<a>` links
- **Forms** return HTML partials instead of redirecting - HTMX swaps the content in place
- **GET handlers** check `htmx.HxRequest(r)` to decide between full page and partial
- **POST handlers** use `htmx.Handle(r, func(){})` to return updated partials with URL push

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

Same structure as the Fluent example - see [../fluent/README.md](../fluent/README.md).
