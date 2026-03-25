# Tether App - Kanban Board

A collaborative kanban board built with a single [Tether](https://github.com/jpl-au/tether) handler. This is a **real-world example** - where the [Feature Explorer](../tether/) splits every concept into its own handler for teaching purposes, this application shows how you would actually build something with tether: one handler, one state struct, shared board state across all connected browsers.

## What it demonstrates

- **One handler** - a single `tether.Handler` serving the entire application
- **WebSocket + SSE failover** - `mode.Both` tries WebSocket first, falls back to SSE+POST automatically
- **Identity** - a landing page asks for your name before showing the board
- **Drag and drop** - cards are draggable between and within columns via `bind.Draggable` and `bind.Sortable`
- **Within-column reordering** - drag cards up/down within a column to prioritise; drop index calculated from cursor position
- **Cross-session sync** - drag a card in one tab, every other tab updates instantly via `Group.Broadcast`
- **Toast notifications** - when someone moves, creates, updates, or deletes a card, all other sessions get a toast
- **Presence** - "James is viewing this" appears on cards when someone has them open, tracked via `tether.Presence[ViewInfo]`
- **Typing indicator** - "James is editing..." appears in green when someone is typing in a card detail, using `bind.OnInput` with debounce
- **Card ownership** - cards show who created them; the first user to join claims the seed data
- **Activity log** - every card tracks its history: created, moved, updated, reordered
- **Relative timestamps** - cards show "2 hours ago", "just now" etc.
- **URL routing** - `/` for the board, `/card/<id>` for detail, `/new` for creating; browser back/forward works via `OnNavigate`
- **Hotkeys** - Escape closes the card detail view via `bind.Hotkey`
- **Overflow menu** - three-dot menu with delete action, toggled via `bind.ToggleClass`
- **SPA-style region swapping** - clicking a card replaces the board with a detail view, no page reloads, no dialogs
- **Reactive online count** - `Group.Count()` + `WatchValue` keeps the badge accurate across all sessions
- **Shimmer animation** - rotating conic-gradient border on card hover

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
  handler.go         <- tether.Handler constructor (mode.Both, Group, Count)
  state.go           <- per-session State (name, view mode, selected card)
  handle.go          <- event dispatch (name, create, move, update, delete)
  view.go            <- render function (landing, board, or detail from store)
  viewers.go         <- tether.Presence[ViewInfo] wrapper for card presence
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
    card/             <- draggable card with presence indicators
    detail/           <- card edit form with back, save, overflow menu
static/
  app.css            <- flat dark theme, shimmer hover, styled scrollbars
```

## How the sync works

The board data lives in `store.Board`, a shared, mutex-protected struct. When any session creates, moves, updates, or deletes a card:

1. The handler mutates the store
2. `Group.Broadcast` triggers a re-render on every connected session
3. Each session's render function reads from the store, so all browsers see the same board
4. `Group.BroadcastOthers` sends a toast to everyone except the actor

Per-session state tracks only the user's name, what they are looking at (board view vs card detail), and their session ID for presence exclusion. The board itself is never duplicated into session state.

## Extension script loading

Tether auto-includes extension scripts (like `tether-drag-and-drop.js`) when their marker attribute appears in the rendered HTML. The runtime also lazy-loads extensions after morphs, so if the marker first appears after a page transition (e.g. landing page to board), the script loads dynamically without a page reload.

As a belt-and-braces measure, the landing page includes a hidden marker element to ensure the DnD script loads on the initial page:

```go
bind.Apply(div.New().Class("sr-only"), bind.Draggable())
```
