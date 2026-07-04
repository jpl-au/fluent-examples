package playwright_test

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/bind"
	pw "github.com/playwright-community/playwright-go"

	"github.com/jpl-au/fluent-examples/tether/store"
)

// The reconnect suite drives a real browser through the three
// recovery paths a production deploy exercises:
//
//   - transport drop and reattach on the same server (state preserved)
//   - server restart with a SessionStore (state restored from disk)
//   - server restart without a store (stale client replaced by a
//     fresh session via a full morph)
//
// These are exactly the flows a unit test cannot prove end to end:
// the assertion is that the DOM in a real browser re-syncs.

// reconnectState is the per-session state for the reconnect fixture.
type reconnectState struct {
	Count int
}

func renderReconnect(s reconnectState) node.Node {
	return div.New(
		span.Text("Count: "+strconv.Itoa(s.Count)).Dynamic("count"),
		bind.Apply(button.Text("Increment"), bind.OnClick("inc")).ID("inc"),
	)
}

func handleReconnect(_ tether.Session, s reconnectState, ev tether.Event) reconnectState {
	if ev.Action == "inc" {
		s.Count++
	}
	return s
}

// restartableServer serves the reconnect fixture on a fixed address
// so a "process restart" can bring up a brand-new handler - with
// empty in-memory pools - at the same URL the browser reconnects to.
type restartableServer struct {
	t     *testing.T
	addr  string
	store tether.SessionStore // nil = no crash recovery
	srv   *http.Server
	h     *tether.Handler[reconnectState]
}

func startRestartable(t *testing.T, sessionStore tether.SessionStore) *restartableServer {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	s := &restartableServer{t: t, addr: ln.Addr().String(), store: sessionStore}
	s.start(ln)
	t.Cleanup(func() { s.stop() })
	return s
}

func (s *restartableServer) start(ln net.Listener) {
	s.h = tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[reconnectState]{
		InitialState: func(*http.Request) reconnectState { return reconnectState{} },
		Render:       renderReconnect,
		Handle:       handleReconnect,
		SessionStore: s.store,
	})
	s.srv = &http.Server{Handler: s.h}
	go s.srv.Serve(ln)
}

func (s *restartableServer) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Handler shutdown first: it persists session state to the store
	// (when one is configured) before the listener goes away.
	s.h.Shutdown(ctx)
	s.srv.Close()
}

// restart simulates a deploy: the old process shuts down (persisting
// to the store when configured) and a fresh one binds the same
// address with no in-memory sessions.
func (s *restartableServer) restart() {
	s.t.Helper()
	s.stop()

	// The port is only free once the old listener has fully closed;
	// the browser may also be hammering it with reconnect attempts.
	var ln net.Listener
	var err error
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		ln, err = net.Listen("tcp", s.addr)
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		s.t.Fatalf("rebind %s: %v", s.addr, err)
	}
	s.start(ln)
}

// clickIncrement clicks the increment button n times and waits for
// the expected count so events are never sent faster than asserted.
func clickIncrement(t *testing.T, page pw.Page, n, want int) {
	t.Helper()
	for i := 0; i < n; i++ {
		if err := page.Locator("#inc").Click(); err != nil {
			t.Fatalf("click: %v", err)
		}
	}
	if err := expect(page.Locator("[data-fluent-key='count']")).
		ToHaveText("Count: " + strconv.Itoa(want)); err != nil {
		t.Fatalf("count did not reach %d: %v", want, err)
	}
}

// sessionID reads the session ID the client is currently bound to.
func sessionID(t *testing.T, page pw.Page) string {
	t.Helper()
	v, err := page.Locator("[data-tether-root]").GetAttribute("data-tether-session")
	if err != nil {
		t.Fatalf("read session attribute: %v", err)
	}
	return v
}

// TestReconnectReattachPreservesState drops the transport and lets
// the client reconnect to the same server. The session must reattach:
// server-side state survives and the page keeps working.
func TestReconnectReattachPreservesState(t *testing.T) {
	s := startRestartable(t, nil)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto("http://" + s.addr + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)
	before := sessionID(t, page)
	clickIncrement(t, page, 3, 3)

	// Drop the transport (dev-mode helper). The client reconnects
	// with its exponential backoff and a fresh connect ticket.
	if _, err := page.Evaluate("Tether.disconnect()"); err != nil {
		t.Fatalf("disconnect: %v", err)
	}
	waitForConnected(t, page)

	// Same session, state intact, still interactive.
	if after := sessionID(t, page); after != before {
		t.Errorf("session changed across reattach: %q != %q", after, before)
	}
	if err := expect(page.Locator("[data-fluent-key='count']")).
		ToHaveText("Count: 3"); err != nil {
		t.Fatalf("state lost across reattach: %v", err)
	}
	clickIncrement(t, page, 1, 4)
}

// TestReconnectRestoresFromSessionStore restarts the server behind
// the browser's back. With a SessionStore configured, the new process
// must restore the session from disk: the DOM re-syncs with the
// pre-restart state and the same session ID.
func TestReconnectRestoresFromSessionStore(t *testing.T) {
	s := startRestartable(t, store.NewFileSessionStore(t.TempDir()))
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto("http://" + s.addr + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)
	before := sessionID(t, page)
	clickIncrement(t, page, 5, 5)

	s.restart()
	waitForConnected(t, page)

	if after := sessionID(t, page); after != before {
		t.Errorf("session changed across restore: %q != %q", after, before)
	}
	if err := expect(page.Locator("[data-fluent-key='count']")).
		ToHaveText("Count: 5"); err != nil {
		t.Fatalf("state not restored from store: %v", err)
	}
	clickIncrement(t, page, 1, 6)
}

// TestReconnectStaleClientGetsFreshSession restarts the server with
// no SessionStore. The browser holds a session ID the new process has
// never seen; it must receive a brand-new session and a full morph
// replacing the stale DOM - without the user clicking anything.
func TestReconnectStaleClientGetsFreshSession(t *testing.T) {
	s := startRestartable(t, nil)
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto("http://" + s.addr + "/"); err != nil {
		t.Fatalf("goto: %v", err)
	}
	waitForConnected(t, page)
	before := sessionID(t, page)
	clickIncrement(t, page, 3, 3)

	s.restart()
	waitForConnected(t, page)

	// The proactive full morph must reset the stale DOM to the fresh
	// session's state with no user interaction.
	if err := expect(page.Locator("[data-fluent-key='count']")).
		ToHaveText("Count: 0"); err != nil {
		t.Fatalf("stale DOM was not replaced: %v", err)
	}
	// The client must have adopted the new session ID.
	if after := sessionID(t, page); after == before {
		t.Error("client kept a session ID the server does not know")
	}
	clickIncrement(t, page, 1, 1)
}
