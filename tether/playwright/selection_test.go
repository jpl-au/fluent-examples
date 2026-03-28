package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestSelectionPageRenders verifies the selection demo loads.
func TestSelectionPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/selection/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Show Selected"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("button not visible: %v", err)
	}
}

// TestSelectionClickSelect clicks an item and verifies the server
// receives its ID when the action button is clicked.
func TestSelectionClickSelect(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/selection/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click the first item (Alpha, id=1).
	item := page.Locator("[data-tether-data-id='1']")
	if err := item.Click(); err != nil {
		t.Fatalf("click item: %v", err)
	}

	// Item should have tether-selected class.
	if err := expect(item).ToHaveClass("list-item tether-selected"); err != nil {
		cls, _ := item.GetAttribute("class")
		t.Errorf("item class = %q, want tether-selected: %v", cls, err)
	}

	// Click Show Selected.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Show Selected"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click button: %v", err)
	}

	result := page.GetByText("Selected 1 items: Alpha")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("selection result not visible: %v", err)
	}
}

// TestSelectionNone clicks the button without selecting anything and
// verifies the "no items selected" result.
func TestSelectionNone(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/selection/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Show Selected"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click button: %v", err)
	}

	result := page.GetByText("No items selected")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("empty selection result not visible: %v", err)
	}
}

// TestSelectionCtrlClick ctrl+clicks two items and verifies both
// are included in the selected set.
func TestSelectionCtrlClick(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/selection/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Click item 1.
	item1 := page.Locator("[data-tether-data-id='1']")
	if err := item1.Click(); err != nil {
		t.Fatalf("click item 1: %v", err)
	}

	// Ctrl+click item 3.
	item3 := page.Locator("[data-tether-data-id='3']")
	if err := item3.Click(pw.LocatorClickOptions{
		Modifiers: []pw.KeyboardModifier{"Control"},
	}); err != nil {
		t.Fatalf("ctrl+click item 3: %v", err)
	}

	// Click Show Selected.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Show Selected"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click button: %v", err)
	}

	result := page.GetByText("Selected 2 items:")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("ctrl+click selection result not visible: %v", err)
	}
}
