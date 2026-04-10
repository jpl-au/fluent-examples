package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestClientActionsPageRenders verifies the client-side actions demo
// loads and shows the copy button and source text.
func TestClientActionsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/client-actions/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	source := page.Locator("#copy-source")
	if err := expect(source).ToBeVisible(); err != nil {
		t.Fatalf("copy source not visible: %v", err)
	}
	if err := expect(source).ToContainText("sk_live_abc123def456"); err != nil {
		t.Fatalf("copy source text mismatch: %v", err)
	}

	exact := true
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copy", Exact: &exact})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("copy button not visible: %v", err)
	}
}

// TestClientActionsCopy clicks the copy button and verifies the
// clipboard contains the expected text. Clipboard read requires the
// "clipboard-read" permission.
func TestClientActionsCopy(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t, WithPermissions(srv, "clipboard-read", "clipboard-write"))
	defer cleanup()

	_, err := page.Goto(srv + "/client-actions/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	exact := true
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copy", Exact: &exact})
	if err := btn.Click(); err != nil {
		t.Fatalf("click copy: %v", err)
	}

	// Read the clipboard via the browser API.
	result, err := page.Evaluate("() => navigator.clipboard.readText()")
	if err != nil {
		t.Fatalf("clipboard read: %v", err)
	}

	text, ok := result.(string)
	if !ok || text != "sk_live_abc123def456" {
		t.Errorf("clipboard = %q, want %q", text, "sk_live_abc123def456")
	}
}

// TestClientActionsFlashText clicks the copy button and verifies the
// button text temporarily changes to "Copied!" via FlashText.
func TestClientActionsFlashText(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t, WithPermissions(srv, "clipboard-read", "clipboard-write"))
	defer cleanup()

	_, err := page.Goto(srv + "/client-actions/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	exact := true
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copy", Exact: &exact})
	if err := btn.Click(); err != nil {
		t.Fatalf("click copy: %v", err)
	}

	// FlashText should change the button text to "Copied!".
	copied := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copied!"})
	if err := expect(copied).ToBeVisible(); err != nil {
		t.Errorf("flash text did not appear: %v", err)
	}
}

// TestClientActionsCopyLink verifies the second panel's copy link
// button copies the URL to the clipboard.
func TestClientActionsCopyLink(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t, WithPermissions(srv, "clipboard-read", "clipboard-write"))
	defer cleanup()

	_, err := page.Goto(srv + "/client-actions/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// The second panel has #link-source and a "Copy Link" button.
	linkSource := page.Locator("#link-source")
	if err := expect(linkSource).ToBeVisible(); err != nil {
		t.Fatalf("link source not visible: %v", err)
	}
	if err := expect(linkSource).ToContainText("https://example.com/share/abc123"); err != nil {
		t.Fatalf("link source text mismatch: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copy Link"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click copy link: %v", err)
	}

	result, err := page.Evaluate("() => navigator.clipboard.readText()")
	if err != nil {
		t.Fatalf("clipboard read: %v", err)
	}

	text, ok := result.(string)
	if !ok || text != "https://example.com/share/abc123" {
		t.Errorf("clipboard = %q, want %q", text, "https://example.com/share/abc123")
	}
}

// TestClientActionsCopyLinkFlashClass verifies that clicking "Copy
// Link" temporarily adds the btn-flashed CSS class via FlashClass.
func TestClientActionsCopyLinkFlashClass(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t, WithPermissions(srv, "clipboard-read", "clipboard-write"))
	defer cleanup()

	_, err := page.Goto(srv + "/client-actions/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Copy Link"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click copy link: %v", err)
	}

	// FlashClass should add "btn-flashed" to the button element.
	flashed := page.Locator("button.btn-flashed")
	if err := expect(flashed).ToBeVisible(); err != nil {
		t.Errorf("flash class not applied to button: %v", err)
	}
}
