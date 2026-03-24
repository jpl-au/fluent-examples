// Package playwright provides end-to-end browser tests for the tether
// example application. These tests start the real example server via
// [app.New], open pages in a headless browser, and verify the full
// stack: initial HTML render, WebSocket connection, event handling,
// DOM morphing, and signal bindings.
//
// The test protocol is controlled by the TETHER_PROTO environment
// variable:
//
//	go test -v ./playwright/...                      # default: HTTP1 (HTTP/1.1)
//	TETHER_PROTO=HTTP2 go test -v ./playwright/...   # HTTP/2 over TLS
//
// HTTP/1.1 uses a real http.Server (not httptest) because httptest's
// internal connection tracking stalls SSE EventSource retries for
// ~30s. HTTP/2 uses httptest with TLS and EnableHTTP2, matching
// modern production defaults. Both modes verify the same test suite.
//
// Uses the system-installed Chrome. Requires Chrome or Chromium to be
// installed - no Playwright browser download is performed.
package playwright_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	tether "github.com/jpl-au/tether"
	pw "github.com/playwright-community/playwright-go"

	"github.com/jpl-au/fluent-examples/tether/app"
)

// useHTTP2 reports whether the test suite is running in HTTP/2 mode.
// Set TETHER_PROTO=HTTP2 to enable.
// ServerMode selects the server type for startApp.
type ServerMode int

const (
	// HTTP1 starts httptest with TLS (HTTP/1.1 over TLS).
	HTTP1 ServerMode = iota + 1

	// HTTP2 starts httptest with TLS and HTTP/2 enabled.
	HTTP2

	// RealHTTP1 starts a real http.Server over plain HTTP/1.1
	// (not httptest). Used for the SSE test because httptest's
	// internal WaitGroup blocks EventSource retries for ~30s.
	// Service workers and push are not available without TLS.
	RealHTTP1
)

// serverMode returns HTTP1 or HTTP2 based on the TETHER_PROTO
// environment variable. Default is HTTP1.
func serverMode() ServerMode {
	if os.Getenv("TETHER_PROTO") == "HTTP2" {
		return HTTP2
	}
	return HTTP1
}

// startApp creates the full example application and returns its base
// URL. The server is shut down when the test ends.
func startApp(t *testing.T, mode ServerMode) string {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	assets := &tether.Asset{
		FS:     os.DirFS("../static"),
		Prefix: "/static/",
	}
	mux, _ := app.New(ctx, assets)

	switch mode {
	case HTTP2:
		srv := httptest.NewUnstartedServer(mux)
		srv.EnableHTTP2 = true
		srv.StartTLS()
		t.Cleanup(srv.Close)
		return srv.URL

	case RealHTTP1:
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		srv := &http.Server{Handler: mux}
		go srv.Serve(ln)
		t.Cleanup(func() { srv.Close() })
		return "http://" + ln.Addr().String()

	default: // HTTP1
		srv := httptest.NewTLSServer(mux)
		t.Cleanup(srv.Close)
		return srv.URL
	}
}

// PageOption configures a browser page created by [newPage].
type PageOption func(ctx pw.BrowserContext) error

// WithPermissions grants browser permissions for the given origin.
// Use for tests that need notification permission or other
// browser-gated APIs:
//
//	page, cleanup := newPage(t, WithPermissions(srv, "notifications"))
func WithPermissions(origin string, permissions ...string) PageOption {
	return func(ctx pw.BrowserContext) error {
		return ctx.GrantPermissions(permissions, pw.BrowserContextGrantPermissionsOptions{
			Origin: &origin,
		})
	}
}

// newPage starts Playwright and opens a headless Chromium browser.
// HTTPS certificate errors are ignored for the self-signed cert from
// httptest (HTTP/2 mode). Pass [WithPermissions] to grant browser
// permissions. Returns a page and a cleanup function.
func newPage(t *testing.T, opts ...PageOption) (pw.Page, func()) {
	t.Helper()

	playwright, err := pw.Run()
	if err != nil {
		t.Skipf("Playwright driver not installed, skipping browser tests.\n" +
			"Install it with:\n\n" +
			"  go run github.com/playwright-community/playwright-go/cmd/playwright@latest install\n")
	}

	launchOpts := pw.BrowserTypeLaunchOptions{
		Headless: pw.Bool(true),
		Args:     []string{"--ignore-certificate-errors"},
	}
	// Find a Chromium-based browser across platforms. The Channel
	// option only works with Chrome installed at known vendor paths.
	// Snap, Homebrew, and distro packages put the binary elsewhere,
	// so we fall back to an explicit PATH lookup.
	if path := findChromium(); path != "" {
		launchOpts.ExecutablePath = pw.String(path)
	} else {
		launchOpts.Channel = pw.String("chrome")
	}
	browser, err := playwright.Chromium.Launch(launchOpts)
	if err != nil {
		playwright.Stop()
		t.Fatalf("browser launch: %v", err)
	}

	ctx, err := browser.NewContext(pw.BrowserNewContextOptions{
		IgnoreHttpsErrors: pw.Bool(true),
	})
	if err != nil {
		browser.Close()
		playwright.Stop()
		t.Fatalf("new context: %v", err)
	}

	for _, opt := range opts {
		if err := opt(ctx); err != nil {
			ctx.Close()
			browser.Close()
			playwright.Stop()
			t.Fatalf("page option: %v", err)
		}
	}

	page, err := ctx.NewPage()
	if err != nil {
		ctx.Close()
		browser.Close()
		playwright.Stop()
		t.Fatalf("new page: %v", err)
	}

	cleanup := func() {
		page.Close()
		ctx.Close()
		browser.Close()
		playwright.Stop()
	}
	return page, cleanup
}

// expect returns a Playwright assertion helper for locator assertions
// with automatic waiting and retrying.
func expect(l pw.Locator) pw.LocatorAssertions {
	return pw.NewPlaywrightAssertions().Locator(l)
}

// waitForConnected waits until the tether root element has
// data-tether-state="connected", indicating the WebSocket or SSE
// transport is open and the client is ready to send/receive events.
func waitForConnected(t *testing.T, page pw.Page) {
	t.Helper()
	connected := page.Locator("[data-tether-state='connected']")
	if err := connected.WaitFor(); err != nil {
		t.Fatalf("tether did not connect: %v", err)
	}
}

// jsNumber converts a value returned by Playwright's Evaluate into a
// float64. JS numbers arrive as either int or float64 depending on
// whether they have a fractional part. This helper handles both so
// callers don't need type switches at every call site.
func jsNumber(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}

// findChromium searches PATH for a Chromium-based browser. Returns
// the absolute path if found, or empty string to fall back to
// Playwright's built-in channel lookup. This handles snap, Homebrew,
// distro packages, and other non-standard install locations.
func findChromium() string {
	for _, name := range []string{
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"chrome",
	} {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}
