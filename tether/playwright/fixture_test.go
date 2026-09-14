package playwright_test

import (
	"bytes"
	"fmt"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

// A closed HTTP listener does not stop a tether handler's pending-session
// reaper. Exercise the real application fixture and check that each instance
// releases its handlers, even when no browser ever connected.
func TestServerFixtureReleasesHandlers(t *testing.T) {
	reapers := func(t *testing.T) int {
		t.Helper()
		var stacks bytes.Buffer
		if err := pprof.Lookup("goroutine").WriteTo(&stacks, 2); err != nil {
			t.Fatal(err)
		}
		count := 0
		for line := range strings.SplitSeq(stacks.String(), "\n") {
			if strings.HasPrefix(line, "github.com/jpl-au/tether.") && strings.Contains(line, ").reapPending(") {
				count++
			}
		}
		return count
	}
	before := reapers(t)
	for _, mode := range []ServerMode{HTTP1, HTTP2, RealHTTP1} {
		t.Run(fmt.Sprint(mode), func(t *testing.T) {
			startApp(t, mode)
			if count := reapers(t); count <= before {
				t.Fatalf("fixture did not expose its running reapers: before=%d during=%d", before, count)
			}
		})
		// Shutdown cancels the reapers; allow their goroutines to observe it.
		deadline := time.NewTimer(cleanupTimeout)
		tick := time.NewTicker(50 * time.Millisecond)
		for {
			if count := reapers(t); count <= before {
				break
			}
			select {
			case <-tick.C:
			case <-deadline.C:
				tick.Stop()
				t.Fatalf("fixture retained pending-session reapers: before=%d after=%d", before, reapers(t))
			}
		}
		tick.Stop()
		deadline.Stop()
	}
}
