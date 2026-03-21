package playwright_test

import "testing"

// TestHealthEndpoint verifies the server is up and the health
// endpoint responds.
func TestHealthEndpoint(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/health")
	if err != nil {
		t.Fatalf("goto /health: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("health status = %d, want 200", resp.Status())
	}
}

// TestHTTPPageRenders verifies a stateless page (tether.Stateless) loads.
func TestHTTPPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}
}
