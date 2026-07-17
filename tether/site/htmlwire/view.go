package htmlwire

import (
	"strconv"

	"github.com/jpl-au/fluent/html5/code"
	"github.com/jpl-au/fluent/html5/pre"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// curlExample is the copy-paste command shown on the page. The
// endpoint answers with readable HTML rather than a JSON envelope.
const curlExample = `curl -si -X POST localhost:8080/html-wire/ \
  -H 'Content-Type: application/json' \
  -d '{"type":"click","action":"htmlwire.fragment","data":{"count":"4"}}'`

// Render builds the HTML wire format demo.
func Render(s State) node.Node {
	counterStr := strconv.Itoa(s.Counter)
	return cpage.New(
		panel.Card("Plain-HTML Responses",
			"This page sets WireFormat: wire.HTML, so POST event responses are "+
				"plain HTML instead of the JSON envelope. 'Targeted fragment' calls "+
				"sess.Morph(\"hw-count\") - the response body is just the keyed "+
				"<span> plus a Tether-Morph: keyed header. 'Full page' sends the "+
				"whole rendered page with the toast riding in a JSON effects island "+
				"appended to the HTML. Watch the Network tab: both responses are "+
				"readable HTML.",
			"wire.HTML · sess.Morph", panel.HTTP,
			layout.Row(
				button.PrimaryAction("Targeted fragment", "htmlwire.fragment",
					bind.EventData("count", counterStr),
				),
				button.SecondaryAction("Full page + toast", "htmlwire.full",
					bind.EventData("count", counterStr),
				),
				span.Text("Count: "+counterStr).Dynamic("hw-count"),
			),
		),

		panel.Card("Inspect It With curl",
			"HTML responses need no client at all - the endpoint is plain HTTP. "+
				"Run this from a terminal and read the fragment straight off the wire:",
			"", panel.HTTP,
			pre.New(code.Text(curlExample)).Class("demo-raw"),
			hint.Text("The initial GET also sends Cache-Control: public, max-age=30 "+
				"via StatelessConfig.CacheControl - stateless pages carry no session "+
				"token, so browsers and CDNs may cache them. (Dev mode forces "+
				"no-store, so you will see the cache header in production builds.)"),
		),
	)
}
