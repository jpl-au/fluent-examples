package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestEventsPageRenders verifies the events page loads and the Click
// Events demo card is visible.
func TestEventsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}

	heading := page.GetByText("Click Events")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestEventsClick clicks the button and verifies the counter
// increments via stateless HTTP POST.
func TestEventsClick(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.Locator("[data-tether-click='events.click']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The click counter text is inline, not Dynamic-keyed, so we
	// check the full page for the updated text after morph.
	result := page.GetByText("Clicked 1 times")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("click counter did not update: %v", err)
	}
}

// TestEventsFormSubmit fills the name field, submits, and verifies
// the success result appears.
func TestEventsFormSubmit(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Two inputs share name="name" (Submit form and Bind form), so
	// scope the locator to the submit form.
	formScope := page.Locator("[data-tether-submit='events.submit']")
	nameField := formScope.Locator("input[name='name']")
	if err := nameField.Fill("Alice"); err != nil {
		t.Fatalf("fill: %v", err)
	}

	submit := formScope.Locator("button[type='submit']")
	if err := submit.Click(); err != nil {
		t.Fatalf("submit: %v", err)
	}

	result := page.Locator("[data-tether-key='submit-result']")
	if err := expect(result).ToContainText("Hello, Alice!"); err != nil {
		text, _ := result.TextContent()
		t.Errorf("result = %q, want to contain 'Hello, Alice!'", text)
	}
}

// TestEventsFormSubmitEmpty submits the form with an empty name and
// verifies the validation error appears.
func TestEventsFormSubmitEmpty(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	formScope := page.Locator("[data-tether-submit='events.submit']")
	submit := formScope.Locator("button[type='submit']")
	if err := submit.Click(); err != nil {
		t.Fatalf("submit: %v", err)
	}

	result := page.Locator("[data-tether-key='submit-result']")
	if err := expect(result).ToContainText("Name is required"); err != nil {
		text, _ := result.TextContent()
		t.Errorf("result = %q, want to contain 'Name is required'", text)
	}
}

// TestEventsChangeDropdown selects a colour from the dropdown and
// verifies the server receives the change event.
func TestEventsChangeDropdown(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	sel := page.Locator("select[name='colour']")
	if _, err := sel.SelectOption(pw.SelectOptionValues{Values: &[]string{"red"}}); err != nil {
		t.Fatalf("select: %v", err)
	}

	result := page.Locator("[data-tether-key='colour-result']")
	if err := expect(result).ToContainText("red"); err != nil {
		text, _ := result.TextContent()
		t.Errorf("result = %q, want to contain 'red'", text)
	}
}

// TestEventsKeydown focuses the key field and presses Enter, verifying
// the server receives the filtered keydown event.
func TestEventsKeydown(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	field := page.Locator("input[name='key']")
	if err := field.Click(); err != nil {
		t.Fatalf("focus: %v", err)
	}
	if err := field.Press("Enter"); err != nil {
		t.Fatalf("press: %v", err)
	}

	result := page.Locator("[data-tether-key='key-result']")
	if err := expect(result).ToContainText("Enter"); err != nil {
		text, _ := result.TextContent()
		t.Errorf("result = %q, want to contain 'Enter'", text)
	}
}

// TestEventsPaste dispatches a paste event and verifies the server
// receives the pasted text.
func TestEventsPaste(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Dispatch a paste event with text via JS since Playwright
	// doesn't have a native paste API.
	_, err = page.Evaluate(`() => {
		var el = document.querySelector('[data-tether-paste="events.paste"]');
		if (!el) return;
		var dt = new DataTransfer();
		dt.setData('text', 'hello from clipboard');
		var ev = new ClipboardEvent('paste', {clipboardData: dt, bubbles: true});
		el.dispatchEvent(ev);
	}`)
	if err != nil {
		t.Fatalf("dispatch paste: %v", err)
	}

	result := page.GetByText("hello from clipboard")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("paste result not visible: %v", err)
	}
}

// TestEventsContextMenu right-clicks the context menu target and
// verifies the server receives the event (browser default suppressed).
func TestEventsContextMenu(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/events")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Dispatch a contextmenu event via JS. Playwright's right-click
	// may not trigger the DOM event consistently on all platforms.
	_, err = page.Evaluate(`() => {
		var el = document.querySelector('[data-tether-contextmenu]');
		if (!el) return;
		el.dispatchEvent(new MouseEvent('contextmenu', {bubbles: true}));
	}`)
	if err != nil {
		t.Fatalf("dispatch contextmenu: %v", err)
	}

	result := page.GetByText("Context menu intercepted!")
	if err := expect(result).ToBeVisible(); err != nil {
		t.Errorf("context menu result not visible: %v", err)
	}
}
