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

// TestRenderingCounterDecrement increments once then decrements,
// verifying the counter returns to zero via targeted patch.
func TestRenderingCounterDecrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/rendering/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	inc := page.Locator("[data-tether-click='rendering.increment']")
	if err := inc.Click(); err != nil {
		t.Fatalf("click increment: %v", err)
	}

	counter := page.Locator("[data-tether-key='rendering-counter']")
	if err := expect(counter).ToContainText("1"); err != nil {
		t.Fatalf("counter did not reach 1: %v", err)
	}

	dec := page.Locator("[data-tether-click='rendering.decrement']")
	if err := dec.Click(); err != nil {
		t.Fatalf("click decrement: %v", err)
	}

	if err := expect(counter).ToContainText("0"); err != nil {
		text, _ := counter.TextContent()
		t.Errorf("counter = %q, want to contain 0", text)
	}
}

// TestRenderingDecrementLowerBound verifies that decrementing at
// zero does not go negative - the counter stays at 0.
func TestRenderingDecrementLowerBound(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/rendering/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	dec := page.Locator("[data-tether-click='rendering.decrement']")
	if err := dec.Click(); err != nil {
		t.Fatalf("click decrement: %v", err)
	}

	counter := page.Locator("[data-tether-key='rendering-counter']")
	if err := expect(counter).ToContainText("0"); err != nil {
		text, _ := counter.TextContent()
		t.Errorf("counter = %q, want to contain 0 after decrement at lower bound", text)
	}
}

// TestRenderingRemoveItem adds two items then removes one, verifying
// the list shrinks correctly.
func TestRenderingRemoveItem(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/rendering/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	addBtn := page.Locator("[data-tether-click='rendering.add-item']")

	// Add two items.
	if err := addBtn.Click(); err != nil {
		t.Fatalf("click add 1: %v", err)
	}
	item0 := page.Locator("[data-tether-key='item-0']")
	if err := expect(item0).ToBeVisible(); err != nil {
		t.Fatalf("item-0 not visible after first add: %v", err)
	}

	if err := addBtn.Click(); err != nil {
		t.Fatalf("click add 2: %v", err)
	}
	item1 := page.Locator("[data-tether-key='item-1']")
	if err := expect(item1).ToBeVisible(); err != nil {
		t.Fatalf("item-1 not visible after second add: %v", err)
	}

	// Remove the last item.
	removeBtn := page.Locator("[data-tether-click='rendering.remove-item']")
	if err := removeBtn.Click(); err != nil {
		t.Fatalf("click remove: %v", err)
	}

	// item-0 should remain, item-1 should be gone.
	if err := expect(item0).ToBeVisible(); err != nil {
		t.Errorf("item-0 should still be visible: %v", err)
	}

	itemList := page.Locator("[data-tether-key='item-list']")
	if err := expect(itemList).Not().ToContainText("Item 2"); err != nil {
		t.Errorf("item-1 should have been removed: %v", err)
	}
}
