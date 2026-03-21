package playwright_test

import "testing"

// TestErrorsPageRenders verifies the error boundaries page loads.
func TestErrorsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/errors")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}

	// The error boundary catches a deliberate panic and shows
	// the fallback content.
	fallback := page.GetByText("Caught by error boundary")
	if err := expect(fallback).ToBeVisible(); err != nil {
		t.Errorf("error boundary fallback not visible: %v", err)
	}
}
