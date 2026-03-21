# Styles Audit - tether Feature Explorer

Assessment of CSS/component consistency, missing styles, dead code,
duplication, and UI polish issues ahead of publishing.

---

## 1. Missing CSS Definitions

Classes applied in Go code with **no corresponding CSS rule** in `app.css`.
These elements render with browser defaults only.

| Class | Used in | Impact |
|---|---|---|
| `msg-list` | `broadcasting/view.go:73` | Message container has no layout styling - items stack by default but lack the padding/border treatment of `.activity-feed` |
| `msg-item` | `broadcasting/view.go:71` | Message items have no padding or bottom border |
| `msg-user` | `broadcasting/view.go:69` | User label has no font-weight or colour |
| `msg-text` | `broadcasting/view.go:70` | Message text has no colour |
| `shoutbox` | `component/shoutbox.go:62` | Container wrapper - no styling (children have rules but the wrapper doesn't) |
| `shout-text` | `component/shoutbox.go:103` | Shout body text has no colour - `.shout-user` and `.shout-time` are styled but `.shout-text` was missed |
| `counter-display` | `valuestore/view.go:34,54,67,71` | Counter value renders in base font - likely intended to use the existing `.stat-value` rule (2rem, 700 weight, tabular-nums) |
| `monitor-charts` | `realtime/view.go:35` | Chart grid container has no layout rule - charts stack vertically by default instead of side by side |
| `monitor-chart` | `realtime/view.go:46` | Individual chart wrapper has no sizing (inline `style` attribute provides dimensions, so this is cosmetic) |
| `viewport-list` | `events/view.go:284` | Viewport demo list has no styling |
| `viewport-item` | `events/view.go:277` | Viewport list items have no styling |
| `viewport-sentinel` | `events/view.go:279` | Sentinel element has no styling (intentionally invisible, but a min-height of 1px would be defensive) |

### Progress element class mismatch

`uploads/view.go:58` applies `Class("progress")` to a `<progress>` element,
but CSS defines the progress bar styling on `.upload-progress` (lines
751-773). The pseudo-element selectors (`::-webkit-progress-bar`,
`::-webkit-progress-value`, `::-moz-progress-bar`) never match, so the
**upload progress bar is completely unstyled** - it falls back to the
browser's default `<progress>` rendering.

---

## 2. Visually Indistinguishable Badge Variants

`.badge-green` and `.badge-indigo` are **identical** (lines 347-355):

```css
.badge-green {
  background: var(--secondary);
  color: var(--secondary-fg);
}
.badge-indigo {
  background: var(--secondary);
  color: var(--secondary-fg);
}
```

Both use the same background and foreground - there is no visual
distinction. The indigo variant presumably needs its own colour.

---

## 3. Duplicate `.upload-status` Definition

`.upload-status` is defined **twice** with conflicting values:

- **Line 782**: `font-size: 0.875rem; color: var(--muted-fg);`
- **Line 989**: `font-size: 0.8125rem; font-weight: 500; margin-left: 0.75rem; text-transform: uppercase; letter-spacing: 0.03em;`

The second definition silently overrides the first. One should be removed
or they should be consolidated.

---

## 4. Dead Component Code

### Unused functions in `component/` package

| Function | File | Status |
|---|---|---|
| `component.Card()` | `component/card.go` | Never imported - exact duplicate of `components/composite/card.New()` |
| `component.Badge()` | `component/badge.go` | Never imported - the simple badge component is also never imported (see below) |
| `component.Demo()` | `component/demo.go` | Never called from any view or handler - `components/composite/demo.New()` is used everywhere |

The `component/` package duplicates `Demo`, `Card`, and `Badge` from the
`components/` tree. Only `Counter`, `CounterGroup`, and `Shoutbox`
(which are genuine `tether.Component` implementations with state) are
actually used.

The `Transport` type and constants are also duplicated - both
`component.Transport` and `demo.Transport` exist with identical definitions.

### Unused simple component

`components/simple/badge/badge.go` is **never imported** anywhere. All
badge usage is inline: `span.New().Class("badge badge-green")` in layout
and view files.

---

## 5. Dead CSS (Defined but Never Used)

| Selector | Lines | Notes |
|---|---|---|
| `.table`, `.table th/td/tr:hover` | 588-613 | No tables in the example app |
| `.tabs`, `.tab`, `.tab-panel` | 802-834 | No tab UI components exist |
| `.flash-target` | 788-796 | Flash demo uses ID selector `#flash-target`, not this class |
| `.upload-area` | 745-749 | Upload demos use `row.New()` instead |
| `.upload-progress` | 751-773 | Intended for `<progress>` but code applies `"progress"` class instead (see section 1) |
| `.progress-area` | 775-780 | Never used |
| `.room-badge` | 1005 | Groups page uses `badge badge-green` directly |
| `.room-none` | 1006 | Never used |
| `.chain-log` | 1009 | Never used |
| `.stat-value` | 733-739 | Never used - value store uses `counter-display` (which also has no CSS, see section 1) |
| `.feature-item/title/desc` | 974-976 | Never used |

Intentionally unused (infrastructure):
- `.sr-only` - accessibility utility
- `.tether-fade-enter/leave`, `.tether-slide-enter/leave` - transition classes for `bind.Transition()`

---

## 6. Raw Class Usage in Views (Bypassing Component Layer)

The component system exists but is only partially adopted. Many views
import raw `fluent/html5/*` elements and apply CSS classes directly.

### No component exists for these patterns

| Pattern | Raw class | Occurrences | Files |
|---|---|---|---|
| Internal link-button | `"btn btn-secondary"` on `<a>` + `bind.Link()` | ~15 | `navigation/view.go` |
| Submit button | `"btn btn-primary"` on `<button type=submit>` | ~8 | `events/view.go`, `valuestore/view.go`, `uploads/view.go`, `uploads/filtered/view.go` |
| Text input | `"input"` on `<input>` | ~12 | `events/view.go`, `signals/ws_view.go`, `broadcasting/view.go`, `groups/view.go`, `valuestore/view.go`, `uploads/view.go` |
| Form group | `"form-group"` on `<div>` | ~6 | `events/view.go` |
| Label | `"label"` on `<label>` | ~6 | `events/view.go` |
| Activity feed | `"activity-feed"`, `"activity-item"`, `"activity-user/text/time"` | ~10 | `live/ws_view.go`, `groups/view.go` |
| Client-only button | `"btn btn-secondary"` on `<button>` + `bind.SetSignal/ToggleSignal` | ~7 | `signals/ws_view.go` |

### `button.Link` is external-only

`button.Link()` applies `target=_blank` and `rel=noopener`, making it
suitable only for external URLs. Navigation demos need internal
link-buttons with `bind.Link()` but there is no component variant for this,
forcing raw `a.New().Class("btn btn-secondary")` throughout
`navigation/view.go`.

### `component/` tether components use raw classes

`counter.go`, `counter_group.go`, and `shoutbox.go` all build buttons with
raw `button.New().Class("btn btn-primary btn-sm")` instead of using the
`components/simple/button` package. These are `tether.Component`
implementations, so they can't easily take a dependency on the example
app's component tree - but it means CSS class names are scattered across
two unrelated packages.

---

## 7. Button Crowding & Touch Targets

### Navigation page

Up to 4 secondary buttons sit side by side in a `demo-row` (`flex-wrap:
wrap; gap: 0.75rem`). At intermediate viewport widths (roughly 900-1100px
depending on sidebar), the buttons are tight before wrapping kicks in.
The text-heavy labels ("Tab: settings, page 3") make the buttons wide
enough to crowd without quite triggering the wrap.

### Counter components

The `−`, `+`, and `Reset` buttons use `btn-sm` (`padding: 0.25rem 0.5rem;
font-size: 0.75rem`). With single-character labels the buttons are very
small - potentially below the 44x44px recommended touch target. The
`CounterGroup` nests two sets of these inside `demo-columns` (a
`1fr 1fr` grid), further compressing horizontal space.

### Form submit buttons

Submit buttons (`btn btn-primary`) sit directly below inputs in forms via
the `form > gap: 0.75rem` rule. There is no visual separation between the
form fields and the action button - a subtle border, divider, or extra
spacing would improve clarity.

---

## 8. CSS Inconsistencies

### `.item-list` vs `.list-item`

CSS styles list items via the tag selector `.item-list li` (line 848), but
the Go component applies a `.list-item` class to each `<li>` (in
`list/item.go:10`). The class is redundant - styling comes from the
parent/tag combination. Either the CSS should target `.list-item` or the
class should be removed from the Go component. Currently the `.list-item`
class name is just noise.

### `.result-success` on result blocks

`events/view.go:219,264,294` applies `Class("result-block result-success")`
to `<pre>` elements, but `.result-success` (line 940) only sets
`color: oklch(0.72 0.17 142)`. On a `.result-block` (which already has
`color: var(--foreground)` and `background: var(--muted)`), this turns the
monospace text green but doesn't change the block's background. It works
visually but looks accidental - if success blocks should look different, a
dedicated `.result-block-success` modifier with a tinted background would
be more intentional.

### `.signal-panel` used for two purposes

`.signal-panel` serves as both a content container (via `panel.Signal()`)
and a text paragraph (`panel.SignalText()` applies it to a `<p>`). It also
appears directly in view files on various elements. The same class on both
block and inline elements means the padding/border treatment may look
inconsistent depending on content.

### `demo-columns` used outside the demo component

`demo-columns` is conceptually part of the demo card layout but is applied
directly in `notifications/view.go:81,106` and
`component/counter_group.go:64`. If this layout is useful outside demos, it
should be a general-purpose component (or at minimum renamed).

---

## 9. Recommendations (Ordered by Impact)

### Must-fix before publishing

1. **Add missing `msg-*` CSS** - broadcasting messages are completely
   unstyled. Either add rules mirroring `activity-*` or reuse the
   activity-feed pattern directly.

2. **Fix progress bar class** - change `uploads/view.go:58` from
   `Class("progress")` to `Class("upload-progress")` so the styled
   pseudo-elements actually apply.

3. **Add `shout-text` CSS** - single rule, easy fix, currently missing
   colour on message body text.

4. **Differentiate `badge-indigo`** - give it a distinct background/colour
   so the two badge variants are visually distinguishable.

5. **Add `counter-display` CSS** (or rename usages to `stat-value`) - the
   value store counters render in the default body font, making them hard
   to read as primary numeric indicators.

6. **Add `monitor-charts` layout CSS** - the real-time dashboard charts
   need a grid or flex rule to sit side by side instead of stacking.

### Should-fix (component hygiene)

7. **Delete dead `component/` duplicates** - remove `component.Demo`,
   `component.Card`, `component.Badge`, and the duplicate `Transport`
   type. Only `Counter`, `CounterGroup`, and `Shoutbox` belong there.

8. **Delete the unused `components/simple/badge/` package** - or wire it
   into the views that currently use raw `badge badge-green` strings.

9. **Remove dead CSS** - the `table`, `tabs`, `flash-target`,
   `upload-area`, `room-badge`, `room-none`, `chain-log`, `stat-value`,
   `feature-*`, and `progress-area` rules are unused weight.

10. **Consolidate `.upload-status`** - remove the duplicate definition.

### Nice-to-have (polish)

11. Add components for common raw-class patterns: `Submit` button, internal
    `LinkButton`, `Input`, `FormGroup`, `ActivityFeed`.

12. Increase `btn-sm` touch targets or add minimum dimensions.

13. Add `viewport-*` CSS for the infinite scroll demo (even minimal styling
    would clarify the list boundary and sentinel).

14. Rename or scope `demo-columns` if it's used outside demo cards.

15. Resolve the `.list-item` / `.item-list li` inconsistency - pick one
    selector strategy and stick with it.
