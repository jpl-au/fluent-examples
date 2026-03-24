// Package panel provides content panels for the example application.
// Card is the primary composite panel with title, badges, description,
// and content. The remaining functions provide simpler bordered panels
// for signal demos, toggle targets, and hook demonstrations.
package panel

import (
	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/html5/h3"
	"github.com/jpl-au/fluent/html5/p"
	"github.com/jpl-au/fluent/node"
)

// Card renders a content panel with a title, API label badge,
// transport compatibility badges, optional description, and content.
func Card(title, description, api string, transport Transport, children ...node.Node) node.Node {
	nodes := []node.Node{
		div.New(
			h3.New().Class("demo-title").Text(title),
			badges(api, transport),
		).Class("demo-header"),
	}
	if description != "" {
		nodes = append(nodes, p.Text(description).Class("demo-description"))
	}
	nodes = append(nodes, div.New(children...).Class("demo-content"))
	return div.New(nodes...).Class("demo")
}

// Signal renders a bordered panel for signal content. Returns the
// concrete element so callers can chain bind.Apply or .Dynamic().
func Signal(children ...node.Node) *div.Element {
	return div.New(children...).Class("signal-panel")
}

// SignalText renders a signal panel with a single text message.
func SignalText(s string) *p.Element {
	return p.Text(s).Class("signal-panel")
}

// SignalSuccess renders a signal panel with success styling.
func SignalSuccess(s string) *p.Element {
	return p.Text(s).Class("signal-panel result-success")
}

// SignalMuted renders a signal panel with muted styling.
func SignalMuted(s string) *p.Element {
	return p.Text(s).Class("signal-panel result-muted")
}

// SignalFocusTrap renders a signal panel configured as a focus trap
// container for bind.FocusTrap chaining.
func SignalFocusTrap(children ...node.Node) *div.Element {
	return div.New(children...).Class("signal-panel focus-trap-demo")
}

// ToggleDemo renders a panel for CSS class toggling via bind.BindClass.
func ToggleDemo(children ...node.Node) *div.Element {
	return div.New(children...).Class("toggle-demo")
}

// HookTarget renders a styled target element for hook demonstrations.
func HookTarget(children ...node.Node) *div.Element {
	return div.New(children...).Class("hook-target")
}
