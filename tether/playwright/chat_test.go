package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestChatPageRenders verifies the chat page loads.
func TestChatPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/chat/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	input := page.Locator("#shout-input")
	if err := expect(input).ToBeVisible(); err != nil {
		t.Fatalf("shout input not visible: %v", err)
	}
}

// TestChatSendMessage types a message, clicks Send, and verifies
// it appears in the sender's feed.
func TestChatSendMessage(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/chat/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	input := page.Locator("#shout-input")
	if err := input.Fill("hello chat"); err != nil {
		t.Fatalf("fill: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	feed := page.Locator("[data-tether-key='shout-feed']")
	if err := expect(feed).ToContainText("hello chat"); err != nil {
		t.Errorf("message not in feed: %v", err)
	}
}

// TestChatCrossSession opens two pages, sends from one, and verifies
// the other receives the message via the shared bus.
func TestChatCrossSession(t *testing.T) {
	srv := startApp(t, serverMode())

	sender, cleanupSender := newPage(t)
	defer cleanupSender()

	receiver, cleanupReceiver := newPage(t)
	defer cleanupReceiver()

	_, err := sender.Goto(srv + "/chat/")
	if err != nil {
		t.Fatalf("sender goto: %v", err)
	}
	_, err = receiver.Goto(srv + "/chat/")
	if err != nil {
		t.Fatalf("receiver goto: %v", err)
	}

	waitForConnected(t, sender)
	waitForConnected(t, receiver)

	input := sender.Locator("#shout-input")
	if err := input.Fill("cross-session chat"); err != nil {
		t.Fatalf("fill: %v", err)
	}

	btn := sender.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	receiverFeed := receiver.Locator("[data-tether-key='shout-feed']")
	if err := expect(receiverFeed).ToContainText("cross-session chat"); err != nil {
		t.Errorf("message not received by other session: %v", err)
	}
}
