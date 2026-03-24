package playwright_test

import (
	"testing"
)

// TestHotkeyPageRenders verifies the hotkey demo loads and shows the
// initial hint text.
func TestHotkeyPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/hotkey/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	hint := page.GetByText("No hotkey triggered yet.")
	if err := expect(hint).ToBeVisible(); err != nil {
		t.Fatalf("hint not visible: %v", err)
	}
}

// TestHotkeyCtrlK presses Ctrl+K and verifies the server receives
// the hotkey event and updates the page.
func TestHotkeyCtrlK(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/hotkey/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Press Ctrl+K anywhere on the page.
	if err := page.Keyboard().Press("Control+k"); err != nil {
		t.Fatalf("press ctrl+k: %v", err)
	}

	result := page.GetByText("Last hotkey: ctrl-k")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("ctrl+k not reflected: %v", err)
	}
}

// TestHotkeyEscape presses Escape and verifies the server receives it.
func TestHotkeyEscape(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/hotkey/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	if err := page.Keyboard().Press("Escape"); err != nil {
		t.Fatalf("press escape: %v", err)
	}

	result := page.GetByText("Last hotkey: escape")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("escape not reflected: %v", err)
	}
}
