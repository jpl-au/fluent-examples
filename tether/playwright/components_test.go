package playwright_test

import "testing"

// TestComponentsPageRenders verifies the components page loads and
// the counter components are interactive.
func TestComponentsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/components/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByText("+").First()
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("increment button not visible: %v", err)
	}
}

// TestComponentsLikesIncrement clicks the Likes counter + button
// and verifies the count updates.
func TestComponentsLikesIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/components/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The Likes counter is in the likes-section Dynamic container.
	section := page.Locator("[data-tether-key='likes-section']")
	btn := section.GetByText("+")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	count := section.Locator("text=/^Count:/")
	if err := expect(count).ToHaveText("Count: 2"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("likes count = %q, want %q", text, "Count: 2")
	}
}

// TestComponentsLikesDecrement verifies the − button decreases the
// count and floors at zero.
func TestComponentsLikesDecrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/components/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	section := page.Locator("[data-tether-key='likes-section']")
	plus := section.GetByText("+")
	minus := section.GetByText("−")

	// Increment twice, then decrement once.
	if err := plus.Click(); err != nil {
		t.Fatalf("click +: %v", err)
	}
	if err := plus.Click(); err != nil {
		t.Fatalf("click +: %v", err)
	}
	if err := minus.Click(); err != nil {
		t.Fatalf("click −: %v", err)
	}

	count := section.Locator("text=/^Count:/")
	if err := expect(count).ToHaveText("Count: 1"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("likes count = %q, want %q", text, "Count: 1")
	}

	// Decrement to zero - should not go negative.
	if err := minus.Click(); err != nil {
		t.Fatalf("click −: %v", err)
	}
	if err := minus.Click(); err != nil {
		t.Fatalf("click − again: %v", err)
	}

	if err := expect(count).ToHaveText("Count: 0"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("likes count = %q, want %q (should not go negative)", text, "Count: 0")
	}
}

// TestComponentsReset clicks Reset and verifies the counter returns
// to zero.
func TestComponentsReset(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/components/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	section := page.Locator("[data-tether-key='likes-section']")
	plus := section.GetByText("+")
	reset := section.GetByText("Reset")

	if err := plus.Click(); err != nil {
		t.Fatalf("click +: %v", err)
	}
	if err := plus.Click(); err != nil {
		t.Fatalf("click +: %v", err)
	}
	if err := reset.Click(); err != nil {
		t.Fatalf("click reset: %v", err)
	}

	count := section.Locator("text=/^Count:/")
	if err := expect(count).ToHaveText("Count: 0"); err != nil {
		text, _ := count.TextContent()
		t.Errorf("likes count = %q, want %q", text, "Count: 0")
	}
}

// TestComponentsIndependent verifies that Likes and Stars counters
// operate independently - incrementing one does not affect the other.
func TestComponentsIndependent(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/components/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	likes := page.Locator("[data-tether-key='likes-section']")
	stars := page.Locator("[data-tether-key='stars-section']")

	// Increment Likes three times.
	likesPlus := likes.GetByText("+")
	for range 3 {
		if err := likesPlus.Click(); err != nil {
			t.Fatalf("click likes +: %v", err)
		}
	}

	// Increment Stars once.
	starsPlus := stars.GetByText("+")
	if err := starsPlus.Click(); err != nil {
		t.Fatalf("click stars +: %v", err)
	}

	// Verify they're independent.
	likesCount := likes.Locator("text=/^Count:/")
	starsCount := stars.Locator("text=/^Count:/")

	if err := expect(likesCount).ToHaveText("Count: 3"); err != nil {
		text, _ := likesCount.TextContent()
		t.Errorf("likes = %q, want Count: 3", text)
	}
	if err := expect(starsCount).ToHaveText("Count: 1"); err != nil {
		text, _ := starsCount.TextContent()
		t.Errorf("stars = %q, want Count: 1", text)
	}
}
