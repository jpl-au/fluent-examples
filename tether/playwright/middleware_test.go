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

	resp, err := page.Goto(srv + "/middleware")
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

	_, err := page.Goto(srv + "/middleware")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='mw.ping']").First()
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The chain log should show Outer → Inner → Inner ← Outer ←
	log := page.GetByText("Outer →")
	if err := expect(log).ToBeVisible(); err != nil {
		t.Errorf("chain log not visible: %v", err)
	}
}

// TestMiddlewareGuardBlock clicks the Blocked Action button and
// verifies the guard middleware short-circuits the chain.
func TestMiddlewareGuardBlock(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/middleware")
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
