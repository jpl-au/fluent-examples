module github.com/jpl-au/fluent-examples/tether

go 1.25.0

require (
	github.com/alicebob/miniredis/v2 v2.38.0
	github.com/go-echarts/go-echarts/v2 v2.7.2
	github.com/jpl-au/chain v0.1.1
	github.com/jpl-au/fluent v0.5.0
	github.com/jpl-au/fluent-jit v0.3.3
	github.com/jpl-au/fluent-security v0.1.0
	github.com/jpl-au/tether v0.3.3
	github.com/jpl-au/tether/tetheredis v0.1.0
	github.com/playwright-community/playwright-go v0.5700.1
	// Pinned to v9.20.0 like tetheredis: v9.21.0 deadlocks Subscribe (go-redis #3839).
	github.com/redis/go-redis/v9 v9.20.0
)

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/deckarep/golang-set/v2 v2.8.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.4 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/lxzan/gws v1.9.0 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yuin/gopher-lua v1.1.2 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
)
