package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestFreezePageRenders verifies the freeze demo page loads and the
// counter is visible.
func TestFreezePageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/freeze/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("Frozen Counter")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestFreezeIncrement clicks increment and verifies the count updates.
func TestFreezeIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/freeze/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='freeze.increment']")
	for range 3 {
		if err := btn.Click(); err != nil {
			t.Fatalf("click: %v", err)
		}
	}

	count := page.Locator("[data-tether-key='count']")
	if err := expect(count).ToHaveText("Count: 3"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("count = %q, want %q", text, "Count: 3")
	}
}

// TestFreezeStateSurvivesDisconnect increments the counter,
// disconnects the network (simulating a transport drop), reconnects,
// and verifies the counter value is preserved through the freeze/thaw
// cycle.
func TestFreezeStateSurvivesDisconnect(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/freeze/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Increment to 5.
	btn := page.Locator("[data-tether-click='freeze.increment']")
	for range 5 {
		if err := btn.Click(); err != nil {
			t.Fatalf("click: %v", err)
		}
	}

	count := page.Locator("[data-tether-key='count']")
	if err := expect(count).ToHaveText("Count: 5"); err != nil {
		text, _ := count.TextContent()
		t.Fatalf("count before disconnect = %q, want %q", text, "Count: 5")
	}

	// Close the transport via the DevMode test hook. This triggers a
	// clean disconnect on the server, which freezes the session.
	if _, err := page.Evaluate("Tether.disconnect()"); err != nil {
		t.Fatalf("disconnect: %v", err)
	}

	// Wait for the client to detect the disconnect.
	disconnected := page.Locator("[data-tether-state='disconnected']")
	if err := disconnected.WaitFor(pw.LocatorWaitForOptions{
		Timeout: pw.Float(5000),
	}); err != nil {
		t.Fatalf("did not disconnect: %v", err)
	}

	// The client auto-reconnects - wait for the connection to restore.

	waitForConnected(t, page)

	// The counter should still be 5 after the freeze/thaw cycle.
	if err := expect(count).ToHaveText("Count: 5"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("count after reconnect = %q, want %q", text, "Count: 5")
	}
}
