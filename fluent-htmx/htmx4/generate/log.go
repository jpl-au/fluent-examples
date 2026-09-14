package generate

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// template defines a log entry prototype. The generator picks one at
// random, stamps it with the current time, and fills in dynamic
// attribute values from the attrs function.
type template struct {
	level   string
	message string
	attrs   func() string
}

// pool is the set of log entry prototypes the generator picks from.
// Each entry mimics a different subsystem so the feed looks like a
// real application rather than a single repeated message.
var pool = []template{
	{"INFO", "request completed", requestAttrs},
	{"WARN", "slow query detected", queryAttrs},
	{"DEBUG", "cache operation", cacheAttrs},
	{"ERROR", "connection refused", connectionAttrs},
	{"INFO", "job completed", jobAttrs},
	{"WARN", "rate limit approaching", rateLimitAttrs},
	{"DEBUG", "auth token issued", authAttrs},
}

// LogEntry creates a single fake log entry with the current timestamp.
func LogEntry() Entry {
	t := pool[rand.IntN(len(pool))]
	return Entry{
		Time:    time.Now().Format("15:04:05.000"),
		Level:   t.level,
		Message: t.message,
		Attrs:   t.attrs(),
	}
}

// Jitter returns a random duration between 2 and 5 seconds. Use it
// as the delay between entries so the feed has an organic rhythm
// rather than a fixed tick.
func Jitter() time.Duration {
	return time.Duration(2000+rand.IntN(3000)) * time.Millisecond
}

// --- Attribute generators ---
//
// Each function produces realistic key-value pairs for a specific log
// category. Values are randomised within plausible ranges.

func requestAttrs() string {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	paths := []string{"/api/users", "/api/orders", "/api/products", "/api/health", "/api/auth/login"}
	statuses := []int{200, 200, 200, 201, 204, 301, 400, 404, 500}
	return fmt.Sprintf("method=%s path=%s duration=%dms status=%d",
		methods[rand.IntN(len(methods))],
		paths[rand.IntN(len(paths))],
		rand.IntN(200)+1,
		statuses[rand.IntN(len(statuses))],
	)
}

func queryAttrs() string {
	tables := []string{"orders", "users", "products", "sessions", "audit_log"}
	return fmt.Sprintf("table=%s duration=%dms rows=%d",
		tables[rand.IntN(len(tables))],
		rand.IntN(800)+200,
		rand.IntN(5000)+1,
	)
}

func cacheAttrs() string {
	ops := []string{"hit", "miss", "evict", "set"}
	keys := []string{"user:1234", "session:abc", "product:99", "config:main", "rate:10.0.0.1"}
	return fmt.Sprintf("op=%s key=%s ttl=%ds",
		ops[rand.IntN(len(ops))],
		keys[rand.IntN(len(keys))],
		rand.IntN(3600),
	)
}

func connectionAttrs() string {
	hosts := []string{"db-primary", "db-replica-2", "redis-01", "queue-broker", "auth-service"}
	return fmt.Sprintf("host=%s retry=%d backoff=%dms",
		hosts[rand.IntN(len(hosts))],
		rand.IntN(5)+1,
		rand.IntN(5000)+500,
	)
}

func jobAttrs() string {
	jobs := []string{"email-digest", "report-gen", "cleanup", "sync-inventory", "recalc-scores"}
	return fmt.Sprintf("job=%s items=%d duration=%dms",
		jobs[rand.IntN(len(jobs))],
		rand.IntN(1000)+1,
		rand.IntN(10000)+100,
	)
}

func rateLimitAttrs() string {
	return fmt.Sprintf("client=10.0.%d.%d requests=%d limit=%d window=60s",
		rand.IntN(255), rand.IntN(255),
		rand.IntN(50)+80,
		100,
	)
}

func authAttrs() string {
	users := []string{"alice", "bob", "carol", "dave", "eve"}
	return fmt.Sprintf("user=%s scope=read,write ttl=%dm",
		users[rand.IntN(len(users))],
		rand.IntN(60)+15,
	)
}
