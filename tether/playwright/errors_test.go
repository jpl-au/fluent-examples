package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestErrorsPageRenders verifies the error boundaries page loads.
func TestErrorsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/errors/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}

	// The error boundary catches a deliberate panic and shows
	// the fallback content.
	fallback := page.GetByText("Caught by error boundary")
	if err := expect(fallback).ToBeVisible(); err != nil {
		t.Errorf("error boundary fallback not visible: %v", err)
	}
}

// TestErrorsFallbackMessage verifies the full fallback message
// rendered by the error boundary, including the recovery text.
func TestErrorsFallbackMessage(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/errors/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The fallback is rendered by panel.SignalSuccess which wraps
	// the full recovery message.
	fallback := page.GetByText("component recovered gracefully")
	if err := expect(fallback).ToBeVisible(); err != nil {
		t.Errorf("full fallback message not visible: %v", err)
	}
}

// TestErrorsPanelRendersAlongsideBoundary verifies that the
// containing panel card still renders correctly even though the
// inner component panicked. The error boundary isolates the failure
// so the rest of the page is unaffected.
func TestErrorsPanelRendersAlongsideBoundary(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/errors/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The panel card title "Error Boundary" should still render.
	// Use the demo-title class to scope to the panel heading, since
	// "Error Boundaries" (plural) also appears in the sidebar nav.
	title := page.Locator("h3.demo-title", pw.PageLocatorOptions{HasText: "Error Boundary"})
	if err := expect(title).ToBeVisible(); err != nil {
		t.Errorf("panel title not visible: %v", err)
	}

	// The panel description text should also be visible, proving
	// the panel rendered correctly around the caught panic.
	desc := page.GetByText("deliberately panics during render")
	if err := expect(desc).ToBeVisible(); err != nil {
		t.Errorf("panel description not visible: %v", err)
	}

	// The tether.Catch badge should be visible in the panel header.
	badge := page.Locator("span.api-label", pw.PageLocatorOptions{HasText: "tether.Catch"})
	if err := expect(badge).ToBeVisible(); err != nil {
		t.Errorf("API badge not visible: %v", err)
	}
}
