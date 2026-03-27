package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestWindowingPageRenders verifies the windowing demo loads with
// the table and navigation buttons visible.
func TestWindowingPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/windowing/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// First row should be visible.
	row := page.Locator("#wrow-1")
	if err := expect(row).ToBeVisible(); err != nil {
		t.Fatalf("first row not visible: %v", err)
	}

	// Position indicator should show the initial range.
	position := page.Locator("[data-tether-key='position']")
	if err := expect(position).ToContainText("Showing rows 1"); err != nil {
		t.Fatalf("position indicator missing: %v", err)
	}
}

// TestWindowingNextPage verifies that clicking Next advances the
// visible window and shows different rows.
func TestWindowingNextPage(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/windowing/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Click Next.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Next"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// Row 31 should now be visible (second page).
	row := page.Locator("#wrow-31")
	if err := expect(row).ToBeAttached(); err != nil {
		t.Fatalf("row 31 not in DOM after Next: %v", err)
	}

	// Row 1 should no longer be in the DOM.
	row1 := page.Locator("#wrow-1")
	if err := expect(row1).Not().ToBeAttached(); err != nil {
		t.Errorf("row 1 should not be in DOM after Next: %v", err)
	}
}

// TestWindowingURLPagination verifies that navigating directly to
// a page via the URL query parameter shows the correct rows.
func TestWindowingURLPagination(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	// Navigate directly to page 5 (rows 121-150).
	_, err := page.Goto(srv + "/windowing/?page=5")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Row 121 should be visible (first row of page 5).
	row := page.Locator("#wrow-121")
	if err := expect(row).ToBeVisible(); err != nil {
		t.Fatalf("row 121 not visible on page 5: %v", err)
	}

	// Position indicator should reflect page 5.
	position := page.Locator("[data-tether-key='position']")
	if err := expect(position).ToContainText("Showing rows 121"); err != nil {
		t.Fatalf("position should show rows 121: %v", err)
	}

	// Row 1 should not be in the DOM.
	row1 := page.Locator("#wrow-1")
	if err := expect(row1).Not().ToBeAttached(); err != nil {
		t.Errorf("row 1 should not be in DOM on page 5: %v", err)
	}
}
