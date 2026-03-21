// Package monitor provides layout elements for the real-time system
// monitor dashboard - a chart grid container and individual chart panels.
package monitor

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
)

// Charts creates the chart grid container. Returns the concrete
// element so callers can chain .Dynamic().
func Charts(children ...node.Node) *div.Element {
	return div.New(children...).Class("monitor-charts")
}

// Chart creates a single chart panel with the given ID. Returns the
// concrete element so callers can chain .SetAttribute() and .SetData()
// for the echarts hook wiring.
func Chart(id string) *div.Element {
	return div.New().ID(id).Class("monitor-chart")
}
