package list

import (
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/node"
)

// Item creates a single list item with text content.
func Item(text string) node.Node {
	return li.New().Class("list-item").Text(text)
}

// ItemNode creates a list item wrapping arbitrary child nodes.
func ItemNode(children ...node.Node) node.Node {
	return li.New(children...).Class("list-item")
}
