package playwright_test

import "testing"

// TestNavigationPageRenders verifies the navigation page loads.
func TestNavigationPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/navigation/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("Client-Side Navigation")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestNavigationLinkToTarget clicks a bind.Link and verifies the
// browser navigates to the target page without a full reload.
func TestNavigationLinkToTarget(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/navigation/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	link := page.GetByText("Go to Target Page")
	if err := link.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The target page should render with a "Back to Navigation Demos" link.
	back := page.GetByText("Back to Navigation Demos")
	if err := expect(back).ToBeVisible(); err != nil {
		t.Errorf("target page not rendered: %v", err)
	}
}

// TestNavigationQueryParams clicks a link with query parameters and
// verifies OnNavigate extracts them into state.
func TestNavigationQueryParams(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/navigation/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	link := page.GetByText("Tab: settings, page 3")
	if err := link.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	result := page.Locator("[data-tether-key='query-params']")
	if err := expect(result).ToContainText("tab    = settings"); err != nil {
		t.Errorf("tab not extracted: %v", err)
	}
	if err := expect(result).ToContainText("page   = 3"); err != nil {
		t.Errorf("page not extracted: %v", err)
	}
}

// TestNavigationMultiValue clicks a link with repeated query keys
// and verifies Params.Strings collects them.
func TestNavigationMultiValue(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/navigation/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	link := page.GetByText("Tags: go, web, sse")
	if err := link.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	result := page.Locator("[data-tether-key='multi-value-params']")
	if err := expect(result).ToContainText("go, web, sse"); err != nil {
		t.Errorf("tags not extracted: %v", err)
	}
}

// TestNavigationTypedParams clicks a link with boolean and float
// query values and verifies BoolDefault and Float64Default parse them.
func TestNavigationTypedParams(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/navigation/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	link := page.GetByText("Active, price 9.99")
	if err := link.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	result := page.Locator("[data-tether-key='typed-params']")
	if err := expect(result).ToContainText("active = true"); err != nil {
		t.Errorf("active not parsed: %v", err)
	}
	if err := expect(result).ToContainText("price  = 9.99"); err != nil {
		t.Errorf("price not parsed: %v", err)
	}
}

// TestNavigationNumericMultiValue clicks a link with repeated numeric
// query keys and verifies Ints and Float64s collect them.
func TestNavigationNumericMultiValue(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/navigation/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	link := page.GetByText("Quantities: 1, 2, 5")
	if err := link.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	result := page.Locator("[data-tether-key='numeric-multi-value-params']")
	if err := expect(result).ToContainText("1, 2, 5"); err != nil {
		t.Errorf("quantities not parsed: %v", err)
	}
}
