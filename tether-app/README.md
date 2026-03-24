# Tether App - Kanban Board

A collaborative kanban board built with a single [Tether](https://github.com/jpl-au/tether) handler. This is a **real-world example** - where the [Feature Explorer](../tether/) splits every concept into its own handler for teaching purposes, this application shows how you would actually build something with tether: one handler, one state struct, shared board state across all connected browsers.

## What it demonstrates

- **One handler** - a single `tether.Handler` serving the entire application
- **WebSocket + SSE failover** - `mode.Both` tries WebSocket first, falls back to SSE+POST automatically
- **Identity** - a landing page asks for your name before showing the board
- **Drag and drop** - cards are draggable between columns via `bind.Draggable` and `bind.DropTarget`
- **Cross-session sync** - drag a card in one tab, every other tab updates instantly via `Group.Broadcast`
- **SPA-style region swapping** - clicking a card replaces the board with a detail view, no page reloads, no dialogs
- **Shared state** - the board lives in a thread-safe store; each session holds only view state (name, current view, selected card)
- **Reactive online count** - `tether.Value` + `WatchValue` push the badge update to all sessions

## Run

```bash
go run .
# Visit http://localhost:8080
```

Open two browser tabs. Enter a name in each, then drag a card in one - it moves in the other immediately.

## Structure

```
main.go              <- entry point, embedded assets, ListenAndServe
store/
  store.go           <- Board, Column, Card types, thread-safe operations
handler/
  doc.go             <- package documentation
  handler.go         <- tether.Handler constructor (mode.Both, Group, Value)
  state.go           <- per-session State (name, view mode, selected card)
  handle.go          <- event dispatch (name, create, move, update, delete)
  view.go            <- render function (landing, board, or detail from store)
layout/
  layout.go          <- page shell (header, user name, online badge, content)
components/
  simple/
    button/           <- Primary, Secondary, Danger, Submit
    badge/            <- Todo, Progress, Done column indicators
    field/            <- Text, TextValue, Area, Label, Inline
  composite/
    board/            <- three-column grid
    column/           <- swimlane with header and card list
    card/             <- draggable card with title and description
    detail/           <- card edit form with back, save, delete
static/
  app.css            <- flat dark theme, shimmer hover animation
```

## How the sync works

The board data lives in `store.Board`, a shared, mutex-protected struct. When any session creates, moves, updates, or deletes a card:

1. The handler mutates the store
2. `Group.Broadcast` triggers a re-render on every connected session
3. Each session's render function reads from the store, so all browsers see the same board

Per-session state tracks only the user's name and what they are looking at (board view vs card detail). The board itself is never duplicated into session state.

## Extension script loading caveat

Tether auto-includes extension scripts (like `tether-drag-and-drop.js`) only when their marker attribute appears in the **initial page render**. If the first view rendered does not contain any draggable elements (e.g. this app shows a landing page first), the DnD script would never load.

The workaround is a hidden marker element on the initial page:

```go
bind.Apply(div.New().Class("sr-only"), bind.Draggable())
```

This ensures the extension script loads upfront so drag-and-drop works when the board view appears later. See `handler/view.go` for the implementation.
