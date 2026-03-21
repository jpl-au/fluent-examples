package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestBroadcastingPageRenders verifies the broadcasting page loads.
func TestBroadcastingPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/broadcasting/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("Cross-Session Events")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestBroadcastingSendMessage types a message, clicks Send, and
// verifies it appears in the sender's message list.
func TestBroadcastingSendMessage(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/broadcasting/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	input := page.Locator("#broadcast-input")
	if err := input.Fill("hello from test"); err != nil {
		t.Fatalf("fill: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The sender sees their own message immediately.
	messages := page.Locator("[data-tether-key='messages']")
	if err := expect(messages).ToContainText("hello from test"); err != nil {
		t.Errorf("message not visible in sender: %v", err)
	}
}

// TestBroadcastingCrossSession opens two browser pages, sends a
// message from one, and verifies the other receives it.
func TestBroadcastingCrossSession(t *testing.T) {
	srv := startApp(t, serverMode())

	sender, cleanupSender := newPage(t)
	defer cleanupSender()

	receiver, cleanupReceiver := newPage(t)
	defer cleanupReceiver()

	_, err := sender.Goto(srv + "/broadcasting/")
	if err != nil {
		t.Fatalf("sender goto: %v", err)
	}
	_, err = receiver.Goto(srv + "/broadcasting/")
	if err != nil {
		t.Fatalf("receiver goto: %v", err)
	}

	waitForConnected(t, sender)
	waitForConnected(t, receiver)

	// Send a message from the sender.
	input := sender.Locator("#broadcast-input")
	if err := input.Fill("cross-session test"); err != nil {
		t.Fatalf("fill: %v", err)
	}
	btn := sender.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The receiver should see the message via WatchBus.
	receiverMessages := receiver.Locator("[data-tether-key='messages']")
	if err := expect(receiverMessages).ToContainText("cross-session test"); err != nil {
		t.Errorf("message not received by other session: %v", err)
	}
}

// TestBroadcastingMessageCounter sends a message and verifies the
// shared message counter increments.
func TestBroadcastingMessageCounter(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/broadcasting/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	input := page.Locator("#broadcast-input")
	if err := input.Fill("counter test"); err != nil {
		t.Fatalf("fill: %v", err)
	}
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The message counter should show at least 1.
	counter := page.Locator("[data-tether-key='message-count']")
	if err := expect(counter).Not().ToContainText("Total messages: 0"); err != nil {
		text, _ := counter.TextContent()
		t.Errorf("counter should have incremented, got %q", text)
	}
}
