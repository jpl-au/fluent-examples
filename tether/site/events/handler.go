package events

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
	"github.com/jpl-au/fluent-examples/tether/middleware"
)

// State holds per-request state for the events and forms demo.
// Each field captures the result of a specific demo interaction
// so the page can render feedback after each event.
type State struct {
	ClickCount        int
	InputValue        string
	SubmitResult      string
	SubmitError       string
	AutoFocusResult   string
	AutoFocusError    string
	ChangeValue       string
	LastKey           string
	ThrottleHits      int
	EventDataResult   string
	TypedResult       string
	BindResult        string
	ViewportPage      int
	CustomEventResult string
	ResetResult       string
	FocusBlurResult   string
	RawDataResult     string
	PasteResult       string
	ContextMenuResult string
	ValidatedResult   string
	EditableResult    string
}

// New creates a stateless page handler for the events demo.
func New(app tether.App, assets *tether.Asset) http.Handler {
	return tether.Stateless(app, tether.StatelessConfig[State]{
		InitialState: func(_ *http.Request) State { return State{} },
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionHTTP, "/events", 0, Render(s))
		},
		Handle:     Handle,
		Middleware: []tether.Middleware[State]{middleware.Logging[State]},
		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Events & Forms"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},
	})
}

// Handle processes events on the events page, dispatching each
// action to the appropriate state mutation. Because the HTTP section
// is stateless, counters and lists are round-tripped via EventData.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "events.click":
		s.ClickCount = eventCount(ev) + 1
	case "events.input":
		s.InputValue = ev.Value()
	case "events.submit":
		name, _ := ev.Get("name")
		if name == "" {
			s.SubmitError = "Name is required"
			s.SubmitResult = ""
		} else {
			s.SubmitResult = "Hello, " + name + "!"
			s.SubmitError = ""
		}
	case "events.autofocus":
		email, _ := ev.Get("email")
		if email == "" {
			s.AutoFocusError = "Email is required"
			s.AutoFocusResult = ""
		} else {
			s.AutoFocusResult = "Submitted: " + email
			s.AutoFocusError = ""
		}
	case "events.change":
		s.ChangeValue = ev.Value()
	case "events.keydown":
		s.LastKey = "Enter"
	case "events.focus":
		s.FocusBlurResult = "Field focused"
	case "events.blur":
		s.FocusBlurResult = "Field blurred"
	case "events.confirm":
		sess.Toast("Confirmed! (nothing was actually deleted)")
	case "events.throttle":
		s.ThrottleHits = eventCount(ev) + 1
	case "events.data":
		itemID, _ := ev.Get("item-id")
		category, _ := ev.Get("category")
		s.EventDataResult = fmt.Sprintf("item-id=%s, category=%s", itemID, category)
	case "events.raw-data":
		status, _ := ev.Get("data-status")
		s.RawDataResult = fmt.Sprintf("data-status=%s", status)
	case "events.typed":
		qty, qerr := ev.Int("qty")
		price, perr := ev.Float64("price")
		urgent := ev.Bool("urgent")
		if qerr != nil {
			s.TypedResult = fmt.Sprintf("quantity parse error: %v", qerr)
		} else if perr != nil {
			s.TypedResult = fmt.Sprintf("price parse error: %v", perr)
		} else {
			s.TypedResult = fmt.Sprintf("quantity=%d, price=%.2f, urgent=%v", qty, price, urgent)
		}
	case "events.bind":
		var form struct {
			Name  string `tether:"name"`
			Email string `tether:"email"`
		}
		if err := ev.Bind(&form); err != nil {
			s.BindResult = fmt.Sprintf("bind error: %v", err)
		} else {
			s.BindResult = fmt.Sprintf("name=%q, email=%q", form.Name, form.Email)
		}
	case "events.indicator":
		time.Sleep(time.Second)
		sess.Toast("Loading complete!")
	case "events.custom":
		s.CustomEventResult = "Double-click received!"
	case "events.reset":
		msg, _ := ev.Get("message")
		if msg != "" {
			s.ResetResult = fmt.Sprintf("Sent: %q", msg)
		}
	case "events.viewport":
		page, _ := ev.Int("page")
		s.ViewportPage = page + 1
	case "events.paste":
		s.PasteResult = fmt.Sprintf("Pasted: %q", ev.Value())
	case "events.contextmenu":
		s.ContextMenuResult = "Context menu intercepted!"
	case "events.validated":
		name, _ := ev.Get("validated-name")
		s.ValidatedResult = fmt.Sprintf("Validated: %s", name)
	case "events.editable":
		s.EditableResult = fmt.Sprintf("Edited to: %q", ev.Value())
	}
	return s
}

// eventCount reads a "count" key from the event data, returning 0 if
// absent or unparseable. Used by stateless counters to carry state
// across requests via bind.EventData.
func eventCount(ev tether.Event) int {
	s, _ := ev.Get("count")
	n, _ := strconv.Atoi(s)
	return n
}
