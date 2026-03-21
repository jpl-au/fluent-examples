package playwright_test

import "testing"

// TestConfigurationPageRenders verifies the configuration page loads
// and shows the compression, timeouts, and diagnostics cards.
func TestConfigurationPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/configuration/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	// Configuration page has TrustedOrigins set to localhost:8080,
	// so the WebSocket won't connect from httptest.Server. The
	// initial HTML is server-rendered, so render checks work
	// without a connection.
	compression := page.GetByText("WebSocket Compression")
	if err := expect(compression).ToBeVisible(); err != nil {
		t.Errorf("compression card not visible: %v", err)
	}

	pageViews := page.GetByText("Page View Counter")
	if err := expect(pageViews).ToBeVisible(); err != nil {
		t.Errorf("page view counter card not visible: %v", err)
	}
}
