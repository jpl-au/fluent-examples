package playwright_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/alicebob/miniredis/v2"
	pw "github.com/playwright-community/playwright-go"
	"github.com/redis/go-redis/v9"

	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/tetheredis"

	"github.com/jpl-au/fluent-examples/tether/site/clustertest"
)

// startClusterServer creates a tether server with the cluster test
// handler mounted at /_test/cluster/. It starts a plain HTTP server
// and returns its base URL. The server is shut down when the test ends.
func startClusterServer(t *testing.T, cluster tether.Cluster) string {
	t.Helper()

	app := tether.App{
		DevMode: true,
		Cluster: cluster,
	}

	bus := tether.NewBus[clustertest.Message](tether.BusConfig{Topic: "cluster-test"})
	handler := clustertest.New(app, bus)

	mux := http.NewServeMux()
	mux.Handle("/_test/cluster/", handler)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	srv := &http.Server{Handler: mux}
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close() })

	return "http://" + ln.Addr().String()
}

// TestClusterCrossSession starts a server with a Redis-backed cluster,
// opens two browser tabs, sends a message from one, and verifies the
// other receives it. This proves the full stack from cluster
// configuration through bus delivery to DOM update.
func TestClusterCrossSession(t *testing.T) {
	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })

	cluster := tetheredis.New(rdb)
	srv := startClusterServer(t, cluster)

	sender, cleanupSender := newPage(t)
	defer cleanupSender()

	receiver, cleanupReceiver := newPage(t)
	defer cleanupReceiver()

	if _, err := sender.Goto(srv + "/_test/cluster/"); err != nil {
		t.Fatalf("sender goto: %v", err)
	}
	if _, err := receiver.Goto(srv + "/_test/cluster/"); err != nil {
		t.Fatalf("receiver goto: %v", err)
	}

	waitForConnected(t, sender)
	waitForConnected(t, receiver)

	// Type a message in the sender and click Send.
	input := sender.Locator("#cluster-input")
	if err := input.Fill("cluster hello"); err != nil {
		t.Fatalf("fill: %v", err)
	}
	btn := sender.GetByRole("button", pw.PageGetByRoleOptions{Name: "Send"})
	if err := btn.Click(); err != nil {
		t.Fatalf("click: %v", err)
	}

	// The sender sees their own message immediately.
	senderMessages := sender.Locator("[data-tether-key='cluster-messages']")
	if err := expect(senderMessages).ToContainText("cluster hello"); err != nil {
		t.Errorf("sender did not see own message: %v", err)
	}

	// The receiver should see the message via WatchBus.
	receiverMessages := receiver.Locator("[data-tether-key='cluster-messages']")
	if err := expect(receiverMessages).ToContainText("cluster hello"); err != nil {
		t.Errorf("receiver did not see message from other session: %v", err)
	}
}

// TestClusterPageRenders verifies the cluster test page loads and
// connects without a cluster backend. This confirms the handler
// works in isolation.
func TestClusterPageRenders(t *testing.T) {
	app := tether.App{DevMode: true}
	bus := tether.NewBus[clustertest.Message](tether.BusConfig{Topic: "cluster-test-render"})
	handler := clustertest.New(app, bus)

	mux := http.NewServeMux()
	mux.Handle("/_test/cluster/", handler)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := &http.Server{Handler: mux}
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close() })

	base := "http://" + ln.Addr().String()

	page, cleanup := newPage(t)
	defer cleanup()

	if _, err := page.Goto(base + "/_test/cluster/"); err != nil {
		t.Fatalf("goto: %v", err)
	}

	waitForConnected(t, page)

	// The page should show the placeholder text.
	placeholder := page.Locator("[data-tether-key='cluster-messages']")
	if err := expect(placeholder).ToContainText("No messages yet"); err != nil {
		t.Errorf("placeholder not visible: %v", err)
	}
}
