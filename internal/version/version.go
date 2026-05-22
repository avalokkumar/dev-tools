// Package version exposes build-time stamped version metadata.
package version

// These vars are overridden via -ldflags at release build time.
// Defaults reflect a local development build.
var (
	Version = "0.0.0-dev"
	Commit  = "none"
	Date    = "unknown"
)

// Info is the structured representation returned to users.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Get returns the current build's Info.
func Get() Info {
	return Info{Version: Version, Commit: Commit, Date: Date}
}
