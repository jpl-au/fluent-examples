package playwright_test

import "testing"

// TestRealtimePageRenders verifies the real-time dashboard loads and
// the chart containers are present in the DOM.
func TestRealtimePageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/realtime/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	heading := page.GetByText("System Monitor")
	if err := expect(heading).ToBeVisible(); err != nil {
		t.Fatalf("heading not visible: %v", err)
	}

	// Verify the three chart containers exist - the echarts hook
	// initialises them after mount.
	for _, id := range []string{"chartcpu", "chartheap", "chartgoroutines"} {
		chart := page.Locator("#" + id)
		if err := expect(chart).ToBeVisible(); err != nil {
			t.Errorf("chart %q not visible: %v", id, err)
		}
	}
}
