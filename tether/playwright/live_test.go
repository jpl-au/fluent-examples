package playwright_test

import (
	"regexp"
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestLiveWSPageRenders verifies the WebSocket live updates page loads.
func TestLiveWSPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/live/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("Uptime Ticker")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestLiveSSEPageRenders verifies the SSE live updates page loads
// and the EventSource connection establishes. On HTTP/1.1, uses a
// real http.Server because httptest stalls SSE retries for ~30s. On
// HTTP/2, uses the default httptest server.
func TestLiveSSEPageRenders(t *testing.T) {
	mode := serverMode()
	if mode == HTTP1 {
		mode = RealHTTP1
	}
	srv := startApp(t, mode)
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/live/sse/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("Uptime Ticker")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}
}

// TestLiveSSEBroadcast clicks the broadcast button on the SSE variant
// and verifies the message appears. Uses RealHTTP1 to avoid httptest
// stalling EventSource retries.
func TestLiveSSEBroadcast(t *testing.T) {
	mode := serverMode()
	if mode == HTTP1 {
		mode = RealHTTP1
	}
	srv := startApp(t, mode)
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/live/sse/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-event-click='live.broadcast']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	result := page.Locator("#broadcast")
	if err := expect(result).ToContainText("Last: broadcast at"); err != nil {
		t.Errorf("SSE broadcast message not visible: %v", err)
	}
}

// TestLiveWSBroadcast clicks the broadcast button and verifies
// the message appears on the page.
func TestLiveWSBroadcast(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/live/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-event-click='live.broadcast']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The broadcast result shows "Last: broadcast at HH:MM:SS".
	result := page.Locator("#broadcast")
	if err := expect(result).ToContainText("Last: broadcast at"); err != nil {
		t.Errorf("broadcast message not visible: %v", err)
	}
}

// TestLiveWSSetTitle clicks "Set Title" and verifies the document
// title changes via Session.SetTitle.
func TestLiveWSSetTitle(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/live/ws/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.Locator("[data-tether-event-click='live.set-title']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	if err := pw.NewPlaywrightAssertions().Page(page).ToHaveTitle(regexp.MustCompile(`^tether - Title set at \d{2}:\d{2}:\d{2}$`)); err != nil {
		t.Errorf("SetTitle did not update the page title: %v", err)
	}
}
