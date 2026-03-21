// session_test.go demonstrates every Session side-effect assertion:
// Toast, Flash, Announce, Signal, Signals (batch and raw map),
// Title, URL, Replaced, Navigate with OnNavigate, and the
// Connect/Disconnect lifecycle hooks.
package tethertest_test

import (
	"strings"
	"testing"

	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/tethertest"

	"github.com/jpl-au/fluent-examples/tether/site/navigation"
	"github.com/jpl-au/fluent-examples/tether/site/notifications"
	"github.com/jpl-au/fluent-examples/tether/site/signals"
)

// ---------- Toast, Flash, Announce ----------

// TestToastAndFlash demonstrates HasToast and HasFlash - both are
// from the notifications handler.
func TestToastAndFlash(t *testing.T) {
	h := tethertest.New(tethertest.Config[notifications.State]{
		Handle: notifications.Handle,
	})

	h.Send("notify.toast")
	if !h.HasToast("This is a toast notification!") {
		t.Fatalf("Toast = %q", h.Toast())
	}

	h.Send("notify.flash")
	if !h.HasFlash("#flash-target", "Flashed!") {
		t.Fatalf("Flash = %v", h.Flash())
	}
}

// TestAnnounce demonstrates HasAnnounce for accessibility
// announcements - the handler fires both an announcement and a
// companion toast in the same event.
func TestAnnounce(t *testing.T) {
	h := tethertest.New(tethertest.Config[notifications.State]{
		Handle: notifications.Handle,
	})

	h.Send("notify.announce")
	if !h.HasAnnounce("New item added to your feed") {
		t.Fatalf("Announce = %q", h.Announce())
	}
	if !h.HasToast("Announcement sent to ARIA live region") {
		t.Fatalf("Toast = %q", h.Toast())
	}
}

// ---------- Signals ----------

// TestSignalAndHasSignal demonstrates HasSignal for asserting
// individual signal values with their original Go types.
func TestSignalAndHasSignal(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Handle: signals.Handle,
	})

	h.Send("signals.increment")
	if !h.HasSignal("signals.counter", 1) {
		t.Fatal("expected signals.counter = 1")
	}
}

// TestSignalsRawMap demonstrates Signals() - the raw map of all
// signal values pushed during the most recent Send. Useful when
// a handler pushes multiple signals via sess.Signals().
func TestSignalsRawMap(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Handle: signals.Handle,
	})

	h.Send("signals.reset-all")
	m := h.Signals()
	if m == nil {
		t.Fatal("Signals() should not be nil after batch reset")
	}
	if m["signals.counter"] != 0 {
		t.Fatalf("signals.counter = %v, want 0", m["signals.counter"])
	}
	if m["signals.panel_visible"] != false {
		t.Fatalf("signals.panel_visible = %v, want false", m["signals.panel_visible"])
	}
}

// ---------- Title ----------

// TestTitle demonstrates Title() - the handler sets the document
// title via sess.SetTitle and the harness captures it.
func TestTitle(t *testing.T) {
	h := tethertest.New(tethertest.Config[signals.State]{
		Handle: func(sess tether.Session, s signals.State, ev tether.Event) signals.State {
			sess.SetTitle("Updated Title")
			return s
		},
	})

	h.Send("anything")
	if h.Title() != "Updated Title" {
		t.Fatalf("Title = %q, want %q", h.Title(), "Updated Title")
	}
}

// ---------- Navigation ----------

// TestNavigateAndURL demonstrates Navigate() with OnNavigate - the
// harness parses the URL and delivers typed query parameters.
func TestNavigateAndURL(t *testing.T) {
	h := tethertest.New(tethertest.Config[navigation.State]{
		Handle: navigation.Handle,
		OnNavigate: func(_ tether.Session, s navigation.State, p tether.Params) navigation.State {
			s.Page = p.Path
			s.Tab = p.Get("tab")
			s.QueryPage = p.IntDefault("page", 1)
			s.Active = p.BoolDefault("active", false)
			s.Price = p.Float64Default("price", 0)
			s.Tags = p.Strings("tag")
			return s
		},
	})

	h.Navigate("/products?tab=shoes&page=3&active=true&price=29.99")
	s := h.State()
	if s.Page != "/products" {
		t.Fatalf("Page = %q", s.Page)
	}
	if s.Tab != "shoes" {
		t.Fatalf("Tab = %q", s.Tab)
	}
	if s.QueryPage != 3 {
		t.Fatalf("QueryPage = %d", s.QueryPage)
	}
	if !s.Active {
		t.Fatal("Active should be true")
	}
	if s.Price != 29.99 {
		t.Fatalf("Price = %f", s.Price)
	}
}

// TestURLAndReplaceURL demonstrates URL() and Replaced()  -
// the handler pushes URLs via Navigate and ReplaceURL.
func TestURLAndReplaceURL(t *testing.T) {
	h := tethertest.New(tethertest.Config[navigation.State]{
		Handle: navigation.Handle,
		OnNavigate: func(_ tether.Session, s navigation.State, p tether.Params) navigation.State {
			s.Page = p.Path
			return s
		},
	})

	h.Send("nav.goto-target")
	if h.URL() != "/navigation/target" {
		t.Fatalf("URL = %q", h.URL())
	}
	if h.Replaced() {
		t.Fatal("Navigate should not be a replacement")
	}

	h.Send("nav.replace-a")
	if !strings.Contains(h.URL(), "tab=a") {
		t.Fatalf("URL = %q, expected tab=a", h.URL())
	}
	if !h.Replaced() {
		t.Fatal("ReplaceURL should be flagged as a replacement")
	}
}

// ---------- Connect / Disconnect ----------

// TestConnectAndDisconnect demonstrates the lifecycle hooks - the
// harness calls OnConnect and OnDisconnect with a test session.
func TestConnectAndDisconnect(t *testing.T) {
	var connected, disconnected bool

	h := tethertest.New(tethertest.Config[signals.State]{
		Handle: signals.Handle,
		OnConnect: func(_ tether.Session) {
			connected = true
		},
		OnDisconnect: func(_ tether.Session) {
			disconnected = true
		},
	})

	h.Connect()
	if !connected {
		t.Fatal("OnConnect was not called")
	}

	h.Disconnect()
	if !disconnected {
		t.Fatal("OnDisconnect was not called")
	}
}
