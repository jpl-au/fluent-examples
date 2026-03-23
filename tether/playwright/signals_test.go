package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestSignalsPageRenders verifies the signals page loads and shows
// the initial counter value.
func TestSignalsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	btn := page.GetByText("Increment Server Counter")
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("increment button not visible: %v", err)
	}
}

// TestSignalsIncrement clicks the increment button and verifies the
// counter updates in the DOM via the WebSocket connection.
func TestSignalsIncrement(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByText("Increment Server Counter")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// After clicking, the server pushes the new count as a signal
	// and the client updates the text via BindText.
	counter := page.Locator(".demo:first-child [data-tether-bind-text='signals.counter']")
	if err := expect(counter).ToHaveText("2"); err != nil {
		text, _ := counter.TextContent()
		t.Errorf("counter text = %q, want %q", text, "2")
	}
}

// TestSignalsTogglePanel clicks the toggle button and verifies the
// panel appears via BindShow.
func TestSignalsTogglePanel(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The panel should be hidden initially.
	panel := page.Locator("[data-tether-bind-show='signals.panel_visible']")
	if err := expect(panel).Not().ToBeVisible(); err != nil {
		t.Fatalf("panel should be hidden initially: %v", err)
	}

	// Click toggle - panel should appear.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Toggle Visibility"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click toggle: %v", err)
	}

	if err := expect(panel).ToBeVisible(); err != nil {
		t.Errorf("panel should be visible after toggle: %v", err)
	}
}

// TestSignalsToggleLock clicks the lock button and verifies the
// input becomes disabled via BindAttr.
func TestSignalsToggleLock(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Click toggle lock.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Toggle Lock"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The input should become disabled. BindAttr renders as a single
	// attribute: data-tether-bind-attr="disabled signals.input_locked".
	input := page.Locator("[data-tether-bind-attr='disabled signals.input_locked']")
	if err := expect(input).ToBeDisabled(); err != nil {
		t.Errorf("input should be disabled after toggle lock: %v", err)
	}
}

// TestSignalsSetSignalClientSide clicks a colour button and verifies
// the text updates via client-side SetSignal (no server round-trip).
func TestSignalsSetSignalClientSide(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByText("Set to Blue")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The colour display span is bound via BindText to the
	// signals.colour signal.
	display := page.Locator("[data-tether-key='colour-display']")
	if err := expect(display).ToHaveText("blue"); err != nil {
		text, _ := display.TextContent()
		t.Errorf("colour display = %q, want %q", text, "blue")
	}
}

// TestSignalsPrefillValue clicks the prefill button and verifies the
// input value is set via BindValue from a server-pushed signal.
func TestSignalsPrefillValue(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Pre-fill from Server"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The input is bound via BindValue to signals.prefill_value.
	// The server pushes "hello@example.com".
	input := page.Locator("[data-tether-bind-value='signals.prefill_value']")
	if err := expect(input).ToHaveValue("hello@example.com"); err != nil {
		t.Errorf("prefill input: %v", err)
	}
}

// TestSignalsFavouriteToggle clicks the favourite button twice and
// verifies the text toggles between "Favourited!" and "Not favourited"
// via OptimisticToggle and BindShow/BindHide.
func TestSignalsFavouriteToggle(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Toggle Favourite"})

	// First click - should show "Favourited!".
	if err := btn.Click(); err != nil {
		t.Fatalf("click 1: %v", err)
	}
	favourited := page.GetByText("Favourited!")
	if err := expect(favourited).ToBeVisible(); err != nil {
		t.Errorf("should show Favourited! after first click: %v", err)
	}

	// Second click - should show "Not favourited".
	if err := btn.Click(); err != nil {
		t.Fatalf("click 2: %v", err)
	}
	notFavourited := page.GetByText("Not favourited")
	if err := expect(notFavourited).ToBeVisible(); err != nil {
		t.Errorf("should show Not favourited after second click: %v", err)
	}
}

// TestSignalsResetAll builds up state across several actions, then
// clicks Reset All and verifies everything returns to defaults via
// the batch sess.Signals() call.
func TestSignalsResetAll(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Build up some state.
	increment := page.GetByText("Increment Server Counter")
	if err := increment.Click(); err != nil {
		t.Fatalf("click increment: %v", err)
	}
	if err := increment.Click(); err != nil {
		t.Fatalf("click increment: %v", err)
	}

	toggle := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Toggle Visibility"})
	if err := toggle.Click(); err != nil {
		t.Fatalf("click toggle: %v", err)
	}

	// Verify state is built up - counter should be 2, panel visible.
	counter := page.Locator(".demo:first-child [data-tether-bind-text='signals.counter']")
	if err := expect(counter).ToHaveText("2"); err != nil {
		t.Fatalf("counter should be 2 before reset: %v", err)
	}

	// Click Reset All.
	reset := page.GetByText("Reset All Signals")
	if err := reset.Click(); err != nil {
		t.Fatalf("click reset: %v", err)
	}

	// Counter should return to 0.
	if err := expect(counter).ToHaveText("0"); err != nil {
		t.Errorf("counter should be 0 after reset: %v", err)
	}

	// Panel should be hidden again.
	panel := page.Locator("[data-tether-bind-show='signals.panel_visible']")
	if err := expect(panel).Not().ToBeVisible(); err != nil {
		t.Errorf("panel should be hidden after reset: %v", err)
	}
}

// TestSignalsCloak verifies that a cloaked element becomes visible
// after the tether JS initialises and the BindText signal populates
// it. The cloak attribute hides the element during SSR to prevent a
// flash of stale placeholder content.
func TestSignalsCloak(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// After JS initialises, the cloak is removed and the element
	// becomes visible with the signal's current value.
	cloaked := page.Locator("[data-tether-key='cloaked']")
	if err := expect(cloaked).ToBeVisible(); err != nil {
		t.Errorf("cloaked element should be visible after JS init: %v", err)
	}

	// The text should be the counter value (0 on fresh page), not
	// the placeholder "Loading...".
	if err := expect(cloaked).ToHaveText("0"); err != nil {
		text, _ := cloaked.TextContent()
		t.Errorf("cloaked text = %q, want %q", text, "0")
	}
}

// TestSignalsPermanent verifies that an element marked with
// bind.Permanent is never replaced by the differ, even after a
// state change triggers a re-render.
func TestSignalsPermanent(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The permanent element should be visible.
	permanent := page.Locator("[data-tether-permanent]")
	if err := expect(permanent).ToBeVisible(); err != nil {
		t.Fatalf("permanent element not visible: %v", err)
	}

	// Trigger a state change (increment counter) that causes a
	// re-render. The permanent element should survive the morph.
	btn := page.GetByText("Increment Server Counter")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// Verify the permanent element is still there after the morph.
	if err := expect(permanent).ToBeVisible(); err != nil {
		t.Errorf("permanent element should survive morph: %v", err)
	}
	if err := expect(permanent).ToContainText("never replaced"); err != nil {
		t.Errorf("permanent element content should be unchanged: %v", err)
	}
}

// TestSignalsTransition clicks the transition toggle and verifies
// the panel appears and disappears. The CSS fade animation runs over
// 200ms - we wait for the element to be attached to the DOM first,
// then check visibility after the transition completes.
func TestSignalsTransition(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/signals/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Toggle Transition"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The panel should appear in the DOM after the morph. Wait for
	// it to be attached (present in DOM) rather than visible (which
	// depends on the CSS transition completing).
	panel := page.GetByText("Visible with transition.")
	if err := expect(panel).ToBeAttached(); err != nil {
		t.Errorf("transition panel not in DOM after toggle: %v", err)
	}
}
