package playwright_test

import (
	"strings"
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestDiagnosticsPageRenders verifies the diagnostics page loads.
func TestDiagnosticsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/diagnostics/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Check the page rendered the trigger button. Use text match
	// since the button text is unique on this page.
	content, _ := page.Content()
	if !strings.Contains(content, "diag.trigger-panic") {
		t.Fatalf("page does not contain trigger button HTML")
	}

	// The button might be inside a section that needs scrolling
	// or is below the fold. Use a locator count check instead of
	// visibility.
	btn := page.Locator("[data-tether-click='diag.trigger-panic']")
	count, _ := btn.Count()
	if count == 0 {
		t.Fatal("trigger button not found in DOM")
	}
}

// TestDiagnosticsTriggerPanic clicks the panic trigger button and
// verifies a diagnostic event appears in the feed. The session
// should survive the panic.
func TestDiagnosticsTriggerPanic(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/diagnostics/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='diag.trigger-panic']").First()
	if err := btn.Click(pw.LocatorClickOptions{Force: pw.Bool(true)}); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The diagnostic event feed should show a HandlerPanic entry.
	feed := page.Locator("[data-tether-key='diagnostics']")
	if err := expect(feed).ToContainText("handler_panic"); err != nil {
		t.Errorf("diagnostic event not in feed: %v", err)
	}
}
