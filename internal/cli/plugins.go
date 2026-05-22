package cli

import (
	"context"
	"log"
	"os"

	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/internal/plugin"
)

// loadPluginsInto scans dir (or the default if empty) for plugins and
// registers their Operations into reg. Errors loading individual plugins are
// logged to stderr but do not abort.
//
// Returns the loaded *Plugin slice so the caller can Close them on shutdown.
func loadPluginsInto(ctx context.Context, reg *mcpserver.Registry, dir string) []*plugin.Plugin {
	if dir == "" {
		dir = plugin.DefaultDir()
	}
	if dir == "" {
		return nil
	}
	logger := log.New(os.Stderr, "[devforge] ", log.LstdFlags|log.Lmsgprefix)
	plugins, err := plugin.LoadAll(ctx, dir, reg, logger)
	if err != nil {
		logger.Printf("plugin scan: %v", err)
	}
	return plugins
}
