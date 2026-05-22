package webserver

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distEmbed embed.FS

// distFS returns the embedded SPA, with the `dist` prefix stripped.
func distFS() (fs.FS, error) {
	return fs.Sub(distEmbed, "dist")
}
