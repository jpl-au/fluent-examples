package playwright_test

import (
	"regexp"
	"testing"

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
	row := page.Locator("#counter-0")
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

	value := page.Locator("#counter-0 .data-value")
	if err := expect(value).ToHaveText(regexp.MustCompile(`^[1-9]\d*$`)); err != nil {
		t.Errorf("counter-0 did not increment: %v", err)
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

	for _, selector := range []string{"#counter-0 .data-value", "#counter-1 .data-value"} {
		if err := expect(page.Locator(selector)).ToHaveText(regexp.MustCompile(`^[1-9]\d*$`)); err != nil {
			t.Errorf("%s did not increment: %v", selector, err)
		}
	}
}

// TestPatchResetAll clicks the Reset All button and verifies counters
// return to zero. Wait for the first two counters to increment, then
// reset them. The ticker advances independently of reset and will not
// revisit them until it has cycled through the remaining 18 counters.
func TestPatchResetAll(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/patch/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	values := []pw.Locator{
		page.Locator("#counter-0 .data-value"),
		page.Locator("#counter-1 .data-value"),
	}
	for i, value := range values {
		if err := expect(value).ToHaveText(regexp.MustCompile(`^[1-9]\d*$`)); err != nil {
			t.Fatalf("counter-%d did not increment before reset: %v", i, err)
		}
	}

	// Click Reset All. This triggers a full Update that zeroes the
	// entire counters array and re-renders every Dynamic row.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Reset All"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click reset: %v", err)
	}

	for i, value := range values {
		if err := expect(value).ToHaveText("0"); err != nil {
			t.Errorf("counter-%d did not reset to zero: %v", i, err)
		}
	}
}
