package playwright_test

import "testing"

// TestMorphPageRenders verifies the full-page morph demo loads.
func TestMorphPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/morph")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}

	heading := page.GetByText("Counter Without Dynamic Keys")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestMorphIncrement clicks + and verifies the counter updates via
// the full-page morph fallback (no Dynamic keys on this page).
func TestMorphIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/morph")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.Locator("[data-tether-click='morph.increment']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// No Dynamic key - the morph updates the entire page. Verify the
	// counter text changed by looking for "Count: 1" anywhere on the
	// page.
	result := page.GetByText("Count: 1")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("counter did not update via full-page morph: %v", err)
	}
}
