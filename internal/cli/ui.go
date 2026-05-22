package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/internal/tools"
	"github.com/devforge/devforge/internal/webserver"
)

func newUICmd(_ Stdio) *cobra.Command {
	var (
		addr      string
		pluginDir string
	)
	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Launch the local web UI",
		Long:  "Starts a localhost-only HTTP server hosting the DevForge web UI.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			reg := mcpserver.NewRegistry()
			if err := tools.Register(reg); err != nil {
				return err
			}
			plugins := loadPluginsInto(cmd.Context(), reg, pluginDir)
			defer func() {
				for _, p := range plugins {
					_ = p.Close()
				}
			}()
			s, err := webserver.New(webserver.Config{Addr: addr, Registry: reg})
			if err != nil {
				return err
			}
			errCh := make(chan error, 1)
			go func() { errCh <- s.ListenAndServe(cmd.Context()) }()

			// Print the bound address once available. ListenAndServe assigns
			// Addr inside the goroutine; race-tolerant short wait via select
			// is acceptable here for a developer-facing CLI.
			waitForAddr(s)
			if a := s.Addr(); a != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "DevForge UI listening on http://%s\n", a)
			}
			return <-errCh
		},
	}
	cmd.Flags().StringVar(&addr, "addr", "127.0.0.1:0", "bind address")
	cmd.Flags().StringVar(&pluginDir, "plugin-dir", "",
		"directory to scan for plugins (default $DEVFORGE_PLUGIN_DIR or ~/.devforge/plugins)")
	return cmd
}

func waitForAddr(s *webserver.Server) {
	for i := 0; i < 50; i++ {
		if s.Addr() != nil {
			return
		}
		// 10ms * 50 = 500ms ceiling.
		sleep10ms()
	}
}
