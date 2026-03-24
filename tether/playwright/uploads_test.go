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

	// Create a temporary file to upload. Use the working directory
	// instead of t.TempDir() because snap-installed Chromium cannot
	// read files from /tmp due to sandbox isolation.

	dir, _ := os.Getwd()
	path := filepath.Join(dir, "test-upload.txt")
	if err := os.WriteFile(path, []byte("hello from playwright"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })

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

// TestUploadsDownload uploads a file and then clicks the Download
// button, verifying that sess.Download triggers a file download via
// normal HTTP. The download event is intercepted by Playwright to
// confirm the file was received.
func TestUploadsDownload(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	_, err := page.Goto(srv + "/uploads/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// Upload a file first.
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "test-download.txt")
	if err := os.WriteFile(path, []byte("download test content"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })

	input := page.Locator("#upload-input")
	if err := input.SetInputFiles(path); err != nil {
		t.Fatalf("set input files: %v", err)
	}

	uploadBtn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Upload"})
	if err := uploadBtn.Click(); err != nil {
		t.Fatalf("click upload: %v", err)
	}

	// Wait for the file to appear in the list.
	list := page.Locator("[data-tether-key='uploads']")
	if err := expect(list).ToContainText("test-download.txt"); err != nil {
		t.Fatalf("uploaded file not in list: %v", err)
	}

	// Click the Download button and expect a download event.
	download, err := page.ExpectDownload(func() error {
		dlBtn := page.GetByRole("button", pw.PageGetByRoleOptions{Name: "Download"})
		return dlBtn.Click()
	})
	if err != nil {
		t.Fatalf("expected download: %v", err)
	}

	if download.SuggestedFilename() != "test-download.txt" {
		t.Errorf("filename = %q, want test-download.txt", download.SuggestedFilename())
	}
}
