package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/devforge/devforge/internal/version"
)

// Server is a thin adapter over mcp-go's MCPServer + StdioServer.
// It wires every Operation in the Registry as an MCP-Tool whose handler
// calls the Operation's Handler.
type Server struct {
	mcp *server.MCPServer
}

// New builds an MCP server from a Registry.
// The returned *Server is ready to Listen on stdio (or any io.Reader/Writer).
func New(reg *Registry) *Server {
	v := version.Get()
	// WithToolCapabilities(false): advertise tools BUT NOT tools/listChanged
	// notifications. DevForge registers tools once at startup (built-ins +
	// any out-of-process plugins discovered via internal/plugin); the
	// catalog is static for the process lifetime. If we ever add runtime
	// plugin hot-reload, flip this to true and call
	// MCPServer.SendNotificationToAllClients("notifications/tools/list_changed").
	mcpSrv := server.NewMCPServer(
		"devforge",
		v.Version,
		server.WithToolCapabilities(false),
	)

	for _, op := range reg.List() {
		op := op // capture
		schema := op.InputSchema
		if len(schema) == 0 {
			schema = json.RawMessage(`{"type":"object","properties":{}}`)
		}
		tool := mcp.NewToolWithRawSchema(op.Name(), op.Description, schema)
		mcpSrv.AddTool(tool, makeHandler(op))
	}

	return &Server{mcp: mcpSrv}
}

func makeHandler(op Operation) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := json.Marshal(req.GetArguments())
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
		}
		raw, err := op.Handler(ctx, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(string(raw)), nil
	}
}

// Listen serves the MCP protocol on the given streams until ctx is cancelled
// or stdin EOFs. Use os.Stdin / os.Stdout for production.
func (s *Server) Listen(ctx context.Context, in io.Reader, out io.Writer) error {
	stdio := server.NewStdioServer(s.mcp)
	return stdio.Listen(ctx, in, out)
}
