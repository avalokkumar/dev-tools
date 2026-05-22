package cli

import (
	"github.com/spf13/cobra"

	"github.com/devforge/devforge/internal/mcpserver"
	"github.com/devforge/devforge/internal/tools"
)

func newMcpCmd(io Stdio) *cobra.Command {
	var pluginDir string
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Run the MCP server over stdio",
		Long: "Starts an MCP (Model Context Protocol) server speaking JSON-RPC " +
			"over stdio. Intended for AI agent clients (Claude Code, Cursor, " +
			"Claude Desktop). Exits when stdin closes.",
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
			srv := mcpserver.New(reg)
			return srv.Listen(cmd.Context(), io.In, io.Out)
		},
	}
	cmd.Flags().StringVar(&pluginDir, "plugin-dir", "",
		"directory to scan for plugins (default $DEVFORGE_PLUGIN_DIR or ~/.devforge/plugins)")
	return cmd
}
