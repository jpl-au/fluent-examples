package security

import (
	fsec "github.com/jpl-au/fluent-security"
	"github.com/jpl-au/fluent/html5/code"
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/html5/pre"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"

	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the security demo page. Two panels: security.HTML
// against a short list of attack fixtures showing before/after output,
// and a Nonce/CSP panel quoting the per-request nonce carried in State.
func Render(s State) node.Node {
	return cpage.New(
		panel.Card(
			"security.HTML - sanitising untrusted HTML",
			"Each row below pairs a raw input (what a user typed, or an attacker tried) with the result of running it through fluent-security's HTML helper. It wraps bluemonday's UGCPolicy: safe formatting tags survive, scripts and dangerous URIs are stripped before the content ever reaches the DOM.",
			"security.HTML", panel.AllTransports,
			cleanRows(),
		),
		panel.Card(
			"Nonce + Content-Security-Policy",
			"This page generates a fresh CSP nonce for every request and uses it twice: once in the Content-Security-Policy header, and once on the inline <script> in the page <head>. Because both values come from the same security.Nonce() call, the browser permits the inline script and logs a message to the console. Open devtools to verify.",
			"security.Nonce", panel.AllTransports,
			nonceSection(s.Nonce),
		),
	)
}

// fixtures drive the demo. Each fixture is a short title, the raw
// input, and a brief note on what the sanitiser is expected to do.
// The rendering is produced live by security.HTML, so the demo
// reflects whatever bluemonday's UGCPolicy actually does today.
type fixture struct {
	Title string
	Raw   string
	Note  string
}

var fixtures = []fixture{
	{
		Title: "Plain text",
		Raw:   "A boring comment with no HTML.",
		Note:  "Nothing to strip; the text passes through unchanged.",
	},
	{
		Title: "Safe formatting",
		Raw:   `Working on the <strong>fluent</strong> docs - see <a href="https://example.com">example.com</a> for details.`,
		Note:  "Formatting tags and http/https links survive the UGC policy.",
	},
	{
		Title: "Inline script",
		Raw:   `Before <script>alert('xss')</script> after.`,
		Note:  "The <script> tag is removed entirely; surrounding text remains.",
	},
	{
		Title: "Event handler attribute",
		Raw:   `<a href="/ok" onclick="steal()">click me</a>`,
		Note:  "The onclick handler is stripped; the link keeps its safe href and text.",
	},
	{
		Title: "javascript: URI",
		Raw:   `<a href="javascript:alert(1)">click</a>`,
		Note:  "The javascript: scheme is rejected and the href (and often the tag) is dropped.",
	},
	{
		Title: "SVG-smuggled script",
		Raw:   `<svg><script>alert('smuggled')</script></svg>`,
		Note:  "UGCPolicy does not allow <svg>; the whole block reduces to inert text.",
	},
}

// cleanRows renders each fixture as a three-column row: title,
// raw source (in a <pre><code> block), and sanitised output. Using
// <pre><code> for the raw source shows what was typed without
// rendering it; the sanitised column is live HTML.
func cleanRows() node.Node {
	rows := make([]node.Node, 0, len(fixtures))
	for _, f := range fixtures {
		rows = append(rows, div.New(
			span.Text(f.Title).Class("demo-row-title"),
			div.New(
				span.Text("Raw input").Class("demo-col-label"),
				pre.New(code.Text(f.Raw)).Class("demo-raw"),
			).Class("demo-col"),
			div.New(
				span.Text("Sanitised output").Class("demo-col-label"),
				div.New(fsec.HTML(f.Raw)).Class("demo-rendered"),
			).Class("demo-col"),
			p.Text(f.Note).Class("demo-row-note"),
		).Class("demo-row"))
	}
	return div.New(rows...).Class("demo-grid")
}

// nonceSection quotes the per-request nonce so the user can see the
// value that was stamped on both the inline script and the CSP
// header. Devtools -> Network -> Response Headers shows the header
// side; devtools -> Console shows the "permitted by CSP nonce" log
// that fired because the values matched.
func nonceSection(nonce string) node.Node {
	return div.New(
		p.Text("Nonce generated for this request:").Class("demo-label"),
		pre.New(code.Text(nonce)).Class("demo-raw"),
		p.Text("The same value is in the Content-Security-Policy response header under 'script-src' and on the <script nonce=\"...\"> in the page <head>. Reload to get a fresh one.").Class("demo-description"),
	)
}
