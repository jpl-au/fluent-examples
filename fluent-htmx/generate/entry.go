package generate

// Entry is a single fake log line with a timestamp, severity level,
// human-readable message, and structured key-value attributes.
type Entry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
	Attrs   string `json:"attrs"`
}
