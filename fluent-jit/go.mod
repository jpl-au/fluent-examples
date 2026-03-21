module github.com/jpl-au/fluent-examples/fluent-jit

go 1.25.0

require (
	github.com/jpl-au/chain v0.1.1
	github.com/jpl-au/fluent v0.2.1
	github.com/jpl-au/fluent-jit v0.2.1
	github.com/lxzan/gws v1.9.0
)

require github.com/klauspost/compress v1.17.9 // indirect

replace (
	github.com/jpl-au/fluent => ../../fluent
	github.com/jpl-au/fluent-jit => ../../fluent-jit
)
