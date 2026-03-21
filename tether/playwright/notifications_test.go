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
