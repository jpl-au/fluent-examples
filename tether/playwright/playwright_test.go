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
// Both modes use httptest with TLS. Tests that need plain HTTP/1.1
// use RealHTTP1 because httptest's internal connection tracking can
// stall SSE EventSource retries for ~30s. HTTP/2 enables EnableHTTP2.
//
// Uses the system-installed Chrome. Requires Chrome or Chromium to be
// installed - no Playwright browser download is performed.
//
// Shared helpers for server and browser setup, cleanup, timing, connection
// readiness, and assertions live in this file. Add new shared helpers here so
// they stay discoverable and are not duplicated; keep per-feature app scaffolding
// (each feature's start/render/handle functions and its state type) in that
// feature's own test file.
package playwright_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"testing"
	"time"

	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/wire"
	pw "github.com/playwright-community/playwright-go"

	"github.com/jpl-au/fluent-examples/tether/app"
)

// primaryModifier is the platform's primary command modifier key: Meta (Cmd)
// on macOS, Control on Windows and Linux. It is what bind.Hotkey's "mod" alias
// resolves to on the client (which reads navigator.platform) and what selects
// multiple items with Cmd/Ctrl+click. Playwright's headless Chromium reports the
// host OS, so any test that drives a mod-based shortcut must press the modifier
// matching the machine the suite runs on, or it passes on one OS and fails on
// another.
func primaryModifier() pw.KeyboardModifier {
	if runtime.GOOS == "darwin" {
		return "Meta"
	}
	return "Control"
}

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
	finish := testPhase(t, "application setup")
	defer finish()
	ctx, cancel := context.WithCancel(context.Background())
	// Also cancel subscriptions if application setup fails before startServer.
	t.Cleanup(cancel)

	assets := &tether.Asset{
		FS:     os.DirFS("../static"),
		Prefix: "/static/",
	}
	mux, handlers := app.New(ctx, assets, wire.JSON)
	shutdown := make([]func(context.Context) error, 0, len(handlers)+1)
	shutdown = append(shutdown, func(context.Context) error { cancel(); return nil })
	for _, handler := range handlers {
		shutdown = append(shutdown, handler.Shutdown)
	}
	return startServer(t, mode, mux, shutdown...)
}

// Cleanup has its own allowance; browser actions and assertions retain their
// existing timeouts. Tests stay sequential, with a fresh server and browser
// for every fixture, including when selected individually with -run.
const cleanupTimeout = 30 * time.Second

// testPhase reports setup and teardown timings when TETHER_TEST_TIMINGS=1.
// Start messages identify a stalled phase; goroutine counts help expose
// resources retained across sequential tests. They are diagnostic, not assertions.
func testPhase(t *testing.T, name string) func() {
	t.Helper()
	if os.Getenv("TETHER_TEST_TIMINGS") != "1" {
		return func() {}
	}
	started := time.Now()
	t.Logf("%s: start (goroutines=%d)", name, runtime.NumGoroutine())
	return func() {
		t.Helper()
		t.Logf("%s: %s (goroutines=%d)", name, time.Since(started).Round(time.Millisecond), runtime.NumGoroutine())
	}
}

// shutdownHandlers uses one deadline for the entire fixture. Even after a
// timeout, every handler gets Shutdown so all of them begin releasing resources.
func shutdownHandlers(t *testing.T, shutdown ...func(context.Context) error) bool {
	t.Helper()
	defer testPhase(t, "handler shutdown")()
	ctx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
	defer cancel()
	ok := true
	for i, stop := range shutdown {
		if err := stop(ctx); err != nil {
			t.Errorf("handler %d shutdown: %v", i, err)
			ok = false
		}
	}
	return ok
}

// waitCleanup waits for a cleanup operation owned by this test. The caller
// must unblock and join an operation that times out before returning.
func waitCleanup(t *testing.T, name string, done <-chan error) bool {
	t.Helper()
	timer := time.NewTimer(cleanupTimeout)
	defer timer.Stop()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("%s: %v", name, err)
		}
		return true
	case <-timer.C:
		t.Errorf("%s did not finish within %s", name, cleanupTimeout)
		return false
	}
}

// Once forced shutdown has also failed, abort this test process instead of
// starting another test alongside resources that could not be released.
func requireCleanup(t *testing.T, name string, done <-chan error) {
	t.Helper()
	if !waitCleanup(t, name, done) {
		panic(name + " is still running after forced shutdown")
	}
}

// serveHTTP owns the serving goroutine as well as the listener. Its returned
// cleanup is idempotent so restart fixtures can stop early and still register
// automatic cleanup for a failed test.
func serveHTTP(t *testing.T, srv *http.Server, ln net.Listener) func() {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		err := srv.Serve(ln)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		done <- err
	}()
	return sync.OnceFunc(func() {
		defer testPhase(t, "HTTP server close")()
		if err := srv.Close(); err != nil {
			t.Errorf("HTTP server close: %v", err)
		}
		requireCleanup(t, "HTTP serving goroutine", done)
	})
}

