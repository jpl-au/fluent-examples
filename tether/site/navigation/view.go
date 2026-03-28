package navigation

import (
	"strconv"
	"strings"

	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
	"github.com/jpl-au/fluent-examples/tether/components/simple/result"
)

// Render builds the navigation demo page.
func Render(s State) node.Node {
	if s.Page == "/navigation/target/" {
		return renderTarget()
	}
	return renderMain(s)
}

// renderMain builds the main navigation page with all demos.
func renderMain(s State) node.Node {
	return page.New(
		panel.Card("Client-Side Navigation", "Click any link to navigate without a full page reload. The browser URL updates via pushState and the new page HTML loads via a fetch POST. This example includes a target page to navigate to and back.", "bind.Link", panel.AllTransports,
			layout.Row(
				button.Nav("Go to Target Page", "/navigation/target/", bind.Link()),
				button.Nav("Stay Here (reload)", "/navigation/", bind.Link()),
			),
		),

		panel.Card("Query Parameters", "Click a link - the URL changes and OnNavigate reads the query string into state using typed helpers. The server uses p.Get for strings, p.IntDefault for optional integers with a default, and displays what it received. This is how you derive state from the URL: filters, pagination, search terms.", "Params.Get · Params.IntDefault", panel.AllTransports,
			layout.Stack(
				layout.Row(
					button.Nav("Tab: overview, page 1", "/navigation/?tab=overview&page=1", bind.Link()),
					button.Nav("Tab: settings, page 3", "/navigation/?tab=settings&page=3", bind.Link()),
					button.Nav("Search: hello, page 2", "/navigation/?q=hello&page=2", bind.Link()),
					button.Nav("No params (defaults)", "/navigation/", bind.Link()),
				),
				queryParamsResult(s),
			),
		),

		panel.Card("Multi-Value Parameters", "Click a link with repeated query keys. OnNavigate uses p.Strings to collect all values for the same key into a slice.", "Params.Strings", panel.AllTransports,
			layout.Stack(
				layout.Row(
					button.Nav("Tags: go, web, sse", "/navigation/?tag=go&tag=web&tag=sse", bind.Link()),
					button.Nav("Tag: react", "/navigation/?tag=react", bind.Link()),
					button.Nav("No tags", "/navigation/", bind.Link()),
				),
				multiValueResult(s.Tags),
			),
		),

		panel.Card("Typed Parameters", "Click a link with boolean and float query values. OnNavigate uses p.BoolDefault and p.Float64Default to parse them with defaults when absent.", "Params.BoolDefault · Params.Float64Default", panel.AllTransports,
			layout.Stack(
				layout.Row(
					button.Nav("Active, price 9.99", "/navigation/?active=true&price=9.99", bind.Link()),
					button.Nav("Inactive, price 4.50", "/navigation/?active=false&price=4.50", bind.Link()),
					button.Nav("No params (defaults)", "/navigation/", bind.Link()),
				),
				typedParamsResult(s),
			),
		),

		panel.Card("Multi-Value Numeric Parameters", "Click a link with repeated numeric query keys. OnNavigate uses p.Ints and p.Float64s to collect all values into typed slices - the numeric counterparts of p.Strings.", "Params.Ints · Params.Float64s", panel.AllTransports,
			layout.Stack(
				layout.Row(
					button.Nav("Quantities: 1, 2, 5", "/navigation/?qty=1&qty=2&qty=5", bind.Link()),
					button.Nav("Prices: 9.99, 24.50, 3.00", "/navigation/?price=9.99&price=24.50&price=3.00", bind.Link()),
					button.Nav("No params", "/navigation/", bind.Link()),
				),
				numericMultiValueResult(s),
			),
		),

		panel.Card("ReplaceURL", "Click a button and watch the browser URL bar - the query string changes, but pressing Back will not return to the previous URL. This uses replaceState instead of pushState, which is useful for filters and tab selections that shouldn't pollute navigation history.", "Session.ReplaceURL", panel.WS|panel.SSE,
			layout.Row(
				button.PrimaryAction("Set ?tab=a", "nav.replace-a"),
				button.PrimaryAction("Set ?tab=b", "nav.replace-b"),
			),
			span.Text("Check the browser URL bar after clicking."),
		),

		panel.Card("Programmatic Navigation", "Click the button - the server decides where to navigate and tells the browser. Unlike bind.Link where the destination is in the HTML, here the navigation target is determined server-side in the event handler.", "Session.Navigate", panel.WS|panel.SSE,
			button.PrimaryAction("Navigate to Target Page", "nav.goto-target"),
		),
	)
}

// renderTarget is a simple landing page for bind.Link and Navigate
// demos. It exists so there is somewhere to navigate to and back.
func renderTarget() node.Node {
	return page.New(
		panel.Card("Target Page", "You navigated here via bind.Link or Session.Navigate. Click below to go back.", "bind.Link", panel.AllTransports,
			layout.Stack(
				button.NavPrimary("Back to Navigation Demos", "/navigation/", bind.Link()),
			),
		),
	)
}

// queryParamsResult renders the Get/IntDefault extraction results so the
// user can see how each typed getter parses the query string.
func queryParamsResult(s State) node.Node {
	tab := s.Tab
	if tab == "" {
		tab = "(empty)"
	}
	search := s.Search
	if search == "" {
		search = "(empty)"
	}
	return result.BlockDynamic("query-params",
		"tab    = "+tab+"\n"+
			"search = "+search+"\n"+
			"page   = "+strconv.Itoa(s.QueryPage)+" (default: 1)",
	)
}

// multiValueResult renders the Params.Strings extraction result,
// showing all values collected from repeated query keys.
func multiValueResult(tags []string) node.Node {
	txt := "(none)"
	if len(tags) > 0 {
		txt = strings.Join(tags, ", ")
	}
	return result.BlockDynamic("multi-value-params", "tags = ["+txt+"]")
}

// typedParamsResult renders the BoolDefault and Float64Default values so the
// user can see how defaults apply when the query key is absent.
func typedParamsResult(s State) node.Node {
	return result.BlockDynamic("typed-params",
		"active = "+strconv.FormatBool(s.Active)+" (default: false)\n"+
			"price  = "+strconv.FormatFloat(s.Price, 'f', 2, 64)+" (default: 0.00)",
	)
}

// numericMultiValueResult renders the Ints and Float64s extraction
// results, showing typed slices collected from repeated query keys.
func numericMultiValueResult(s State) node.Node {
	qtys := "(none)"
	if len(s.Quantities) > 0 {
		parts := make([]string, len(s.Quantities))
		for i, q := range s.Quantities {
			parts[i] = strconv.Itoa(q)
		}
		qtys = strings.Join(parts, ", ")
	}
	prices := "(none)"
	if len(s.Prices) > 0 {
		parts := make([]string, len(s.Prices))
		for i, p := range s.Prices {
			parts[i] = strconv.FormatFloat(p, 'f', 2, 64)
		}
		prices = strings.Join(parts, ", ")
	}
	return result.BlockDynamic("numeric-multi-value-params",
		"quantities = ["+qtys+"]\n"+
			"prices     = ["+prices+"]",
	)
}
