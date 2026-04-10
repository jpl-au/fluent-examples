package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestHealthEndpoint verifies the server is up and the health
// endpoint responds.
func TestHealthEndpoint(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/health")
	if err != nil {
		t.Fatalf("goto /health: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("health status = %d, want 200", resp.Status())
	}
}

// TestHTTPPageRenders verifies a stateless page (tether.Stateless) loads.
func TestHTTPPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}
}

// TestHTTPOverviewContent verifies the root page contains the
// Welcome panel and table of contents links.
func TestHTTPOverviewContent(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The Welcome panel heading is rendered as demo-title in the main
	// content area. Scope the locator to avoid matching sidebar nav.
	welcome := page.Locator("h3.demo-title", pw.PageLocatorOptions{HasText: "Welcome"})
	if err := expect(welcome).ToBeVisible(); err != nil {
		t.Errorf("Welcome heading not visible: %v", err)
	}

	// The overview page should contain links to feature sections.
	// The toc renders links with class toc-title; scope to that to
	// avoid matching the same link in the sidebar nav.
	eventsLink := page.Locator("a.toc-title", pw.PageLocatorOptions{HasText: "Events & Forms"})
	if err := expect(eventsLink).ToBeVisible(); err != nil {
		t.Errorf("Events & Forms link not visible: %v", err)
	}
}

// TestHTTPNotFoundPage navigates to a non-existent URL under the
// root handler and verifies the 404 page content renders. The HTTP
// section uses router.NotFound to render a custom empty state.
func TestHTTPNotFoundPage(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	// Navigate to a path that does not match any registered route
	// but falls through to the catch-all "/" handler.
	_, err := page.Goto(srv + "/this-page-does-not-exist")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The NotFound renderer shows "Page not found" as the title.
	title := page.GetByText("Page not found")
	if err := expect(title).ToBeVisible(); err != nil {
		t.Errorf("404 title not visible: %v", err)
	}

	// The description text should also be present.
	desc := page.GetByText("The page you're looking for doesn't exist.")
	if err := expect(desc).ToBeVisible(); err != nil {
		t.Errorf("404 description not visible: %v", err)
	}

	// The back link to the overview should be present.
	backLink := page.GetByText("Back to overview")
	if err := expect(backLink).ToBeVisible(); err != nil {
		t.Errorf("back to overview link not visible: %v", err)
	}
}
