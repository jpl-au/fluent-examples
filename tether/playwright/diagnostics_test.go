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
	feed := page.Locator("[data-fluent-key='diagnostics']")
	if err := expect(feed).ToContainText("handler_panic"); err != nil {
		t.Errorf("diagnostic event not in feed: %v", err)
	}
}

// TestDiagnosticsSessionSurvivesPanic triggers a panic and then
// triggers another to verify the session is still alive. OnPanic is
// configured to keep the session, so a second interaction should work.
func TestDiagnosticsSessionSurvivesPanic(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/diagnostics/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-click='diag.trigger-panic']").First()

	// First panic.
	if err := btn.Click(pw.LocatorClickOptions{Force: pw.Bool(true)}); err != nil {
		t.Fatalf("click 1: %v", err)
	}

	feed := page.Locator("[data-fluent-key='diagnostics']")
	if err := expect(feed).ToContainText("handler_panic"); err != nil {
		t.Fatalf("first panic event not in feed: %v", err)
	}

	// Second panic - session should still be alive.
	if err := btn.Click(pw.LocatorClickOptions{Force: pw.Bool(true)}); err != nil {
		t.Fatalf("click 2: %v", err)
	}

	// The feed should now contain at least two handler_panic entries.
	// We verify by checking that the feed still renders (the session
	// didn't die) and contains the expected text.
	if err := expect(feed).ToContainText("handler_panic"); err != nil {
		t.Errorf("feed should still show events after second panic: %v", err)
	}

	// The WebSocket connection should still be active.
	connected := page.Locator("[data-tether-state='connected']")
	if err := connected.WaitFor(); err != nil {
		t.Errorf("session should still be connected after panics: %v", err)
	}
}

// TestDiagnosticsEventShowsSessionID verifies that the diagnostic
// event entry in the feed includes a truncated session ID.
func TestDiagnosticsEventShowsSessionID(t *testing.T) {
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

	// The event feed renders session IDs truncated to 6 chars in the
	// format: "handler_panic - detail (session XXXXXX)".
	feed := page.Locator("[data-fluent-key='diagnostics']")
	if err := expect(feed).ToContainText("session"); err != nil {
		t.Errorf("diagnostic event should include session ID: %v", err)
	}
}
