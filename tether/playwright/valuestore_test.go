package playwright_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestValuestorePageRenders verifies the value store page loads and
// the Update demo card is visible.
func TestValuestorePageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/valuestore/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("Update (Read-Modify-Write)")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestValuestoreIncrement clicks Increment and verifies the shared
// counter increases by one.
func TestValuestoreIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/valuestore/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Read the initial count before clicking.
	count := page.Locator("#update-count")
	before := readCountText(t, count)

	btn := page.Locator("[data-tether-event-click='value.increment']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	expected := fmt.Sprintf("Count: %d", before+1)
	if err := expect(count).ToHaveText(expected); err != nil {
		text, _ := count.TextContent()
		t.Errorf("count = %q, want %q", text, expected)
	}
}

// TestValuestoreReset increments the counter then resets it to zero
// via the Store (Direct Set) demo.
func TestValuestoreReset(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/valuestore/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Increment first to ensure the counter is non-zero.
	count := page.Locator("#store-count")
	before := readCountText(t, count)
	inc := page.Locator("[data-tether-event-click='value.increment']")
	if err := inc.Click(); err != nil {
		t.Fatalf("increment: %v", err)
	}

	// Wait for the increment to register before resetting.
	if err := expect(count).ToHaveText(fmt.Sprintf("Count: %d", before+1)); err != nil {
		t.Fatalf("counter did not increment before reset: %v", err)
	}

	reset := page.Locator("[data-tether-event-click='value.reset']")
	if err := reset.Click(); err != nil {
		t.Fatalf("reset: %v", err)
	}

	if err := expect(count).ToHaveText("Count: 0"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("count = %q, want Count: 0", text)
	}
}

// TestValuestoreLocalIncrement verifies that clicking Increment Local
// updates the local counter. The shared counter may have a non-zero
// value from previous tests (it's a global tether.Value), so we only
// assert the local count changes.
func TestValuestoreLocalIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/valuestore/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-event-click='value.local-inc']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	local := page.Locator("#local-count")
	if err := expect(local).ToHaveText("Local: 1"); err != nil {
		text, _ := local.TextContent()
		t.Errorf("local = %q, want Local: 1", text)
	}
}

// readCountText extracts the numeric value from a "Count: N" or
// "Shared: N" text element. Fails the test if the value cannot be read.
func readCountText(t *testing.T, locator pw.Locator) int {
	t.Helper()
	text, err := locator.TextContent()
	if err != nil {
		t.Fatalf("read count: %v", err)
	}
	// Extract the number after the last space or colon.
	text = strings.TrimSpace(text)
	parts := strings.Fields(text)
	if len(parts) == 0 {
		t.Fatal("count text is empty")
	}
	n, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		t.Fatalf("invalid count %q: %v", text, err)
	}
	return n
}
