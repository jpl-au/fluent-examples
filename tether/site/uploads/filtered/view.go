package filtered

import (
	"fmt"

	"github.com/jpl-au/fluent/html5/attr/accept"
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

// Render builds the filtered upload page, demonstrating MIME type
// filtering via UploadConfig.Accept.
func Render(s State) node.Node {
	return cpage.New(
		panel.Card(
			"Filtered Upload",
			"Two layers of filtering protect this endpoint. "+
				"The HTML accept attribute on the file input restricts the browser's file picker to images and PDFs. "+
				"If a disallowed file gets through (e.g. via drag-and-drop or browser override), "+
				"the server's UploadConfig.Accept rejects it with 415 Unsupported Media Type before the handler callback runs.",
			"tether.UploadConfig · Accept", panel.AllTransports,
			layout.Row(
				field.File("filtered-upload-input", "file").
					Accept(accept.ImageWildcard, accept.MimePDF),
				button.Primary("Upload",
					bind.Upload("upload"),
					bind.UploadInput("#filtered-upload-input"),
				),
			),
		),

		panel.Card(
			"Accepted Files",
			"Files that pass the MIME type filter appear here. "+
				"The content type is shown alongside each entry so you can verify the filter is working.",
			"sess.Update", panel.WS|panel.SSE,
			layout.Stack(
				fileList(s.Uploads),
				button.SmallAction("Clear List", "clear"),
			),
		),
	)
}

// fileList renders uploaded files or a placeholder.
func fileList(files []FileEntry) node.Node {
	if len(files) == 0 {
		return layout.Container(
			hint.Text("No files uploaded yet. Try an image or PDF."),
		).Dynamic("filtered-uploads")
	}
	nodes := make([]node.Node, len(files))
	for i, f := range files {
		nodes[i] = upload.ItemWithType(f.Name, formatSize(f.Size), f.ContentType)
	}
	return layout.Container(list.New(nodes...)).Dynamic("filtered-uploads")
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
