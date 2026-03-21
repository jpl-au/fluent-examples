package playwright_test

import (
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestFilteredUploadsPageRenders verifies the filtered uploads page
// loads and the upload button is visible.
func TestFilteredUploadsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/uploads/filtered/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Upload"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Errorf("upload button not visible: %v", err)
	}
}
