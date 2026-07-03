package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestNotificationsPageRenders verifies the SW push page loads and
// the subscribe button is visible (VAPID keys are available).
func TestNotificationsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/sw/push")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Enable Push Notifications"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("subscribe button not visible: %v", err)
	}
}

// TestNotificationsPushSubscribe grants notification permission,
// clicks the subscribe button, and verifies the send buttons become
// visible (subscription succeeded).
func TestNotificationsPushSubscribe(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t, WithPermissions(srv, "notifications"))
	defer cleanup()

	_, err := page.Goto(srv + "/sw/push")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Click subscribe - notification permission was granted above.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Enable Push Notifications"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click subscribe: %v", err)
	}

	// After subscribing, the "Send Test Push" button should be visible.
	send := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send Test Push"})
	if err := expect(send).ToBeVisible(); err != nil {
		t.Errorf("send button not visible after subscribe: %v", err)
	}
}

// TestNotificationsMainPageRenders verifies the main notifications
// page loads and shows all demo buttons.
func TestNotificationsMainPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/notifications/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Show Toast"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("Show Toast button not visible: %v", err)
	}
}

// TestNotificationsToast clicks the toast button and verifies a
// toast notification appears in the DOM.
func TestNotificationsToast(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/notifications/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Show Toast"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The toast container is created lazily by tether.js and toast
	// elements have the class "tether-toast".
	toast := page.Locator(".tether-toast")
	if err := expect(toast).ToBeAttached(); err != nil {
		t.Errorf("toast element not found in DOM: %v", err)
	}
	if err := expect(toast).ToContainText("This is a toast notification!"); err != nil {
		t.Errorf("toast text mismatch: %v", err)
	}
}

// TestNotificationsFlash clicks the flash button and verifies the
// target element text changes temporarily.
func TestNotificationsFlash(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/notifications/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Flash Message"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The flash target's text content should be replaced with "Flashed!".
	target := page.Locator("#flash-target")
	if err := expect(target).ToHaveText("Flashed!"); err != nil {
		text, _ := target.TextContent()
		t.Errorf("flash target text = %q, want %q", text, "Flashed!")
	}
}

// TestNotificationsAnnounce clicks the announce button and verifies
// the ARIA live region is populated and a visual echo appears.
func TestNotificationsAnnounce(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/notifications/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Announce"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The handler sends a signal that echoes the announcement text
	// visually via bind.Text and bind.Show.
	echo := page.Locator("[data-tether-bind-text='notify.announced']")
	if err := expect(echo).ToContainText("Screen readers heard"); err != nil {
		t.Errorf("announce echo not visible: %v", err)
	}

	// The ARIA live region should contain the announcement text.
	liveRegion := page.Locator("[aria-live='polite'][aria-atomic='true']")
	if err := expect(liveRegion).ToHaveText("New item added to your feed"); err != nil {
		text, _ := liveRegion.TextContent()
		t.Errorf("live region text = %q, want announcement text", text)
	}

	// A toast is also sent alongside the announcement.
	toast := page.Locator(".tether-toast")
	if err := expect(toast).ToBeAttached(); err != nil {
		t.Errorf("announce should also fire a toast: %v", err)
	}
}

// TestNotificationsIndicator clicks the loading indicator button and
// verifies the spinner becomes visible during the simulated delay.
func TestNotificationsIndicator(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/notifications/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The spinner should be hidden initially.
	spinner := page.Locator("#notify-spinner")
	if err := expect(spinner).Not().ToBeVisible(); err != nil {
		t.Fatalf("spinner should be hidden initially: %v", err)
	}

	// Click the indicator button - the bind.Indicator directive
	// shows the spinner while the action is in flight.
	btn := page.Locator("[data-tether-click='notify.indicator']").First()
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The spinner should become visible during the simulated delay.
	if err := expect(spinner).ToBeVisible(); err != nil {
		t.Errorf("spinner should be visible during loading: %v", err)
	}
}

// TestNotificationsFlashCompare clicks the Flash button in the
// comparison panel and verifies the target updates.
func TestNotificationsFlashCompare(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/notifications/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The "Flash vs Signal" panel has two buttons. Click the Flash one.
	btn := page.Locator("[data-tether-click='notify.flash-compare']").First()
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	target := page.Locator("#flash-compare-target")
	if err := expect(target).ToHaveText("Saved!"); err != nil {
		text, _ := target.TextContent()
		t.Errorf("flash compare target text = %q, want %q", text, "Saved!")
	}
}

// TestNotificationsSignalFlash clicks the Signal button in the
// comparison panel and verifies the "Saved!" text appears via
// bind.Show.
func TestNotificationsSignalFlash(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/notifications/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Click the signal button in the comparison panel.
	btn := page.Locator("[data-tether-click='notify.signal-flash']").First()
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The "Saved!" text should become visible via bind.Show.
	saved := page.Locator("[data-tether-bind-show='notify.saved']")
	if err := expect(saved).ToBeVisible(); err != nil {
		t.Errorf("signal saved text should be visible: %v", err)
	}
}
