package playwright_test

import (
	"strings"
	"testing"
)

// TestSecurityPageRenders verifies the security demo page loads and
// the two panel headings are visible.
func TestSecurityPageRenders(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/security/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	if resp.Status() != 200 {
		t.Errorf("status = %d, want 200", resp.Status())
	}

	clean := page.GetByText("security.HTML - sanitising untrusted HTML")
	if err := expect(clean).ToBeVisible(); err != nil {
		t.Errorf("security.HTML panel heading not visible: %v", err)
	}

	nonce := page.GetByText("Nonce + Content-Security-Policy")
	if err := expect(nonce).ToBeVisible(); err != nil {
		t.Errorf("Nonce panel heading not visible: %v", err)
	}
}

// TestSecuritySanitisedColumnsDropDangerousConstructs walks every
// sanitised-output column on the page and verifies no <script> tag,
// javascript: URI, or inline event handler survived the sanitiser.
// If any attack fixture makes it through, bluemonday's UGC policy has
// changed (or our integration has regressed) and the demo is lying.
func TestSecuritySanitisedColumnsDropDangerousConstructs(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(srv + "/security/"); err != nil {
		t.Fatalf("goto: %v", err)
	}

	cols := page.Locator(".demo-rendered")
	count, err := cols.Count()
	if err != nil {
		t.Fatalf("count .demo-rendered: %v", err)
	}
	if count == 0 {
		t.Fatal("no .demo-rendered columns found; view may have regressed")
	}

	for i := range count {
		inner, err := cols.Nth(i).InnerHTML()
		if err != nil {
			t.Fatalf("inner html col %d: %v", i, err)
		}
		low := strings.ToLower(inner)
		if strings.Contains(low, "<script") {
			t.Errorf("sanitised column %d still contains <script>: %s", i, inner)
		}
		if strings.Contains(low, "javascript:") {
			t.Errorf("sanitised column %d still contains javascript: URI: %s", i, inner)
		}
		if strings.Contains(low, "onclick=") || strings.Contains(low, "onerror=") {
			t.Errorf("sanitised column %d still contains inline event handler: %s", i, inner)
		}
	}
}

// TestSecurityNonceMatchesCSPHeader is the integration assertion that
// the whole nonce/CSP wiring actually holds together: the inline
// <script nonce="..."> served in the HTTP response body must carry
// the same value as the 'nonce-...' token in the Content-Security-
// Policy response header. If they diverge, the browser would silently
// refuse to execute the inline script - the demo would look fine but
// the point of the demo (inline script permitted by nonce) would be
// broken.
//
// The assertion reads the nonce from the raw response body rather
// than the live DOM because browsers hide the nonce attribute from
// scripts once CSP has been applied (a W3C mitigation to stop CSS
// attribute selectors from exfiltrating the value). Any
// GetAttribute("nonce") call against the rendered page would return
// an empty string even when the attribute was present in the wire
// HTML, so the only reliable source is the server's output.
func TestSecurityNonceMatchesCSPHeader(t *testing.T) {
	srv := startApp(t, serverMode())
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/security/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}

	csp, err := resp.HeaderValue("content-security-policy")
	if err != nil {
		t.Fatalf("header value: %v", err)
	}
	if csp == "" {
		t.Fatal("no Content-Security-Policy header on /security/ response")
	}

	headerNonce := extractNonce(csp)
	if headerNonce == "" {
		t.Fatalf("could not find nonce in CSP header: %q", csp)
	}

	body, err := resp.Body()
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	scriptNonce := extractScriptNonce(string(body))
	if scriptNonce == "" {
		t.Fatalf("no <script nonce=...> found in server response body")
	}

	if scriptNonce != headerNonce {
		t.Errorf("nonce mismatch:\n  header = %q\n  script = %q\nbrowser would reject the inline script", headerNonce, scriptNonce)
	}
}

// extractScriptNonce finds the first <script ... nonce="value">
// token in the raw HTML response body. Deliberately simple string
// search - the security demo's HTML is well-formed and only has one
// nonce'd script; nothing fancier is needed.
func extractScriptNonce(html string) string {
	i := strings.Index(html, "<script")
	for i >= 0 {
		// Find the end of this tag.
		end := strings.Index(html[i:], ">")
		if end < 0 {
			return ""
		}
		tag := html[i : i+end]
		if _, after, ok := strings.Cut(tag, ` nonce="`); ok {
			rest := after
			if before, _, ok := strings.Cut(rest, "\""); ok {
				return before
			}
		}
		// Move past this tag and look for the next <script.
		html = html[i+end+1:]
		i = strings.Index(html, "<script")
	}
	return ""
}

// TestSecurityNonceIsFreshPerRequest verifies that two requests to
// the same URL produce different nonces. Reuse would defeat the
// freshness property CSP relies on.
func TestSecurityNonceIsFreshPerRequest(t *testing.T) {
	srv := startApp(t, serverMode())

	first := fetchNonceFromCSP(t, srv)
	second := fetchNonceFromCSP(t, srv)

	if first == "" || second == "" {
		t.Fatalf("expected nonces on both requests, got %q and %q", first, second)
	}
	if first == second {
		t.Errorf("nonce reused across requests: %q - CSP freshness defeated", first)
	}
}

// extractNonce pulls the nonce value out of a CSP string of the form
// "...; script-src 'self' 'nonce-<value>'; ...". Returns "" if no
// 'nonce-' token is present. Deliberately simple: the demo CSP is
// fixed-format, so a strings.Index is enough.
func extractNonce(csp string) string {
	const prefix = "'nonce-"
	_, after, ok := strings.Cut(csp, prefix)
	if !ok {
		return ""
	}
	rest := after
	before, _, ok := strings.Cut(rest, "'")
	if !ok {
		return ""
	}
	return before
}

// fetchNonceFromCSP opens /security/ in a fresh page and returns the
// nonce value lifted from the response's CSP header.
func fetchNonceFromCSP(t *testing.T, srv string) string {
	t.Helper()
	page, cleanup := newPage(t)
	defer cleanup()

	resp, err := page.Goto(srv + "/security/")
	if err != nil {
		t.Fatalf("goto: %v", err)
	}
	csp, err := resp.HeaderValue("content-security-policy")
	if err != nil {
		t.Fatalf("header value: %v", err)
	}
	return extractNonce(csp)
}
