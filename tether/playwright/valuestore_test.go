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
	count := page.Locator("[data-tether-key='update-count']")
	before := readCountText(t, count)

	btn := page.Locator("[data-tether-click='value.increment']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	expected := fmt.Sprintf("Count: %d", before+1)
	if err := expect(count).ToContainText(strconv.Itoa(before + 1)); err != nil {
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
	inc := page.Locator("[data-tether-click='value.increment']")
	if err := inc.Click(); err != nil {
		t.Fatalf("increment: %v", err)
	}

	// Wait for the increment to register before resetting.
	count := page.Locator("[data-tether-key='store-count']")
	before := readCountText(t, count)
	if before == 0 {
		// The watcher may not have propagated yet - wait briefly.
		page.WaitForTimeout(500)
	}

	reset := page.Locator("[data-tether-click='value.reset']")
	if err := reset.Click(); err != nil {
		t.Fatalf("reset: %v", err)
	}

	if err := expect(count).ToContainText("0"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("count = %q, want to contain 0", text)
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

	btn := page.Locator("[data-tether-click='value.local-inc']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	local := page.Locator("[data-tether-key='local-count']")
	if err := expect(local).ToContainText("1"); err != nil {
		text, _ := local.TextContent()
		t.Errorf("local = %q, want to contain 1", text)
	}
}

// readCountText extracts the numeric value from a "Count: N" or
// "Shared: N" text element. Returns 0 if unparseable.
func readCountText(t *testing.T, locator pw.Locator) int {
	t.Helper()
	text, err := locator.TextContent()
	if err != nil {
		return 0
	}
	// Extract the number after the last space or colon.
	text = strings.TrimSpace(text)
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(parts[len(parts)-1])
	return n
}
