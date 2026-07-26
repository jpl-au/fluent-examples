package playwright_test

import (
	"strconv"
	"testing"
)

// TestMorphPageRenders verifies the full-page morph demo loads.
func TestMorphPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/morph/")
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

	_, err := page.Goto(srv + "/morph/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.Locator("[data-tether-event-click='morph.increment']")
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

// TestMorphDecrement increments once then decrements, verifying the
// counter returns to zero via full-page morph.
func TestMorphDecrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/morph/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	inc := page.Locator("[data-tether-event-click='morph.increment']")
	if err := inc.Click(); err != nil {
		t.Fatalf("click increment: %v", err)
	}

	after1 := page.GetByText("Count: 1")
	if err := expect(after1).ToBeVisible(); err != nil {
		t.Fatalf("counter did not reach 1: %v", err)
	}

	dec := page.Locator("[data-tether-event-click='morph.decrement']")
	if err := dec.Click(); err != nil {
		t.Fatalf("click decrement: %v", err)
	}

	after0 := page.GetByText("Count: 0")
	if err := expect(after0).ToBeVisible(); err != nil {
		t.Errorf("counter did not decrement to 0: %v", err)
	}
}

// TestMorphDecrementLowerBound verifies that decrementing at zero
// does not go negative - the counter stays at 0.
func TestMorphDecrementLowerBound(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/morph/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	dec := page.Locator("[data-tether-event-click='morph.decrement']")
	if err := dec.Click(); err != nil {
		t.Fatalf("click decrement: %v", err)
	}

	counter := page.GetByText("Count: 0")
	if err := expect(counter).ToBeVisible(); err != nil {
		t.Errorf("counter went below 0: %v", err)
	}
}

// TestMorphMultipleIncrements clicks + several times and verifies
// the counter round-trips correctly through full-page morphs.
func TestMorphMultipleIncrements(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/morph/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.Locator("[data-tether-event-click='morph.increment']")

	for i := 1; i <= 5; i++ {
		if err := btn.Click(); err != nil {
			t.Fatalf("click %d: %v", i, err)
		}

		expected := page.GetByText("Count: " + strconv.Itoa(i))
		if err := expect(expected).ToBeVisible(); err != nil {
			t.Fatalf("counter did not reach %d: %v", i, err)
		}
	}
}
