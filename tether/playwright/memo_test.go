package playwright_test

import (
	"strconv"
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestMemoPageRenders verifies the memo demo loads with the initial
// table and counter.
func TestMemoPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memo/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	tbl := page.Locator("#memo-table")
	if err := expect(tbl).ToBeVisible(); err != nil {
		t.Fatalf("table not visible: %v", err)
	}

	count := page.Locator("#memo-count")
	if err := expect(count).ToHaveText("0"); err != nil {
		t.Fatalf("counter should start at 0: %v", err)
	}
}

// TestMemoIncrementDoesNotChangeTable verifies that clicking the
// counter button updates the count but does not re-render the table.
// The table row count stays the same (memo hit - skipped).
func TestMemoIncrementDoesNotChangeTable(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memo/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Count initial rows.
	rows := page.Locator("#memo-table tbody tr")
	initialCount, err := rows.Count()
	if err != nil {
		t.Fatalf("count rows: %v", err)
	}

	// Click increment button (use role to avoid matching description text).
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Increment Counter"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// Counter should update.
	count := page.Locator("#memo-count")
	if err := expect(count).ToHaveText("1"); err != nil {
		t.Fatalf("counter should be 1: %v", err)
	}

	// Table row count should be unchanged.
	afterCount, err := rows.Count()
	if err != nil {
		t.Fatalf("count rows after: %v", err)
	}
	if afterCount != initialCount {
		t.Errorf("table rows changed from %d to %d after increment (memo should skip)", initialCount, afterCount)
	}
}

// TestMemoAddItemUpdatesTable verifies that clicking Add Item adds a
// new row to the table (memo miss - re-rendered).
func TestMemoAddItemUpdatesTable(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memo/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Count initial rows.
	rows := page.Locator("#memo-table tbody tr")
	initialCount, err := rows.Count()
	if err != nil {
		t.Fatalf("count rows: %v", err)
	}

	// Click Add Item button (use role to avoid matching description text).
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Add Item"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// New row should appear.
	newRow := page.Locator("#row-" + itoa(initialCount+1))
	if err := expect(newRow).ToBeAttached(); err != nil {
		t.Fatalf("new row not in DOM: %v", err)
	}

	// Table should have one more row.
	afterCount, err := rows.Count()
	if err != nil {
		t.Fatalf("count rows after: %v", err)
	}
	if afterCount != initialCount+1 {
		t.Errorf("expected %d rows after add, got %d", initialCount+1, afterCount)
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

// TestMemoRealtimePageRenders verifies the memoised real-time
// dashboard loads and the chart containers are present.
func TestMemoRealtimePageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memo/realtime/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("System Monitor")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}

	for _, id := range []string{"memocpu", "memoheap", "memogoroutines"} {
		chart := page.Locator("#" + id)
		if err := expect(chart).ToBeVisible(); err != nil {
			t.Errorf("chart %q not visible: %v", id, err)
		}
	}
}
