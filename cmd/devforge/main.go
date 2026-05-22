// Command devforge is the unified developer toolkit binary.
//
// It exposes Tools (utilities) through three Surfaces:
//   - CLI (default invocation)
//   - Web UI (`devforge ui`)
//   - MCP server over stdio (`devforge mcp`)
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/devforge/devforge/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	os.Exit(cli.Execute(ctx, os.Args[1:], cli.Stdio{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}))
}
