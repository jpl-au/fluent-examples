package playwright_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
	"github.com/jpl-au/tether/mode"
	pw "github.com/playwright-community/playwright-go"
)

// These tests exercise the three developer-experience features -
// opt-in View Transitions, the Prefetch effect, and the error-slug
// channel - against a self-contained tether server rather than the
// full example app, because View Transitions is an app-level switch and
// the tests need it both on and off.

type devxState struct{ N int }

func devxRender(s devxState) node.Node {
	return div.New(
		div.Text(strconv.Itoa(s.N)).ID("counter").Dynamic("counter"),
		bind.Apply(button.Text("Bump"), bind.OnClick("bump")).ID("bump"),
		bind.Apply(button.Text("Warm"), bind.OnClick("warm")).ID("warm"),
		bind.Apply(button.Text("BadFlash"), bind.OnClick("badflash")).ID("badflash"),
	)
}

func devxHandle(sess tether.Session, s devxState, ev tether.Event) devxState {
	switch ev.Action {
	case "bump":
		s.N++
	case "warm":
		// Repeating the same URL exercises the client's per-page dedupe.
		sess.Prefetch("/devx-next")
	case "badflash":
		// An unparseable selector drives the client's safeQuery down its
		// error path, emitting the "invalid-selector" slug.
		sess.Flash("[", "boom")
	}
	return s
}

func devxLayout(_ devxState, content node.Node) node.Node {
	return html.New(
		head.New(
			meta.UTF8(),
			meta.Viewport("width=device-width, initial-scale=1"),
			title.Static("Tether - DevX"),
		),
		body.New(content),
	).Lang("en")
}

// startDevx starts a self-contained stateful tether server at the root
// path, with View Transitions on or off. The handler serves its own
// /_tether/ client assets, so no surrounding mux is needed.
func startDevx(t *testing.T, viewTransitions bool) string {
	t.Helper()
	app := tether.App{
		DevMode: true,
		Client:  tether.Client{ViewTransitions: viewTransitions},
	}
	h := tether.Stateful(app, tether.StatefulConfig[devxState]{
		Name:         "devx",
		Mode:         mode.WebSocket,
		InitialState: func(_ *http.Request) devxState { return devxState{} },
		Render:       devxRender,
		Handle:       devxHandle,
		Layout:       devxLayout,
	})
	srv := httptest.NewTLSServer(h)
	t.Cleanup(srv.Close)
	return srv.URL
}

// vtStub replaces document.startViewTransition with a recorder that
// still invokes the mutation callback, so the DOM updates exactly as a
// real transition would while the test can observe whether the runtime
// took the transition path.
const vtStub = `
window.__vtCalls = 0;
document.startViewTransition = function (cb) {
  window.__vtCalls++;
  var p = Promise.resolve(cb ? cb() : undefined);
  return { updateCallbackDone: p, ready: p, finished: p, skipTransition: function () {} };
};
`

func addInit(t *testing.T, page pw.Page, script string) {
	t.Helper()
	if err := page.AddInitScript(pw.Script{Content: pw.String(script)}); err != nil {
		t.Fatalf("add init script: %v", err)
	}
}

// TestViewTransitionsUsed verifies that when View Transitions are
// enabled the runtime routes the morph through startViewTransition and
// the DOM still updates.
func TestViewTransitionsUsed(t *testing.T) {
	srv := startDevx(t, true)
	page, cleanup := newPage(t)
	defer cleanup()
	addInit(t, page, vtStub)

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#bump").Click(); err != nil {
		t.Fatalf("click bump: %v", err)
	}
	if err := expect(page.Locator("#counter")).ToContainText("1"); err != nil {
		t.Fatalf("counter did not update: %v", err)
	}

	calls, err := page.Evaluate("() => window.__vtCalls")
	if err != nil {
		t.Fatalf("evaluate vtCalls: %v", err)
	}
	if jsNumber(calls) < 1 {
		t.Errorf("startViewTransition was not used, __vtCalls = %v", calls)
	}
}

