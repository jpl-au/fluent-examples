package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestGroupsPageRenders verifies the groups page loads.
func TestGroupsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/groups/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Verify the page content loaded by checking for a button.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Join Alpha"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("Join Alpha button not visible: %v", err)
	}
}

// TestGroupsJoinRoom clicks "Join Alpha" and verifies the room
// status updates.
func TestGroupsJoinRoom(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/groups/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Join Alpha"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	status := page.Locator("[data-tether-key='room-status']")
	if err := expect(status).ToContainText("alpha"); err != nil {
		t.Errorf("room status should show alpha: %v", err)
	}
}

// TestGroupsLeaveRoom joins a room then leaves and verifies the
// status resets.
func TestGroupsLeaveRoom(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/groups/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	join := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Join Alpha"})
	if err := join.Click(); err != nil {
		t.Fatalf("click join: %v", err)
	}

	leave := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Leave Room"})
	if err := leave.Click(); err != nil {
		t.Fatalf("click leave: %v", err)
	}

	hint := page.Locator("[data-tether-key='leave-hint']")
	if err := expect(hint).ToContainText("not in a room"); err != nil {
		t.Errorf("should show 'not in a room' after leaving: %v", err)
	}
}

// TestGroupsBroadcast joins a room, sends a message, and verifies
// it appears.
func TestGroupsBroadcast(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/groups/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	join := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Join Alpha"})
	if err := join.Click(); err != nil {
		t.Fatalf("click join: %v", err)
	}

	input := page.Locator("#message-input")
	if err := input.Fill("hello room"); err != nil {
		t.Fatalf("fill: %v", err)
	}

	send := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send to Room"})
	if err := send.Click(); err != nil {
		t.Fatalf("click send: %v", err)
	}

	msg := page.Locator("[data-tether-key='room-message']")
	if err := expect(msg).ToContainText("hello room"); err != nil {
		t.Errorf("message not visible: %v", err)
	}
}
