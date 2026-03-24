// Package result provides styled blocks for displaying server-returned
// data - monospace pre-formatted panels for structured output.
package result

import (
	"github.com/jpl-au/fluent/html5/pre"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// Block renders a monospace pre-formatted panel for server data.
func Block(s string) node.Node {
	return pre.New().Class("result-block").Text(s)
}

// BlockDynamic renders a result block with a Tether Dynamic key so
// the differ can track it across re-renders.
func BlockDynamic(key, s string) node.Node {
	return pre.New().Class("result-block").Text(s).Dynamic(key)
}

// Success renders a green result block for positive outcomes.
func Success(s string) node.Node {
	return pre.New().Class("result-block result-success").Text(s)
}

// Danger renders a red result block for errors or warnings.
func Danger(s string) node.Node {
	return pre.New().Class("result-block result-danger").Text(s)
}

// Blue renders a blue result block.
func Blue(s string) node.Node {
	return pre.New().Class("result-block result-blue").Text(s)
}

// Label renders a small muted heading above a result block.
func Label(s string) node.Node {
	return span.Text(s).Class("result-label")
}