// startServer registers handler shutdown before closing the HTTP server, so
// persistent transports cannot keep httptest.Close waiting on a live handler.
// The shutdown callbacks belong only to this fixture.
func startServer(t *testing.T, mode ServerMode, handler http.Handler, shutdown ...func(context.Context) error) string {
	t.Helper()
	defer testPhase(t, "server setup")()

	if mode == RealHTTP1 {
		var closeServer func()
		t.Cleanup(func() {
			stopped := shutdownHandlers(t, shutdown...)
			if closeServer != nil {
				closeServer()
			}
			if !stopped {
				panic("fixture handlers did not shut down")
			}
		})
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		closeServer = serveHTTP(t, &http.Server{Handler: handler}, ln)
		return "http://" + ln.Addr().String()
	}

	srv := httptest.NewUnstartedServer(handler)
	srv.EnableHTTP2 = mode == HTTP2
	t.Cleanup(func() {
		stopped := shutdownHandlers(t, shutdown...)
		defer testPhase(t, "HTTPS server close")()
		done := make(chan error, 1)
		go func() { srv.Close(); done <- nil }()
		if !waitCleanup(t, "HTTPS server close", done) {
			srv.CloseClientConnections()
			requireCleanup(t, "HTTPS server close after closing connections", done)
		}
		if !stopped {
			panic("fixture handlers did not shut down")
		}
	})
	srv.StartTLS()
	return srv.URL
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

	finish := testPhase(t, "Playwright driver startup")
	playwright, err := pw.Run()
	finish()
	if err != nil {
		t.Skipf("Playwright driver unavailable: %v\n"+
			"Install the driver matching go.mod (no browser download) with:\n\n"+
			"  go run github.com/playwright-community/playwright-go/cmd/playwright --version\n", err)
	}
	var browser pw.Browser
	var ctx pw.BrowserContext
	var page pw.Page
	cleanup := sync.OnceFunc(func() {
		defer testPhase(t, "browser cleanup")()
		closed := make(chan error, 1)
		go func() {
			var errs []error
			collect := func(name string, err error) {
				if err != nil {
					errs = append(errs, fmt.Errorf("%s: %w", name, err))
				}
			}
			if page != nil {
				collect("page close", page.Close())
			}
			if ctx != nil {
				collect("context close", ctx.Close())
			}
			if browser != nil {
				collect("browser close", browser.Close())
			}
			closed <- errors.Join(errs...)
		}()
		browserClosed := waitCleanup(t, "browser cleanup", closed)

		// Stop owns driver shutdown even when a browser close operation stalls.
		// Closing its pipe releases outstanding Playwright calls. Only this
		// fixture's driver may be killed if it fails to exit after that.
		stopped := make(chan error, 1)
		pid := playwright.Pid()
		finish := testPhase(t, "Playwright driver stop")
		go func() { stopped <- playwright.Stop() }()
		if !waitCleanup(t, "Playwright driver stop", stopped) {
			if pid > 0 {
				process, err := os.FindProcess(pid)
				if err == nil {
					err = process.Kill()
					if releaseErr := process.Release(); releaseErr != nil {
						t.Errorf("release Playwright driver process handle: %v", releaseErr)
					}
				}
				if err != nil && !errors.Is(err, os.ErrProcessDone) {
					t.Errorf("kill stalled Playwright driver %d: %v", pid, err)
				}
			}
			requireCleanup(t, "Playwright driver exit after kill", stopped)
		}
		finish()
		if !browserClosed {
			requireCleanup(t, "browser cleanup after driver stop", closed)
		}
	})
	t.Cleanup(cleanup)

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
	finish = testPhase(t, "Chrome launch")
	browser, err = playwright.Chromium.Launch(launchOpts)
	finish()
	if err != nil {
		t.Fatalf("browser launch: %v", err)
	}

	finish = testPhase(t, "browser context setup")
	ctx, err = browser.NewContext(pw.BrowserNewContextOptions{
		IgnoreHttpsErrors: pw.Bool(true),
	})
	finish()
	if err != nil {
		t.Fatalf("new context: %v", err)
	}

	for _, opt := range opts {
		if err := opt(ctx); err != nil {
			t.Fatalf("page option: %v", err)
		}
	}

	finish = testPhase(t, "browser page setup")
	page, err = ctx.NewPage()
	finish()
	if err != nil {
		t.Fatalf("new page: %v", err)
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
	defer testPhase(t, "tether connection")()
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
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case uint:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
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
