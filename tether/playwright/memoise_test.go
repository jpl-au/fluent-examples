package playwright_test

import (
	"strconv"
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestMemoisedPageRenders verifies the memoisation demo loads with
// the initial table and counter.
func TestMemoisedPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memoise/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	tbl := page.Locator("#memoise-table")
	if err := expect(tbl).ToBeVisible(); err != nil {
		t.Fatalf("table not visible: %v", err)
	}

	count := page.Locator("#memoise-count")
	if err := expect(count).ToHaveText("0"); err != nil {
		t.Fatalf("counter should start at 0: %v", err)
	}
}

// TestMemoisedIncrementDoesNotChangeTable verifies that clicking the
// counter button updates the count but does not re-render the table.
// The table row count stays the same (memoiser hit - skipped).
func TestMemoisedIncrementDoesNotChangeTable(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memoise/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Count initial rows.
	rows := page.Locator("#memoise-table tbody tr")
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
	count := page.Locator("#memoise-count")
	if err := expect(count).ToHaveText("1"); err != nil {
		t.Fatalf("counter should be 1: %v", err)
	}

	// Table row count should be unchanged.
	afterCount, err := rows.Count()
	if err != nil {
		t.Fatalf("count rows after: %v", err)
	}
	if afterCount != initialCount {
		t.Errorf("table rows changed from %d to %d after increment (memoiser should skip)", initialCount, afterCount)
	}
}

// TestMemoisedAddItemUpdatesTable verifies that clicking Add Item
// adds a new row to the table (memoiser miss - re-rendered).
func TestMemoisedAddItemUpdatesTable(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memoise/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Count initial rows.
	rows := page.Locator("#memoise-table tbody tr")
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

// TestMemoisedRealtimePageRenders verifies the memoised real-time
// dashboard loads and the chart containers are present.
func TestMemoisedRealtimePageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/memoise/realtime/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("System Monitor")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}

	for _, id := range []string{"memoise-cpu", "memoise-heap", "memoise-goroutines"} {
		chart := page.Locator("#" + id)
		if err := expect(chart).ToBeVisible(); err != nil {
			t.Errorf("chart %q not visible: %v", id, err)
		}
	}
}
