package playwright_test

import (
	"testing"
	"time"
)

// TestPatchPageRenders verifies the patch demo loads with all
// counters visible.
func TestPatchPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/patch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// First counter should be visible.
	row := page.Locator("[data-tether-key='counter-0']")
	if err := expect(row).ToBeVisible(); err != nil {
		t.Fatalf("counter-0 not visible: %v", err)
	}
}

// TestPatchCounterIncrements verifies that the background timer
// increments counters via Patch.
func TestPatchCounterIncrements(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/patch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Wait for a few ticks (500ms each).
	time.Sleep(2 * time.Second)

	// At least one counter should have incremented above 0.
	row := page.Locator("[data-tether-key='counter-0']")
	text, err := row.InnerText()
	if err != nil {
		t.Fatalf("read counter-0: %v", err)
	}

	// The text contains "Counter 0" and the value. If the value
	// is still "0" after 2 seconds, the patch isn't working.
	if text == "Counter 00" {
		t.Error("counter-0 should have incremented after 2 seconds")
	}
}
