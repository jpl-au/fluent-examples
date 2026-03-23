# Tether Feature Explorer

Interactive example application for the [Tether](https://github.com/jpl-au/tether) framework. Each demo is a self-contained package that exercises a specific Tether feature with real, runnable code.

```bash
go run .             # defaults to :8080
PORT=3000 go run .   # or set a custom port
```

Open multiple browser tabs to see real-time features in action.

## Why so many handlers?

This application creates a separate `tether.Handler` (or `tether.Stateless`)
for every demo. That is a **deliberate teaching choice**, not an architectural
recommendation.

Each handler isolates a single feature - signals, uploads, groups, chat,
etc. - so you can open one package, read one handler, and understand that
feature without any other demo's state or events getting in the way. The
wiring in `app/app.go` reflects this: it stitches together many independent
demos into a single navigable site.

**A real application would look different.** Most production apps need only
one `tether.Handler` (or a small handful if distinct sections have genuinely
different lifecycle requirements). You would define your state, events, and
views in that single handler rather than splitting every concern into its
own package. Think of this explorer as a reference catalogue - pick the
feature you need, study it here, then integrate the pattern into your own
application.

## Demos

### HTTP (stateless)

| Demo | Package | What it shows |
|------|---------|---------------|
| Rendering | `site/rendering` | Dynamic keys, dynamic lists, error boundaries (`tether.Catch`), components, nested routing (`RouteTyped`) |
| Events | `site/events` | Every event binding: click, submit, input, change, keydown, focus, blur, viewport, throttle, debounce, confirm, custom events, typed extraction (`ev.Int`, `ev.Bool`, `ev.Bind`) |
| Errors | `site/errors` | `tether.Catch` error boundaries recovering from panics in child components |
| Morph | `site/morph` | Full-page morph fallback when no Dynamic keys are present |
| Navigation | `site/navigation` | `bind.Link`, `OnNavigate`, query parameter extraction (`Params`), `ReplaceURL`, server-driven `Navigate` |

### WebSocket + SSE (stateful)

| Demo | Package | What it shows |
|------|---------|---------------|
| Live Updates | `site/live` | Uptime ticker via `sess.Go`, activity feeds via `Bus`, online count via `Value`, Group operations, `SetTitle`, `State()`, `Close()` |
| Signals | `site/signals` | Every signal binding: `BindText`, `BindShow`, `BindHide`, `BindClass`, `BindAttr`, `BindValue`, `SetSignal`, `ToggleSignal`, `Optimistic`, `OptimisticToggle`, `Cloak`, `Permanent`, `Hook`, `Transition`, `FocusTrap` |
| Chat | `site/chat` | Real-time cross-session chat using `tether.Component`, `Mounter`, `Bus`, and `WatchBus` |
| Broadcasting | `site/broadcasting` | `Bus.Emit`, `Bus.Publish`, shared counter with `tether.Value`, `WatchBus`, `WatchValue`, async subscribers |
| Groups | `site/groups` | Room membership with `Group.Add`/`Remove`, `Broadcast`, `BroadcastOthers`, `OnJoin`/`OnLeave` callbacks |
| Value Store | `site/valuestore` | `tether.Value` for shared observable state, `Store`, `Update`, `WatchValue`, shared vs local state |
| Components | `site/components` | `tether.Component`, `StatefulConfig.Components`, `Mounter`, `Event.Target`, multiple independent instances |
| Middleware | `site/mw` | `tether.Middleware` chain: timing, guard, counting, and ordered wrappers showing onion execution |
| Notifications | `site/notifications` | Server-push side effects: `Toast`, `Flash`, `Announce`, `Signal` |
| Diagnostics | `site/diagnostics` | `Handler.Diagnostics` bus, live event feed via `WatchBus`, triggerable panics, diagnostic kind reference |
| Configuration | `site/configuration` | `Timeouts`, `Limits`, `Security`, compression, `SessionStore`, `DiffStore`, `OnRestore` |
| Uploads | `site/uploads` | File uploads via `bind.Upload` with real-time feedback through `sess.Update` |
| Filtered Uploads | `site/uploads/filtered` | `UploadConfig.Accept` MIME-type filtering |
| Freeze | `site/freeze` | `FreezeWithConnect`, `SessionStore` persistence, zero-memory disconnected sessions, state restoration |
| Real-time Monitor | `site/realtime` | Live Go runtime metrics (CPU, heap, goroutines) pushed every second via `sess.Go`, rendered as go-echarts line charts |

### Service Worker

| Demo | Package | What it shows |
|------|---------|---------------|
| Overview | `site/sw/page` | Service worker lifecycle, caching strategies, push notifications |

## Project structure

```
app/           Application wiring - creates all handlers, mounts routes
component/     Reusable tether.Component implementations (counter, shoutbox)
components/    UI rendering helpers (buttons, panels, cards, layouts)
layout/        Shared page shell and navigation
middleware/    Example middleware implementations
playwright/    End-to-end browser tests
site/          One package per demo (listed above)
static/        CSS and client-side JS
store/         File-based SessionStore and DiffStore implementations
```

Each `site/` package follows a consistent pattern:
- `doc.go` - what the demo exercises
- `handler.go` - state type, handler constructor, event handling
- `view.go` - render function using UI components

## Where to start

1. **New to Tether?** Start with `site/rendering` and `site/events` for the basics, then `site/live` for real-time features.
2. **Building a real app?** Look at `site/chat` for a complete component-based feature, `site/configuration` for production settings, and `site/diagnostics` for observability.
3. **Testing?** See `site/tethertest` for unit testing patterns without a browser.

## Playwright tests

End-to-end browser tests live in `playwright/`. They use the system-installed Google Chrome via the Playwright Go driver - no bundled Chromium download is required.

### Prerequisites

1. **Google Chrome** must be installed at the default system path.
2. **Playwright Go driver** - installed automatically as a Go module dependency.

### Running

```bash
go test ./playwright/...                      # HTTP/1.1 (default)
TETHER_PROTO=HTTP2 go test ./playwright/...   # HTTP/2 over TLS
```

### How it works

The tests start a full `httptest` server per test, open a headless Chrome page via `Channel: "chrome"` (system Chrome, not a Playwright-managed download), and interact with the application using Playwright's Go API. The `IgnoreHttpsErrors` context option handles the self-signed certificate from `httptest.NewTLSServer`.

No Playwright browser binaries are downloaded or cached. If Chrome is not installed, the tests skip with a descriptive message.
