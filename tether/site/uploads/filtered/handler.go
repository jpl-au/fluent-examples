package filtered

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jpl-au/fluent/html5/body"
	"github.com/jpl-au/fluent/html5/head"
	"github.com/jpl-au/fluent/html5/html"
	"github.com/jpl-au/fluent/html5/meta"
	"github.com/jpl-au/fluent/html5/title"
	"github.com/jpl-au/fluent/node"
	tether "github.com/jpl-au/tether"
	"github.com/jpl-au/tether/mode"
	wsupgrade "github.com/jpl-au/tether/ws"

	"github.com/jpl-au/fluent-examples/tether/layout"
	"github.com/jpl-au/fluent-examples/tether/site/shared"
)

// State is the per-session state for the filtered upload demo.
type State struct {
	Uploads     []FileEntry
	OnlineCount int
}

// FileEntry records a successfully uploaded file.
type FileEntry struct {
	Name        string
	Size        int64
	ContentType string
}

var filteredPresence = shared.NewPresenceCountOnly()

// New creates a handler that only accepts image and PDF uploads.
// The Accept field on UploadConfig tells the framework to reject
// anything else with 415 Unsupported Media Type - the Handle
// callback never sees rejected files.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "uploads/filtered",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: filteredPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/uploads/filtered/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Filtered Upload"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Upload: &tether.UploadConfig[State]{
			Handle:  handleUpload,
			MaxSize: 5 << 20, // 5 MB
			Accept:  []string{"image/*", "application/pdf"},
		},

		Watchers: shared.Watchers[State](filteredPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("filtered: connected", "id", sess.ID())
			shared.TrackPresence(filteredPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("filtered: disconnected", "id", sess.ID())
			shared.UntrackPresence(filteredPresence, sess.ID())
		},
	})
}

// handleUpload is called only for files that pass the Accept filter.
func handleUpload(sess *tether.StatefulSession[State], upload tether.Upload) error {
	slog.Info("filtered upload received",
		"name", upload.Name,
		"size", upload.Size,
		"content_type", upload.ContentType,
	)

	sess.Update(func(s State) State {
		s.Uploads = append(s.Uploads, FileEntry{
			Name:        upload.Name,
			Size:        upload.Size,
			ContentType: upload.ContentType,
		})
		return s
	})

	sess.Toast(fmt.Sprintf("Accepted %s (%s)", upload.Name, upload.ContentType))
	return nil
}

// Handle processes events on the filtered upload page.
func Handle(_ tether.Session, s State, ev tether.Event) State {
	if ev.Action == "clear" {
		s.Uploads = nil
	}
	return s
}
