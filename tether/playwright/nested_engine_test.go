package playwright_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
)

type nestedBrowserState struct {
	keys []string
	a    string
}

func nestedBrowserRow(key, a string) node.Node {
	value := key
	if key == "a" {
		value = a
	}
	return span.Text(value).Dynamic(key)
}

func nestedBrowserList(s nestedBrowserState) node.Node {
	var rows []node.Node
	for _, key := range s.keys {
		rows = append(rows, nestedBrowserRow(key, s.a))
	}
	return div.New(rows...).Dynamic("list")
}

func TestNestedEngineBrowser(t *testing.T) {
	for _, transport := range []mode.Transport{mode.WebSocket, mode.ServerSentEvents} {
		for _, memo := range []bool{false, true} {
			t.Run(fmt.Sprintf("transport=%v/memo=%v", transport, memo), func(t *testing.T) {
				sessions := make(chan *tether.StatefulSession[nestedBrowserState], 1)
				h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[nestedBrowserState]{
					Mode: transport, Memoise: memo,
					InitialState: func(*http.Request) nestedBrowserState { return nestedBrowserState{[]string{"a", "b"}, "0"} },
					Render: func(s nestedBrowserState) node.Node {
						return div.New(span.Text("outside").ID("outside"), div.New(nestedBrowserList(s)).Dynamic("panel"))
					},
					Handle:    func(_ tether.Session, s nestedBrowserState, _ tether.Event) nestedBrowserState { return s },
					OnConnect: func(s *tether.StatefulSession[nestedBrowserState]) { sessions <- s },
				})
				srv := httptest.NewServer(h)
				defer srv.Close()
				defer h.Shutdown(context.Background())
				page, cleanup := newPage(t)
				defer cleanup()
				if _, err := page.Goto(srv.URL); err != nil {
					t.Fatal(err)
				}
				waitForConnected(t, page)
				var session *tether.StatefulSession[nestedBrowserState]
				select {
				case session = <-sessions:
				case <-time.After(5 * time.Second):
					t.Fatal("no session")
				}
				js := func(script string) any {
					t.Helper()
					result, err := page.Evaluate(script)
					if err != nil {
						t.Fatal(err)
					}
					return result
				}
				wantA := func(value string) {
					t.Helper()
					if err := expect(page.Locator("#a")).ToHaveText(value); err != nil {
						t.Fatal(err)
					}
				}
				js("document.getElementById('outside').setAttribute('data-client', 'preserve')")
				session.Update(func(s nestedBrowserState) nestedBrowserState { s.keys = []string{"a", "b", "c"}; return s })
				if err := expect(page.Locator("#c")).ToHaveText("c"); err != nil {
					t.Fatal(err)
				}
				if js("document.getElementById('outside').getAttribute('data-client')") != "preserve" {
					t.Fatal("nested addition morphed the root")
				}
				session.Update(func(s nestedBrowserState) nestedBrowserState { s.keys = []string{"c", "b", "a"}; return s })
				if _, err := page.WaitForFunction("Array.from(document.querySelector('#list').children).map(e => e.id).join(',') === 'c,b,a'", nil); err != nil {
					t.Fatal(err)
				}
				patchList := func(value string) {
					session.Patch("list", func(s nestedBrowserState) (nestedBrowserState, node.Node) {
						s.a = value
						return s, nestedBrowserList(s)
					})
				}
				patchA := func(value string) {
					session.Patch("a", func(s nestedBrowserState) (nestedBrowserState, node.Node) {
						s.a = value
						return s, nestedBrowserRow("a", value)
					})
				}
				patchList("1")
				wantA("1")
				patchA("0")
				wantA("0")
				patchA("2")
				wantA("2")
				patchList("0")
				wantA("0")
				patchA("3")
				wantA("3")
				js("document.getElementById('list').setAttribute('data-client', 'child-only')")
				session.Update(func(s nestedBrowserState) nestedBrowserState { s.a = "0"; return s })
				wantA("0")
				if js("document.getElementById('list').getAttribute('data-client')") != "child-only" {
					t.Fatal("child correction unnecessarily morphed its parent")
				}
				session.Update(func(s nestedBrowserState) nestedBrowserState { s.keys = []string{"b", "a"}; return s })
				if err := expect(page.Locator("#c")).Not().ToBeAttached(); err != nil {
					t.Fatal(err)
				}
				if js("document.getElementById('outside').getAttribute('data-client')") != "preserve" {
					t.Fatal("nested removal morphed the root")
				}
			})
		}
	}
}

func TestNestedMoveKeepsDOMIdentity(t *testing.T) {
	for _, memo := range []bool{false, true} {
		t.Run(fmt.Sprint(memo), func(t *testing.T) {
			sessions := make(chan *tether.StatefulSession[bool], 1)
			h := tether.Stateful(tether.App{DevMode: true}, tether.StatefulConfig[bool]{
				Memoise:      memo,
				InitialState: func(*http.Request) bool { return false },
				Render: func(right bool) node.Node {
					var leftChildren, rightChildren []node.Node
					row := span.Text("row").Dynamic("row")
					if right {
						rightChildren = append(rightChildren, row)
					} else {
						leftChildren = append(leftChildren, row)
					}
					return div.New(div.New(leftChildren...).Dynamic("left"), div.New(rightChildren...).Dynamic("right")).Dynamic("panel")
				},
				Handle:    func(_ tether.Session, s bool, _ tether.Event) bool { return s },
				OnConnect: func(s *tether.StatefulSession[bool]) { sessions <- s },
			})
			srv := httptest.NewServer(h)
			defer srv.Close()
			defer h.Shutdown(context.Background())
			page, cleanup := newPage(t)
			defer cleanup()
			if _, err := page.Goto(srv.URL); err != nil {
				t.Fatal(err)
			}
			waitForConnected(t, page)
			var session *tether.StatefulSession[bool]
			select {
			case session = <-sessions:
			case <-time.After(5 * time.Second):
				t.Fatal("no session")
			}
			if _, err := page.Evaluate("window.originalRow = document.getElementById('row')"); err != nil {
				t.Fatal(err)
			}
			session.Update(func(bool) bool { return true })
			if err := expect(page.Locator("#right #row")).ToBeAttached(); err != nil {
				t.Fatal(err)
			}
			same, err := page.Evaluate("window.originalRow === document.getElementById('row')")
			if err != nil {
				t.Fatal(err)
			}
			if same != true {
				t.Fatal("moving a keyed row replaced the DOM element")
			}
		})
	}
}
