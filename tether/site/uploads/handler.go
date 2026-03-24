package uploads

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

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

// State is the per-session state for the uploads demo.
type State struct {
	// Uploads is the list of files received via the upload handler.
	Uploads []FileEntry
	// OnlineCount tracks connected sessions for the badge.
	OnlineCount int
}

// FileEntry records an uploaded file with its data so it can be
// downloaded again.
type FileEntry struct {
	ID          string
	Name        string
	ContentType string
	Size        int64
	Time        time.Time
	Data        []byte
}

// fileStore holds uploaded files keyed by ID so they can be served
// back for download. Shared across sessions.
var fileStore = struct {
	mu    sync.RWMutex
	files map[string]FileEntry
	seq   int
}{files: make(map[string]FileEntry)}

func storeFile(f FileEntry) string {
	fileStore.mu.Lock()
	defer fileStore.mu.Unlock()
	fileStore.seq++
	f.ID = fmt.Sprintf("f%d", fileStore.seq)
	fileStore.files[f.ID] = f
	return f.ID
}

func loadFile(id string) (FileEntry, bool) {
	fileStore.mu.RLock()
	defer fileStore.mu.RUnlock()
	f, ok := fileStore.files[id]
	return f, ok
}

// ServeFile is an HTTP handler that serves uploaded files for download.
// Mount it at /uploads/files/ so sess.Download can reference them.
func ServeFile(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	f, ok := loadFile(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, f.Name))
	w.Header().Set("Content-Type", f.ContentType)
	w.Write(f.Data)
}

var uploadPresence = shared.NewPresenceCountOnly()

// New creates a handler demonstrating file uploads with real-time
// feedback. The upload arrives via HTTP POST; sess.Update and
// sess.Toast push the result to the client over WebSocket.
func New(app tether.App, assets *tether.Asset) *tether.Handler[State] {
	return tether.Stateful(app, tether.StatefulConfig[State]{
		Name:    "uploads",
		Mode:    mode.WebSocket,
		Upgrade: wsupgrade.Upgrade(),

		InitialState: func(_ *http.Request) State {
			return State{OnlineCount: uploadPresence.OnlineCount.Load()}
		},
		Render: func(s State) node.Node {
			return layout.Shell(layout.SectionLive, "/uploads/", s.OnlineCount, Render(s))
		},
		Handle: Handle,

		Layout: func(_ State, content node.Node) node.Node {
			return html.New(
				head.New(
					meta.UTF8(),
					meta.Viewport("width=device-width, initial-scale=1"),
					title.Static("Tether - Uploads"),
					assets.Stylesheet("app.css"),
				),
				body.New(content),
			).Lang("en")
		},

		Upload: &tether.UploadConfig[State]{
			Handle:  handleUpload,
			MaxSize: 10 << 20, // 10 MB
		},

		Watchers: shared.Watchers[State](uploadPresence,
			func(n int, s State) State { s.OnlineCount = n; return s },
			nil,
		),

		OnConnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("uploads: connected", "id", sess.ID())
			shared.TrackPresence(uploadPresence, sess.ID())
		},
		OnDisconnect: func(sess *tether.StatefulSession[State]) {
			slog.Info("uploads: disconnected", "id", sess.ID())
			shared.UntrackPresence(uploadPresence, sess.ID())
		},
	})
}

// handleUpload appends the file to session state and shows a
// confirmation toast. The upload arrives via HTTP POST - sess.Update
// and sess.Toast push the updated list and notification to the client.
func handleUpload(sess *tether.StatefulSession[State], upload tether.Upload) error {
	slog.Info("upload received", "action", upload.Action, "name", upload.Name, "size", upload.Size)

	f, err := upload.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	entry := FileEntry{
		Name:        upload.Name,
		ContentType: upload.ContentType,
		Size:        upload.Size,
		Time:        time.Now(),
		Data:        data,
	}
	id := storeFile(entry)
	entry.ID = id

	sess.Update(func(s State) State {
		s.Uploads = append(s.Uploads, entry)
		return s
	})

	sess.Toast("Uploaded " + upload.Name)
	return nil
}

// Handle processes events on the uploads page.
func Handle(sess tether.Session, s State, ev tether.Event) State {
	switch ev.Action {
	case "uploads.clear":
		s.Uploads = nil
	case "uploads.download":
		id, _ := ev.Get("id")
		sess.Download("/uploads/files/?id=" + id)
	}
	return s
}
