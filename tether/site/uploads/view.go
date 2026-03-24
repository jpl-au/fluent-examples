package uploads

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/div"
	"github.com/jpl-au/fluent/node"
	"github.com/jpl-au/tether/bind"

	"github.com/jpl-au/fluent-examples/tether/components/composite/layout"
	"github.com/jpl-au/fluent-examples/tether/components/composite/list"
	cpage "github.com/jpl-au/fluent-examples/tether/components/composite/page"
	"github.com/jpl-au/fluent-examples/tether/components/composite/upload"
	"github.com/jpl-au/fluent-examples/tether/components/simple/button"
	"github.com/jpl-au/fluent-examples/tether/components/simple/field"
	"github.com/jpl-au/fluent-examples/tether/components/simple/hint"
	"github.com/jpl-au/fluent-examples/tether/components/simple/panel"
)

// Render builds the file uploads page, demonstrating bind.Upload,
// bind.UploadProgress, and the UploadConfig handler.
func Render(s State) node.Node {
	return cpage.New(
		panel.Card(
			"File Upload",
			"File uploads use plain HTTP POST - they are completely transport-agnostic. "+
				"The upload handler uses the live session for real-time feedback: sess.Toast() confirms the upload and sess.Update() adds the file to the list below. "+
				"Both are opt-in choices, not requirements. A minimal handler could simply save the file and return.",
			"bind.Upload · bind.UploadInput · tether.UploadConfig", panel.AllTransports,
			layout.Row(
				field.File("upload-input", "file"),
				button.Primary("Upload",
					bind.Upload("upload"),
					bind.UploadInput("#upload-input"),
				),
			),
		),

		panel.Card(
			"Uploaded Files",
			"The upload handler appends each file to the session state via sess.Update(), which triggers an immediate re-render - the list updates without a page reload. "+
				"Click Download to retrieve the file via sess.Download - the browser fetches it over normal HTTP.",
			"sess.Update · sess.Download", panel.WS|panel.SSE,
			layout.Stack(
				uploadList(s.Uploads),
				button.SmallAction("Clear List", "uploads.clear"),
			),
		),

		panel.Card(
			"Upload Progress",
			"UploadProgress binds a progress element's value to the upload progress signal. "+
				"The client JS updates the signal as bytes are sent, so the bar advances in real time without any server round-trips. "+
				"The signal name is derived from the upload action: upload:{action}:progress.",
			"bind.UploadProgress", panel.AllTransports,
			bind.Apply(upload.Progress(),
				bind.UploadProgress("upload"),
			),
		),
	)
}

// uploadList renders uploaded files or a placeholder.
func uploadList(files []FileEntry) node.Node {
	if len(files) == 0 {
		return layout.Container(
			hint.Text("No files uploaded yet."),
		).Dynamic("uploads")
	}
	nodes := make([]node.Node, len(files))
	for i, f := range files {
		nodes[i] = div.New(
			upload.ItemWithTime(f.Name, formatSize(f.Size), f.Time.Format("15:04:05")),
			button.SmallAction("Download", "uploads.download", bind.EventData("id", f.ID)),
		).Class("layout-row")
	}
	return layout.Container(list.New(nodes...)).Dynamic("uploads")
}

// formatSize converts a byte count to a human-readable string.
func formatSize(b int64) string {
	switch {
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
