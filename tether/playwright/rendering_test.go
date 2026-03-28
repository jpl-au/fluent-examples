package playwright_test

import "testing"

// TestRenderingPageRenders verifies the rendering demo page loads.
func TestRenderingPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/rendering/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}

	waitForConnected(t, page)

	heading := page.GetByText("Dynamic Keys")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Errorf("heading not visible: %v", err)
	}
}

// TestRenderingCounterIncrement clicks + and verifies the counter
// updates via a stateless HTTP POST round-trip.
func TestRenderingCounterIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/rendering/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The + button in the Dynamic Keys section.
	btn := page.Locator("[data-tether-click='rendering.increment']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	counter := page.Locator("[data-tether-key='rendering-counter']")
	if err := expect(counter).ToContainText("1"); err != nil {
		text, _ := counter.TextContent()
		t.Errorf("counter = %q, want to contain 1", text)
	}
}

// TestRenderingAddItem clicks "Add Item" and verifies the list grows.
func TestRenderingAddItem(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/rendering/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='rendering.add-item']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// Individual items have Dynamic keys item-0, item-1, etc.
	item := page.Locator("[data-tether-key='item-0']")
	if err := expect(item).ToContainText("Item 1"); err != nil {
		t.Errorf("item not added: %v", err)
	}
}
