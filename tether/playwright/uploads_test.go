package playwright_test

import (
	"os"
	"path/filepath"
	"testing"

	pw "github.com/playwright-community/playwright-go"
)

// TestUploadsPageRenders verifies the uploads page loads with
// the file input and upload button.
func TestUploadsPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/uploads/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	input := page.Locator("#upload-input")
	if err := expect(input).ToBeAttached(); err != nil {
		t.Fatalf("file input not found: %v", err)
	}

	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Upload"})
	if err := expect(btn).ToBeVisible(); err != nil {
		t.Fatalf("upload button not visible: %v", err)
	}
}

// TestUploadsFileUpload selects a file, clicks Upload, and verifies
// the filename appears in the uploaded files list.
func TestUploadsFileUpload(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/uploads/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Create a temporary file to upload.
	dir := t.TempDir()
	path := filepath.Join(dir, "test-upload.txt")
	if err := os.WriteFile(path, []byte("hello from playwright"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	// Set the file on the input element.
	input := page.Locator("#upload-input")
	if err := input.SetInputFiles(path); err != nil {
		t.Fatalf("set input files: %v", err)
	}

	// Click the upload button.
	btn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Upload"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click upload: %v", err)
	}

	// The filename should appear in the uploaded files list.
	list := page.Locator("[data-tether-key='uploads']")
	if err := expect(list).ToContainText("test-upload.txt"); err != nil {
		t.Errorf("uploaded file not in list: %v", err)
	}
}
