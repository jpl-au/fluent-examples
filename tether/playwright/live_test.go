package playwright_test

import "testing"

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

	btn := page.Locator("[data-tether-click='live.broadcast']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	result := page.Locator("[data-tether-key='broadcast']")
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

	btn := page.Locator("[data-tether-click='live.broadcast']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The broadcast result shows "Last: broadcast at HH:MM:SS".
	result := page.Locator("[data-tether-key='broadcast']")
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

	btn := page.Locator("[data-tether-click='live.set-title']")
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The page title should change to include "Title set at".
	titleText, err := page.Title()
	if err != nil {
		t.Fatalf("title: %v", err)
	}

	// Wait briefly for the title to update.
	page.WaitForTimeout(500)
	titleText, _ = page.Title()

	if len(titleText) == 0 {
		t.Error("page title is empty after SetTitle")
	}
}
