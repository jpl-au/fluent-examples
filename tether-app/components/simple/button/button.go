// Package button provides styled button components. Each variant
// returns a fully wired node.Node - call sites never apply raw CSS.
package button

import (
	el "github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"
)

// Primary creates a primary button with optional bind options.
func Primary(label string, opts ...bind.Option) node.Node {
	n := el.New().Class("btn btn-primary").Text(label)
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// PrimaryAction creates a primary button that fires action on click.
func PrimaryAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(
		el.New().Class("btn btn-primary").Text(label),
		prepend(bind.OnClick(action), opts)...,
	)
}

// Secondary creates a secondary button with optional bind options.
func Secondary(label string, opts ...bind.Option) node.Node {
	n := el.New().Class("btn btn-secondary").Text(label)
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// SecondaryAction creates a secondary button that fires action on click.
func SecondaryAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(
		el.New().Class("btn btn-secondary").Text(label),
		prepend(bind.OnClick(action), opts)...,
	)
}

// DangerAction creates a danger button that fires action on click.
func DangerAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(
		el.New().Class("btn btn-danger").Text(label),
		prepend(bind.OnClick(action), opts)...,
	)
}

// SmallAction creates a small secondary button that fires action on click.
func SmallAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(
		el.New().Class("btn btn-secondary btn-sm").Text(label),
		prepend(bind.OnClick(action), opts)...,
	)
}

// Submit creates a primary submit button for forms.
func Submit(label string, opts ...bind.Option) node.Node {
	n := el.Submit(label).Class("btn btn-primary")
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

func prepend(first bind.Option, rest []bind.Option) []bind.Option {
	out := make([]bind.Option, 0, len(rest)+1)
	return append(append(out, first), rest...)
}
