package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestSignalsSSEPageRenders verifies the SSE signals page loads.
// Uses RealHTTP1 because httptest's internal connection tracking
// stalls EventSource retries for ~30s.
func TestSignalsSSEPageRenders(t *testing.T) {
	srv := startApp(t, RealHTTP1)
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/sse/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Increment Server Counter"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("increment button not visible: %v", err)
	}
}

// TestSignalsSSEIncrement clicks the increment button and verifies
// the counter updates via the SSE transport.
func TestSignalsSSEIncrement(t *testing.T) {
	srv := startApp(t, RealHTTP1)
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/sse/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='signals.increment']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The counter is updated via BindText signal - look for the text
	// content changing rather than a specific Dynamic key.
	counter := page.Locator("[data-tether-bind-text='signals.counter']")
	if err := expect(counter).ToContainText("2"); err != nil {
		text, _ := counter.TextContent()
		t.Errorf("counter = %q, want to contain 2", text)
	}
}
