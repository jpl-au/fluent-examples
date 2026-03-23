# Fluent-JIT Example - Contact Manager

The same contact manager as the [Fluent example](../fluent/), enhanced with [Fluent-JIT](https://github.com/jpl-au/fluent-jit) optimisation strategies. Same functionality, same UI - the JIT wrappers are the only code change.

This is a **teaching example** that applies every JIT strategy to a single
app so you can compare them side by side. A real application would choose
one or two strategies based on its rendering profile rather than using all
three.

## JIT Strategies Demonstrated

### Compile - pre-render static portions

Used for the contact list and detail pages. The tree structure stays the same across renders - names, emails, and note content change, but the surrounding HTML (divs, classes, labels) is frozen after the first render.

```go
// Global API - string-keyed registry
jit.Compile("contacts-list", content, w)

// Instance API - fine-grained control
var compiler = jit.NewCompiler()
compiler.Render(content, w)
```

### Tune - adaptive buffer sizing

Used for the page layout shell and edit forms. Tune learns the optimal buffer size over repeated renders without analysing the tree structure.

```go
jit.Tune("page", doc, w)
```

### Flatten - fully static to raw bytes

Used for the new contact form and the 404 page. These pages contain no dynamic content - every element uses Static(). Flatten pre-renders the entire tree to a []byte once and returns it directly on every subsequent call.

```go
// Global API
jit.Flatten("not-found", content, w)

// Instance API
flattener, err := jit.NewFlattener(content)
flattener.Render(w)
```

## Run

```bash
go run .
# Visit http://localhost:8080
```

## Structure

Same structure as the Fluent example - see [../fluent/README.md](../fluent/README.md). The only code differences are in `layout/layout.go` and `handler/contacts.go` where JIT wrappers are applied.
