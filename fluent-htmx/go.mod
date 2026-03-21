module github.com/jpl-au/fluent-examples/fluent-htmx

go 1.25.0

require (
	github.com/jpl-au/chain v0.1.1
	github.com/jpl-au/fluent v0.2.1
	github.com/jpl-au/fluent-htmx v0.1.0
)

require (
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lxzan/gws v1.9.0 // indirect
)

replace (
	github.com/jpl-au/fluent => ../../fluent
	github.com/jpl-au/fluent-htmx => ../../fluent-htmx
)
