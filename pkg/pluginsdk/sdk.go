package pluginsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// SDKVersion is the manifest schema version this SDK speaks.
// Plugins must keep this in sync with the host's manifest `sdk` field.
const SDKVersion = 1

// OpDecl declares a single Operation a plugin contributes.
type OpDecl struct {
	Tool        string          `json:"tool"`
	Op          string          `json:"op"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
	Handler     OpHandler       `json:"-"`
}

// OpHandler is the user-supplied implementation invoked per call.
type OpHandler func(ctx context.Context, args json.RawMessage) (json.RawMessage, error)

// Plugin describes the plugin and its Operations.
type Plugin struct {
	Name       string
	Version    string
	Operations []OpDecl
}

// Serve runs the plugin's RPC loop on os.Stdin / os.Stdout until stdin closes.
// Logs go to os.Stderr.
func Serve(p Plugin) error { return ServeStreams(p, os.Stdin, os.Stdout) }

// ServeStreams is the streams-explicit variant, for tests.
func ServeStreams(p Plugin, r io.Reader, w io.Writer) error {
	if p.Name == "" {
		return fmt.Errorf("pluginsdk: Plugin.Name required")
	}
	byName := make(map[string]OpHandler, len(p.Operations))
	decls := make([]map[string]any, 0, len(p.Operations))
	for _, o := range p.Operations {
		byName[o.Tool+"_"+o.Op] = o.Handler
		decls = append(decls, map[string]any{
			"tool":        o.Tool,
			"op":          o.Op,
			"description": o.Description,
			"inputSchema": o.InputSchema,
		})
	}
	conn := NewConn(r, w, func(ctx context.Context, f Frame) (json.RawMessage, *FrameError) {
		switch f.Method {
		case "initialize":
			b, err := json.Marshal(map[string]any{
				"plugin":     p.Name,
				"operations": decls,
			})
			if err != nil {
				return nil, &FrameError{Code: -32603, Message: err.Error()}
			}
			return b, nil
		case "invoke":
			var ip struct {
				Name string          `json:"name"`
				Args json.RawMessage `json:"args"`
			}
			if err := json.Unmarshal(f.Params, &ip); err != nil {
				return nil, &FrameError{Code: -32602, Message: "invalid params: " + err.Error()}
			}
			h, ok := byName[ip.Name]
			if !ok {
				return nil, &FrameError{Code: -32601, Message: "unknown operation: " + ip.Name}
			}
			out, err := h(ctx, ip.Args)
			if err != nil {
				return nil, &FrameError{Code: -32000, Message: err.Error()}
			}
			return out, nil
		case "shutdown":
			return json.RawMessage(`null`), nil
		default:
			return nil, &FrameError{Code: -32601, Message: "method not found: " + f.Method}
		}
	})
	return conn.Serve(context.Background())
}

// rawMarshal is a convenience for plugin handlers.
func rawMarshal(v any) (json.RawMessage, error) {
	b, err := json.Marshal(v)
	return b, err
}

// Helper exported for handler authors.
var Marshal = rawMarshal
