package panel

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Transport is a bitmask encoding which Tether transports support a
// given feature. Combine constants with | to express multi-transport
// availability.
type Transport uint8

const (
	HTTP Transport = 1 << iota
	WS
	SSE
	AllTransports = HTTP | WS | SSE
)

// badges builds the badge group: optional API label followed by the
// three transport indicators.
func badges(api string, t Transport) node.Node {
	var items []node.Node
	if api != "" {
		items = append(items, span.Text(api).Class("api-label"))
	}
	items = append(items,
		transportBadge("HTTP", t&HTTP != 0),
		transportBadge("WS", t&WS != 0),
		transportBadge("SSE", t&SSE != 0),
	)
	return div.New(items...).Class("demo-badges")
}

// transportBadge renders a single transport indicator. Active badges
// use a distinct colour per transport; inactive badges are dimmed.
func transportBadge(label string, active bool) node.Node {
	class := "transport-badge "
	switch label {
	case "HTTP":
		class += "transport-http"
	case "WS":
		class += "transport-ws"
	default:
		class += "transport-sse"
	}
	if !active {
		class += " transport-off"
	}
	return span.Text(label).Class(class)
}
