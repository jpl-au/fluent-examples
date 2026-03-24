// Package upload provides styled elements for file upload displays  -
// list items showing file name, size, and metadata, plus a progress bar.
package upload

import (
	"github.com/jpl-au/fluent/html5/li"
	"github.com/jpl-au/fluent/html5/progress"
	"github.com/jpl-au/fluent/html5/span"
	"github.com/jpl-au/fluent/node"
)

// ItemWithTime renders an upload list entry showing name, size, and timestamp.
func ItemWithTime(name, size, time string) node.Node {
	return li.New(
		span.Text(name).Class("upload-name"),
		span.Text(size).Class("upload-size"),
		span.Text(time).Class("upload-time"),
	).Class("upload-item")
}

// ItemWithType renders an upload list entry showing name, size, and content type.
func ItemWithType(name, size, contentType string) node.Node {
	return li.New(
		span.Text(name).Class("upload-name"),
		span.Text(size).Class("upload-size"),
		span.Text(contentType).Class("upload-type"),
	).Class("upload-item")
}

// Progress creates a styled progress bar for upload tracking. Returns
// the concrete element so callers can chain bind.UploadProgress.
func Progress() *progress.Element {
	return progress.New().Max(100).Class("upload-progress")
}
