// Package button provides reusable button components for the example
// application. Each variant returns a fully styled node.Node with
// bind wiring - call sites never apply raw CSS classes.
//
// Base variants (Primary, Secondary) accept a label and optional bind
// options. Action variants (PrimaryAction, SecondaryAction, etc.)
// additionally wire a bind.OnClick handler for the given action.
package button

import (
	"github.com/jpl-au/fluent/html5/a"
	"github.com/jpl-au/fluent/html5/attr/rel"
	"github.com/jpl-au/fluent/html5/attr/target"
	el "github.com/jpl-au/fluent/html5/button"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"
)

// Primary creates a primary button with no default event binding.
func Primary(label string, opts ...bind.Option) node.Node {
	n := el.New().Class("btn btn-primary").Text(label)
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// PrimaryAction creates a primary button that fires the given action on click.
func PrimaryAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-primary").Text(label), prepend(bind.OnClick(action), opts)...)
}

// Secondary creates a secondary button with no default event binding.
func Secondary(label string, opts ...bind.Option) node.Node {
	n := el.New().Class("btn btn-secondary").Text(label)
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// SecondaryAction creates a secondary button that fires the given action on click.
func SecondaryAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-secondary").Text(label), prepend(bind.OnClick(action), opts)...)
}

// DangerAction creates a danger button that fires the given action on click.
func DangerAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-danger").Text(label), prepend(bind.OnClick(action), opts)...)
}

// SmallAction creates a small secondary button that fires the given action.
func SmallAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-secondary btn-sm").Text(label), prepend(bind.OnClick(action), opts)...)
}

// SmallPrimaryAction creates a small primary button that fires the given action.
func SmallPrimaryAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-primary btn-sm").Text(label), prepend(bind.OnClick(action), opts)...)
}

// SmallDangerAction creates a small danger button that fires the given action.
func SmallDangerAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-danger btn-sm").Text(label), prepend(bind.OnClick(action), opts)...)
}

// SmallOutlineAction creates a small outline button that fires the given action.
func SmallOutlineAction(label, action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-outline btn-sm").Text(label), prepend(bind.OnClick(action), opts)...)
}

// Link creates a button-styled anchor that opens href in a new tab.
func Link(label, href string) node.Node {
	return a.New().Href(href).Target(target.Blank).Rel(rel.Rel("noopener")).Class("btn btn-secondary").Text(label)
}

// Nav creates a secondary button-styled anchor for internal navigation.
func Nav(label, href string, opts ...bind.Option) node.Node {
	n := a.New().Href(href).Class("btn btn-secondary").Text(label)
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// NavPrimary creates a primary button-styled anchor for internal navigation.
func NavPrimary(label, href string, opts ...bind.Option) node.Node {
	n := a.New().Href(href).Class("btn btn-primary").Text(label)
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// Submit creates a primary submit button for forms.
func Submit(label string, opts ...bind.Option) node.Node {
	n := el.Submit(label).Class("btn btn-primary")
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// Increment creates a compact "+" button with no default event binding.
func Increment(opts ...bind.Option) node.Node {
	n := el.New().Class("btn btn-primary btn-sm").Text("+")
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// IncrementAction creates a compact "+" button that fires the given action.
func IncrementAction(action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-primary btn-sm").Text("+"), prepend(bind.OnClick(action), opts)...)
}

// Decrement creates a compact "-" button with no default event binding.
func Decrement(opts ...bind.Option) node.Node {
	n := el.New().Class("btn btn-secondary btn-sm").Text("-")
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// DecrementAction creates a compact "-" button that fires the given action.
func DecrementAction(action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-secondary btn-sm").Text("-"), prepend(bind.OnClick(action), opts)...)
}

// Reset creates a compact "Reset" button with no default event binding.
func Reset(opts ...bind.Option) node.Node {
	n := el.New().Class("btn btn-outline btn-sm").Text("Reset")
	if len(opts) == 0 {
		return n
	}
	return bind.Apply(n, opts...)
}

// ResetAction creates a compact "Reset" button that fires the given action.
func ResetAction(action string, opts ...bind.Option) node.Node {
	return bind.Apply(el.New().Class("btn btn-outline btn-sm").Text("Reset"), prepend(bind.OnClick(action), opts)...)
}

// prepend inserts first before the rest of the options.
func prepend(first bind.Option, rest []bind.Option) []bind.Option {
	out := make([]bind.Option, 0, len(rest)+1)
	return append(append(out, first), rest...)
}
