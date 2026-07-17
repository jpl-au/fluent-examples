package playwright_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
)

// The modifiers suite drives the declarative event modifiers through a
// real browser: Outside (click-outside closes a panel) and Once (a
// binding fires exactly once). Both run server events so the test can
// observe the effect in re-rendered state.

type modState struct {
	Open  bool
	Count int
}

func renderMods(s modState) node.Node {
	// The panel carries OnClick("close") + Outside(): the close action
	// fires only when a click lands outside the panel. A click inside
	// the panel does nothing.
	var panel node.Node = span.Text("") // empty placeholder when closed
	if s.Open {
		panel = bind.Apply(
			div.New(span.Text("PANEL BODY")),
			bind.OnClick("close"),
			bind.Outside(),
		).ID("panel")
	}

	// The bump button fires OnClick("bump") + Once(): only the first
	// click increments the count. It sits outside the Dynamic region so
	// morphs never replace it (which would reset the Once guard).
	bump := bind.Apply(button.Text("bump"), bind.OnClick("bump"), bind.Once()).ID("bump")

	return div.New(
		div.New(span.Text("outside area")).ID("outside"),
		bind.Apply(button.Text("open"), bind.OnClick("open")).ID("open"),
		bump,
		div.New(
			panel,
			span.Text(fmt.Sprintf("count:%d", s.Count)).ID("count"),
		).Dynamic("m"),
	)
}

func handleMods(_ tether.Session, s modState, ev tether.Event) modState {
	switch ev.Action {
	case "open":
		s.Open = true
	case "close":
		s.Open = false
	case "bump":
		s.Count++
	}
	return s
}

func startModsServer(t *testing.T) string {
	t.Helper()
	h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[modState]{
		InitialState: func(*http.Request) modState { return modState{} },
		Render:       renderMods,
		Handle:       handleMods,
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

// TestModifierOutside opens a panel, clicks inside it (which must NOT
// close it), then clicks outside it (which must close it via the
// click-outside binding).
func TestModifierOutside(t *testing.T) {
	srv := startModsServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#open").Click(); err != nil {
		t.Fatalf("click open: %v", err)
	}
	if err := expect(page.Locator("#panel")).ToBeVisible(); err != nil {
		t.Fatalf("panel should open: %v", err)
	}

	// A click inside the panel must not fire the close action.
	if err := page.Locator("#panel").Click(); err != nil {
		t.Fatalf("click panel: %v", err)
	}
	if err := expect(page.Locator("#panel")).ToBeVisible(); err != nil {
		t.Errorf("panel closed on an inside click: %v", err)
	}

	// A click outside the panel fires close and the panel disappears.
	if err := page.Locator("#outside").Click(); err != nil {
		t.Fatalf("click outside: %v", err)
	}
	if err := expect(page.Locator("#panel")).ToHaveCount(0); err != nil {
		t.Errorf("panel did not close on an outside click: %v", err)
	}
}

// TestModifierOnce clicks a Once-guarded button three times and asserts
// the server saw only the first click.
func TestModifierOnce(t *testing.T) {
	srv := startModsServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	bump := page.Locator("#bump")
	if err := bump.Click(); err != nil {
		t.Fatalf("first click: %v", err)
	}
	// The first click must register.
	if err := expect(page.Locator("#count")).ToHaveText("count:1"); err != nil {
		t.Fatalf("first click did not register: %v", err)
	}

	// Two further clicks must be ignored by the Once guard.
	if err := bump.Click(); err != nil {
		t.Fatalf("second click: %v", err)
	}
	if err := bump.Click(); err != nil {
		t.Fatalf("third click: %v", err)
	}

	// Give any erroneous extra events time to land, then confirm the
	// count is still 1.
	if err := expect(page.Locator("#count")).ToHaveText("count:1"); err != nil {
		text, _ := page.Locator("#count").TextContent()
		t.Errorf("Once fired more than once: count = %q, want count:1", text)
	}
}
