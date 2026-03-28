package mw

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/layout"
)

// State holds per-request state for the middleware demo. Each field
// captures the result of a specific middleware interaction so the
// page can render feedback after each event.
type State struct {
	LastAction    string
	TimingResult  string
	BlockedResult string
	ChainLog      string
	EventCount    int
}

// New creates a stateless page handler for the middleware demo.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State { return State{} },
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/middleware/", 0, Render(s))
		},
		Handle: Handle,
		Middleware: []tether.Middleware[State]{
			ordered("Outer"),
			guard,
			timing,
			counting,
			ordered("Inner"),
		},
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Middleware"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes events on the middleware page, dispatching each
// action to the appropriate state mutation. Because the HTTP section
// is stateless, the event count is round-tripped via EventData.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "mw.ping":
		s.LastAction = "ping"
	case "mw.slow":
		time.Sleep(100 * time.Millisecond)
		s.LastAction = "slow"
	case "mw.blocked":
		// Guard middleware intercepts this before it reaches here.
		s.LastAction = "blocked (should not appear)"
	}
	// Read event count from round-trip data (stateless page).
	if c, _ := ev.Get("count"); c != "" {
		if n, err := strconv.Atoi(c); err == nil {
			s.EventCount = n
		}
	}
	return s
}

// timing measures how long the inner handler takes and records the
// duration in state so the view can display it.
func timing(next tether.HandleFunc[State]) tether.HandleFunc[State] {
	return func(sess tether.Session, s State, ev tether.Event) State {
		start := time.Now()
		s = next(sess, s, ev)
		s.TimingResult = fmt.Sprintf("Handled in %s", time.Since(start).Round(time.Microsecond))
		return s
	}
}

// counting increments the event count after the inner handler runs,
// regardless of which action was dispatched.
func counting(next tether.HandleFunc[State]) tether.HandleFunc[State] {
	return func(sess tether.Session, s State, ev tether.Event) State {
		s = next(sess, s, ev)
		s.EventCount++
		return s
	}
}

// guard short-circuits the middleware chain when it sees the
// "mw.blocked" action, preventing the inner handler from running.
func guard(next tether.HandleFunc[State]) tether.HandleFunc[State] {
	return func(sess tether.Session, s State, ev tether.Event) State {
		if ev.Action == "mw.blocked" {
			s.BlockedResult = "Blocked by guard middleware"
			sess.Toast("This action was intercepted before reaching the handler")
			return s
		}
		return next(sess, s, ev)
	}
}

// ordered returns middleware that appends its name to ChainLog on
// entry and exit, making the onion-like execution order visible.
func ordered(name string) tether.Middleware[State] {
	return func(next tether.HandleFunc[State]) tether.HandleFunc[State] {
		return func(sess tether.Session, s State, ev tether.Event) State {
			s.ChainLog += name + " → "
			s = next(sess, s, ev)
			s.ChainLog += name + " ← "
			return s
		}
	}
}
