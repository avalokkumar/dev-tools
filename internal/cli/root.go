// Package cli implements the cobra command tree for the `devforge` binary.
//
// The root command and all subcommands share a single Stdio so tests can
// inject buffers. The single entry point is Execute, returning a process
// exit code.
package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	clitools "github.com/devforge/devforge/internal/cli/tools"
)

// Execute parses args, runs the matched subcommand, and returns an exit code.
func Execute(ctx context.Context, args []string, io Stdio) int {
	root := newRootCmd(io)
	root.SetArgs(args)
	root.SetIn(io.In)
	root.SetOut(io.Out)
	root.SetErr(io.Err)
	if err := root.ExecuteContext(ctx); err != nil {
		// cobra has already printed the message via SetErr above unless
		// SilenceErrors is set; we leave the default behavior and just
		// translate the error into a non-zero exit code.
		var ec exitCoder
		if errors.As(err, &ec) {
			return ec.ExitCode()
		}
		return 1
	}
	return 0
}

type exitCoder interface{ ExitCode() int }

func newRootCmd(io Stdio) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "devforge",
		Short:         "DevForge — unified developer toolkit",
		Long:          "DevForge bundles daily-driver developer utilities behind a CLI, a local web UI, and an MCP stdio server.",
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	// Global --json flag so any subcommand can render machine-readable output.
	cmd.PersistentFlags().Bool("json", false, "emit machine-readable JSON output where supported")

	cmd.AddCommand(newVersionCmd(io))
	cmd.AddCommand(newMcpCmd(io))
	cmd.AddCommand(newUICmd(io))
	cmd.AddCommand(clitools.NewUUIDCmd())
	cmd.AddCommand(clitools.NewJSONCmd())
	cmd.AddCommand(newRunCmd(io))
	cmd.AddCommand(newUpdateCmd(io))
	cmd.AddCommand(newTelemetryCmd(io))
	cmd.AddCommand(newDocsCmd(io))
	return cmd
}
