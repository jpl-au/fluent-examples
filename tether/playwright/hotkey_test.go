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

// TestHotkeyCtrlSlash presses Ctrl+/ which contains a character that
// is special in CSS selectors. Without CSS.escape in the hotkey
// handler, querySelector would throw a SyntaxError.
func TestHotkeyCtrlSlash(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/hotkey/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	if err := page.Keyboard().Press("Control+/"); err != nil {
		t.Fatalf("press ctrl+/: %v", err)
	}

	result := page.GetByText("Last hotkey: ctrl-/")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("ctrl+/ not reflected: %v", err)
	}
}

// TestHotkeyShiftQuestion presses Shift+? which produces a key name
// with a special CSS character. Verifies the CSS.escape fix handles
// punctuation keys correctly.
func TestHotkeyShiftQuestion(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/hotkey/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	if err := page.Keyboard().Press("Shift+?"); err != nil {
		t.Fatalf("press shift+?: %v", err)
	}

	result := page.GetByText("Last hotkey: shift-?")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("shift+? not reflected: %v", err)
	}
}

// TestHotkeySequence presses multiple hotkeys in sequence and verifies
// each one updates the display, confirming the handler stays functional
// after processing special character keys.
func TestHotkeySequence(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/hotkey/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Press Ctrl+/ first (special character)
	if err := page.Keyboard().Press("Control+/"); err != nil {
		t.Fatalf("press ctrl+/: %v", err)
	}
	result := page.GetByText("Last hotkey: ctrl-/")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("ctrl+/ not reflected: %v", err)
	}

	// Then press Ctrl+K (normal character) to confirm handler still works
	if err := page.Keyboard().Press("Control+k"); err != nil {
		t.Fatalf("press ctrl+k: %v", err)
	}
	result = page.GetByText("Last hotkey: ctrl-k")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("ctrl+k after ctrl+/ not reflected: %v", err)
	}

	// Then Escape
	if err := page.Keyboard().Press("Escape"); err != nil {
		t.Fatalf("press escape: %v", err)
	}
	result = page.GetByText("Last hotkey: escape")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("escape after sequence not reflected: %v", err)
	}
}

// TestHotkeyUnregisteredKeyIgnored presses a key that has no binding
// and verifies it does not trigger a hotkey event or cause errors.
func TestHotkeyUnregisteredKeyIgnored(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/hotkey/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Press an unregistered key
	if err := page.Keyboard().Press("Control+z"); err != nil {
		t.Fatalf("press ctrl+z: %v", err)
	}

	// The hint should still show the initial text
	hint := page.GetByText("No hotkey triggered yet.")
	if err := expect(hint).ToBeVisible(); err != nil {
		t.Errorf("unregistered key should not trigger hotkey: %v", err)
	}
}
