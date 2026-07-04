package playwright_test

import (
	"testing"
	"time"

	pw "github.com/playwright-community/playwright-go"
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
	row := page.Locator("[data-fluent-key='counter-0']")
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
	row := page.Locator("[data-fluent-key='counter-0']")
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

// TestPatchMultipleCountersIncrement verifies that the background
// timer cycles through multiple counters, not just counter-0. After
// enough ticks, counter-1 should also have a non-zero value.
func TestPatchMultipleCountersIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/patch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The ticker increments one counter every 500ms, cycling through
	// all 20. After 1.5s, counters 0, 1, and 2 should each have been
	// hit at least once.
	time.Sleep(2 * time.Second)

	row1 := page.Locator("[data-fluent-key='counter-1']")
	text, err := row1.InnerText()
	if err != nil {
		t.Fatalf("read counter-1: %v", err)
	}

	if text == "Counter 10" {
		t.Error("counter-1 should have incremented after 2 seconds")
	}
}

// TestPatchResetAll clicks the Reset All button and verifies counters
// return to zero. The background ticker increments one counter every
// 500ms, so after a reset the low-index counters can race with the
// ticker and re-increment. To avoid the race, this test waits for
// several counters to increment, then clicks Reset and checks a
// counter in the upper range (counter-15) that the ticker cannot
// possibly reach for several seconds after the reset fires.
func TestPatchResetAll(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/patch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Wait long enough for counters 0-15 to have incremented at least
	// once. 8 seconds at 500ms per tick = 16 ticks, so by this point
	// counter-15 will have value 1 or more.
	time.Sleep(8 * time.Second)

	// Confirm counter-15 has incremented so we know reset will do
	// something observable.
	row15 := page.Locator("[data-fluent-key='counter-15']")
	before, _ := row15.InnerText()
	if before == "Counter 150" {
		t.Fatalf("counter-15 before reset = %q, expected non-zero after 8 seconds", before)
	}

	// Click Reset All. This triggers a full Update that zeroes the
	// entire counters array and re-renders every Dynamic row.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Reset All"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click reset: %v", err)
	}

	// After reset, counter-15 must show 0. The ticker is at index ~16
	// so it won't touch counter-15 again for ~2.5s (until index wraps
	// past 15 on the second pass), giving us a safe window to verify.
	if err := expect(row15).ToHaveText("Counter 150"); err != nil {
		text, _ := row15.InnerText()
		t.Errorf("counter-15 after reset = %q, want %q", text, "Counter 150")
	}
}
