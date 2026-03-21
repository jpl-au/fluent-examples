# fluent-examples

Getting-started projects for the Fluent ecosystem. Each directory is a
standalone Go module demonstrating one framework area. The first three
examples build on the same contact manager application, each layering a
different concern.

| Directory | Framework | Description |
|-----------|-----------|-------------|
| [fluent](fluent/) | [fluent](https://github.com/jpl-au/fluent) | Server-rendered contact manager - pure Go, pure HTML |
| [fluent-jit](fluent-jit/) | [fluent-jit](https://github.com/jpl-au/fluent-jit) | Same contact manager with JIT optimisation strategies |
| [fluent-htmx](fluent-htmx/) | [fluent-htmx](https://github.com/jpl-au/fluent-htmx) | Same contact manager with HTMX partial page updates |
| [tether](tether/) | [tether](https://github.com/jpl-au/tether) | Reactive server-driven UI - Feature Explorer |

## Prerequisites

- Go 1.25+

## Running

Each example is a standalone module. To run any example:

```bash
cd fluent        # or fluent-jit, fluent-htmx, tether
go run .         # starts on :8080 by default
```

## Testing

Unit tests run with no extra dependencies:

```bash
cd tether
go test ./...
```

### Playwright browser tests

The tether example includes end-to-end browser tests in `tether/playwright/`.
These require the Playwright driver and Chrome or Chromium. If either is
missing, the tests skip automatically.

To install the Playwright driver:

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install
```

To run the browser tests:

```bash
cd tether
go test -v ./playwright/...
```

To run under HTTP/2 instead of the default HTTP/1.1:

```bash
TETHER_PROTO=HTTP2 go test -v ./playwright/...
```

## Licence

MIT - see [LICENSE](LICENSE).
