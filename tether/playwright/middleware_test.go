package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestMiddlewarePageRenders verifies the middleware demo page loads.
func TestMiddlewarePageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/middleware/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}

	waitForConnected(t, page)

	heading := page.GetByText("Middleware Chain")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Errorf("heading not visible: %v", err)
	}
}

// TestMiddlewareChainOrder clicks Send Event and verifies the
// middleware chain log shows the onion-like execution order.
func TestMiddlewareChainOrder(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/middleware/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='mw.ping']").First()
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The chain log should show Outer -> Inner -> Inner <- Outer <-
	log := page.GetByText("Outer →")
	if err := expect(log).ToBeVisible(); err != nil {
		t.Errorf("chain log not visible: %v", err)
	}
}

// TestMiddlewareChainFullOrder verifies the complete onion execution
// order including both Inner markers in the chain log.
func TestMiddlewareChainFullOrder(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/middleware/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='mw.ping']").First()
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The chain log shows: "Outer → Inner → Inner ← Outer ←"
	// Verify the Inner markers appear in the chain result.
	chainResult := page.Locator("[data-fluent-key='chain-result']")
	if err := expect(chainResult).ToContainText("Inner →"); err != nil {
		t.Errorf("chain log missing Inner entry marker: %v", err)
	}
	if err := expect(chainResult).ToContainText("Inner ←"); err != nil {
		t.Errorf("chain log missing Inner exit marker: %v", err)
	}
	if err := expect(chainResult).ToContainText("Outer ←"); err != nil {
		t.Errorf("chain log missing Outer exit marker: %v", err)
	}
}

// TestMiddlewareGuardBlock clicks the Blocked Action button and
// verifies the guard middleware short-circuits the chain.
func TestMiddlewareGuardBlock(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/middleware/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Blocked Action"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	result := page.GetByText("Blocked by guard middleware")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("guard result not visible: %v", err)
	}
}

// TestMiddlewareSlowEvent clicks Slow Event and verifies the timing
// middleware records a measurable duration.
func TestMiddlewareSlowEvent(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/middleware/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='mw.slow']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click slow: %v", err)
	}

	// The timing middleware records the handler duration.
	timingResult := page.Locator("[data-fluent-key='timing-result']")
	if err := expect(timingResult).ToContainText("Handled in"); err != nil {
		t.Errorf("timing result not visible: %v", err)
	}
}

// TestMiddlewareEventCount clicks the Click Me button multiple times
// and verifies the event counter increments.
func TestMiddlewareEventCount(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/middleware/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The "Click Me" button in the Event Counting card.
	btn := page.Locator("[data-tether-click='mw.ping']").Last()
	if err := btn.Click(); err != nil {
		t.Fatalf("click 1: %v", err)
	}

	// After one click, the counting middleware increments to 1.
	count1 := page.GetByText("Events processed: 1")
	if err := expect(count1).ToBeVisible(); err != nil {
		t.Errorf("event count did not reach 1: %v", err)
	}

	// Click again - count should reach 2.
	if err := btn.Click(); err != nil {
		t.Fatalf("click 2: %v", err)
	}

	count2 := page.GetByText("Events processed: 2")
	if err := expect(count2).ToBeVisible(); err != nil {
		t.Errorf("event count did not reach 2: %v", err)
	}
}
