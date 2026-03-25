// tether Feature Explorer - a comprehensive example application
// demonstrating every feature of the tether framework, organised
// by feature: each demo is a standalone package under site/.
package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os/signal"
	"syscall"

	tether "github.com/jpl-au/tether"

	"github.com/jpl-au/fluent-examples/tether/app"
)

//go:embed static
var staticEmbed embed.FS

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	staticFS, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		log.Fatal(err)
	}
	assets := &tether.Asset{
		FS:       staticFS,
		Prefix:   "/static/",
		Precache: []string{"app.css", "hooks.js"},
	}

	mux, drainables := app.New(ctx, assets)

	if err := tether.ListenAndServe("", mux, drainables...); err != nil {
		log.Fatal(err)
	}
}
