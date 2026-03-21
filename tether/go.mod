module github.com/jpl-au/fluent-examples/tether

go 1.25.0

require (
	github.com/go-echarts/go-echarts/v2 v2.7.1
	github.com/jpl-au/chain v0.1.1
	github.com/jpl-au/fluent v0.2.1
	github.com/jpl-au/tether v0.0.0
)

require (
	github.com/deckarep/golang-set/v2 v2.8.0 // indirect
	github.com/dolthub/maphash v0.1.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.4 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/jpl-au/fluent-jit v0.2.1 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/lxzan/gws v1.9.0 // indirect
	github.com/playwright-community/playwright-go v0.5700.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.49.0 // indirect
)

replace (
	github.com/jpl-au/fluent => ../../fluent
	github.com/jpl-au/fluent-jit => ../../fluent-jit
	github.com/jpl-au/tether => ../../tether
)
