package plugin

import (
	"os"
	"path/filepath"
)

// DefaultDir returns the directory the loader scans by default.
// Honors $DEVFORGE_PLUGIN_DIR; falls back to ~/.devforge/plugins.
func DefaultDir() string {
	if d := os.Getenv("DEVFORGE_PLUGIN_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".devforge", "plugins")
}
