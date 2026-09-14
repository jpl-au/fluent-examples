package playwright_test

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jpl-au/fluent-examples/tether/store"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
	pw "github.com/playwright-community/playwright-go"
)

func TestAuditReconnectNavigation(t *testing.T) {
	for _, transport := range []mode.Transport{mode.WebSocket, mode.ServerSentEvents} {
		for _, scenario := range []string{"held", "offline", "thaw", "restore"} {
			t.Run(fmt.Sprintf("%v/%s", transport, scenario), func(t *testing.T) {
				sessions := make(chan *tether.StatefulSession[string], 8)
				cfg := tether.StatefulConfig[string]{
					Mode:         transport,
					InitialState: func(*http.Request) string { return "/" },
					Render:       func(s string) node.Node { return div.Text(s).Dynamic("path") },
					Handle:       func(_ tether.Session, s string, _ tether.Event) string { return s },
					OnNavigate:   func(_ tether.Session, _ string, p tether.Params) string { return p.Path },
					OnConnect:    func(s *tether.StatefulSession[string]) { sessions <- s },
					Timeouts:     tether.Timeouts{Retry: 500 * time.Millisecond, DisableJitter: true},
				}
				if scenario == "thaw" || scenario == "restore" {
					if scenario == "thaw" {
						cfg.Freeze = tether.FreezeWithConnect
					}
					cfg.SessionStore = store.NewFileSessionStore(t.TempDir())
				} else {
					cfg.OnDisconnect = func(s *tether.StatefulSession[string]) { s.Navigate("/held") }
				}
				h := tether.Stateful(tether.App{DevMode: true}, cfg)
				var current atomic.Pointer[tether.Handler[string]]
				current.Store(h)
				base := startServer(t, RealHTTP1,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { current.Load().ServeHTTP(w, r) }),
					func(ctx context.Context) error { return current.Load().Shutdown(ctx) },
				)
				page, cleanup := newPage(t)
				defer cleanup()
				if _, err := page.Goto(base); err != nil {
					t.Fatal(err)
				}
				waitForConnected(t, page)
				var s *tether.StatefulSession[string]
				select {
				case s = <-sessions:
				case <-time.After(5 * time.Second):
					t.Fatal("no session")
				}
				s.Close()
				if err := page.Locator("[data-tether-state='disconnected']").WaitFor(pw.LocatorWaitForOptions{Timeout: pw.Float(3000)}); err != nil {
					t.Fatal(err)
				}
				want := "/held"
				if scenario != "held" {
					want = "/offline"
					if _, err := page.Evaluate(`history.replaceState({}, '', '/offline'); dispatchEvent(new PopStateEvent('popstate'));`); err != nil {
						t.Fatal(err)
					}
				}
				if scenario == "restore" {
					if !shutdownHandlers(t, h.Shutdown) {
						t.FailNow()
					}
					current.Store(tether.Stateful(tether.App{DevMode: true}, cfg))
				}
				if _, err := page.WaitForFunction(`want => document.querySelector('[data-tether-root]').getAttribute('data-tether-state') === 'connected' && location.pathname === want && document.getElementById('path').textContent === want`, want, pw.PageWaitForFunctionOptions{Timeout: pw.Float(4000)}); err != nil {
					state, _ := page.Evaluate(`({url:location.pathname, dom:document.getElementById('path').textContent})`)
					t.Fatalf("reconnect did not settle at %s: %v", want, state)
				}
			})
		}
	}
}
