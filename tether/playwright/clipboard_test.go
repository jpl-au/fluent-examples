package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestClipboardPageRenders verifies the clipboard demo loads and
// shows the copy button and source text.
func TestClipboardPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/clipboard")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	source := page.Locator("#copy-source")
	if err := expect(source).ToBeVisible(); err != nil {
		t.Fatalf("copy source not visible: %v", err)
	}
	if err := expect(source).ToContainText("tether-secret-key-abc123"); err != nil {
		t.Fatalf("copy source text mismatch: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copy"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("copy button not visible: %v", err)
	}
}

// TestClipboardCopy clicks the copy button and verifies the clipboard
// contains the expected text. Clipboard read requires the
// "clipboard-read" permission.
func TestClipboardCopy(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t, WithPermissions(srv, "clipboard-read", "clipboard-write"))
	defer cleanup()

	_, err := page.Goto(srv + "/clipboard")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copy"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click copy: %v", err)
	}

	// Read the clipboard via the browser API.
	result, err := page.Evaluate("() => navigator.clipboard.readText()")
	if err != nil {
		t.Fatalf("clipboard read: %v", err)
	}

	text, ok := result.(string)
	if !ok || text != "tether-secret-key-abc123" {
		t.Errorf("clipboard = %q, want %q", text, "tether-secret-key-abc123")
	}
}
