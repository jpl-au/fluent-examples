package playwright_test

import "testing"

// TestScrollPageRenders verifies the scroll demo loads.
func TestScrollPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/scroll/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByText("Scroll to Target")
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("scroll button not visible: %v", err)
	}
}

// TestScrollToClient clicks the client-side ScrollTo button and
// verifies the target element is scrolled into the viewport.
func TestScrollToClient(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/scroll/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Target starts off-screen.
	target := page.Locator("#scroll-target")
	if err := expect(target).ToBeAttached(); err != nil {
		t.Fatalf("scroll target not in DOM: %v", err)
	}

	// Click scroll button.
	btn := page.GetByText("Scroll to Target")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// Verify target is now visible (scrolled into view).
	if err := expect(target).ToBeVisible(); err != nil {
		t.Errorf("scroll target not visible after ScrollTo: %v", err)
	}
}

// TestScrollToServer clicks the server-side ScrollTo button and
// verifies the target is scrolled into view via the WebSocket.
func TestScrollToServer(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/scroll/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByText("Server Scroll")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	target := page.Locator("#scroll-target")
	if err := expect(target).ToBeVisible(); err != nil {
		t.Errorf("scroll target not visible after server ScrollTo: %v", err)
	}
}

// TestPreserveScroll verifies that scroll position survives a re-render.
func TestPreserveScroll(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/scroll/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Scroll the preserve-scroll list down.
	list := page.Locator("[data-tether-preserve-scroll]")
	if err := expect(list).ToBeVisible(); err != nil {
		t.Fatalf("preserve list not visible: %v", err)
	}

	// Scroll to bottom of the list.
	scrollInfo, err := list.Evaluate("el => { el.scrollTop = el.scrollHeight; return {scrollTop: el.scrollTop, scrollHeight: el.scrollHeight, clientHeight: el.clientHeight}; }", nil)
	if err != nil {
		t.Fatalf("scroll list: %v", err)
	}
	info := scrollInfo.(map[string]any)
	t.Logf("Scroll Info: %+v", info)

	// Read scroll position before adding items.
	beforeRaw, err := list.Evaluate("el => el.scrollTop", nil)
	if err != nil {
		t.Fatalf("read scrollTop: %v", err)
	}
	t.Logf("beforeRaw: %v (%T)", beforeRaw, beforeRaw)
	before := jsNumber(beforeRaw)
	if before == 0 {
		t.Fatalf("scrollTop should be non-zero after scrolling (info: %+v, beforeRaw: %v %T)", info, beforeRaw, beforeRaw)
	}

	// Add items (triggers re-render).
	btn := page.GetByText("Add 5 Items")
	if err := btn.Click(); err != nil {
		t.Fatalf("click add: %v", err)
	}

	// Wait for the new item to appear.
	item15 := page.Locator("#item-15")
	if err := expect(item15).ToBeAttached(); err != nil {
		t.Fatalf("item 15 not in DOM after add: %v", err)
	}

	// Read scroll position after re-render.
	afterRaw, err := list.Evaluate("el => el.scrollTop", nil)
	if err != nil {
		t.Fatalf("read scrollTop after: %v", err)
	}
	after := jsNumber(afterRaw)

	// The scroll position should be preserved (not reset to 0).
	if after == 0 {
		t.Errorf("scrollTop reset to 0 after re-render, expected preserved position (was %.0f)", before)
	}
}
