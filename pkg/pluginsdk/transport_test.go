package pluginsdk

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

// TestPluginTransport_RoundtripFrame — B9: Conn pair round-trips a request and
// receives the matching response.
func TestPluginTransport_RoundtripFrame(t *testing.T) {
	t.Parallel()

	// "Server" side of the pipe pair receives requests via handler.
	srvIn, cliOut := io.Pipe()  // client writes -> server reads
	cliIn, srvOut := io.Pipe()  // server writes -> client reads

	server := NewConn(srvIn, srvOut, func(_ context.Context, f Frame) (json.RawMessage, *FrameError) {
		if f.Method != "echo" {
			return nil, &FrameError{Code: -32601, Message: "method not found"}
		}
		return f.Params, nil
	})
	client := NewConn(cliIn, cliOut, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go server.Serve(ctx)
	go client.Serve(ctx)

	res, err := client.Call(ctx, "echo", map[string]any{"hello": "world"})
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if !strings.Contains(string(res), `"hello":"world"`) {
		t.Fatalf("res = %s", res)
	}

	cancel()
	server.Close()
	client.Close()
	_ = cliOut.Close()
	_ = srvOut.Close()
}

// TestPluginTransport_UnknownMethod — B9: handler error propagates.
func TestPluginTransport_UnknownMethod(t *testing.T) {
	t.Parallel()
	srvIn, cliOut := io.Pipe()
	cliIn, srvOut := io.Pipe()

	server := NewConn(srvIn, srvOut, func(_ context.Context, _ Frame) (json.RawMessage, *FrameError) {
		return nil, &FrameError{Code: -32601, Message: "method not found"}
	})
	client := NewConn(cliIn, cliOut, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go server.Serve(ctx)
	go client.Serve(ctx)

	_, err := client.Call(ctx, "missing", nil)
	if err == nil || !strings.Contains(err.Error(), "method not found") {
		t.Fatalf("expected method-not-found error, got %v", err)
	}
	cancel()
	server.Close()
	client.Close()
	_ = cliOut.Close()
	_ = srvOut.Close()
}
