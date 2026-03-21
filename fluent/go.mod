module github.com/jpl-au/fluent-examples/fluent

go 1.25.0

replace github.com/jpl-au/fluent => ../../fluent

replace github.com/jpl-au/chain => ../../../chain

require (
	github.com/jpl-au/chain v0.0.0-00010101000000-000000000000
	github.com/jpl-au/fluent v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lxzan/gws v1.9.0 // indirect
)
