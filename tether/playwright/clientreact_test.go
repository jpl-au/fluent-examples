package playwright_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/input"
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/html5/template"
	"github.com/jpl-au/fluent/html5/ul"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
	pw "github.com/playwright-community/playwright-go"
)

// The client-reactivity suite drives the four hyperscript-inspired
// features through a real browser, proving they run entirely on the
// client:
//
//   - conditional bindings (ShowWhen / ClassWhen) react to a
//     server-pushed signal without the server sending booleans
//   - client-side filtering hides list items by text match
//   - client events (Emit / OnClientEvent) clear the filter box and
//     re-run the filter with no round-trip
//   - signal-driven templates render a JSON-array signal into a list
//
// A stateful server is used so the test can push signals in response to
// clicks; the reactive behaviour it verifies is all client-side.

type person struct {
	Name string `json:"name"`
}

type reactState struct {
	Count int
}

func renderReact(s reactState) node.Node {
	// A <template> whose content is the per-item markup. li.Text keeps
	// the {{name}} placeholder verbatim (no HTML-special characters to
	// escape), so the extension sees "<li>{{name}}</li>".
	itemTemplate := bind.Apply(
		template.New(li.Text("{{name}}")),
		bind.Template("people", "#people-list"),
	).ID("item-template")

	return div.New(
		// Conditional bindings driven by the "count" signal.
		bind.Apply(span.Text("HIGH"), bind.ShowWhen("count", ">", 2)).ID("high"),
		bind.Apply(span.Text("box"), bind.ClassWhen("danger", "count", ">=", 3)).ID("box"),
		bind.Apply(button.Text("inc"), bind.OnClick("inc")).ID("inc"),
		bind.Apply(button.Text("load"), bind.OnClick("load")).ID("load"),

		// Client-side filter plus a Clear button that clears it via a
		// client event - no server involvement for either.
		bind.Apply(input.Text("q", ""),
			bind.Filter("#items"),
			bind.Value("query"),
			bind.OnClientEvent("clear", bind.SetSignal("query", "")),
		).ID("search"),
		bind.Apply(button.Text("Clear"), bind.Emit("clear", "#search")).ID("clear"),
		ul.New(
			bind.Apply(li.Text("apple"), bind.FilterItem()),
			bind.Apply(li.Text("banana"), bind.FilterItem()),
			bind.Apply(li.Text("cherry"), bind.FilterItem()),
		).ID("items"),

		// Target the template renders into.
		ul.New().ID("people-list"),
		itemTemplate,
	)
}

func handleReact(sess tether.Session, s reactState, ev tether.Event) reactState {
	switch ev.Action {
	case "inc":
		s.Count++
		sess.Signal("count", s.Count)
	case "load":
		sess.Signal("people", []person{{Name: "Ada"}, {Name: "Alan"}, {Name: "Grace"}})
	}
	return s
}

func startReactServer(t *testing.T) string {
	t.Helper()
	h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[reactState]{
		InitialState: func(*http.Request) reactState { return reactState{} },
		Render:       renderReact,
		Handle:       handleReact,
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := &http.Server{Handler: h}
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close() })
	return "http://" + ln.Addr().String()
}

// TestConditionalBindingsReactToSignal clicks an increment button that
// pushes a "count" signal and asserts ShowWhen and ClassWhen
// derive their state from the value - no boolean signal is ever sent.
func TestConditionalBindingsReactToSignal(t *testing.T) {
	srv := startReactServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	inc := page.Locator("#inc")

	// One click: count = 1, which is not > 2, so HIGH stays hidden.
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	if err := expect(page.Locator("#high")).ToBeHidden(); err != nil {
		t.Errorf("HIGH should be hidden at count=1: %v", err)
	}

	// Two more clicks: count = 3, which is > 2 and >= 3.
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	if err := inc.Click(); err != nil {
		t.Fatalf("click inc: %v", err)
	}
	if err := expect(page.Locator("#high")).ToBeVisible(); err != nil {
		t.Errorf("HIGH should be visible at count=3: %v", err)
	}
	if err := expect(page.Locator("#box.danger")).ToHaveCount(1); err != nil {
		t.Errorf("box should carry the danger class at count=3: %v", err)
	}
}

// TestClientFilter types into the search box and asserts only matching
// items remain visible, with no server round-trip.
func TestClientFilter(t *testing.T) {
	srv := startReactServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#search").Fill("ban"); err != nil {
		t.Fatalf("fill search: %v", err)
	}

	if err := expect(page.Locator("#items li").Filter(pw.LocatorFilterOptions{HasText: "banana"})).ToBeVisible(); err != nil {
		t.Errorf("banana should stay visible: %v", err)
	}
	if err := expect(page.Locator("#items li").Filter(pw.LocatorFilterOptions{HasText: "apple"})).ToBeHidden(); err != nil {
		t.Errorf("apple should be filtered out: %v", err)
	}
}

// TestClientEventClearsFilter clicks a Clear button that emits a client
// event; the search input receives it, clears its bound value, and the
// filter re-runs so every item reappears - all client-side.
func TestClientEventClearsFilter(t *testing.T) {
	srv := startReactServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#search").Fill("ban"); err != nil {
		t.Fatalf("fill search: %v", err)
	}
	if err := expect(page.Locator("#items li").Filter(pw.LocatorFilterOptions{HasText: "apple"})).ToBeHidden(); err != nil {
		t.Fatalf("apple should be filtered before clear: %v", err)
	}

	if err := page.Locator("#clear").Click(); err != nil {
		t.Fatalf("click clear: %v", err)
	}

	// The bound input is emptied and every item shows again.
	if err := expect(page.Locator("#search")).ToHaveValue(""); err != nil {
		t.Errorf("search should be cleared: %v", err)
	}
	if err := expect(page.Locator("#items li").Filter(pw.LocatorFilterOptions{HasText: "apple"})).ToBeVisible(); err != nil {
		t.Errorf("apple should reappear after clear: %v", err)
	}
}

// TestSignalTemplate clicks a button that pushes a JSON-array signal and
// asserts the template extension renders one list item per element.
func TestSignalTemplate(t *testing.T) {
	srv := startReactServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#load").Click(); err != nil {
		t.Fatalf("click load: %v", err)
	}

	list := page.Locator("#people-list li")
	if err := expect(list).ToHaveCount(3); err != nil {
		t.Fatalf("template should render three items: %v", err)
	}
	if err := expect(page.Locator("#people-list")).ToContainText("Ada"); err != nil {
		t.Errorf("rendered list should contain Ada: %v", err)
	}
	if err := expect(page.Locator("#people-list")).ToContainText("Grace"); err != nil {
		t.Errorf("rendered list should contain Grace: %v", err)
	}
}
