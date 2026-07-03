package playwright_test

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
	pw "github.com/playwright-community/playwright-go"
)

// The auto-fragments browser test proves the whole protocol through a
// real client: the initial GET seeds the hash island, the runtime
// echoes the map with each event, and the response carries only the
// fragment that changed - asserted on the actual wire bytes, not just
// the resulting DOM.

// autofragState has two regions plus a deliberately large static
// section, so a fragment response is visibly smaller than the page.
type autofragState struct {
	A, B int
}

func renderAutofrag(s autofragState) node.Node {
	// Stateless mode reconstructs state each request, so the current
	// count travels with the event as data (the standard pattern -
	// see the morph demo).
	return div.New(
		div.New(
			span.Text("Region A: "+strconv.Itoa(s.A)).Dynamic("af-a"),
			span.Text("Region B: "+strconv.Itoa(s.B)).Dynamic("af-b"),
		),
		// The button carries the current count as event data, so it
		// lives in its own fragment to refresh alongside region A.
		div.New(
			bind.Apply(button.Text("Bump A"), bind.OnClick("bump-a"),
				bind.EventData("a", strconv.Itoa(s.A))).ID("bump-a"),
		).Dynamic("af-btn"),
		div.New(span.Text(strings.Repeat("static filler content ", 200))).Class("filler"),
	)
}

func startAutofragServer(t *testing.T) string {
	t.Helper()
	h := tether.Stateless(tether.App{}, tether.StatelessConfig[autofragState]{
		AutoFragments: true,
		InitialState: func(*http.Request) autofragState {
			return autofragState{A: 1, B: 2}
		},
		Render: renderAutofrag,
		Handle: func(_ tether.Session, s autofragState, ev tether.Event) autofragState {
			if ev.Action == "bump-a" {
				if raw, ok := ev.Get("a"); ok {
					if n, err := strconv.Atoi(raw); err == nil {
						s.A = n
					}
				}
				s.A++
			}
			return s
		},
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

// TestAutoFragmentsBrowserSendsOnlyChangedRegion drives the protocol
// through the real runtime and asserts on the POST response body:
// only the changed fragment travels, and the untouched region and
// static filler stay off the wire.
func TestAutoFragmentsBrowserSendsOnlyChangedRegion(t *testing.T) {
	srv := startAutofragServer(t)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Capture the event POST's response while clicking.
	resp, err := page.ExpectResponse(srv+"/", func() error {
		return page.Locator("#bump-a").Click()
	}, pw.PageExpectResponseOptions{})
	if err != nil {
		t.Fatalf("no event response captured: %v", err)
	}
	body, err := resp.Text()
	if err != nil {
		t.Fatalf("response body: %v", err)
	}

	if !strings.Contains(body, "Region A: 2") {
		t.Errorf("response should carry the changed fragment, got %s", body)
	}
	if strings.Contains(body, "Region B") {
		t.Error("unchanged region travelled on the wire")
	}
	if strings.Contains(body, "static filler content") {
		t.Error("static content travelled on the wire - not a fragment response")
	}
	if !strings.Contains(body, `"hashes"`) {
		t.Error("response should carry the refreshed hash map")
	}

	// And the DOM applied it.
	if err := expect(page.Locator("[data-tether-key='af-a']")).
		ToHaveText("Region A: 2"); err != nil {
		t.Fatalf("fragment not applied to the DOM: %v", err)
	}
	if err := expect(page.Locator("[data-tether-key='af-b']")).
		ToHaveText("Region B: 2"); err != nil {
		t.Fatalf("untouched region should be intact: %v", err)
	}

	// A second event proves the refreshed map round-trips: the next
	// response still targets only region A.
	resp2, err := page.ExpectResponse(srv+"/", func() error {
		return page.Locator("#bump-a").Click()
	}, pw.PageExpectResponseOptions{})
	if err != nil {
		t.Fatalf("no second response: %v", err)
	}
	body2, _ := resp2.Text()
	if !strings.Contains(body2, "Region A: 3") || strings.Contains(body2, "Region B") {
		t.Errorf("second event should still be a targeted fragment, got %s", body2)
	}
	if err := expect(page.Locator("[data-tether-key='af-a']")).
		ToHaveText("Region A: 3"); err != nil {
		t.Fatalf("second fragment not applied: %v", err)
	}
}
