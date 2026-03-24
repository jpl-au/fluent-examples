// Package main is the entry point for the kanban board application,
// a real-world tether example demonstrating a single handler with
// WebSocket transport and SSE+POST failover.
package main

import (
	"embed"
	"io/fs"
	"log"

	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether-app/handler"
	"github.com/jpl-au/fluent-examples/tether-app/store"
)

//go:embed static
var staticEmbed embed.FS

func main() {
	staticFS, _ := fs.Sub(staticEmbed, "static")
	assets := &tether.Asset{
		FS:       staticFS,
		Prefix:   "/static/",
		Precache: []string{"app.css"},
	}

	board := store.New()
	h := handler.New(board, assets)

	if err := tether.ListenAndServe("", h, h); err != nil {
		log.Fatal(err)
	}
}