// TestViewTransitionsDisabled verifies that with the switch off the
// runtime never calls startViewTransition, yet the DOM still updates.
func TestViewTransitionsDisabled(t *testing.T) {
	srv := startDevx(t, false)
	page, cleanup := newPage(t)
	defer cleanup()
	addInit(t, page, vtStub)

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#bump").Click(); err != nil {
		t.Fatalf("click bump: %v", err)
	}
	if err := expect(page.Locator("#counter")).ToContainText("1"); err != nil {
		t.Fatalf("counter did not update: %v", err)
	}

	calls, err := page.Evaluate("() => window.__vtCalls")
	if err != nil {
		t.Fatalf("evaluate vtCalls: %v", err)
	}
	if jsNumber(calls) != 0 {
		t.Errorf("startViewTransition was called with the switch off, __vtCalls = %v", calls)
	}
}

// hintCount counts the speculation-rules scripts and prefetch links in
// the head that reference the given URL. Either form satisfies the
// Prefetch contract depending on browser support.
const hintCountFn = `() => {
  var url = "/devx-next";
  var scripts = Array.from(document.head.querySelectorAll('script[type="speculationrules"]'))
    .filter(function (s) { return s.textContent.indexOf(url) !== -1; });
  var links = Array.from(document.head.querySelectorAll('link[rel="prefetch"]'))
    .filter(function (l) { return l.getAttribute("href").indexOf(url) !== -1; });
  return scripts.length + links.length;
}`

// TestPrefetchEffect verifies the Prefetch effect lands a speculation
// rule (or fallback link) in the head, and that repeating the same URL
// does not emit a duplicate.
func TestPrefetchEffect(t *testing.T) {
	srv := startDevx(t, false)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if err := page.Locator("#warm").Click(); err != nil {
		t.Fatalf("click warm: %v", err)
	}
	if _, err := page.WaitForFunction("() => ("+hintCountFn+")() >= 1 ? true : null", nil); err != nil {
		t.Fatalf("prefetch hint did not appear: %v", err)
	}

	first, err := page.Evaluate("(" + hintCountFn + ")()")
	if err != nil {
		t.Fatalf("evaluate hint count: %v", err)
	}
	if jsNumber(first) != 1 {
		t.Fatalf("expected exactly 1 prefetch hint, got %v", first)
	}

	// Repeat the same URL, then bump so we can wait on an observable
	// change that is ordered after the duplicate prefetch was processed.
	if err := page.Locator("#warm").Click(); err != nil {
		t.Fatalf("click warm again: %v", err)
	}
	if err := page.Locator("#bump").Click(); err != nil {
		t.Fatalf("click bump: %v", err)
	}
	if err := expect(page.Locator("#counter")).ToContainText("1"); err != nil {
		t.Fatalf("counter did not update: %v", err)
	}

	after, err := page.Evaluate("(" + hintCountFn + ")()")
	if err != nil {
		t.Fatalf("evaluate hint count after repeat: %v", err)
	}
	if jsNumber(after) != 1 {
		t.Errorf("duplicate prefetch emitted, hint count = %v (want 1)", after)
	}
}

// TestErrorSlug verifies that a dev warning reaches Tether.onError with
// the catalogued slug alongside the type and message.
func TestErrorSlug(t *testing.T) {
	srv := startDevx(t, false)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)

	if _, err := page.Evaluate(`() => {
		window.__errs = [];
		window.Tether.onError = function (e) { window.__errs.push(e); };
	}`); err != nil {
		t.Fatalf("install onError: %v", err)
	}

	if err := page.Locator("#badflash").Click(); err != nil {
		t.Fatalf("click badflash: %v", err)
	}
	if _, err := page.WaitForFunction("() => window.__errs && window.__errs.length > 0 ? true : null", nil); err != nil {
		t.Fatalf("onError never fired: %v", err)
	}

	res, err := page.Evaluate("() => window.__errs[0]")
	if err != nil {
		t.Fatalf("read error payload: %v", err)
	}
	payload, ok := res.(map[string]any)
	if !ok {
		t.Fatalf("error payload was not an object: %#v", res)
	}
	if payload["slug"] != "invalid-selector" {
		t.Errorf("slug = %v, want invalid-selector", payload["slug"])
	}
	if payload["type"] != "render" {
		t.Errorf("type = %v, want render", payload["type"])
	}
	if _, hasMsg := payload["message"]; !hasMsg {
		t.Errorf("payload missing message field: %#v", payload)
	}
}
